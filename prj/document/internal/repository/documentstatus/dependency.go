package documentstatus

type (
	Metrics interface {
		IncDocumentStatusRepositoryUpdates(status string, result string)
	}
)
