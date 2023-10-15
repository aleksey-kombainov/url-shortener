package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
)

type Options struct {
	ServerListenAddr string `env:"SERVER_ADDRESS"`
	BaseURL          string `env:"BASE_URL"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn      string `env:"DATABASE_DSN"`
}

var defaultOptions = Options{
	ServerListenAddr: ":8080",
	BaseURL:          "http://localhost:8080",
	FileStoragePath:  "/tmp/short-url-db.json",
	//DatabaseDsn:      "postgres://shortener_user:12345@localhost:5432/shortener",
	DatabaseDsn: "postgres://shortener_user:12345@localhost:5432/shortener",
}

func GetOptions(logger zerolog.Logger) (opts Options) {
	flag.StringVar(&opts.ServerListenAddr, "a", defaultOptions.ServerListenAddr, "server listen address")
	flag.StringVar(&opts.BaseURL, "b", defaultOptions.BaseURL, "url for shortcuts")
	flag.StringVar(&opts.FileStoragePath, "f", defaultOptions.FileStoragePath, "file storage path")
	flag.StringVar(&opts.DatabaseDsn, "d", defaultOptions.DatabaseDsn, "db dsn")
	flag.Parse()

	envOptions := Options{}
	if err := env.Parse(&envOptions); err == nil {
		if envOptions.ServerListenAddr != "" {
			opts.ServerListenAddr = envOptions.ServerListenAddr
		}
		if envOptions.BaseURL != "" {
			opts.BaseURL = envOptions.BaseURL
		}
		if envOptions.FileStoragePath != "" {
			opts.FileStoragePath = envOptions.FileStoragePath
		}
		if envOptions.DatabaseDsn != "" {
			opts.DatabaseDsn = envOptions.DatabaseDsn
		}
	} else {
		logger.Error().Msg("Can't parse env vars: " + err.Error())
	}
	return opts
}
