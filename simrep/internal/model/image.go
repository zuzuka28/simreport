package model

type Image struct {
	ClipImageVector []float32
	Fname           string
	Hashes          *HashImage
	Sha256          string
	SourceBytes     string
}
