package document

import "time"

type document struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ImageIDs    []string  `json:"imageIDs"`
	TextID      string    `json:"textID"`
	LastUpdated time.Time `json:"lastUpdated"`
}
