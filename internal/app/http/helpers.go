package http

import (
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/logger"
	"net/http"
	"strings"
)

const (
	ErrorHTTPCode = http.StatusBadRequest
)

func ExtractMIMETypeFromStr(str string) string {
	mtypeSlice := strings.Split(str, ";")
	return strings.TrimSpace(mtypeSlice[0])
}

func httpError(res http.ResponseWriter, errStr string) {
	logger.Logger.Error().
		Msg("http error: " + errStr)
	http.Error(res, errStr, ErrorHTTPCode)
}
