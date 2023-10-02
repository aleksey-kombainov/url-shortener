package http

import (
	"compress/gzip"
	"github.com/go-http-utils/headers"
	"github.com/ldez/mimetype"
	"io"
	"net/http"
)

type responseWriter struct {
	writer                    http.ResponseWriter
	gzWriter                  *gzip.Writer
	length                    int
	statusCode                int
	clientSupportsCompression bool
	useCompression            bool
}

func newResponseWriter(w http.ResponseWriter, clientSupportsCompression bool) *responseWriter {
	return &responseWriter{
		writer:                    w,
		gzWriter:                  gzip.NewWriter(w),
		length:                    0,
		statusCode:                0,
		clientSupportsCompression: clientSupportsCompression,
		useCompression:            false,
	}
}

func (w *responseWriter) Write(b []byte) (n int, err error) {

	if w.clientSupportsCompression && w.length == 0 {
		conType := w.writer.Header().Get(headers.ContentType)
		if conType == mimetype.TextHTML || conType == mimetype.ApplicationJSON {
			w.writer.Header().Set(headers.ContentEncoding, acceptableEncodingValue)
			w.useCompression = true
		}
	}
	if w.useCompression {
		n, err = w.gzWriter.Write(b)
	} else {
		n, err = w.writer.Write(b)
	}
	w.length += n
	return
}

func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.writer.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (w *responseWriter) Close() error {
	if w.useCompression {
		return w.gzWriter.Close()
	}
	return nil
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
