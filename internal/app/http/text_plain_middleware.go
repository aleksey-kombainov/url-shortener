package http

import (
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/rs/zerolog"
	"net/http"
)

type TextPlainMiddleware struct {
	logger *zerolog.Logger
}

func NewTextPlainMiddleware(logger *zerolog.Logger) *TextPlainMiddleware {
	return &TextPlainMiddleware{logger: logger}
}

func (m TextPlainMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		if vc := append(ValidEncodedContentTypesForShortener, mimetype.TextPlain); !IsHeaderContainsMIMETypes(request.Header.Values(headers.ContentType), vc) {
			http.Error(respWriter, "Content-type not allowed", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(respWriter, request)
	})
}
