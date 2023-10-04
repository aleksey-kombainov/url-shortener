package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Options struct {
	ServerListenAddr string `env:"SERVER_ADDRESS"`
	BaseURL          string `env:"BASE_URL"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
}

var defaultOptions = Options{
	ServerListenAddr: ":8080",
	BaseURL:          "http://localhost:8080",
	FileStoragePath:  "/tmp/short-url-db.json",
}

var envOptions = new(Options)
var options = new(Options)

func init() {
	flag.StringVar(&options.ServerListenAddr, "a", defaultOptions.ServerListenAddr, "server listen address")
	flag.StringVar(&options.BaseURL, "b", defaultOptions.BaseURL, "url for shortcuts")
	flag.StringVar(&options.FileStoragePath, "f", defaultOptions.FileStoragePath, "file storage path")
}

func GetOptions() Options {
	flag.Parse()

	if err := env.Parse(&envOptions); err != nil {
		if envOptions.ServerListenAddr != "" {
			options.ServerListenAddr = envOptions.ServerListenAddr
		}
		if envOptions.BaseURL != "" {
			options.BaseURL = envOptions.BaseURL
		}
		if envOptions.FileStoragePath != "" {
			options.FileStoragePath = envOptions.FileStoragePath
		}
	}
	return *options
}
