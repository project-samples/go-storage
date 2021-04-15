package app

import (
	"github.com/common-go/log"
	"github.com/common-go/s3"
	"github.com/common-go/storage"
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
