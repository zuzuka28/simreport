package model

type SimilarityMatch struct {
	ID            string
	Rate          float64
	Highlights    []string
	SimilarImages []string
}
