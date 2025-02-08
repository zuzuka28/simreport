package analyzehistory

type (
	Metrics interface {
		IncAnalyzeHistoryRepositoryRequests(op string, status string, dur float64)
	}
)
