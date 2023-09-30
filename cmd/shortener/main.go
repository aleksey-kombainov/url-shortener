package main

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	nethttp "net/http"
)

var options config.Options
var loggerInstance zerolog.Logger

func main() {
	options = config.GetOptions()
	loggerInstance = logger.Init()

	loggerInstance.Info().Msg("Starting server")

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
	mux.Post("/", http.RequestLoggerMiddleware(http.ShortenerHandler, &loggerInstance))
	mux.Get("/{shortcut}", http.RequestLoggerMiddleware(http.ExpanderHandler, &loggerInstance))
	mux.NotFound(http.ErrorHandler)
	mux.MethodNotAllowed(http.ErrorHandler)

	return mux
}
