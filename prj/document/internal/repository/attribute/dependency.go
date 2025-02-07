package attribute

type (
	Metrics interface {
		IncAttributeRepositoryRequests(op string, status string, dur float64)
	}
)
