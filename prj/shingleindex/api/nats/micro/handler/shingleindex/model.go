package shingleindex

type documentSimilarMatch struct {
	ID            string
	Rate          float64
	Highlights    []string
	SimilarImages []string
}
