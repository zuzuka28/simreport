//go:build wireinject

package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	documentnats "github.com/zuzuka28/simreport/prj/document/api/nats/handler/document"
	servernats "github.com/zuzuka28/simreport/prj/document/api/nats/server"
	serverhttp "github.com/zuzuka28/simreport/prj/document/api/rest/server"
	analyzeapi "github.com/zuzuka28/simreport/prj/document/api/rest/server/handler/analyze"
	attributeapi "github.com/zuzuka28/simreport/prj/document/api/rest/server/handler/attribute"
	documentapi "github.com/zuzuka28/simreport/prj/document/api/rest/server/handler/document"
	"github.com/zuzuka28/simreport/prj/document/internal/config"
	"github.com/zuzuka28/simreport/prj/document/internal/model"
	analyzehistoryrepo "github.com/zuzuka28/simreport/prj/document/internal/repository/analyzehistory"
	attributerepo "github.com/zuzuka28/simreport/prj/document/internal/repository/attribute"
	documentrepo "github.com/zuzuka28/simreport/prj/document/internal/repository/document"
	documentstatusrepo "github.com/zuzuka28/simreport/prj/document/internal/repository/documentstatus"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/filestorage"
	fulltextindexrepo "github.com/zuzuka28/simreport/prj/document/internal/repository/fulltextindexclient"
	semanticindexrepo "github.com/zuzuka28/simreport/prj/document/internal/repository/semanticindexclient"
	shingleindexrepo "github.com/zuzuka28/simreport/prj/document/internal/repository/shingleindexclient"
	analyzesrv "github.com/zuzuka28/simreport/prj/document/internal/service/analyze"
	attributesrv "github.com/zuzuka28/simreport/prj/document/internal/service/attribute"
	documentsrv "github.com/zuzuka28/simreport/prj/document/internal/service/document"
	documentparsersrv "github.com/zuzuka28/simreport/prj/document/internal/service/documentparser"
	"github.com/zuzuka28/simreport/prj/document/internal/service/documentpipeline"
	filesavedhandler "github.com/zuzuka28/simreport/prj/document/internal/service/documentpipeline/handler/filesaved"
	documentstatussrv "github.com/zuzuka28/simreport/prj/document/internal/service/documentstatus"
	fulltextindexsrv "github.com/zuzuka28/simreport/prj/document/internal/service/fulltextindex"
	semanticindexsrv "github.com/zuzuka28/simreport/prj/document/internal/service/semanticindex"
	shingleindexsrv "github.com/zuzuka28/simreport/prj/document/internal/service/shingleindex"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/lib/minioutil"
	"github.com/zuzuka28/simreport/lib/tikaclient"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
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

func InitTika(
	_ context.Context,
	_ *config.Config,
) (*tikaclient.Client, error) {
	panic(wire.Build(
		wire.Value(http.DefaultClient),
		wire.FieldsOf(new(*config.Config), "Tika"),
		tikaclient.New,
	))
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

func InitNatsJetstream(
	_ *nats.Conn,
) (jetstream.JetStream, error) {
	panic(wire.Build(
		wire.Value([]jetstream.JetStreamOpt(nil)),
		jetstream.New,
	))
}

func ProvideDocumentStatusJetstreamKV(
	ctx context.Context,
	js jetstream.JetStream,
) (jetstream.KeyValue, error) {
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{ //nolint:exhaustruct
		Bucket: "documentstatus",
	})
	if err != nil {
		return nil, fmt.Errorf("new kv: %w", err)
	}

	return kv, nil
}

func ProvideDocumentStatusJetstreamStream(
	ctx context.Context,
	js jetstream.JetStream,
) (jetstream.Stream, error) {
	s, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{ //nolint:exhaustruct
		Name:      "documentstatus",
		Subjects:  []string{"documentstatus.>"},
		Retention: jetstream.InterestPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("new steream: %w", err)
	}

	return s, nil
}

func InitFilestorageRepository(
	_ *minio.Client,
	_ *config.Config,
) (*filestorage.Repository, error) {
	panic(wire.Build(
		filestorage.NewRepository,
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

func InitDocumentRepository(
	_ *elasticsearch.Client,
	_ *config.Config,
) (*documentrepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "DocumentRepo"),
		documentrepo.NewRepository,
	))
}

func InitAttributeRepository(
	_ *elasticsearch.Client,
	_ *config.Config,
) (*attributerepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "AttributeRepo"),
		attributerepo.NewRepository,
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

func InitDocumentStatusRepository(
	_ context.Context,
	_ jetstream.JetStream,
) (*documentstatusrepo.Repository, error) {
	panic(wire.Build(
		ProvideDocumentStatusJetstreamKV,
		wire.Bind(new(jetstream.Publisher), new(jetstream.JetStream)),
		documentstatusrepo.NewRepository,
	))
}

func InitAttributeService(
	_ *attributerepo.Repository,
) (*attributesrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(attributesrv.Repository), new(*attributerepo.Repository)),
		attributesrv.NewService,
	))
}

func InitDocumentStatusService(
	_ *documentstatusrepo.Repository,
) (*documentstatussrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(documentstatussrv.Repository), new(*documentstatusrepo.Repository)),
		documentstatussrv.NewService,
	))
}

func ProvideAnalyzeServiceOpts() analyzesrv.Opts {
	return analyzesrv.Opts{}
}

func InitAnalyzeService(
	_ *config.Config,
	_ *shingleindexsrv.Service,
	_ *fulltextindexsrv.Service,
	_ *documentsrv.Service,
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

func ProvideDocumentServiceOpts() documentsrv.Opts {
	return documentsrv.Opts{} //nolint:exhaustruct
}

func InitDocumentParserService(
	_ *tikaclient.Client,
) (*documentparsersrv.Service, error) {
	panic(wire.Build(
		documentparsersrv.NewService,
	))
}

func InitDocumentService(
	_ *config.Config,
	_ *tikaclient.Client,
	_ *filestorage.Repository,
	_ *documentrepo.Repository,
) (*documentsrv.Service, error) {
	panic(wire.Build(
		ProvideDocumentServiceOpts,
		InitDocumentParserService,
		wire.Bind(new(documentsrv.FileRepository), new(*filestorage.Repository)),
		wire.Bind(new(documentsrv.Repository), new(*documentrepo.Repository)),
		wire.Bind(new(documentsrv.Parser), new(*documentparsersrv.Service)),
		documentsrv.NewService,
	))
}

func InitDocumentHandler(
	_ *documentsrv.Service,
	_ *documentstatussrv.Service,
) *documentapi.Handler {
	panic(wire.Build(
		wire.Bind(new(documentapi.Service), new(*documentsrv.Service)),
		wire.Bind(new(documentapi.StatusService), new(*documentstatussrv.Service)),
		documentapi.NewHandler,
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

func InitAttributeHandler(
	_ *attributesrv.Service,
) *attributeapi.Handler {
	panic(wire.Build(
		wire.Bind(new(attributeapi.Service), new(*attributesrv.Service)),
		attributeapi.NewHandler,
	))
}

func InitRestAPI(
	_ context.Context,
	_ *config.Config,
) (*serverhttp.Server, error) {
	panic(wire.Build(
		ProvideSpec,
		ProvideS3,
		ProvideElastic,
		ProvideNats,
		InitNatsJetstream,
		InitTika,

		InitDocumentStatusRepository,
		InitDocumentStatusService,

		InitShingleIndexRepository,
		InitShingleIndexService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitSemanticIndexRepository,
		InitSemanticIndexService,

		InitFilestorageRepository,

		InitAttributeRepository,

		InitDocumentRepository,
		InitAnalyzeHistoryRepository,

		InitAttributeService,

		InitDocumentService,

		InitAnalyzeService,

		InitDocumentHandler,
		InitAnalyzeHandler,
		InitAttributeHandler,
		wire.Bind(new(serverhttp.DocumentHandler), new(*documentapi.Handler)),
		wire.Bind(new(serverhttp.AnalyzeHandler), new(*analyzeapi.Handler)),
		wire.Bind(new(serverhttp.AttributeHandler), new(*attributeapi.Handler)),
		wire.FieldsOf(new(*config.Config), "Port"),
		wire.Struct(new(serverhttp.Opts), "*"),
		serverhttp.New,
	))
}

func InitFileSavedHandler(
	_ *documentsrv.Service,
) (*filesavedhandler.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(filesavedhandler.DocumentService), new(*documentsrv.Service)),
		filesavedhandler.NewHandler,
	))
}

func InitDocumentNatsHandler(
	_ *documentsrv.Service,
) *documentnats.Handler {
	panic(wire.Build(
		wire.Bind(new(documentnats.Service), new(*documentsrv.Service)),
		documentnats.NewHandler,
	))
}

func InitNatsAPI(
	_ context.Context,
	_ *config.Config,
) (*servernats.Server, error) {
	panic(wire.Build(
		ProvideS3,
		ProvideElastic,
		ProvideNats,
		InitTika,

		InitFilestorageRepository,

		InitDocumentRepository,

		InitDocumentService,

		InitDocumentNatsHandler,
		wire.Bind(new(servernats.DocumentHandler), new(*documentnats.Handler)),
		servernats.NewServer,
	))
}

func ProvideDocumentPipelineStages(
	fsh *filesavedhandler.Handler,
) []documentpipeline.Stage {
	return []documentpipeline.Stage{
		{
			Trigger: model.DocumentProcessingStatusFileSaved,
			Action:  fsh,
			Next:    model.DocumentProcessingStatusDocumentSaved,
		},
		// {
		// 	Trigger: model.DocumentProcessingStatusDocumentSaved,
		// 	Action:  dsh,
		// 	Next:    model.DocumentProcessingStatusDocumentAnalyzed,
		// },
	}
}

func InitDocumentPipeline(
	_ context.Context,
	_ *config.Config,
) (*documentpipeline.Service, error) {
	panic(wire.Build(
		ProvideS3,
		ProvideElastic,
		ProvideNats,
		InitNatsJetstream,
		InitTika,

		ProvideDocumentStatusJetstreamStream,
		InitDocumentStatusRepository,
		InitDocumentStatusService,

		InitFilestorageRepository,

		InitDocumentRepository,

		InitDocumentService,

		InitFileSavedHandler,

		ProvideDocumentPipelineStages,
		wire.Bind(new(jetstream.ConsumerManager), new(jetstream.Stream)),
		wire.Bind(new(documentpipeline.StatusService), new(*documentstatussrv.Service)),

		documentpipeline.NewService,
	))
}
