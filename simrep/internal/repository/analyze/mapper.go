package analyze

import (
	"encoding/json"
	"fmt"
	"simrep/internal/model"
	"simrep/pkg/elasticutil"
	"slices"
	"time"
)

//nolint:gochecknoglobals
var now = time.Now

func parseDocument(hit *elasticutil.Hit) (model.AnalyzedDocument, error) {
	var doc analyzedDocument

	if err := json.Unmarshal(hit.Source, &doc); err != nil {
		return model.AnalyzedDocument{}, fmt.Errorf("unmarshal analyzed document: %w", err)
	}

	return mapDocumentToModel(doc), nil
}

func mapDocumentToInternal(in model.AnalyzedDocument) analyzedDocument {
	images := make([]analyzedImage, 0, len(in.Images))
	for _, img := range in.Images {
		images = append(images, mapImageToInternal(img))
	}

	return analyzedDocument{
		ID:          in.ID,
		Text:        in.Text,
		TextVector:  in.TextVector,
		Images:      images,
		LastUpdated: now(),
	}
}

func mapImageToInternal(in model.AnalyzedImage) analyzedImage {
	return analyzedImage{
		ID:        in.ID,
		Vector:    in.Vector,
		HashImage: mapHashImageToInternal(in.HashImage),
	}
}

func mapHashImageToInternal(in model.HashImage) hashImage {
	return hashImage{
		Ahash:       in.Ahash,
		AhashVector: in.AhashVector,
		Dhash:       in.Dhash,
		DhashVector: in.DhashVector,
		Phash:       in.Phash,
		PhashVector: in.PhashVector,
		Whash:       in.Whash,
		WhashVector: in.WhashVector,
	}
}

func mapDocumentToModel(in analyzedDocument) model.AnalyzedDocument {
	images := make([]model.AnalyzedImage, 0, len(in.Images))
	for _, img := range in.Images {
		images = append(images, mapImageToModel(img))
	}

	return model.AnalyzedDocument{
		ID:         in.ID,
		Text:       in.Text,
		TextVector: in.TextVector,
		Images:     images,
	}
}

func mapImageToModel(in analyzedImage) model.AnalyzedImage {
	return model.AnalyzedImage{
		ID:        in.ID,
		Vector:    in.Vector,
		HashImage: mapHashImageToModel(in.HashImage),
	}
}

func mapHashImageToModel(in hashImage) model.HashImage {
	return model.HashImage{
		Ahash:       in.Ahash,
		AhashVector: in.AhashVector,
		Dhash:       in.Dhash,
		DhashVector: in.DhashVector,
		Phash:       in.Phash,
		PhashVector: in.PhashVector,
		Whash:       in.Whash,
		WhashVector: in.WhashVector,
	}
}

func mapSearchResponseToMatches(
	query model.AnalyzedDocumentSimilarQuery,
	in *elasticutil.SearchResponse,
) ([]model.AnalyzedDocumentMatch, error) {
	docs := make([]model.AnalyzedDocument, 0, len(in.Hits.Hits))

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

	matches := make([]model.AnalyzedDocumentMatch, 0, len(docs))

	for i, match := range docs {
		matches = append(matches, model.AnalyzedDocumentMatch{
			ID:            match.ID,
			Rate:          in.Hits.Hits[i].Score,
			Highlights:    highlights[i],
			SimilarImages: filterSimilarImages(query, match),
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

func filterSimilarImages(
	query model.AnalyzedDocumentSimilarQuery,
	in model.AnalyzedDocument,
) []string {
	var imgs []string

	for _, img := range in.Images {
		if slices.Contains(query.Item.ImagesIDs(), img.ID) {
			imgs = append(imgs, img.ID)
		}
	}

	return imgs
}

func buildSearchQuery(
	query model.AnalyzedDocumentSimilarQuery,
) ([]byte, error) {
	var criteria []map_

	if len(query.Item.Images) > 0 {
		criteria = append(
			criteria,
			map_{
				"terms": map_{
					"images.id": query.Item.ImagesIDs(),
				},
			},
		)
	}

	if len(query.Item.Text) > 0 {
		criteria = append(
			criteria,
			map_{
				"match": map_{
					"text": query.Item.Text,
				},
			},
			map_{
				"match": map_{
					"text.russian": query.Item.Text,
				},
			},
			map_{
				"match": map_{
					"text.english": query.Item.Text,
				},
			},
		)
	}

	// if len(query.Item.TextVector) > 0 {
	// 	criteria = append(
	// 		criteria,
	// 		map_{
	// 			"script_score": map_{
	// 				"query": map_{
	// 					"match_all": map_{},
	// 				},
	// 				"script": map_{
	// 					"source": `1 + cosineSimilarity(params.query_vector, 'textVector') * 0.5`,
	// 					"params": map_{
	// 						"query_vector": query.Item.TextVector,
	// 					},
	// 				},
	// 			},
	// 		})
	// }
	//
	// if len(query.Item.Images) > 0 {
	// 	for _, vec := range query.Item.Images {
	// 		criteria = append(
	// 			criteria,
	// 			map_{
	// 				"script_score": map_{
	// 					"query": map_{
	// 						"match_all": map_{},
	// 					},
	// 					"script": map_{
	// 						"source": `1 + cosineSimilarity(params.query_vector, 'images.vector') * 0.5`,
	// 						"params": map_{
	// 							"query_vector": vec,
	// 						},
	// 					},
	// 				},
	// 			})
	// 	}
	// }

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

	fmt.Println(string(m))

	return m, nil
}
