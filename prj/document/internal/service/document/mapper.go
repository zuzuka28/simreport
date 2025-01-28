package document

import "github.com/zuzuka28/simreport/prj/document/internal/model"

func mapDocumentWithContentToDocument(in model.Document) model.Document {
	return model.Document{
		ParentID:    in.ParentID,
		Name:        in.Name,
		LastUpdated: in.Source.LastUpdated,
		Version:     in.Version,
		GroupID:     in.GroupID,
		SourceID:    in.SourceID,
		TextID:      in.TextID,
		ImageIDs:    in.ImageIDs,
		WithContent: false,
		Source:      model.File{}, //nolint:exhaustruct
		Text:        model.File{}, //nolint:exhaustruct
		Images:      []model.File{},
	}
}
