package similarity

type (
	Metrics interface {
		IncSimilarityRepositoryRequests(op string, status string, dur float64)
	}
)
