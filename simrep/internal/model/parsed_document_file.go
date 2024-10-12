package model

import "time"

type ParsedDocumentFile struct {
	ID string

	Sha256      string
	RawContent  []byte
	LastUpdated time.Time

	Images      []MediaFile
	TextContent string
}
