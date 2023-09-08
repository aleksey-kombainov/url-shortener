package main

import (
	"errors"
	"fmt"
	"github.com/aleksey-kombainov/url-shortener.git/pkg/memstorage"
	"github.com/aleksey-kombainov/url-shortener.git/pkg/random"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"io"
	"net/http"
	"strings"
)

const (
	routeURI                = `/`
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

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(routeURI, routerHandler)

	return http.ListenAndServe(`:8080`, mux)
}

func routerHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		shortenerHandler(res, req)
	} else if req.Method == http.MethodGet {
		expanderHandler(res, req)
	} else {
		http.Error(res, "bad request", errorHTTPCode)
	}
}

func shortenerHandler(res http.ResponseWriter, req *http.Request) {

	mtypeSlice := strings.Split(req.Header.Get(headers.ContentType), ";")
	mtype := strings.TrimSpace(mtypeSlice[0])
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

	shortcut, err := getAndSaveUniqueShortcut(urlStr)
	if err != nil {
		http.Error(res, err.Error(), errorHTTPCode)
		return
	}

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("http://localhost:8080/" + shortcut))
}

func expanderHandler(res http.ResponseWriter, req *http.Request) {
	shortcut := strings.TrimPrefix(req.RequestURI, routeURI)
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
