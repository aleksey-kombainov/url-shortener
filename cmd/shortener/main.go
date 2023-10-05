package main

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage"
	nethttp "net/http"
)

var options config.Options

func main() {
	initGlobals()
	defer storage.ShortcutStorage.Close()
	logger.Logger.Info().Msg("Starting server")

	if err := run(); err != nil {
		panic(err)
	}
}

func initGlobals() {
	options = config.GetOptions()
	logger.Init()
	logger.Logger.Info().Msg("Shutting down")
	storage.ShortcutStorageFactoryInit(options, &logger.Logger)
}

func run() error {
	mux := http.GetRouter()

	return nethttp.ListenAndServe(options.ServerListenAddr, mux)
}
