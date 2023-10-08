package http

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http/handler"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func GetRouter(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(NewLoggerEncoderMiddleware(logger).Handler)

	mux.Route("/", func(r chi.Router) {
		r.Use(NewTextPlainMiddleware(logger).Handler)

		mux.Post("/", handler.NewShortenerHandler(logger, shortcutService, urlService).ServeHTTP)

		mux.Get("/{shortcut}", handler.NewExpanderHandler(logger, shortcutService, urlService).ServeHTTP)
		mux.Get("/ping", handler.NewPingHandler(logger, shortcutService.Storage).ServeHTTP)
	})

	mux.Route("/api", func(r chi.Router) {
		r.Use(NewAPIMiddleware(logger).Handler)
		r.Post("/shorten", handler.NewShortenerAPIHandler(logger, shortcutService, urlService).ServeHTTP)
	})

	errHandler := handler.NewErrorHandler(logger).ServeHTTP
	mux.NotFound(errHandler)
	mux.MethodNotAllowed(errHandler)

	return mux
}
