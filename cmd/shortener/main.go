package main

import (
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"io"
	"net/http"
	"strings"
)

const (
	routeUri        = `/`
	errorHttpCode   = http.StatusBadRequest
	shortenerResult = "http://localhost:8080/EwHXdJfB"
	expanderResult  = "https://practicum.yandex.ru/"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(routeUri, routerHandler)

	return http.ListenAndServe(`:8080`, mux)
}

func routerHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		shortenerHandler(res, req)
	} else if req.Method == http.MethodGet {
		expanderHandler(res, req)
	} else {
		http.Error(res, "bad request", errorHttpCode)
	}
}

func shortenerHandler(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get(headers.ContentType) != mimetype.TextPlain {
		http.Error(res, fmt.Sprintf("Content-type \"%s\" not allowed", req.Header.Get("Content-Type")), errorHttpCode)
		return
	}
	// @todo
	//url, err := io.ReadAll(req.Body)
	_, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), errorHttpCode)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortenerResult))
}

func expanderHandler(res http.ResponseWriter, req *http.Request) {
	shortcut := strings.TrimPrefix(req.RequestURI, routeUri)
	if len(shortcut) == 0 {
		http.Error(res, "invalid shortcut", errorHttpCode)
		return
	}
	res.Header().Add(headers.Location, expanderResult)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
