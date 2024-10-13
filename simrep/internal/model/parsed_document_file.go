package model

type ParsedDocumentFile struct {
	ID     string
	Name   string
	Source DocumentFile

	Images      []MediaFile
	TextContent string
}
