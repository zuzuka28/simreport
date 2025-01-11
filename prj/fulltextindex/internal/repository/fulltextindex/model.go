package fulltextindex

import "time"

type analyzedDocument struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type similarityHighlight struct {
	Text []string `json:"text.russian"`
}
