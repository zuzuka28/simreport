package fulltextindex

import "time"

const (
	metricsError = "error"
)

//nolint:revive
type map_ map[string]any

type analyzedDocument struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type similarityHighlight struct {
	TextRu []string `json:"text.russian"`
	TextEn []string `json:"text.english"`
}
