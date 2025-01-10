package model

type FileSaveManyCommand struct {
	Bucket string
	Items  []File
}
