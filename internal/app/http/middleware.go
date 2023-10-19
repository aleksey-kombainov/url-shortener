package http

import (
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
	"time"
)

// @todo разбить на две middleware: encoder & logger. проблема: подменяем responsewriter в одном и используем подмену в другом
const (
	acceptableEncodingValue = "gzip"
)

var ValidEncodedContentTypesForShortener []string = []string{
	mimetype.ApplicationGzip, "application/x-gzip",
}

type LoggerEncoderMiddleware struct {
	logger *zerolog.Logger
}

func NewLoggerEncoderMiddleware(logger *zerolog.Logger) *LoggerEncoderMiddleware {
	return &LoggerEncoderMiddleware{logger: logger}
}

func (m LoggerEncoderMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := request.Header.Get(headers.AcceptEncoding)
		supportsGzip := strings.Contains(acceptEncoding, acceptableEncodingValue)
		rw := newResponseWriter(respWriter, supportsGzip)
		defer func() {
			if err := rw.Close(); err != nil {
				m.logger.Error().Msg("can't close response writer")
			}
		}()

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
			defer func() {
				if err := cr.Close(); err != nil {
					m.logger.Error().Msg("can't close response writer")
				}
			}()
		}

		m.logger.Info().
			Str("uri", request.RequestURI).
			Str("method", request.Method).
			Send()

		start := time.Now()
		next.ServeHTTP(rw, request)
		execTimeMicroseconds := time.Since(start).Microseconds()

		m.logger.Info().
			Int64("executionTimeMicroseconds", execTimeMicroseconds).
			Int("responseLength", rw.length).
			Int("responseCode", rw.statusCode).
			Send()
	})
}
