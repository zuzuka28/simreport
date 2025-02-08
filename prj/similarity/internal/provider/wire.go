//go:build wireinject

package provider

import (
	"context"
	"io"
	"os"
	"sync"

	analyzenatsapi "github.com/zuzuka28/simreport/prj/similarity/api/nats/handler/similarity"
	servernats "github.com/zuzuka28/simreport/prj/similarity/api/nats/server"
	serverhttp "github.com/zuzuka28/simreport/prj/similarity/api/rest/server"
	analyzeapi "github.com/zuzuka28/simreport/prj/similarity/api/rest/server/handler/similarity"
	"github.com/zuzuka28/simreport/prj/similarity/internal/config"
	"github.com/zuzuka28/simreport/prj/similarity/internal/metrics"
	analyzehistoryrepo "github.com/zuzuka28/simreport/prj/similarity/internal/repository/analyzehistory"
	documentrepo "github.com/zuzuka28/simreport/prj/similarity/internal/repository/document"
	similarityindexrepo "github.com/zuzuka28/simreport/prj/similarity/internal/repository/similarityindexclient"
	documentsrv "github.com/zuzuka28/simreport/prj/similarity/internal/service/document"
	fulltextindexsrv "github.com/zuzuka28/simreport/prj/similarity/internal/service/fulltextindex"
	semanticindexsrv "github.com/zuzuka28/simreport/prj/similarity/internal/service/semanticindex"
	shingleindexsrv "github.com/zuzuka28/simreport/prj/similarity/internal/service/shingleindex"
	analyzesrv "github.com/zuzuka28/simreport/prj/similarity/internal/service/similarity"

	"github.com/zuzuka28/simreport/lib/elasticutil"

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

func ProvideSpec() ([]byte, error) {
	f, err := os.Open("./api/rest/doc/openapi.yaml")
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	spec, err := io.ReadAll(f)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return spec, nil
}

func InitConfig(_ string) (*config.Config, error) {
	panic(wire.Build(config.New))
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

	return elasticCli, err //nolint:wrapcheck
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

func InitDocumentService(
	_ *documentrepo.Repository,
) (*documentsrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(documentsrv.Repository), new(*documentrepo.Repository)),
		documentsrv.NewService,
	))
}

func InitSimilarityIndexRepository(
	_ similarityindexrepo.Opts,
	_ *nats.Conn,
	_ *metrics.Metrics,
) (*similarityindexrepo.Repository, error) {
	panic(wire.Build(
		wire.Bind(new(similarityindexrepo.Metrics), new(*metrics.Metrics)),
		similarityindexrepo.NewRepository,
	))
}

func InitFulltextIndexService(
	_ *nats.Conn,
	_ *metrics.Metrics,
) (*fulltextindexsrv.Service, error) {
	panic(wire.Build(
		wire.Value(similarityindexrepo.Opts{
			MicroSubject: "similarity_fulltext",
		}),
		InitSimilarityIndexRepository,
		wire.Bind(new(fulltextindexsrv.Repository), new(*similarityindexrepo.Repository)),
		fulltextindexsrv.NewService,
	))
}

func InitShingleIndexService(
	_ *nats.Conn,
	_ *metrics.Metrics,
) (*shingleindexsrv.Service, error) {
	panic(wire.Build(
		wire.Value(similarityindexrepo.Opts{
			MicroSubject: "similarity_shingle",
		}),
		InitSimilarityIndexRepository,
		wire.Bind(new(shingleindexsrv.Repository), new(*similarityindexrepo.Repository)),
		shingleindexsrv.NewService,
	))
}

func InitSemanticIndexService(
	_ *nats.Conn,
	_ *metrics.Metrics,
) (*semanticindexsrv.Service, error) {
	panic(wire.Build(
		wire.Value(similarityindexrepo.Opts{
			MicroSubject: "similarity_semantic",
		}),
		InitSimilarityIndexRepository,
		wire.Bind(new(semanticindexsrv.Repository), new(*similarityindexrepo.Repository)),
		semanticindexsrv.NewService,
	))
}

func InitAnalyzeHistoryRepository(
	_ *elasticsearch.Client,
	_ *config.Config,
	_ *metrics.Metrics,
) (*analyzehistoryrepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "AnalyzeHistoryRepo"),
		wire.Bind(new(analyzehistoryrepo.Metrics), new(*metrics.Metrics)),
		analyzehistoryrepo.NewRepository,
	))
}

func ProvideAnalyzeServiceOpts() analyzesrv.Opts {
	return analyzesrv.Opts{}
}

func InitAnalyzeService(
	_ *config.Config,
	_ *documentsrv.Service,
	_ *shingleindexsrv.Service,
	_ *fulltextindexsrv.Service,
	_ *semanticindexsrv.Service,
	_ *analyzehistoryrepo.Repository,
) (*analyzesrv.Service, error) {
	panic(wire.Build(
		ProvideAnalyzeServiceOpts,
		wire.Bind(new(analyzesrv.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(analyzesrv.ShingleIndexService), new(*shingleindexsrv.Service)),
		wire.Bind(new(analyzesrv.FulltextIndexService), new(*fulltextindexsrv.Service)),
		wire.Bind(new(analyzesrv.SemanticIndexService), new(*semanticindexsrv.Service)),
		wire.Bind(new(analyzesrv.HistoryRepository), new(*analyzehistoryrepo.Repository)),
		analyzesrv.NewService,
	))
}

func InitAnalyzeHandler(
	_ *analyzesrv.Service,
) *analyzeapi.Handler {
	panic(wire.Build(
		wire.Bind(new(analyzeapi.Service), new(*analyzesrv.Service)),
		analyzeapi.NewHandler,
	))
}

func InitRestAPI(
	_ context.Context,
	_ *config.Config,
) (*serverhttp.Server, error) {
	panic(wire.Build(
		ProvideMetrics,
		ProvideSpec,
		ProvideElastic,
		ProvideNats,

		InitDocumentRepository,
		InitDocumentService,

		InitShingleIndexService,

		InitFulltextIndexService,

		InitSemanticIndexService,

		InitAnalyzeHistoryRepository,

		InitAnalyzeService,

		InitAnalyzeHandler,
		wire.Bind(new(serverhttp.SimilarityHandler), new(*analyzeapi.Handler)),
		wire.FieldsOf(new(*config.Config), "Port"),
		wire.Struct(new(serverhttp.Opts), "*"),
		serverhttp.New,
	))
}

func InitAnalyzeNatsHandler(
	_ *analyzesrv.Service,
) *analyzenatsapi.Handler {
	panic(wire.Build(
		wire.Bind(new(analyzenatsapi.Service), new(*analyzesrv.Service)),
		analyzenatsapi.NewHandler,
	))
}

func InitNatsAPI(
	_ context.Context,
	_ *config.Config,
) (*servernats.Server, error) {
	panic(wire.Build(
		ProvideMetrics,
		ProvideElastic,
		ProvideNats,

		InitDocumentRepository,
		InitDocumentService,

		InitShingleIndexService,

		InitFulltextIndexService,

		InitSemanticIndexService,

		InitAnalyzeHistoryRepository,

		InitAnalyzeService,

		InitAnalyzeNatsHandler,

		wire.Bind(new(servernats.SimilarityHandler), new(*analyzenatsapi.Handler)),
		servernats.NewServer,
	))
}
