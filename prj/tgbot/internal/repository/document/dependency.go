package document

type (
	Metrics interface {
		IncDocumentRepositoryRequests(op string, status string, dur float64)
	}
)
