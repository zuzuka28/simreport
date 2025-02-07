package metrics

import (
	"io"
	"net/http"
)

type statusResponseWriter interface {
	http.ResponseWriter
	Size() int
	Status() int
	Written() bool
}

type statusRW struct {
	http.ResponseWriter
	pendingStatus int
	status        int
	size          int
}

func newStatusRW(rw http.ResponseWriter) statusResponseWriter {
	return &statusRW{
		ResponseWriter: rw,
		pendingStatus:  0,
		status:         0,
		size:           0,
	}
}

func (rw *statusRW) WriteHeader(s int) {
	if rw.Written() {
		return
	}

	rw.pendingStatus = s
	rw.status = s

	rw.ResponseWriter.WriteHeader(s)
}

func (rw *statusRW) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}

	size, err := rw.ResponseWriter.Write(b)
	rw.size += size

	return size, err //nolint:wrapcheck
}

func (rw *statusRW) ReadFrom(r io.Reader) (int64, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}

	n, err := io.Copy(rw.ResponseWriter, r)
	rw.size += int(n)

	return n, err //nolint:wrapcheck
}

func (rw *statusRW) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func (rw *statusRW) Status() int {
	if rw.Written() {
		return rw.status
	}

	return rw.pendingStatus
}

func (rw *statusRW) Size() int {
	return rw.size
}

func (rw *statusRW) Written() bool {
	return rw.status != 0
}
