package app

import (
	"context"
	"github.com/core-go/log"
	"github.com/core-go/storage"
	"github.com/core-go/storage/google"
	"github.com/core-go/storage/s3"
)

type ApplicationContext struct {
	FileHandler *storage.FileHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	logError := log.ErrorMsg
	storageService, err := CreateStorageService(ctx, root)
	if err != nil {
		return nil, err
	}

	fileHandler := storage.NewFileHandler(storageService, root.KeyFile, logError)

	return &ApplicationContext{FileHandler: fileHandler}, nil
}

func CreateStorageService(ctx context.Context, root Root) (storage.StorageService, error) {
	if root.Provider == "google" {
		return google.NewGoogleStorageServiceWithCredentials(ctx, []byte(root.GoogleCredentials), root.Storage)
	} else {
		return s3.NewS3ServiceWithConfig(root.AWS, root.Storage)
	}
}