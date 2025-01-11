package model

type DocumentStatusUpdateCommand struct {
	ID     string
	Status DocumentProcessingStatus
}
