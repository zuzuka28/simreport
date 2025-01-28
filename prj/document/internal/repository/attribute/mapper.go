package attribute

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/model"

	"github.com/zuzuka28/simreport/lib/elasticutil"
)

func mapSearchResponseToAttributes(in *elasticutil.SearchResponse) ([]model.Attribute, error) {
	var doc attributeAggs

	if err := json.Unmarshal(in.Aggs, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal aggs: %w", err)
	}

	items := make([]model.Attribute, 0, len(doc.Attr.Buckets))

	for _, hit := range doc.Attr.Buckets {
		items = append(items, model.Attribute{
			Value: hit.Key,
			Label: hit.Key,
		})
	}

	return items, nil
}

func buildSearchQuery(query model.AttributeQuery) ([]byte, error) {
	searchQuery := make(map_)

	searchQuery["aggs"] = map_{
		"attr": map_{
			"terms": map_{
				"field": query.ID + ".keyword",
			},
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
