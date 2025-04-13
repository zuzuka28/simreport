package document

type Handler struct {
	s  Service
	ss StatusService
}

func NewHandler(
	s Service,
	ss StatusService,
) *Handler {
	return &Handler{
		s:  s,
		ss: ss,
	}
}
