package vectorizer

type (
	Metrics interface {
		IncVectorizerRequests(op string, status string, dur float64)
	}
)
