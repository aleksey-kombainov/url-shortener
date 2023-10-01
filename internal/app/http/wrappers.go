package http

import (
	"compress/gzip"
	"io"
	"net/http"
)

type responseWriter struct {
	writer     http.ResponseWriter
	length     int
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		writer:     w,
		length:     0,
		statusCode: 0,
	}
}

func (w *responseWriter) Write(b []byte) (n int, err error) {
	n, err = w.writer.Write(b)

	w.length += n

	return
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.writer.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

//////////////////

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	writer http.ResponseWriter
	zw     *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		writer: w,
		zw:     gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.writer.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.writer.Header().Set("Content-Encoding", "gzip")
	}
	c.writer.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
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
