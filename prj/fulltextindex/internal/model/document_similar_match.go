package model

type DocumentSimilarMatch struct {
	ID            string
	Rate          float64
	Highlights    []string
	SimilarImages []string
}
