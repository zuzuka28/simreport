package semanticindex

type (
	Metrics interface {
		IncSemanticIndexRequests(op string, status string, dur float64)
	}
)
