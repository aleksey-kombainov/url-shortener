package http

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"net/http"
	"strings"
)

func IsHeaderContainsMIMETypes(headerValues []string, searchValues []string) bool {
	for _, headerVal := range headerValues {
		for _, searchVal := range searchValues {
			if strings.Contains(headerVal, searchVal) {
				return true
			}
		}
	}
	return false
}

func httpError(res http.ResponseWriter, errStr string) {
	logger.Logger.Error().
		Msg("http error: " + errStr)
	http.Error(res, errStr, http.StatusBadRequest)
}
