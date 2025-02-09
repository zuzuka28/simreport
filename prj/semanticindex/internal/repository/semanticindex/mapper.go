package semanticindex

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/semanticindex/internal/model"
)

//nolint:gochecknoglobals
var now = time.Now

func mapDocumentToInternal(in model.Document) document {
	return document{
		ID:          in.SourceID,
		TextVector:  in.Vector,
		LastUpdated: now(),
	}
}

func mapSearchResponseToMatches(
	in *elasticutil.SearchResponse,
) []model.DocumentSimilarMatch {
	matches := make([]model.DocumentSimilarMatch, 0, len(in.Hits.Hits))

	for _, match := range in.Hits.Hits {
		matches = append(matches, model.DocumentSimilarMatch{
			ID:            match.ID,
			Rate:          match.Score,
			Highlights:    nil,
			SimilarImages: nil,
		})
	}

	return matches
}

func buildSearchQuery(
	query model.DocumentSimilarQuery,
) ([]byte, error) {
	q := map_{
		"knn": map_{
			"field":        "text_vector",
			"query_vector": query.Item.Vector,
			"k":            100, //nolint:gomnd,mnd
		},
	}

	m, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	return m, nil
}
