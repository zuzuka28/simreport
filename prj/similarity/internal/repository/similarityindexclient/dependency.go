package similarityindexclient

type (
	Metrics interface {
		IncSimilarityIndexRequests(index string, op string, status string, dur float64)
	}
)
