package main

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/config"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/http"
	"github.com/go-chi/chi/v5"
	nethttp "net/http"
)

var Options config.Options

func main() {
	Options = config.GetOptions()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := getRouter()

	return nethttp.ListenAndServe(Options.ServerListenAddr, mux)
}

func getRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Post("/", func(res nethttp.ResponseWriter, req *nethttp.Request) {
		http.ShortenerHandler(res, req)
	})
	mux.Get("/{shortcut}", func(res nethttp.ResponseWriter, req *nethttp.Request) {
		http.ExpanderHandler(res, req)
	})
	mux.NotFound(http.ErrorHandler)
	mux.MethodNotAllowed(http.ErrorHandler)

	return mux
}
