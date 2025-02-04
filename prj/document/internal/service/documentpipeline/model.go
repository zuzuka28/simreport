package documentpipeline

import "github.com/zuzuka28/simreport/prj/document/internal/model"

type Stage struct {
	Trigger model.DocumentProcessingStatus
	Action  Handler
	Next    model.DocumentProcessingStatus
}
