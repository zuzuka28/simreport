package document

import "time"

type document struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ImageIDs    []string  `json:"imageIDs"`
	TextContent string    `json:"textContent"`
	LastUpdated time.Time `json:"lastUpdated"`
}
