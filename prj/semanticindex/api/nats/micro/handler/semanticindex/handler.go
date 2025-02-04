package semanticindex

type Handler struct {
	s  Service
	ds DocumentService
}

func NewHandler(
	s Service,
	ds DocumentService,
) *Handler {
	return &Handler{
		s:  s,
		ds: ds,
	}
}
