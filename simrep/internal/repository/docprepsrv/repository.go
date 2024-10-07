package docprepsrv

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"simrep/internal/model"
	client "simrep/pkg/docprepsrv"
)

var ErrStatusNotOK = errors.New("status not ok")

type Opts struct {
	APIURL     string     `yaml:"apiurl"`
	HTTPClient HTTPClient `yaml:"-"`
}

type Repository struct {
	cli apiClient
}

func NewRepository(opts Opts) (*Repository, error) {
	httpcli := opts.HTTPClient
	if httpcli == nil {
		httpcli = http.DefaultClient
	}

	cli, err := client.NewClientWithResponses(
		opts.APIURL,
		client.WithHTTPClient(httpcli),
	)
	if err != nil {
		return nil, fmt.Errorf("new simrep preprocess api client: %w", err)
	}

	return &Repository{
		cli: cli,
	}, nil
}

func (r *Repository) PreprocessRawDocument(
	ctx context.Context,
	doc []byte,
) (*model.Document, error) {
	body, ftype, err := mapDocToPreprocessRequest(doc)
	if err != nil {
		return nil, fmt.Errorf("map document to request: %w", err)
	}

	raw, err := r.cli.ProcessItemPreprocessDocPostWithBodyWithResponse(ctx, ftype, body)
	if err != nil {
		return nil, fmt.Errorf("request api: %w", err)
	}

	if raw.StatusCode() != http.StatusOK {
		return nil, ErrStatusNotOK
	}

	res := mapDocResponseToModel(raw.JSON200)

	return res, nil
}
