package userstate

type (
	Metrics interface {
		IncUserStateRepositoryRequests(op string, status string, dur float64)
	}
)
