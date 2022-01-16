package app

import (
	"context"
	"github.com/core-go/storage"
	"github.com/core-go/storage/dropbox"
	"github.com/core-go/storage/google"
	"github.com/core-go/storage/google-drive"
	"github.com/core-go/storage/one-drive"
	"github.com/core-go/storage/s3"

	"go-service/internal/handler"
)

type ApplicationContext struct {
	FileHandler *handler.FileHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	cloudService, err := CreateCloudService(ctx, root)
	if err != nil {
		return nil, err
	}

	fileHandler := handler.NewFileHandler(cloudService, root.Provider, root.GeneralDirectory, root.KeyFile, root.Storage.Directory)

	return &ApplicationContext{FileHandler: fileHandler}, nil
}

func CreateCloudService(ctx context.Context, root Root) (storage.StorageService, error) {
	if root.Provider == "google-drive" {
		return google_drive.NewGoogleDriveService([]byte(root.GoogleDriveCredentials))
	} else if root.Provider == "drop-box" {
		return drop_box.NewDropboxService(root.DropboxToken)
	} else if root.Provider == "one-drive" {
		return one_drive.NewOneDriveService(ctx, root.OneDriveToken)
	} else if root.Provider == "google-storage" {
		return google.NewGoogleStorageServiceWithCredentials(ctx, []byte(root.GoogleCredentials), root.Storage)
	} else {
		return s3.NewS3ServiceWithConfig(root.AWS, root.Storage)
	}
}
