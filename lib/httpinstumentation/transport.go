package httpinstumentation

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type LogAttributeProvider func(ctx context.Context) []any

type InstumentedTransport struct {
	*http.Transport
	ExtractAttrs    LogAttributeProvider
	LogRequestBody  bool
	LogResponseBody bool
}

//revive:disable:cyclomatic
func (it *InstumentedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t := time.Now()

	attrs := []any{
		"request_url", req.URL.String(),
	}

	if it.ExtractAttrs != nil {
		attrs = append(attrs, it.ExtractAttrs(req.Context())...)
	}

	if it.LogRequestBody && req != nil && req.Body != nil && req.Body != http.NoBody {
		body, reqbody, _ := duplicateBody(req.Body)
		req.Body = reqbody

		var buf bytes.Buffer

		_, err := io.Copy(&buf, body)
		if err != nil {
			return nil, fmt.Errorf("log request body: %w", err)
		}

		attrs = append(attrs, "request_body", buf.String())
	}

	res, err := it.Transport.RoundTrip(req)

	attrs = append(attrs,
		"elapsed_time", time.Since(t),
		"response_status", res.StatusCode,
		"response_size", res.Body,
	)

	if it.LogResponseBody && res != nil && res.Body != nil && res.Body != http.NoBody {
		defer res.Body.Close()
		body, reqbody, _ := duplicateBody(req.Body)
		res.Body = reqbody

		var buf bytes.Buffer

		_, err := io.Copy(&buf, body)
		if err != nil {
			return nil, fmt.Errorf("log response body: %w", err)
		}

		attrs = append(attrs, "response_body", buf.String())
	}

	if err != nil {
		attrs = append(attrs, "error", err)
	}

	slog.Info(
		"send request",
		attrs...,
	)

	return res, err //nolint:wrapcheck
}

func duplicateBody(body io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	var (
		b1 bytes.Buffer
		b2 bytes.Buffer
		tr = io.TeeReader(body, &b2)
	)

	_, err := b1.ReadFrom(tr)
	if err != nil {
		return io.NopCloser(io.MultiReader(&b1, errorReader{err: err})),
			io.NopCloser(io.MultiReader(&b2, errorReader{err: err})), err
	}

	defer func() { _ = body.Close() }()

	return io.NopCloser(&b1), io.NopCloser(&b2), nil
}

type errorReader struct{ err error }

func (r errorReader) Read([]byte) (int, error) { return 0, r.err }
