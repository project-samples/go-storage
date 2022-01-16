package app

import (
	"context"
	"github.com/gorilla/mux"
)

const (
	POST   = "POST"
	DELETE = "DELETE"
)

func Route(r *mux.Router, ctx context.Context, conf Root) error {
	app, err := NewApp(ctx, conf)
	if err != nil {
		return err
	}

	r.HandleFunc("/upload", app.FileHandler.UploadFile).Methods(POST)
	r.HandleFunc("/delete/{id}", app.FileHandler.DeleteFile).Methods(DELETE)

	return nil
}