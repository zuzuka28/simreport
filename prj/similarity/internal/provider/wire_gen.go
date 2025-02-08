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
	similarity3 "github.com/zuzuka28/simreport/prj/similarity/api/nats/handler/similarity"
	server2 "github.com/zuzuka28/simreport/prj/similarity/api/nats/server"
	"github.com/zuzuka28/simreport/prj/similarity/api/rest/server"
	similarity2 "github.com/zuzuka28/simreport/prj/similarity/api/rest/server/handler/similarity"
	"github.com/zuzuka28/simreport/prj/similarity/internal/config"
	"github.com/zuzuka28/simreport/prj/similarity/internal/metrics"
	"github.com/zuzuka28/simreport/prj/similarity/internal/model"
	"github.com/zuzuka28/simreport/prj/similarity/internal/repository/analyzehistory"
	"github.com/zuzuka28/simreport/prj/similarity/internal/repository/document"
	"github.com/zuzuka28/simreport/prj/similarity/internal/repository/similarityindexclient"
	document2 "github.com/zuzuka28/simreport/prj/similarity/internal/service/document"
	"github.com/zuzuka28/simreport/prj/similarity/internal/service/fulltextindex"
	"github.com/zuzuka28/simreport/prj/similarity/internal/service/semanticindex"
	"github.com/zuzuka28/simreport/prj/similarity/internal/service/shingleindex"
	"github.com/zuzuka28/simreport/prj/similarity/internal/service/similarity"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// Injectors from wire.go:

func InitDocumentRepository(conn *nats.Conn, metricsMetrics *metrics.Metrics) (*document.Repository, error) {
	repository := document.NewRepository(conn, metricsMetrics)
	return repository, nil
}

func InitDocumentService(repository *document.Repository) (*document2.Service, error) {
	service := document2.NewService(repository)
	return service, nil
}

func InitSimilarityIndexRepository(opts similarityindexclient.Opts, conn *nats.Conn, metricsMetrics *metrics.Metrics) (*similarityindexclient.Repository, error) {
	repository := similarityindexclient.NewRepository(opts, conn, metricsMetrics)
	return repository, nil
}

func InitFulltextIndexService(conn *nats.Conn, metricsMetrics *metrics.Metrics) (*fulltextindex.Service, error) {
	opts := _wireOptsValue
	repository, err := InitSimilarityIndexRepository(opts, conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	service := fulltextindex.NewService(repository)
	return service, nil
}

var (
	_wireOptsValue = similarityindexclient.Opts{
		MicroSubject: "similarity_fulltext",
	}
)

func InitShingleIndexService(conn *nats.Conn, metricsMetrics *metrics.Metrics) (*shingleindex.Service, error) {
	opts := _wireSimilarityindexclientOptsValue
	repository, err := InitSimilarityIndexRepository(opts, conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	service := shingleindex.NewService(repository)
	return service, nil
}

var (
	_wireSimilarityindexclientOptsValue = similarityindexclient.Opts{
		MicroSubject: "similarity_shingle",
	}
)

func InitSemanticIndexService(conn *nats.Conn, metricsMetrics *metrics.Metrics) (*semanticindex.Service, error) {
	opts := _wireOptsValue2
	repository, err := InitSimilarityIndexRepository(opts, conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	service := semanticindex.NewService(repository)
	return service, nil
}

var (
	_wireOptsValue2 = similarityindexclient.Opts{
		MicroSubject: "similarity_semantic",
	}
)

func InitAnalyzeHistoryRepository(client *elasticsearch.Client, configConfig *config.Config, metricsMetrics *metrics.Metrics) (*analyzehistory.Repository, error) {
	opts := configConfig.AnalyzeHistoryRepo
	repository, err := analyzehistory.NewRepository(opts, client, metricsMetrics)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func InitAnalyzeService(configConfig *config.Config, service *document2.Service, shingleindexService *shingleindex.Service, fulltextindexService *fulltextindex.Service, semanticindexService *semanticindex.Service, repository *analyzehistory.Repository) (*similarity.Service, error) {
	opts := ProvideAnalyzeServiceOpts()
	similarityService := similarity.NewService(opts, service, shingleindexService, fulltextindexService, semanticindexService, repository)
	return similarityService, nil
}

func InitAnalyzeHandler(service *similarity.Service) *similarity2.Handler {
	handler := similarity2.NewHandler(service)
	return handler
}

func InitRestAPI(contextContext context.Context, configConfig *config.Config) (*server.Server, error) {
	int2 := configConfig.Port
	v, err := ProvideSpec()
	if err != nil {
		return nil, err
	}
	conn, err := ProvideNats(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	metricsMetrics := ProvideMetrics()
	repository, err := InitDocumentRepository(conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	service, err := InitDocumentService(repository)
	if err != nil {
		return nil, err
	}
	shingleindexService, err := InitShingleIndexService(conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	fulltextindexService, err := InitFulltextIndexService(conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	semanticindexService, err := InitSemanticIndexService(conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	client, err := ProvideElastic(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	analyzehistoryRepository, err := InitAnalyzeHistoryRepository(client, configConfig, metricsMetrics)
	if err != nil {
		return nil, err
	}
	similarityService, err := InitAnalyzeService(configConfig, service, shingleindexService, fulltextindexService, semanticindexService, analyzehistoryRepository)
	if err != nil {
		return nil, err
	}
	handler := InitAnalyzeHandler(similarityService)
	opts := server.Opts{
		Port:           int2,
		Spec:           v,
		AnalyzeHandler: handler,
		Metrics:        metricsMetrics,
	}
	serverServer, err := server.New(opts)
	if err != nil {
		return nil, err
	}
	return serverServer, nil
}

func InitAnalyzeNatsHandler(service *similarity.Service) *similarity3.Handler {
	handler := similarity3.NewHandler(service)
	return handler
}

func InitNatsAPI(contextContext context.Context, configConfig *config.Config) (*server2.Server, error) {
	conn, err := ProvideNats(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	metricsMetrics := ProvideMetrics()
	repository, err := InitDocumentRepository(conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	service, err := InitDocumentService(repository)
	if err != nil {
		return nil, err
	}
	shingleindexService, err := InitShingleIndexService(conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	fulltextindexService, err := InitFulltextIndexService(conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	semanticindexService, err := InitSemanticIndexService(conn, metricsMetrics)
	if err != nil {
		return nil, err
	}
	client, err := ProvideElastic(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	analyzehistoryRepository, err := InitAnalyzeHistoryRepository(client, configConfig, metricsMetrics)
	if err != nil {
		return nil, err
	}
	similarityService, err := InitAnalyzeService(configConfig, service, shingleindexService, fulltextindexService, semanticindexService, analyzehistoryRepository)
	if err != nil {
		return nil, err
	}
	handler := InitAnalyzeNatsHandler(similarityService)
	serverServer := server2.NewServer(conn, handler, metricsMetrics)
	return serverServer, nil
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

func ProvideSpec() ([]byte, error) {
	f, err := os.Open("./api/rest/doc/openapi.yaml")
	if err != nil {
		return nil, err
	}

	spec, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return spec, nil
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

func ProvideAnalyzeServiceOpts() similarity.Opts {
	return similarity.Opts{}
}
