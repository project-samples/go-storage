package app

import (
	"github.com/core-go/log"
	"github.com/core-go/storage"
	"github.com/core-go/storage/s3"
)

type Root struct {
	Server  ServerConfig   `mapstructure:"server"`
	Log     log.Config     `mapstructure:"log"`
	AWS     s3.Config      `mapstructure:"aws"`
	Storage storage.Config `mapstructure:"storage"`
	KeyFile string         `mapstructure:"key_file"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}
