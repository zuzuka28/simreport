package document

import "time"

type document struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Version     int       `json:"version"`
	GroupID     []string  `json:"groupIDs"`
	SourceID    string    `json:"sourceID"`
	ImageIDs    []string  `json:"imageIDs"`
	TextID      string    `json:"textID"`
	LastUpdated time.Time `json:"lastUpdated"`
}
