package document

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/zuzuka28/simreport/lib/elasticutil"
)

//nolint:gochecknoglobals
var now = time.Now

func mapDocumentToInternal(src model.Document) document {
	return document{
		ParentID:    src.ParentID,
		Name:        src.Name,
		Version:     src.Version,
		GroupID:     src.GroupID,
		SourceID:    src.SourceID,
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
		ParentID:    in.ParentID,
		Name:        in.Name,
		LastUpdated: in.LastUpdated,
		Version:     in.Version,
		GroupID:     in.GroupID,
		SourceID:    in.SourceID,
		TextID:      in.TextID,
		ImageIDs:    in.ImageIDs,
		WithContent: false,
		Source:      model.File{}, //nolint:exhaustruct
		Text:        model.File{}, //nolint:exhaustruct
		Images:      nil,
	}
}

func buildSearchQuery(query model.DocumentSearchQuery) ([]byte, error) {
	searchQuery := make(map_)

	var filter []map_

	if query.ParentID != "" {
		filter = append(filter, map_{
			"term": map_{
				"id.keyword": query.ParentID,
			},
		})
	}

	if query.Version != "" {
		filter = append(filter, map_{
			"term": map_{
				"version": query.Version,
			},
		})
	}

	if len(query.GroupID) > 0 {
		filter = append(filter, map_{
			"terms": map_{
				"groupID": query.GroupID,
			},
		})
	}

	if len(filter) > 0 {
		searchQuery["query"] = map_{
			"bool": map_{
				"filter": filter,
			},
		}
	}

	searchQuery["sort"] = map_{
		"lastUpdated": map_{
			"order": "desc",
		},
	}

	m, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	return m, nil
}

func mapErrorToModel(err error) error {
	switch {
	case errors.Is(err, elasticutil.ErrInvalid):
		return errors.Join(err, model.ErrInvalid)

	case errors.Is(err, elasticutil.ErrNotFound):
		return errors.Join(err, model.ErrNotFound)

	default:
		return err
	}
}
