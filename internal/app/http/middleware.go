package http

import (
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

func RequestLoggerMiddleware(handler http.HandlerFunc, logger *zerolog.Logger) http.HandlerFunc {

	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		rw := &responseWriterWithLength{responseWriter, 0, 0}

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
