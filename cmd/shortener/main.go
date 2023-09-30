package main

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/go-chi/chi/v5"
	nethttp "net/http"
)

var options config.Options

func main() {
	options = config.GetOptions()
	logger.Init()

	logger.Logger.Info().Msg("Starting server")

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := getRouter()

	return nethttp.ListenAndServe(options.ServerListenAddr, mux)
}

func getRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Post("/", http.RequestLoggerMiddleware(http.ShortenerHandler, &logger.Logger))
	mux.Get("/{shortcut}", http.RequestLoggerMiddleware(http.ExpanderHandler, &logger.Logger))
	mux.NotFound(http.ErrorHandler)
	mux.MethodNotAllowed(http.ErrorHandler)

	return mux
}
