package main

import (
	"errors"
	"fmt"
	"github.com/aleksey-kombainov/url-shortener.git/pkg/memstorage"
	"github.com/aleksey-kombainov/url-shortener.git/pkg/random"
	"github.com/go-chi/chi/v5"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"io"
	"net/http"
	"strings"
)

const (
	errorHTTPCode           = http.StatusBadRequest
	shortcutLength          = 8
	generatorIterationLimit = 10000
)

type Storager interface {
	Put(key string, val string)
	GetValueByKey(key string) (string, error)
	GetKeyByValue(val string) (string, error)
}

var storage Storager = memstorage.NewStorage()
var um *urlManager = newURLManager("http", "localhost:8080", "/")

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := getRouter()

	return http.ListenAndServe(um.host, mux)
}

func getRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Post("/", func(res http.ResponseWriter, req *http.Request) {
		shortenerHandler(res, req)
	})
	mux.Get("/{shortcut}", func(res http.ResponseWriter, req *http.Request) {
		expanderHandler(res, req)
	})
	mux.NotFound(routerErrorHandler)
	mux.MethodNotAllowed(routerErrorHandler)

	return mux
}

func routerErrorHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "", errorHTTPCode)
}

func shortenerHandler(res http.ResponseWriter, req *http.Request) {

	mtype := extractMIMETypeFromStr(req.Header.Get(headers.ContentType))
	if mtype != mimetype.TextPlain {
		http.Error(res, fmt.Sprintf("Content-type \"%s\" not allowed", mtype), errorHTTPCode)
		return
	}
	defer req.Body.Close()
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

	shortcut, err := getAndSaveUniqueShortcut(urlStr)
	if err != nil {
		http.Error(res, err.Error(), errorHTTPCode)
		return
	}
	res.Header().Add(headers.ContentType, mimetype.TextPlain)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(um.buildFullURLByShortcut(shortcut)))
}

func expanderHandler(res http.ResponseWriter, req *http.Request) {
	shortcut := um.getShortcutFromURI(req.RequestURI)
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

/*
один урл генерирует разные шорткаты. из задания не понятно, верно ли такое поведение
*/
func getAndSaveUniqueShortcut(url string) (string, error) {
	var shortcut string
	isGenerated := false
	for i := 0; i < generatorIterationLimit; i++ {
		shortcut = random.GenString(shortcutLength)
		if _, err := storage.GetValueByKey(shortcut); err != nil {
			isGenerated = true
			break
		}
	}
	if !isGenerated {
		return "", errors.New("generator limit exceeded")
	}
	storage.Put(shortcut, url)
	return shortcut, nil
}

func extractMIMETypeFromStr(str string) string {
	mtypeSlice := strings.Split(str, ";")
	return strings.TrimSpace(mtypeSlice[0])
}
