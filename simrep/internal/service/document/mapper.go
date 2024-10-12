package document

import "simrep/internal/model"

func mapParsedDocumentFileToParsedDocument(in model.ParsedDocumentFile) model.ParsedDocument {
	imageIDs := make([]string, len(in.Images))
	for i, img := range in.Images {
		imageIDs[i] = img.Sha256
	}

	return model.ParsedDocument{
		ID:          in.ID,
		Sha256:      in.Sha256,
		ImageIDs:    imageIDs,
		TextContent: in.TextContent,
	}
}
