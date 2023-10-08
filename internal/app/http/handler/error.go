package handler

import (
	"github.com/rs/zerolog"
	"net/http"
)

type ErrorHandler struct {
	logger *zerolog.Logger
}

func NewErrorHandler(logger *zerolog.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (h ErrorHandler) ServeHTTP(res http.ResponseWriter, _ *http.Request) {
	http.Error(res, "", http.StatusBadRequest)
}
