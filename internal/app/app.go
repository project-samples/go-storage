package app

import (
	"context"

	"github.com/common-go/log"
	"github.com/common-go/storage"
	"github.com/common-go/storage/s3"
)

type ApplicationContext struct {
	FileHandler *storage.FileHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	logError := log.ErrorMsg

	s3Service, err := s3.NewS3ServiceWithConfig(root.AWS, root.Storage)
	if err != nil {
		return nil, err
	}

	fileHandler := storage.NewFileHandler(s3Service, root.KeyFile, logError)

	return &ApplicationContext{FileHandler: fileHandler}, nil
}
