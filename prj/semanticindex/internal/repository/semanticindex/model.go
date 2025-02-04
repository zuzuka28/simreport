package semanticindex

import "time"

//nolint:revive
type map_ map[string]any

type document struct {
	ID          string    `json:"id"`
	TextVector  []float64 `json:"text_vector"`
	LastUpdated time.Time `json:"lastUpdated"`
}
