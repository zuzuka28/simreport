package analyze

import "time"

type hashImage struct {
	Ahash       string    `json:"ahash"`
	AhashVector []float64 `json:"ahashVector"`
	Dhash       string    `json:"dhash"`
	DhashVector []float64 `json:"dhashVector"`
	Phash       string    `json:"phash"`
	PhashVector []float64 `json:"phashVector"`
	Whash       string    `json:"whash"`
	WhashVector []float64 `json:"whashVector"`
}

type analyzedImage struct {
	ID        string    `json:"id"`
	Vector    []float64 `json:"vector"`
	HashImage hashImage `json:"hashImage"`
}

type analyzedDocument struct {
	ID          string          `json:"id"`
	Text        string          `json:"text"`
	TextVector  []float64       `json:"textVector"`
	Images      []analyzedImage `json:"images"`
	LastUpdated time.Time       `json:"lastUpdated"`
}

type similarityHighlight struct {
	Text []string `json:"text.russian"`
}
