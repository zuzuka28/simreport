package userstate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
)

type Config struct {
	Index string `yaml:"index"`
}

type Repository struct {
	es    *elasticsearch.Client
	index string
}

func NewRepository(
	cfg Config,
	cli *elasticsearch.Client,
) *Repository {
	return &Repository{
		es:    cli,
		index: cfg.Index,
	}
}

func (r *Repository) Fetch(
	ctx context.Context,
	query model.UserStateQuery,
) (*model.UserState, error) {
	esRes, err := r.es.Get(
		r.index,
		strconv.Itoa(query.UserID),
		r.es.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("get from es: %w", err)
	}

	defer esRes.Body.Close()

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

func (r *Repository) Save(
	ctx context.Context,
	cmd model.UserStateSaveCommand,
) error {
	item := mapUserStateSaveCommandToInternal(cmd)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(item); err != nil {
		return fmt.Errorf("encode internal state: %w", err)
	}

	esRes, err := r.es.Index(
		r.index,
		&buf,
		r.es.Index.WithDocumentID(strconv.Itoa(item.UserID)),
		r.es.Index.WithContext(ctx),
		r.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("execute index request: %w", err)
	}

	defer esRes.Body.Close()

	if err := elasticutil.IsErr(esRes); err != nil {
		return fmt.Errorf("save user state: %w", mapErrorToModel(err))
	}

	return nil
}
