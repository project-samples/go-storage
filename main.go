package main

import (
	"fmt"
	s3 "github.com/common-go/s3"
	"github.com/common-go/storage"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func main() {
	cloud, _ := s3.NewS3ServiceWithConfig( s3.Config{Region: "", AccessKeyID: "", SecretAccessKey: "", Token: ""},
	 									   storage.Config{CredentialsFile: "", BucketName: "", SubDirectory: "", PermissionFileRoleAll: true})
	r := mux.NewRouter()
	handler := storage.UploadFileHandler{CloudService: cloud, KeyFile: "file", LogError: nil}
	r.HandleFunc("/upload", handler.UploadFile).Methods("POST")
	r.HandleFunc("/delete/{id}", handler.DeleteFile).Methods("DELETE")
	http.Handle("/", r)
	server := ":" + strconv.Itoa(8080)
	fmt.Println(server + " started")
	http.ListenAndServe(server, nil)
}
