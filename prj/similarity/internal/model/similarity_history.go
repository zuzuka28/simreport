package model

import "time"

type SimilarityHistory struct {
	Date       time.Time
	DocumentID string
	ID         string
	Matches    []*SimilarityMatch
}

type SimilarityHistoryList = List[*SimilarityHistory]
