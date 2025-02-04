package semanticindex

type Service struct {
	r  Repository
	vs VectorizerService
}

func NewService(
	r Repository,
	vs VectorizerService,
) *Service {
	return &Service{
		r:  r,
		vs: vs,
	}
}
