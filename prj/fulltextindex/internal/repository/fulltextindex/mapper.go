package fulltextindex

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
)

//nolint:gochecknoglobals
var now = time.Now

func parseDocument(hit *elasticutil.Hit) (analyzedDocument, error) {
	var doc analyzedDocument

	if err := json.Unmarshal(hit.Source, &doc); err != nil {
		return analyzedDocument{}, fmt.Errorf("unmarshal analyzed document: %w", err)
	}

	return doc, nil
}

func mapDocumentToInternal(in model.Document) analyzedDocument {
	return analyzedDocument{
		ID:          in.SourceID,
		Text:        string(in.Text),
		LastUpdated: now(),
	}
}

func mapSearchResponseToMatches(
	in *elasticutil.SearchResponse,
) ([]model.DocumentSimilarMatch, error) {
	docs := make([]analyzedDocument, 0, len(in.Hits.Hits))

	for _, hit := range in.Hits.Hits {
		doc, err := parseDocument(&hit)
		if err != nil {
			return nil, fmt.Errorf("parse document: %w", err)
		}

		docs = append(docs, doc)
	}

	highlights := make([][]string, 0, len(in.Hits.Hits))

	for _, hit := range in.Hits.Hits {
		highlight, err := parseSimilarityHighlights(hit.Highlight)
		if err != nil {
			return nil, fmt.Errorf("parse highlight: %w", err)
		}

		highlights = append(highlights, highlight)
	}

	matches := make([]model.DocumentSimilarMatch, 0, len(docs))

	for i, match := range docs {
		matches = append(matches, model.DocumentSimilarMatch{
			ID:            match.ID,
			Rate:          in.Hits.Hits[i].Score,
			Highlights:    highlights[i],
			SimilarImages: nil,
		})
	}

	return matches, nil
}

func parseSimilarityHighlights(in json.RawMessage) ([]string, error) {
	if in == nil {
		return nil, nil
	}

	var doc similarityHighlight

	if err := json.Unmarshal(in, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal highlights: %w", err)
	}

	return doc.TextRu, nil
}

func buildSearchQuery(
	index string,
	query model.DocumentSimilarQuery,
) ([]byte, error) {
	q := map_{
		"query": map_{
			"more_like_this": map_{
				"fields": []string{"text", "text.russian", "text.english"},
				"like": []map_{
					{
						"_index": index,
						"_id":    query.ID,
					},
				},
			},
		},
		"highlight": map_{
			"number_of_fragments": 0,
			"fields": map_{
				"text.russian": map_{},
				"text.english": map_{},
			},
		},
	}

	m, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	return m, nil
}
