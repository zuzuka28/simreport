package analyzehistory

import "time"

type history struct {
	Date       time.Time              `json:"date"`
	DocumentID string                 `json:"documentID"`
	ID         string                 `json:"id"`
	Matches    []*documentSimilarMatch `json:"matches"`
}

type documentSimilarMatch struct {
	ID            string   `json:"id"`
	Rate          float64  `json:"rate"`
	Highlights    []string `json:"highlights"`
	SimilarImages []string `json:"similarImages"`
}
