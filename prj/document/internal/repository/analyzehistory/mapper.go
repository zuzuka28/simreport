package analyzehistory

import (
	"encoding/json"
	"errors"
	"fmt"
	"document/internal/model"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"time"
)

//nolint:gochecknoglobals
var now = time.Now

func parseHistory(hit *elasticutil.Hit) (*model.SimilarityHistory, error) {
	var doc history

	if err := json.Unmarshal(hit.Source, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal document: %w", err)
	}

	return mapHistoryToModel(&doc), nil
}

func mapSearchResponseToHistoryList(in *elasticutil.SearchResponse) (*model.SimilarityHistoryList, error) {
	items := make([]*model.SimilarityHistory, 0, len(in.Hits.Hits))

	for _, hit := range in.Hits.Hits {
		item, err := parseHistory(&hit)
		if err != nil {
			return nil, fmt.Errorf("parse document: %w", err)
		}

		items = append(items, item)
	}

	return &model.SimilarityHistoryList{
		Count: in.Hits.Total.Value,
		Items: items,
	}, nil
}

func mapHistoryToModel(in *history) *model.SimilarityHistory {
	if in == nil {
		return nil
	}

	matches := make([]*model.DocumentSimilarMatch, 0, len(in.Matches))
	for _, v := range in.Matches {
		matches = append(matches, mapMatchToModel(v))
	}

	return &model.SimilarityHistory{
		Date:       in.Date,
		DocumentID: in.DocumentID,
		ID:         in.ID,
		Matches:    matches,
	}
}

func mapMatchToModel(in *documentSimilarMatch) *model.DocumentSimilarMatch {
	if in == nil {
		return nil
	}

	return &model.DocumentSimilarMatch{
		ID:            in.ID,
		Rate:          in.Rate,
		Highlights:    in.Highlights,
		SimilarImages: in.SimilarImages,
	}
}

func mapHistoryToInternal(in model.SimilarityHistory) history {
	matches := make([]*documentSimilarMatch, 0, len(in.Matches))
	for _, v := range in.Matches {
		matches = append(matches, mapMatchToInternal(v))
	}

	return history{
		Date:       in.Date,
		DocumentID: in.DocumentID,
		ID:         in.ID,
		Matches:    matches,
	}
}

func mapMatchToInternal(in *model.DocumentSimilarMatch) *documentSimilarMatch {
	if in == nil {
		return nil
	}

	return &documentSimilarMatch{
		ID:            in.ID,
		Rate:          in.Rate,
		Highlights:    in.Highlights,
		SimilarImages: in.SimilarImages,
	}
}

func buildSearchQuery(query model.SimilarityHistoryQuery) ([]byte, error) {
	searchQuery := make(map_)

	searchQuery["size"] = query.Limit
	searchQuery["from"] = query.Offset

	searchQuery["sort"] = map_{
		"date": map_{
			"order": "desc",
		},
	}

	var filter []map_

	if query.DocumentID != "" {
		filter = append(filter, map_{
			"term": map_{
				"documentID": query.DocumentID,
			},
		})
	}

	if !query.DateFrom.IsZero() || !query.DateTo.IsZero() {
		rangeConstraints := make(map_)

		if !query.DateFrom.IsZero() {
			rangeConstraints["gte"] = query.DateFrom
		}

		if !query.DateTo.IsZero() {
			rangeConstraints["lt"] = query.DateTo
		}

		filter = append(filter, map_{"range": map_{"date": rangeConstraints}})
	}

	searchQuery["query"] = map_{
		"bool": map_{
			"filter": filter,
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
