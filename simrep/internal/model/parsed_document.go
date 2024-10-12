package model

type ParsedDocument struct {
	ID          string
	Sha256      string
	ImageIDs    []string
	TextContent string
}
