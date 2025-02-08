package userstate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

func (r *Repository) Save(
	ctx context.Context,
	cmd model.UserStateSaveCommand,
) error {
	const op = "save"

	item := mapUserStateSaveCommandToInternal(cmd)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(item); err != nil {
		return fmt.Errorf("encode internal state: %w", err)
	}

	t := time.Now()

	esRes, err := r.es.Index(
		r.index,
		&buf,
		r.es.Index.WithDocumentID(strconv.Itoa(item.UserID)),
		r.es.Index.WithContext(ctx),
		r.es.Index.WithRefresh("true"),
	)
	if err != nil {
		r.m.IncUserStateRepositoryRequests(op, metricsError, time.Since(t).Seconds())
		return fmt.Errorf("execute index request: %w", err)
	}

	defer esRes.Body.Close()

	r.m.IncUserStateRepositoryRequests(op, strconv.Itoa(esRes.StatusCode), time.Since(t).Seconds())

	if err := elasticutil.IsErr(esRes); err != nil {
		return fmt.Errorf("save user state: %w", mapErrorToModel(err))
	}

	return nil
}
