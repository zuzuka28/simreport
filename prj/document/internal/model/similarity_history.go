package model

import "time"

type SimilarityHistory struct {
	Date       time.Time
	DocumentID string
	ID         string
	Matches    []*DocumentSimilarMatch
}

type SimilarityHistoryList = List[*SimilarityHistory]
