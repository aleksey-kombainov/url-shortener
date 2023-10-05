package http

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage/postgres"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"io"
	"net/http"
	"time"
)

func ShortenerHandler(res http.ResponseWriter, req *http.Request) {

	if vc := append(validEncodedContentTypesForShortener, mimetype.TextPlain); !IsHeaderContainsMIMETypes(req.Header.Values(headers.ContentType), vc) {
		httpError(res, "Content-type not allowed", 0)
		return
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			logger.Logger.Error().
				Msg("Can not close request.Body(): " + err.Error())
		}
	}()
	url, err := io.ReadAll(req.Body)
	if err != nil {
		httpError(res, err.Error(), 0)
		return
	}
	shortcut, err := app.MakeShortcut(string(url))
	if err != nil {
		httpError(res, err.Error(), 0)
		return
	}
	res.Header().Add(headers.ContentType, mimetype.TextPlain)
	res.WriteHeader(http.StatusCreated)
	if _, err := res.Write([]byte(NewURLManagerFromFullURL(config.GetOptions().BaseURL).BuildFullURLByShortcut(shortcut))); err != nil {
		logger.Logger.Error().
			Msg("Can not Write response: " + err.Error())
	}
}

func ExpanderHandler(res http.ResponseWriter, req *http.Request) {
	shortcut := NewURLManagerFromFullURL(config.GetOptions().BaseURL).GetShortcutFromURI(req.RequestURI)
	if len(shortcut) == 0 {
		httpError(res, "invalid shortcut", 0)
		return
	}
	url, err := storage.ShortcutStorage.GetOriginalURLByShortcut(shortcut)
	if err != nil {
		httpError(res, "shortcut not found", 0)
		return
	}
	res.Header().Add(headers.Location, url) // @todo проверить редирект на самого себя
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func PingHandler(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := postgres.GetConnection(ctx, config.GetOptions().DatabaseDsn, logger.Logger)
	if err != nil {
		httpError(res, err.Error(), http.StatusInternalServerError)
	}

	switch ctx.Err() {
	case context.Canceled:
		httpError(res, "canceled", http.StatusInternalServerError)
	case context.DeadlineExceeded:
		httpError(res, "connection timeout", http.StatusInternalServerError)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = conn.Ping(ctx); err != nil {
		logger.Logger.Error().Msg("unable to ping: " + err.Error())
		httpError(res, err.Error(), http.StatusInternalServerError)
	}

	switch ctx.Err() {
	case context.Canceled:
		httpError(res, "canceled", http.StatusInternalServerError)
	case context.DeadlineExceeded:
		httpError(res, "connection timeout", http.StatusInternalServerError)
	}
}

func ErrorHandler(res http.ResponseWriter, _ *http.Request) {
	httpError(res, "", 0)
}
