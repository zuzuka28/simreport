package model

import "time"

type File struct {
	Content     []byte
	Sha256      string
	LastUpdated time.Time
}
