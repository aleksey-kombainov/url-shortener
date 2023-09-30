package http

import "net/http"

type responseWriterWithLength struct {
	http.ResponseWriter
	length     int
	statusCode int
}

func (w *responseWriterWithLength) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)

	w.length += n

	return
}

func (w *responseWriterWithLength) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}
