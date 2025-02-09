package model

type Document struct {
	ID       string
	SourceID string // used to deduplicate saves
	TextID   string // used to fetch file stream from fs
	Text     []byte // FIXME: use io.Reader
	Vector   Vector
}
