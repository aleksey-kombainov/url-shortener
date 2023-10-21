package http

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http/handler"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func GetRouter(logger *zerolog.Logger, shortcutService *app.ShortcutService, urlService *app.URLManagerService) *chi.Mux {
	authMiddleware := NewAuthMiddleware(logger).Handler

	mux := chi.NewRouter()

	mux.Use(NewLoggerEncoderMiddleware(logger).Handler)

	textPlainMiddleware := NewTextPlainMiddleware(logger).Handler

	mux.Get("/{shortcut}", handler.NewExpanderHandler(logger, shortcutService, urlService).ServeHTTP)

	mux.With(textPlainMiddleware).With(authMiddleware).Post("/", handler.NewShortenerHandler(logger, shortcutService, urlService).ServeHTTP)
	mux.Get("/ping", handler.NewPingHandler(logger, *shortcutService.Storage).ServeHTTP)

	mux.Route("/api", func(r chi.Router) {
		r.Use(NewAPIMiddleware(logger).Handler)
		r.Use(authMiddleware)
		r.Post("/shorten", handler.NewShortenerAPIHandler(logger, shortcutService, urlService).ServeHTTP)
		r.Post("/shorten/batch", handler.NewShortenerBatchAPIHandler(logger, shortcutService, urlService).ServeHTTP)
		r.Get("/user/urls", handler.NewUserURLsAPIHandler(logger, shortcutService, urlService).ServeHTTP)
		r.Delete("/user/urls", handler.NewDeleteBatchAPIHandler(logger, shortcutService, urlService).ServeHTTP)
	})
	//errHandler := handler.NewErrorHandler(logger).ServeHTTP
	//mux.NotFound(errHandler)
	//mux.MethodNotAllowed(errHandler)

	return mux
}
