package analyzehistory

import "time"

//nolint:revive
type map_ map[string]any

type history struct {
	Date       time.Time               `json:"date"`
	DocumentID string                  `json:"github.com/zuzuka28/simreport/prj/similarityID"`
	ID         string                  `json:"id"`
	Matches    []*documentSimilarMatch `json:"matches"`
}

type documentSimilarMatch struct {
	ID            string   `json:"id"`
	Rate          float64  `json:"rate"`
	Highlights    []string `json:"highlights"`
	SimilarImages []string `json:"similarImages"`
}
