package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/core-go/storage"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

const contentTypeHeader = "Content-Type"

type FileHandler struct {
	Service	storage.StorageService
	Provider string
	GeneralDirectory string
	Directory string
	KeyFile	string
}

func NewFileHandler(service storage.StorageService, provider string, generalDirectory string, keyFile string, directory string) *FileHandler {
	return &FileHandler{Service: service, Provider: provider, GeneralDirectory: generalDirectory, KeyFile: keyFile, Directory: directory}
}

func (f FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "not available", http.StatusInternalServerError)
		return
	}

	// FormFile returns the first file for the provided form key
	file, handler, err0 := r.FormFile(f.KeyFile)
	if err0 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// NewBuffer creates and initializes a new Buffer using buf as its initial contents.
	// The new Buffer takes ownership of buf, and the caller should not use buf after this call.
	// NewBuffer is intended to prepare a Buffer to read existing data
	bufferFile := bytes.NewBuffer(nil)
	_, err1 := io.Copy(bufferFile, file)
	if err1 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// Get gets the first value associated with the given key
	bytes := bufferFile.Bytes()
	contentType := handler.Header.Get(contentTypeHeader)
	if len(contentType) == 0 {
		contentType = getExt(handler.Filename)
	}

	var directory string
	if f.Provider == "google-storage" {
		directory = f.Directory
	} else {
		// google-drive or drop_box
		directory = f.GeneralDirectory
	}

	rs, err2 := f.Service.Upload(r.Context(), directory, handler.Filename, bytes, contentType)
	if err2 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, rs)
}

func (f FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	i := strings.LastIndex(r.RequestURI, "/")
	filename := ""
	if i <= 0 {
		http.Error(w, "require id", http.StatusBadRequest)
		return
	}
	filename = r.RequestURI[i+1:]
/*
	var directory string
	if f.Provider == "google-storage" {
		directory = f.Directory
	} else {
		// google-drive or drop_box
		directory = f.GeneralDirectory
	}
*/
	rs, err := f.Service.Delete(r.Context(), filename)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var msg string
	if rs {
		msg = fmt.Sprintf("file '%s' has been deleted successfully!!!", filename)
	} else {
		msg = fmt.Sprintf("delete file '%s' failed!!!", filename)
	}
	respond(w, http.StatusOK, msg)
}

func respond(w http.ResponseWriter, code int, result interface{}) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getExt(file string) string {
	ext := filepath.Ext(file)
	if strings.HasPrefix(ext, ":") {
		ext = ext[1:]
		return ext
	}
	return ext
}
