package fulltextindex

type (
	Metrics interface {
		IncFulltextIndexRequests(op string, status string, dur float64)
	}
)
