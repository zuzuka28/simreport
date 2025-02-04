package model

type Document struct {
	ID     string
	Text   []byte
	Vector Vector
}
