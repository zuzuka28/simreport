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
	analyzehistoryrepo "github.com/zuzuka28/simreport/prj/similarity/internal/repository/analyzehistory"
	documentrepo "github.com/zuzuka28/simreport/prj/similarity/internal/repository/document"
	fulltextindexrepo "github.com/zuzuka28/simreport/prj/similarity/internal/repository/fulltextindexclient"
	semanticindexrepo "github.com/zuzuka28/simreport/prj/similarity/internal/repository/semanticindexclient"
	shingleindexrepo "github.com/zuzuka28/simreport/prj/similarity/internal/repository/shingleindexclient"
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
) (*documentrepo.Repository, error) {
	panic(wire.Build(
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
	_ *nats.Conn,
) (*shingleindexrepo.Repository, error) {
	panic(wire.Build(
		shingleindexrepo.NewRepository,
	))
}

func InitShingleIndexService(
	_ *shingleindexrepo.Repository,
) (*shingleindexsrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(shingleindexsrv.Repository), new(*shingleindexrepo.Repository)),
		shingleindexsrv.NewService,
	))
}

func InitFulltextIndexRepository(
	_ *nats.Conn,
) (*fulltextindexrepo.Repository, error) {
	panic(wire.Build(
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

func InitSemanticIndexRepository(
	_ *nats.Conn,
) (*semanticindexrepo.Repository, error) {
	panic(wire.Build(
		semanticindexrepo.NewRepository,
	))
}

func InitSemanticIndexService(
	_ *semanticindexrepo.Repository,
) (*semanticindexsrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(semanticindexsrv.Repository), new(*semanticindexrepo.Repository)),
		semanticindexsrv.NewService,
	))
}

func InitAnalyzeHistoryRepository(
	_ *elasticsearch.Client,
	_ *config.Config,
) (*analyzehistoryrepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "AnalyzeHistoryRepo"),
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
		ProvideSpec,
		ProvideElastic,
		ProvideNats,

		InitDocumentRepository,
		InitDocumentService,

		InitShingleIndexRepository,
		InitShingleIndexService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitSemanticIndexRepository,
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
		ProvideElastic,
		ProvideNats,

		InitDocumentRepository,
		InitDocumentService,

		InitShingleIndexRepository,
		InitShingleIndexService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitSemanticIndexRepository,
		InitSemanticIndexService,

		InitAnalyzeHistoryRepository,

		InitAnalyzeService,

		InitAnalyzeNatsHandler,

		wire.Bind(new(servernats.SimilarityHandler), new(*analyzenatsapi.Handler)),
		servernats.NewServer,
	))
}
