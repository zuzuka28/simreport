package model

type Document struct {
	ID              string
	Images          []*Image
	SbertTextVector []float32
	Sha256          string
	SourceBytes     string
	TextContent     string
}
