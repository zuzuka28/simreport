//go:build wireinject

package provider

import (
	"context"
	"sync"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	serverevent "github.com/zuzuka28/simreport/prj/semanticindex/api/nats/event"
	indexerevent "github.com/zuzuka28/simreport/prj/semanticindex/api/nats/event/handler/indexer"
	semanticindexmicro "github.com/zuzuka28/simreport/prj/semanticindex/api/nats/micro/handler/semanticindex"
	servermicro "github.com/zuzuka28/simreport/prj/semanticindex/api/nats/micro/server"
	"github.com/zuzuka28/simreport/prj/semanticindex/internal/config"
	"github.com/zuzuka28/simreport/prj/semanticindex/internal/metrics"
	documentrepo "github.com/zuzuka28/simreport/prj/semanticindex/internal/repository/document"
	semanticindexrepo "github.com/zuzuka28/simreport/prj/semanticindex/internal/repository/semanticindex"
	vectorizerrepo "github.com/zuzuka28/simreport/prj/semanticindex/internal/repository/vectorizer"
	documentsrv "github.com/zuzuka28/simreport/prj/semanticindex/internal/service/document"
	semanticindexsrv "github.com/zuzuka28/simreport/prj/semanticindex/internal/service/semanticindex"
	vectorizersrv "github.com/zuzuka28/simreport/prj/semanticindex/internal/service/vectorizer"

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

func InitConfig(path string) (*config.Config, error) {
	panic(wire.Build(config.New))
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

func InitDocumentRepository(
	_ *nats.Conn,
	_ *metrics.Metrics,
) (*documentrepo.Repository, error) {
	panic(wire.Build(
		wire.Bind(new(documentrepo.Metrics), new(*metrics.Metrics)),
		documentrepo.NewRepository,
	))
}

func InitVectorizerRepository(
	_ *config.Config,
	_ *metrics.Metrics,
) (*vectorizerrepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "VectorizerRepo"),
		wire.Bind(new(vectorizerrepo.Metrics), new(*metrics.Metrics)),
		vectorizerrepo.NewRepository,
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

func InitSemanticIndexRepository(
	_ *elasticsearch.Client,
	_ *config.Config,
	_ *metrics.Metrics,
) (*semanticindexrepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "SemanticRepo"),
		wire.Bind(new(semanticindexrepo.Metrics), new(*metrics.Metrics)),
		semanticindexrepo.NewRepository,
	))
}

func InitVectorizerService(
	_ *config.Config,
	_ *vectorizerrepo.Repository,
) (*vectorizersrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(vectorizersrv.Repository), new(*vectorizerrepo.Repository)),
		vectorizersrv.NewService,
	))
}

func InitSemanticIndexService(
	_ *semanticindexrepo.Repository,
	_ *vectorizersrv.Service,
) (*semanticindexsrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(semanticindexsrv.Repository), new(*semanticindexrepo.Repository)),
		wire.Bind(new(semanticindexsrv.VectorizerService), new(*vectorizersrv.Service)),
		semanticindexsrv.NewService,
	))
}

func InitSemanticHandler(
	_ *semanticindexsrv.Service,
	_ *documentsrv.Service,
) (*semanticindexmicro.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(semanticindexmicro.Service), new(*semanticindexsrv.Service)),
		wire.Bind(new(semanticindexmicro.DocumentService), new(*documentsrv.Service)),
		semanticindexmicro.NewHandler,
	))
}

func InitIndexerHandler(
	_ *semanticindexsrv.Service,
	_ *documentsrv.Service,
) (*indexerevent.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(indexerevent.Service), new(*semanticindexsrv.Service)),
		wire.Bind(new(indexerevent.DocumentService), new(*documentsrv.Service)),
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

		InitVectorizerRepository,
		InitVectorizerService,

		InitDocumentRepository,
		InitDocumentService,

		InitSemanticIndexRepository,
		InitSemanticIndexService,

		InitSemanticHandler,

		wire.Bind(new(servermicro.Handler), new(*semanticindexmicro.Handler)),
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

		InitVectorizerRepository,
		InitVectorizerService,

		InitDocumentRepository,
		InitDocumentService,

		InitSemanticIndexRepository,
		InitSemanticIndexService,

		InitIndexerHandler,

		wire.Bind(new(serverevent.IndexerHandler), new(*indexerevent.Handler)),
		serverevent.NewServer,
	))
}
