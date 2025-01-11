package model

type DocumentProcessingStatus string

const (
	DocumentProcessingStatusFileSaved        DocumentProcessingStatus = "file_saved"
	DocumentProcessingStatusDocumentSaved    DocumentProcessingStatus = "document_parsed"
	DocumentProcessingStatusDocumentAnalyzed DocumentProcessingStatus = "document_analyzed"
	DocumentProcessingStatusNotFound         DocumentProcessingStatus = "not_found"
)

type DocumentStatus struct {
	ID     string
	Status DocumentProcessingStatus
}
