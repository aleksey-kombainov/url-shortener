package http

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/storage"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"io"
	"net/http"
)

func ShortenerHandler(res http.ResponseWriter, req *http.Request) {

	if !IsHeaderContainsMIMEType(req.Header.Values(headers.ContentType), mimetype.TextPlain) {
		httpError(res, "Content-type not allowed")
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
		httpError(res, err.Error())
		return
	}
	shortcut, err := app.MakeShortcut(string(url))
	if err != nil {
		httpError(res, err.Error())
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
		httpError(res, "invalid shortcut")
		return
	}
	url, err := storage.ShortcutStorage.GetOriginalURLByShortcut(shortcut)
	if err != nil {
		httpError(res, "shortcut not found")
		return
	}
	res.Header().Add(headers.Location, url) // @todo проверить редирект на самого себя
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func ErrorHandler(res http.ResponseWriter, _ *http.Request) {
	httpError(res, "")
}
