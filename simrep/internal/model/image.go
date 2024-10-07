package model

type ImageHashes struct {
	Ahash       string
	AhashVector []float32
	Dhash       string
	DhashVector []float32
	Phash       string
	PhashVector []float32
	Whash       string
	WhashVector []float32
}

type Image struct {
	ClipImageVector []float32
	Fname           string
	Hashes          *ImageHashes
	Sha256          string
	SourceBytes     string
}
