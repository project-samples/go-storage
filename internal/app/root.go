package app

import (
	"github.com/core-go/storage"
	"github.com/core-go/storage/s3"
)

type Root struct {
	Server                 ServerConfig   `mapstructure:"server"`
	AWS                    s3.Config      `mapstructure:"aws"`
	KeyFile                string         `mapstructure:"key_file"`
	Storage                storage.Config `mapstructure:"google_storage"`
	Provider               string         `mapstructure:"provider"`
	DropboxToken           string         `mapstructure:"dropbox_token"`
	OneDriveToken          string         `mapstructure:"one_drive_token"`
	GeneralDirectory       string         `mapstructure:"general_directory"`
	GoogleCredentials      string         `mapstructure:"google_credentials"`
	GoogleDriveCredentials string         `mapstructure:"google_drive_credentials"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port *int64 `mapstructure:"port"`
}
