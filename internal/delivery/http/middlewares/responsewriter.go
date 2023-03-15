package middlewares

import (
	"bytes"
	"fmt"
	"net/http"

)

// ChatGPT's generated package for handling Response for Cache Middleware

// Custom ResponseWriter for handling data
// after the next.ServeHTTP is done.
type ResponseWriter struct {
	http.ResponseWriter
	status int
	Body   *bytes.Buffer
}

func NewResponseWriter(w http.ResponseWriter, status int) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		status:         status,
		Body:           bytes.NewBuffer(nil),
	}
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.status = status
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.Body.Write(b)

	return rw.ResponseWriter.Write(b) //nolint:wrapcheck
}

// `interfacer` says that http.ResponseWrite can be replaced by io.Writer.
//nolint:interfacer
func (rw *ResponseWriter) WriteTo(w http.ResponseWriter) error {
	_, err := w.Write(rw.Body.Bytes())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (rw *ResponseWriter) Status() int {
	return rw.status
}
