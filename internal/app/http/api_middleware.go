package http

import (
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/rs/zerolog"
	"net/http"
)

type APIMiddleware struct {
	logger *zerolog.Logger
}

func NewAPIMiddleware(logger *zerolog.Logger) *APIMiddleware {
	return &APIMiddleware{logger: logger}
}

func (m APIMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			if vc := append(ValidEncodedContentTypesForShortener, mimetype.ApplicationJSON); !IsHeaderContainsMIMETypes(request.Header.Values(headers.ContentType), vc) {
				http.Error(respWriter, "Content-type not allowed", http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(respWriter, request)
	})
}
