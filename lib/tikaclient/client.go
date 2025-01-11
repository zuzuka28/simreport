package tikaclient

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	cli   *http.Client
	url   string
	langs []string
}

func New(
	cli *http.Client,
	url string,
) *Client {
	return &Client{
		cli:   cli,
		url:   url,
		langs: OCRLanguage,
	}
}

func (c *Client) Parse(ctx context.Context, r io.Reader) (Result, error) {
	header := http.Header{
		"X-Tika-OCRLanguage": []string{strings.Join(c.langs, "+")},
	}

	body, err := c.call(ctx, http.MethodPut, "/tika", header, r)
	if err != nil {
		return Result{}, fmt.Errorf("call api: %w", err)
	}

	b := &bytes.Buffer{}
	if _, err := io.Copy(b, body); err != nil {
		return Result{}, fmt.Errorf("read body: %w", err)
	}

	return Result{
		Content: b.Bytes(),
		Sha256:  sha256hash(b.Bytes()),
		Name:    "",
	}, nil
}

func (c *Client) ParseText(ctx context.Context, r io.Reader) (Result, error) {
	header := http.Header{
		"X-Tika-OCRLanguage": []string{strings.Join(c.langs, "+")},
		"Accept":             []string{"text/plain"},
	}

	body, err := c.call(ctx, http.MethodPut, "/tika", header, r)
	if err != nil {
		return Result{}, fmt.Errorf("call api: %w", err)
	}

	b := &bytes.Buffer{}
	if _, err := io.Copy(b, body); err != nil {
		return Result{}, fmt.Errorf("read body: %w", err)
	}

	return Result{
		Content: b.Bytes(),
		Sha256:  sha256hash(b.Bytes()),
		Name:    "",
	}, nil
}

func (c *Client) ParseEmbedded(ctx context.Context, r io.Reader) ([]Result, error) {
	header := http.Header{
		"X-Tika-OCRLanguage": []string{strings.Join(c.langs, "+")},
		"Accept":             []string{"application/zip"},
	}

	body, err := c.call(ctx, http.MethodPut, "/unpack", header, r)
	if err != nil {
		var cerr ClientError
		if errors.As(err, &cerr); cerr.StatusCode == http.StatusNoContent {
			return nil, nil
		}

		return nil, fmt.Errorf("call api: %w", err)
	}

	b := &bytes.Buffer{}
	if _, err := io.Copy(b, body); err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		return nil, fmt.Errorf("new zip reader: %w", err)
	}

	embedded := make([]Result, 0, len(zr.File))

	for _, file := range zr.File {
		fb, err := readZipFile(file)
		if err != nil {
			return nil, fmt.Errorf("read zipped file: %w", err)
		}

		embedded = append(embedded, Result{
			Content: fb,
			Sha256:  sha256hash(fb),
			Name:    "",
		})
	}

	return embedded, nil
}

func (c *Client) call(
	ctx context.Context,
	method, path string,
	header http.Header,
	input io.Reader,
) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.url+path, input)
	if err != nil {
		return nil, fmt.Errorf("new request with context: %w", err)
	}

	req.Header = header

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, ClientError{resp.StatusCode}
	}

	return resp.Body, nil
}

func readZipFile(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("open zip file: %w", err)
	}

	defer rc.Close()

	res, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read content: %w", err)
	}

	return res, nil
}

func sha256hash(in []byte) string {
	h := sha256.New()
	_, _ = h.Write(in)

	return hex.EncodeToString(h.Sum([]byte(nil)))
}
