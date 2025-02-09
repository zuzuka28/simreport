package model

import "time"

type File struct {
	Name        string
	Content     []byte
	Sha256      string
	LastUpdated time.Time
}
