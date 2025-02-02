package model

type DocumentSaveCommand struct {
	Item Document
}

type DocumentSaveResult struct {
	ID string
}
