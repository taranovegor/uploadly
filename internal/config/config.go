package config

import (
	"github.com/taranovegor/uploadly/pkg/config"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	Version        = "0.2.0"
	StaticHttpPath = "static"
)

type Config struct {
	Debug    bool `yaml:"debug"`
	Database struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"database"`
	Storage struct {
		StaticPath string `yaml:"static_path"`
	} `yaml:"storage"`
	Http struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"http"`
	FileContext config.FileContexts `yaml:"file_context"`
}

var appConfig Config

func IsDebug() bool {
	return appConfig.Debug
}

func Init() Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&appConfig); err != nil {
		panic(err)
	}

	return appConfig
}
