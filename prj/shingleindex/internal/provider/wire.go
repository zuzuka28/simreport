//go:build wireinject

package provider

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/minio/minio-go/v7"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/lib/httpinstumentation"
	"github.com/zuzuka28/simreport/lib/minioutil"
	"github.com/zuzuka28/simreport/prj/shingleindex/internal/config"
	serverevent "github.com/zuzuka28/simreport/prj/shingleindex/internal/handler/nats/event"
	indexerevent "github.com/zuzuka28/simreport/prj/shingleindex/internal/handler/nats/event/handler/indexer"
	shingleindexmicro "github.com/zuzuka28/simreport/prj/shingleindex/internal/handler/nats/micro/handler/shingleindex"
	servermicro "github.com/zuzuka28/simreport/prj/shingleindex/internal/handler/nats/micro/server"
	"github.com/zuzuka28/simreport/prj/shingleindex/internal/metrics"
	"github.com/zuzuka28/simreport/prj/shingleindex/internal/model"
	documentrepo "github.com/zuzuka28/simreport/prj/shingleindex/internal/repository/document"
	"github.com/zuzuka28/simreport/prj/shingleindex/internal/repository/filestorage"
	shingleindexrepo "github.com/zuzuka28/simreport/prj/shingleindex/internal/repository/shingleindex"
	documentsrv "github.com/zuzuka28/simreport/prj/shingleindex/internal/service/document"
	shingleindexsrv "github.com/zuzuka28/simreport/prj/shingleindex/internal/service/shingleindex"

	"github.com/google/wire"
	"github.com/nats-io/nats.go"
)

func ProvideMetrics() *metrics.Metrics {
	return metrics.New()
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

func InitShingleIndexRepository(
	_ *elasticsearch.Client,
	_ *config.Config,
) (*shingleindexrepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "ShingleRepo"),
		shingleindexrepo.NewRepository,
	))
}

func InitShingleIndexService(
	_ *shingleindexrepo.Repository,
	_ *documentsrv.Service,
) (*shingleindexsrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(shingleindexsrv.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(shingleindexsrv.Repository), new(*shingleindexrepo.Repository)),
		shingleindexsrv.NewService,
	))
}

func InitShingleHandler(
	_ *shingleindexsrv.Service,
	_ *documentsrv.Service,
	_ *filestorage.Repository,
) (*shingleindexmicro.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(shingleindexmicro.Service), new(*shingleindexsrv.Service)),
		wire.Bind(new(shingleindexmicro.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(shingleindexmicro.Filestorage), new(*filestorage.Repository)),
		shingleindexmicro.NewHandler,
	))
}

func InitIndexerHandler(
	_ *shingleindexsrv.Service,
	_ *documentsrv.Service,
	_ *filestorage.Repository,
) (*indexerevent.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(indexerevent.Service), new(*shingleindexsrv.Service)),
		wire.Bind(new(indexerevent.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(indexerevent.Filestorage), new(*filestorage.Repository)),
		indexerevent.NewHandler,
	))
}

func InitNatsMicroAPI(
	_ context.Context,
	_ *config.Config,
	_ *metrics.Metrics,
) (*servermicro.Server, error) {
	panic(wire.Build(
		ProvideNats,
		InitElastic,

		ProvideS3,
		InitFilestorageRepository,

		InitDocumentRepository,
		InitDocumentService,

		InitShingleIndexRepository,
		InitShingleIndexService,

		InitShingleHandler,

		wire.Bind(new(servermicro.Handler), new(*shingleindexmicro.Handler)),
		wire.Bind(new(servermicro.Metrics), new(*metrics.Metrics)),
		servermicro.NewServer,
	))
}

func InitNatsEventAPI(
	_ context.Context,
	_ *config.Config,
	_ *metrics.Metrics,
) (*serverevent.Server, error) {
	panic(wire.Build(
		ProvideNats,
		InitElastic,

		ProvideS3,
		InitFilestorageRepository,

		InitDocumentRepository,
		InitDocumentService,

		InitShingleIndexRepository,
		InitShingleIndexService,

		InitIndexerHandler,

		wire.Bind(new(serverevent.IndexerHandler), new(*indexerevent.Handler)),
		serverevent.NewServer,
	))
}
