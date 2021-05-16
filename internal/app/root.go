package app

import (
	"github.com/core-go/log"
	"github.com/core-go/storage"
	"github.com/core-go/storage/s3"
)

type Root struct {
	Server            ServerConfig   `mapstructure:"server"`
	Log               log.Config     `mapstructure:"log"`
	KeyFile           string         `mapstructure:"key_file"`
	Storage           storage.Config `mapstructure:"storage"`
	Provider          string         `mapstructure:"provider"`
	GoogleCredentials string         `mapstructure:"google_credentials"`
	AWS               s3.Config      `mapstructure:"aws"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port *int64 `mapstructure:"port"`
}
