package semanticindexclient

type documentSimilarMatch struct {
	ID            string   `json:"id"`
	Rate          float64  `json:"rate"`
	Highlights    []string `json:"highlights"`
	SimilarImages []string `json:"similar_images"`
}
