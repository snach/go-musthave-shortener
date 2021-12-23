package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"storage.txt"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func MakeConf() (ServerConfig, error) {
	var conf ServerConfig
	if err := env.Parse(&conf); err != nil {
		return ServerConfig{}, err
	}

	address := flag.String("a", conf.ServerAddress, "address server (or env var SERVER_ADDRESS)")
	fileStoragePath := flag.String("f", conf.FileStoragePath, "path to storage file (or env var FILE_STORAGE_PATH)")
	baseURL := flag.String("b", conf.BaseURL, "base url ajh shortened link (or env var BASE_URL)")
	flag.Parse()

	return ServerConfig{
		ServerAddress:   *address,
		FileStoragePath: *fileStoragePath,
		BaseURL:         *baseURL,
	}, nil
}
