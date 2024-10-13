package document

import "simrep/internal/model"

func mapParsedDocumentFileToDocument(in model.ParsedDocumentFile) model.Document {
	imageIDs := make([]string, len(in.Images))
	for i, img := range in.Images {
		imageIDs[i] = img.Sha256
	}

	return model.Document{
		ID:          in.ID,
		Name:        in.Name,
		ImageIDs:    imageIDs,
		TextContent: in.TextContent,
		LastUpdated: in.Source.LastUpdated,
	}
}
