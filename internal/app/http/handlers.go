package http

import (
	"fmt"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/memstorage"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"io"
	"net/http"
	"strings"
)

const (
	errorHTTPCode = http.StatusBadRequest
)

var storage memstorage.Storager = memstorage.NewStorage()

func ShortenerHandler(res http.ResponseWriter, req *http.Request) {

	mtype := ExtractMIMETypeFromStr(req.Header.Get(headers.ContentType))
	if mtype != mimetype.TextPlain {
		http.Error(res, fmt.Sprintf("Content-type \"%s\" not allowed", mtype), errorHTTPCode)
		return
	}
	defer func() {
		if err := req.Body.Close(); err != nil {

		}
	}()
	url, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), errorHTTPCode)
		return
	}
	urlStr := strings.TrimSpace(string(url)) // @todo валидация url
	if urlStr == "" {
		http.Error(res, "empty url", errorHTTPCode)
		return
	}

	shortcut, err := app.GetAndSaveUniqueShortcut(urlStr, storage)
	if err != nil {
		http.Error(res, err.Error(), errorHTTPCode)
		return
	}
	res.Header().Add(headers.ContentType, mimetype.TextPlain)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(NewURLManagerFromFullURL(config.GetOptions().BaseURL).BuildFullURLByShortcut(shortcut)))
}

func ExpanderHandler(res http.ResponseWriter, req *http.Request) {
	shortcut := NewURLManagerFromFullURL(config.GetOptions().BaseURL).GetShortcutFromURI(req.RequestURI)
	if len(shortcut) == 0 {
		http.Error(res, "invalid shortcut", errorHTTPCode)
		return
	}
	url, err := storage.GetValueByKey(shortcut)
	if err != nil {
		http.Error(res, "shortcut not found", errorHTTPCode)
		return
	}
	res.Header().Add(headers.Location, url) // @todo проверить редирект на самого себя
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func ErrorHandler(res http.ResponseWriter, _ *http.Request) {
	http.Error(res, "", errorHTTPCode)
}
