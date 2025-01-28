package document

import "time"

const requestTimeout = 60 * time.Second

type Handler struct {
	s Service
}

func NewHandler(
	s Service,
) *Handler {
	return &Handler{
		s: s,
	}
}
