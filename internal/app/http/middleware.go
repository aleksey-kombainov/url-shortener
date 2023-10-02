package http

import (
	"github.com/go-http-utils/headers"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
	"time"
)

const (
	acceptableEncodingValue = "gzip"
)

func RequestLoggerMiddleware(handler http.HandlerFunc, logger *zerolog.Logger) http.HandlerFunc {

	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := request.Header.Get(headers.AcceptEncoding)
		supportsGzip := strings.Contains(acceptEncoding, acceptableEncodingValue)
		rw := newResponseWriter(respWriter, supportsGzip)
		defer rw.Close()

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := request.Header.Get(headers.ContentEncoding)
		sendsGzip := strings.Contains(contentEncoding, acceptableEncodingValue)
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(request.Body)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			request.Body = cr
			defer cr.Close()
		}

		logger.Info().
			Str("uri", request.RequestURI).
			Str("method", request.Method).
			Send()

		start := time.Now()
		handler(rw, request)
		execTimeMicroseconds := time.Since(start).Microseconds()

		logger.Info().
			Int64("executionTimeMicroseconds", execTimeMicroseconds).
			Int("responseLength", rw.length).
			Int("responseCode", rw.statusCode).
			Send()
	})
}
