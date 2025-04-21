package elasticstorage

import "errors"

var errBadStatusCode = errors.New("bad status code")

type set map[string]struct{}

type storedDocumentSource struct {
	ID      string   `json:"id"`
	Content []string `json:"content"`
}

type fetchResponse struct {
	Source storedDocumentSource `json:"_source"`
}
