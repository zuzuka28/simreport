package metrics

type (
	Metrics interface {
		IncNatsMicroRequest(op string, status string, size int, dur float64)
	}
)
