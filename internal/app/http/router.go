package http

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/go-chi/chi/v5"
)

func GetRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Post("/", RequestLoggerMiddleware(ShortenerHandler, &logger.Logger))
	mux.Post("/api/shorten", RequestLoggerMiddleware(ShortenerAPIHandler, &logger.Logger))
	mux.Get("/{shortcut}", RequestLoggerMiddleware(ExpanderHandler, &logger.Logger))
	mux.Get("/ping", RequestLoggerMiddleware(PingHandler, &logger.Logger))
	mux.NotFound(ErrorHandler)
	mux.MethodNotAllowed(ErrorHandler)

	return mux
}
