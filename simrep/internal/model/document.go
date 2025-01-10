package model

import "time"

type Document struct {
	ID          string
	Name        string
	LastUpdated time.Time

	Version int
	GroupID []string

	SourceID string
	TextID   string
	ImageIDs []string

	WithContent bool
	Source      File
	Text        File
	Images      []File
}
