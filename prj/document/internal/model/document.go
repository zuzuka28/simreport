package model

import (
	"strconv"
	"time"
)

type Document struct {
	ParentID    string
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

func (d *Document) ID() string {
	return d.ParentID + "_" + strconv.Itoa(d.Version)
}
