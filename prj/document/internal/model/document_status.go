package model

type DocumentProcessingStatus string

const (
	DocumentProcessingStatusFileSaved     DocumentProcessingStatus = "file_saved"
	DocumentProcessingStatusDocumentSaved DocumentProcessingStatus = "document_parsed"
	DocumentProcessingStatusNotFound      DocumentProcessingStatus = "not_found"
)

type DocumentStatus struct {
	ID     string
	Status DocumentProcessingStatus
}
