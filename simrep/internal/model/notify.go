package model

type NotifyAction string

const (
	NotifyActionFileSaved        NotifyAction = "file_saved"
	NotifyActionDocumentSaved    NotifyAction = "document_parsed"
	NotifyActionDocumentAnalyzed NotifyAction = "document_analyzed"
)
