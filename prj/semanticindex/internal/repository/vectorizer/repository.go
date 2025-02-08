package vectorizer

import (
	"errors"
	"fmt"

	client "github.com/zuzuka28/simreport/prj/semanticindex/internal/repository/vectorizer/vectorizerclient"
)

var errBadResponse = errors.New("bad response")

type Opts struct {
	Host string `yaml:"host"`
}

type Repository struct {
	cli client.ClientWithResponsesInterface

	m Metrics
}

func NewRepository(
	cfg Opts,
	m Metrics,
) (*Repository, error) {
	cli, err := client.NewClientWithResponses(cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	return &Repository{
		cli: cli,
		m:   m,
	}, nil
}
