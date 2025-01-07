package elasticutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var ErrInvalidIndexCount = errors.New("invalid index count")

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalid     = errors.New("invalid")
	ErrInternal    = errors.New("internal error")
	ErrNilResponse = errors.New("nil body in response")
)

func IsErr(res *esapi.Response) error {
	if !res.IsError() {
		return nil
	}

	var err error

	switch res.StatusCode {
	case http.StatusNotFound:
		err = ErrNotFound

	default:
		err = fmt.Errorf("%w, status: %d", ErrInvalid, res.StatusCode)
	}

	if res.Body == nil {
		return err
	}

	var e ESError
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	errStrs := []string{e.Error.Type + "/" + e.Error.Reason}

	for _, cause := range e.Error.Causes {
		errStrs = append(errStrs, fmt.Sprintf("%s/%s", cause.Index, cause.Reason))
	}

	return fmt.Errorf("%w: %s", err, strings.Join(errStrs, ", "))
}
