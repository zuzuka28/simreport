// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package provider

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/nats-io/nats.go"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/lib/httpinstumentation"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/bot"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/config"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/metrics"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/repository/document"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/repository/similarity"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/repository/userstate"
	document2 "github.com/zuzuka28/simreport/prj/tgbot/internal/service/document"
	similarity2 "github.com/zuzuka28/simreport/prj/tgbot/internal/service/similarity"
	userstate2 "github.com/zuzuka28/simreport/prj/tgbot/internal/service/userstate"
	"net"
	"net/http"
	"sync"
	"time"
)

// Injectors from wire.go:

func ProvideUserStateRepository(configConfig *config.Config, client *elasticsearch.Client, metricsMetrics *metrics.Metrics) *userstate.Repository {
	userstateConfig := configConfig.UserStateRepo
	repository := userstate.NewRepository(userstateConfig, client, metricsMetrics)
	return repository
}

func ProvideDocumentRepository(conn *nats.Conn, metricsMetrics *metrics.Metrics) *document.Repository {
	repository := document.NewRepository(conn, metricsMetrics)
	return repository
}

func ProvideSimilarityRepository(conn *nats.Conn, metricsMetrics *metrics.Metrics) *similarity.Repository {
	repository := similarity.NewRepository(conn, metricsMetrics)
	return repository
}

func ProvideUserStateService(repository *userstate.Repository) *userstate2.Service {
	service := userstate2.NewService(repository)
	return service
}

func ProvideDocumentService(repository *document.Repository) *document2.Service {
	service := document2.NewService(repository)
	return service
}

func ProvideSimilarityService(repository *similarity.Repository) *similarity2.Service {
	service := similarity2.NewService(repository)
	return service
}

func InitBot(contextContext context.Context, configConfig *config.Config) (*bot.Bot, error) {
	botConfig := configConfig.Bot
	client, err := ProvideElastic(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	metricsMetrics := ProvideMetrics()
	repository := ProvideUserStateRepository(configConfig, client, metricsMetrics)
	service := ProvideUserStateService(repository)
	conn, err := ProvideNats(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	documentRepository := ProvideDocumentRepository(conn, metricsMetrics)
	documentService := ProvideDocumentService(documentRepository)
	similarityRepository := ProvideSimilarityRepository(conn, metricsMetrics)
	similarityService := ProvideSimilarityService(similarityRepository)
	botBot, err := bot.New(botConfig, service, documentService, similarityService, metricsMetrics)
	if err != nil {
		return nil, err
	}
	return botBot, nil
}

// wire.go:

//nolint:gochecknoglobals
var (
	metricsS    *metrics.Metrics
	metricsOnce sync.Once
)

func ProvideMetrics() *metrics.Metrics {
	metricsOnce.Do(func() {
		metricsS = metrics.New()
	})

	return metricsS
}

func ProvideConfig(path string) (*config.Config, error) {
	cfg, err := config.New(path)
	if err != nil {
		return nil, err
	}

	defaultTransportDialContext := func(
		dialer *net.Dialer,
	) func(context.Context, string, string) (net.Conn, error) {
		return dialer.DialContext
	}

	transport := &httpinstumentation.InstumentedTransport{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: defaultTransportDialContext(&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}),
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		ExtractAttrs: func(ctx context.Context) []any {
			return []any{"request_id", ctx.Value(model.RequestIDKey)}
		},
		LogRequestBody:  true,
		LogResponseBody: false,
	}

	cfg.Elastic.Transport = transport

	return cfg, nil
}

//nolint:gochecknoglobals
var (
	elasticCli     *elasticsearch.Client
	elasticCliOnce sync.Once
)

func ProvideElastic(
	ctx context.Context,
	cfg *config.Config,
) (*elasticsearch.Client, error) {
	var err error

	elasticCliOnce.Do(func() {
		elasticCli, err = elasticutil.NewClientWithStartup(ctx, cfg.Elastic)
	})

	return elasticCli, err
}

//nolint:gochecknoglobals
var (
	natsCli     *nats.Conn
	natsCliOnce sync.Once
)

func ProvideNats(
	_ context.Context,
	cfg *config.Config,
) (*nats.Conn, error) {
	var err error

	natsCliOnce.Do(func() {
		natsCli, err = nats.Connect(cfg.Nats)
	})

	return natsCli, err
}
