package model

type AnalyzedDocumentMatch struct {
	ID            string
	Rate          float64
	Highlights    []string
	SimilarImages []string
}
