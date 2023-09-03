package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Options struct {
	ServerListenAddr string `env:"SERVER_ADDRESS"`
	BaseUrl          string `env:"BASE_URL"`
}

var defaultOptions = Options{
	ServerListenAddr: ":8080",
	BaseUrl:          "http://localhost:8000",
}

var envOptions = new(Options)
var options = new(Options)

func init() {
	flag.StringVar(&options.ServerListenAddr, "a", defaultOptions.ServerListenAddr, "server listen address")
	flag.StringVar(&options.BaseUrl, "b", defaultOptions.BaseUrl, "url for shortcuts")
}

func GetOptions() Options {
	flag.Parse()

	if err := env.Parse(&envOptions); err != nil {
		if envOptions.ServerListenAddr != "" {
			options.ServerListenAddr = envOptions.ServerListenAddr
		}
		if envOptions.BaseUrl != "" {
			options.BaseUrl = envOptions.BaseUrl
		}
	}
	return *options
}
