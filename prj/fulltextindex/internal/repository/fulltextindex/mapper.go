package fulltextindex

import (
	"encoding/json"
	"fmt"
	"fulltextindex/internal/model"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"time"
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
		ID:          in.ID,
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
	var doc similarityHighlight

	if err := json.Unmarshal(in, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal highlights: %w", err)
	}

	return doc.Text, nil
}

func buildSearchQuery(
	query model.DocumentSimilarQuery,
) ([]byte, error) {
	var criteria []map_

	if len(query.Item.Text) > 0 {
		criteria = append(
			criteria,
			map_{
				"match": map_{
					"text": string(query.Item.Text),
				},
			},
			map_{
				"match": map_{
					"text.russian": string(query.Item.Text),
				},
			},
			map_{
				"match": map_{
					"text.english": string(query.Item.Text),
				},
			},
		)
	}

	q := map_{
		"query": map_{
			"bool": map_{
				"should": criteria,
			},
		},
		"highlight": map_{
			"number_of_fragments": 0,
			"fields": map_{
				"text.russian": map_{},
			},
		},
	}

	m, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	return m, nil
}
