package filestorage

type (
	Metrics interface {
		IncFilestorageRequests(op string, status string, dur float64)
		IncFilestorageUploads(status string)
	}
)
