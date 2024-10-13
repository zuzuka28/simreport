package model

import "time"

type Document struct {
	ID          string
	Name        string
	ImageIDs    []string
	TextContent string
	LastUpdated time.Time
}