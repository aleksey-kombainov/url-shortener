package http

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"net/http"
	"strings"
)

const (
	ErrorHTTPCode = http.StatusBadRequest
)

func IsHeaderContainsMIMEType(headerValues []string, str string) bool {
	for _, val := range headerValues {
		if strings.Contains(val, str) {
			return true
		}
	}
	return false
}

func httpError(res http.ResponseWriter, errStr string) {
	logger.Logger.Error().
		Msg("http error: " + errStr)
	http.Error(res, errStr, ErrorHTTPCode)
}
