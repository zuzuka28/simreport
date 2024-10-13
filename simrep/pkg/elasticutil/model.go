package elasticutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type CountResponse struct {
	Count int `json:"count"`
}

type HitHighlight struct{}

type Hit struct {
	ID        string          `json:"_id"`
	Source    json.RawMessage `json:"_source"`
	Score     float64         `json:"_score"`
	Index     string          `json:"_index"`
	Highlight HitHighlight    `json:"highlight"`
}

type TotalHits struct {
	Value int `json:"value"`
}

type SearchResponseHits struct {
	Total TotalHits `json:"total"`
	Hits  []Hit     `json:"hits"`
}

type SearchResponse struct {
	ScrollID string             `json:"_scroll_id"`
	Hits     SearchResponseHits `json:"hits"`
	Aggs     json.RawMessage    `json:"aggregations"`
}

type BulkAction string

const (
	BulkActionCreate BulkAction = "create"
	BulkActionUpdate BulkAction = "update"
	BulkActionDelete BulkAction = "delete"
	BulkActionIndex  BulkAction = "index"
)

type BulkIndexerRecordResponse esutil.BulkIndexerResponseItem

type BulkIndexerRecord esutil.BulkIndexerItem

type BulkRecord struct {
	ID     string
	Index  string
	Action BulkAction
	Body   any

	OnSuccess func(context.Context, BulkIndexerRecord, BulkIndexerRecordResponse)        // Per item
	OnFailure func(context.Context, BulkIndexerRecord, BulkIndexerRecordResponse, error) // Per item
}

func (br *BulkRecord) ToIndexer() (esutil.BulkIndexerItem, error) {
	//nolint:exhaustruct
	item := esutil.BulkIndexerItem{
		Index:      br.Index,
		Action:     string(br.Action),
		DocumentID: br.ID,
		OnSuccess: func(
			ctx context.Context,
			item esutil.BulkIndexerItem,
			item2 esutil.BulkIndexerResponseItem,
		) {
			if br.OnSuccess == nil {
				return
			}

			br.OnSuccess(ctx, BulkIndexerRecord(item), BulkIndexerRecordResponse(item2))
		},
		OnFailure: func(
			ctx context.Context,
			item esutil.BulkIndexerItem,
			item2 esutil.BulkIndexerResponseItem,
			err error,
		) {
			if br.OnFailure == nil {
				return
			}

			br.OnFailure(ctx, BulkIndexerRecord(item), BulkIndexerRecordResponse(item2), err)
		},
	}

	data := br.Body

	if br.Action == BulkActionUpdate {
		data = map[string]any{"doc": data}
	}

	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			return esutil.BulkIndexerItem{}, fmt.Errorf("marshal record: %w", err)
		}

		br.Body = bytes.NewReader(body)
	}

	return item, nil
}

type BulkResponseItemError struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	IndexUUID string `json:"index_uuid"`
	Shard     string `json:"shard"`
	Index     string `json:"index"`
}

type BulkResponseItem struct {
	Index  string                 `json:"_index"`
	ID     string                 `json:"_id"`
	Result string                 `json:"result"`
	Status int                    `json:"status"`
	Error  *BulkResponseItemError `json:"error"`
}

type BulkResponse struct {
	Took   int                           `json:"took"`
	Errors bool                          `json:"errors"`
	Items  []map[string]BulkResponseItem `json:"items,omitempty"`
}

func (bulk *BulkResponse) ErrorItems() []BulkResponseItem {
	if !bulk.Errors {
		return nil
	}

	var errorItems []BulkResponseItem

	for _, item := range bulk.Items {
		for _, v := range item {
			if v.Error != nil {
				errorItems = append(errorItems, v)
			}
		}
	}

	return errorItems
}

type BcrollBody struct {
	ScrollID string `json:"scroll_id"`
}

type UpdateDocBody struct {
	Doc any `json:"doc"`
}

type ESErrCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	Index  string `json:"index"`
}

type ESErrorValue struct {
	Type   string       `json:"type"`
	Reason string       `json:"reason"`
	Causes []ESErrCause `json:"root_cause"`
}

type ESError struct {
	Error ESErrorValue `json:"error"`
}

type indexConfig struct {
	Settings map[string]any `json:"settings,omitempty"`
	Mappings map[string]any `json:"mappings"`
}
