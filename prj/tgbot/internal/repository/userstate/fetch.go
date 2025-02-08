package userstate

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

func (r *Repository) Fetch(
	ctx context.Context,
	query model.UserStateQuery,
) (*model.UserState, error) {
	const op = "fetch"

	t := time.Now()

	esRes, err := r.es.Get(
		r.index,
		strconv.Itoa(query.UserID),
		r.es.Get.WithContext(ctx),
	)
	if err != nil {
		r.m.IncUserStateRepositoryRequests(op, metricsError, time.Since(t).Seconds())
		return nil, fmt.Errorf("get from es: %w", err)
	}

	defer esRes.Body.Close()

	r.m.IncUserStateRepositoryRequests(op, strconv.Itoa(esRes.StatusCode), time.Since(t).Seconds())

	if err := elasticutil.IsErr(esRes); err != nil {
		return nil, fmt.Errorf("fetch user state: %w", mapErrorToModel(err))
	}

	raw, err := elasticutil.ParseDocResponse(esRes.Body)
	if err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	res, err := parseFetchUserStateResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse user state: %w", err)
	}

	return res, nil
}
