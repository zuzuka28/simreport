package document

import "simrep/internal/model"

func mapParsedDocumentToInternal(src model.ParsedDocument) parsedDocument {
	return parsedDocument{
		ID:          src.ID,
		Sha256:      src.Sha256,
		ImageIDs:    src.ImageIDs,
		TextContent: src.TextContent,
	}
}
