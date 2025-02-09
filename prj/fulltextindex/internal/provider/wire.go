//go:build wireinject

package provider

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/lib/httpinstumentation"
	"github.com/zuzuka28/simreport/lib/minioutil"
	serverevent "github.com/zuzuka28/simreport/prj/fulltextindex/api/nats/event"
	indexerevent "github.com/zuzuka28/simreport/prj/fulltextindex/api/nats/event/handler/indexer"
	fulltextindexmicro "github.com/zuzuka28/simreport/prj/fulltextindex/api/nats/micro/handler/fulltextindex"
	servermicro "github.com/zuzuka28/simreport/prj/fulltextindex/api/nats/micro/server"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/config"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/metrics"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/model"
	documentrepo "github.com/zuzuka28/simreport/prj/fulltextindex/internal/repository/document"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/repository/filestorage"
	fulltextindexrepo "github.com/zuzuka28/simreport/prj/fulltextindex/internal/repository/fulltextindex"
	documentsrv "github.com/zuzuka28/simreport/prj/fulltextindex/internal/service/document"
	fulltextindexsrv "github.com/zuzuka28/simreport/prj/fulltextindex/internal/service/fulltextindex"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
)

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
		return nil, err //nolint:wrapcheck
	}

	defaultTransportDialContext := func(
		dialer *net.Dialer,
	) func(context.Context, string, string) (net.Conn, error) {
		return dialer.DialContext
	}

	//nolint:exhaustruct,gomnd,mnd
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
	cfg.S3.Transport = transport

	return cfg, nil
}

func InitElastic(
	_ context.Context,
	_ *config.Config,
) (*elasticsearch.Client, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "Elastic"),
		elasticutil.NewClientWithStartup,
	))
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

	return natsCli, err //nolint:wrapcheck
}

//nolint:gochecknoglobals
var (
	s3Cli     *minio.Client
	s3CliOnce sync.Once
)

func ProvideS3(
	ctx context.Context,
	cfg *config.Config,
) (*minio.Client, error) {
	var err error

	s3CliOnce.Do(func() {
		s3Cli, err = minioutil.NewClientWithStartup(ctx, cfg.S3)
	})

	return s3Cli, err //nolint:wrapcheck
}

func InitFilestorageRepository(
	_ *minio.Client,
	_ *config.Config,
	_ *metrics.Metrics,
) (*filestorage.Repository, error) {
	panic(wire.Build(
		wire.Bind(new(filestorage.Metrics), new(*metrics.Metrics)),
		filestorage.NewRepository,
	))
}

func InitDocumentRepository(
	_ *nats.Conn,
	_ *metrics.Metrics,
) (*documentrepo.Repository, error) {
	panic(wire.Build(
		wire.Bind(new(documentrepo.Metrics), new(*metrics.Metrics)),
		documentrepo.NewRepository,
	))
}

func InitDocumentService(
	_ *documentrepo.Repository,
) (*documentsrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(documentsrv.Repository), new(*documentrepo.Repository)),
		documentsrv.NewService,
	))
}

func InitFulltextIndexRepository(
	_ *elasticsearch.Client,
	_ *config.Config,
	_ *metrics.Metrics,
) (*fulltextindexrepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "FulltextRepo"),
		wire.Bind(new(fulltextindexrepo.Metrics), new(*metrics.Metrics)),
		fulltextindexrepo.NewRepository,
	))
}

func InitFulltextIndexService(
	_ *fulltextindexrepo.Repository,
) (*fulltextindexsrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(fulltextindexsrv.Repository), new(*fulltextindexrepo.Repository)),
		fulltextindexsrv.NewService,
	))
}

func InitFulltextHandler(
	_ *fulltextindexsrv.Service,
	_ *documentsrv.Service,
	_ *filestorage.Repository,
) (*fulltextindexmicro.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(fulltextindexmicro.Service), new(*fulltextindexsrv.Service)),
		wire.Bind(new(fulltextindexmicro.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(fulltextindexmicro.Filestorage), new(*filestorage.Repository)),
		fulltextindexmicro.NewHandler,
	))
}

func InitIndexerHandler(
	_ *fulltextindexsrv.Service,
	_ *documentsrv.Service,
	_ *filestorage.Repository,
) (*indexerevent.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(indexerevent.Service), new(*fulltextindexsrv.Service)),
		wire.Bind(new(indexerevent.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(indexerevent.Filestorage), new(*filestorage.Repository)),
		indexerevent.NewHandler,
	))
}

func InitNatsMicroAPI(
	_ context.Context,
	_ *config.Config,
) (*servermicro.Server, error) {
	panic(wire.Build(
		ProvideMetrics,
		InitElastic,
		ProvideNats,

		ProvideS3,
		InitFilestorageRepository,

		InitDocumentRepository,
		InitDocumentService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitFulltextHandler,

		wire.Bind(new(servermicro.Handler), new(*fulltextindexmicro.Handler)),
		wire.Bind(new(servermicro.Metrics), new(*metrics.Metrics)),
		servermicro.NewServer,
	))
}

func InitNatsEventAPI(
	_ context.Context,
	_ *config.Config,
) (*serverevent.Server, error) {
	panic(wire.Build(
		ProvideMetrics,
		InitElastic,
		ProvideNats,

		ProvideS3,
		InitFilestorageRepository,

		InitDocumentRepository,
		InitDocumentService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitIndexerHandler,

		wire.Bind(new(serverevent.IndexerHandler), new(*indexerevent.Handler)),
		serverevent.NewServer,
	))
}
