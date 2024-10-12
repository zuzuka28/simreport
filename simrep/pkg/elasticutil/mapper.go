package elasticutil

import (
	"encoding/json"
	"fmt"
	"io"
)

func ParseCountResponse(resp io.ReadCloser) (*CountResponse, error) {
	defer func() { _ = resp.Close() }()

	if resp == nil {
		return nil, ErrNilResponse
	}

	res := new(CountResponse)
	if err := json.NewDecoder(resp).Decode(res); err != nil {
		return nil, fmt.Errorf("unmarshal raw into count response: %w", err)
	}

	return res, nil
}

func ParseSearchResponse(r io.ReadCloser) (*SearchResponse, error) {
	defer func() { _ = r.Close() }()

	if r == nil {
		return nil, ErrNilResponse
	}

	res := new(SearchResponse)
	if err := json.NewDecoder(r).Decode(res); err != nil {
		return nil, fmt.Errorf("unmarshal raw into search response: %w", err)
	}

	return res, nil
}

func ParseDocResponse(r io.ReadCloser) (*Hit, error) {
	defer func() { _ = r.Close() }()

	if r == nil {
		return nil, ErrNilResponse
	}

	res := new(Hit)
	if err := json.NewDecoder(r).Decode(res); err != nil {
		return nil, fmt.Errorf("unmarshal raw into hit: %w", err)
	}

	return res, nil
}
