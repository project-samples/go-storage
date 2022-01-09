package service

import (
	"context"
)

type CloudService interface {
	Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, directory string, fileName string) (bool, error)
}