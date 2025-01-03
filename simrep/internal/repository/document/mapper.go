package document

import (
	"encoding/json"
	"fmt"
	"simrep/internal/model"
	"simrep/pkg/elasticutil"
	"time"
)

//nolint:gochecknoglobals
var now = time.Now

func mapDocumentToInternal(src model.Document) document {
	return document{
		ID:          src.ID,
		Name:        src.Name,
		ImageIDs:    src.ImageIDs,
		TextID:      src.TextID,
		LastUpdated: src.LastUpdated,
	}
}

func parseDocument(hit *elasticutil.Hit) (model.Document, error) {
	var doc document

	if err := json.Unmarshal(hit.Source, &doc); err != nil {
		return model.Document{}, fmt.Errorf("unmarshal document: %w", err)
	}

	return mapDocument(doc), nil
}

func mapSearchResponseToDocuments(in *elasticutil.SearchResponse) ([]model.Document, error) {
	items := make([]model.Document, 0, len(in.Hits.Hits))

	for _, hit := range in.Hits.Hits {
		item, err := parseDocument(&hit)
		if err != nil {
			return nil, fmt.Errorf("parse document: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func mapDocument(in document) model.Document {
	return model.Document{
		ID:          in.ID,
		Name:        in.Name,
		LastUpdated: in.LastUpdated,
		ImageIDs:    in.ImageIDs,
		TextID:      in.TextID,
		WithContent: false,
		Source:      model.File{}, //nolint:exhaustruct
		Text:        model.File{}, //nolint:exhaustruct
		Images:      nil,
	}
}

func buildSearchQuery(query model.DocumentSearchQuery) ([]byte, error) {
	searchQuery := make(map_)

	if query.Name != "" {
		searchQuery = map_{
			"query": map_{
				"query_string": map_{
					"query": fmt.Sprintf("name: %s*", query.Name),
				},
			},
		}
	}

	m, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	return m, nil
}
