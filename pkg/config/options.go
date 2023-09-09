package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Options struct {
	ServerListenAddr string `env:"SERVER_ADDRESS"`
	BaseURL          string `env:"BASE_URL"`
}

var defaultOptions = Options{
	ServerListenAddr: ":8080",
	BaseURL:          "http://localhost:8080",
}

var envOptions = new(Options)
var options = new(Options)

func init() {
	flag.StringVar(&options.ServerListenAddr, "a", defaultOptions.ServerListenAddr, "server listen address")
	flag.StringVar(&options.BaseURL, "b", defaultOptions.BaseURL, "url for shortcuts")
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
	}
	return *options
}
