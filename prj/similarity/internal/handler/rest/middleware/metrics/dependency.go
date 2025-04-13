package metrics

type (
	Metrics interface {
		IncHTTPRequest(op string, status string, size int, dur float64)
	}
)
