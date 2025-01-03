package model

import "time"

type Document struct {
	ID          string
	Name        string
	LastUpdated time.Time

	ImageIDs []string
	TextID   string

	WithContent bool
	Source      File
	Text        File
	Images      []File
}
