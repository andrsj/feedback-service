package middlewares

import (
	"bytes"
	"fmt"
	"net/http"

)

// Custom ResponeWriter for handling data
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
