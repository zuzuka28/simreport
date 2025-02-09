package fulltextindex

import (
	"time"
)

const requestTimeout = 60 * time.Second

type Handler struct {
	s  Service
	ds DocumentService
	fs Filestorage
}

func NewHandler(
	s Service,
	ds DocumentService,
	fs Filestorage,
) *Handler {
	return &Handler{
		s:  s,
		ds: ds,
		fs: fs,
	}
}
