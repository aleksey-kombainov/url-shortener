package config

import (
	"flag"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
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

var options = new(Options)

func init() {
	flag.StringVar(&options.ServerListenAddr, "a", defaultOptions.ServerListenAddr, "server listen address")
	flag.StringVar(&options.BaseURL, "b", defaultOptions.BaseURL, "url for shortcuts")
	flag.StringVar(&options.FileStoragePath, "f", defaultOptions.FileStoragePath, "file storage path")
}

func GetOptions() Options {
	flag.Parse()

	envOptions := Options{}
	if err := env.Parse(&envOptions); err == nil {
		if envOptions.ServerListenAddr != "" {
			options.ServerListenAddr = envOptions.ServerListenAddr
		}
		if envOptions.BaseURL != "" {
			options.BaseURL = envOptions.BaseURL
		}
		if envOptions.FileStoragePath != "" {
			options.FileStoragePath = envOptions.FileStoragePath
		}
	} else {
		logger.Logger.Error().Msg("Can't parse env vars: " + err.Error())
	}
	return *options
}
