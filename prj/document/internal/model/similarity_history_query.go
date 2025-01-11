package model

import "time"

type SimilarityHistoryQuery struct {
	DocumentID string

	Limit  int
	Offset int

	DateFrom time.Time
	DateTo   time.Time
}
