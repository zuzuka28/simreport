//go:build wireinject

package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	documentnats "simrep/api/nats/handler/document"
	servernats "simrep/api/nats/server"
	serverhttp "simrep/api/rest/server"
	analyzeapi "simrep/api/rest/server/handler/analyze"
	anysaveapi "simrep/api/rest/server/handler/anysave"
	documentapi "simrep/api/rest/server/handler/document"
	"simrep/internal/config"
	"simrep/internal/model"
	anysaverepo "simrep/internal/repository/anysave"
	documentrepo "simrep/internal/repository/document"
	documentstatusrepo "simrep/internal/repository/documentstatus"
	fulltextindexrepo "simrep/internal/repository/fulltextindexclient"
	semanticindexrepo "simrep/internal/repository/semanticindexclient"
	shingleindexrepo "simrep/internal/repository/shingleindexclient"
	analyzesrv "simrep/internal/service/analyze"
	anysavesrv "simrep/internal/service/anysave"
	documentsrv "simrep/internal/service/document"
	"simrep/internal/service/documentparser"
	documentparsersrv "simrep/internal/service/documentparser"
	"simrep/internal/service/documentpipeline"
	filesavedhandler "simrep/internal/service/documentpipeline/handler/filesaved"
	documentstatussrv "simrep/internal/service/documentstatus"
	fulltextindexsrv "simrep/internal/service/fulltextindex"
	semanticindexsrv "simrep/internal/service/semanticindex"
	shingleindexsrv "simrep/internal/service/shingleindex"
	"simrep/pkg/elasticutil"
	"simrep/pkg/minioutil"
	"simrep/pkg/tikaclient"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

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

func InitConfig(path string) (*config.Config, error) {
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

	return elasticCli, err
}

//nolint:gochecknoglobals
var (
	natsCli     *nats.Conn
	natsCliOnce sync.Once
)

func ProvideNats(
	ctx context.Context,
	cfg *config.Config,
) (*nats.Conn, error) {
	var err error

	natsCliOnce.Do(func() {
		natsCli, err = nats.Connect(cfg.Nats)
	})

	return natsCli, err
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

	return s3Cli, err
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

func InitShingleIndexRepository(
	_ *nats.Conn,
) (*shingleindexrepo.Repository, error) {
	panic(wire.Build(
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

func InitDocumentFileRepository(
	_ *minio.Client,
	_ *config.Config,
) (*anysaverepo.Repository, error) {
	panic(wire.Build(
		anysaverepo.NewRepository,
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

func InitDocumentStatusRepository(
	ctx context.Context,
	js jetstream.JetStream,
) (*documentstatusrepo.Repository, error) {
	panic(wire.Build(
		ProvideDocumentStatusJetstreamKV,
		wire.Bind(new(jetstream.Publisher), new(jetstream.JetStream)),
		documentstatusrepo.NewRepository,
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
	return analyzesrv.Opts{} //nolint:exhaustruct
}

func InitAnalyzeService(
	_ *config.Config,
	_ *shingleindexsrv.Service,
	_ *fulltextindexsrv.Service,
	_ *documentsrv.Service,
	_ *semanticindexsrv.Service,
) (*analyzesrv.Service, error) {
	panic(wire.Build(
		ProvideAnalyzeServiceOpts,
		wire.Bind(new(analyzesrv.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(analyzesrv.ShingleIndexService), new(*shingleindexsrv.Service)),
		wire.Bind(new(analyzesrv.FulltextIndexService), new(*fulltextindexsrv.Service)),
		wire.Bind(new(analyzesrv.SemanticIndexService), new(*semanticindexsrv.Service)),
		analyzesrv.NewService,
	))
}

func ProvideAnysaveServiceOpts() anysavesrv.Opts {
	return anysavesrv.Opts{} //nolint:exhaustruct
}

func InitAnysaveService(
	_ *minio.Client,
	_ *config.Config,
) (*anysavesrv.Service, error) {
	panic(wire.Build(
		InitDocumentFileRepository,
		ProvideAnysaveServiceOpts,
		wire.Bind(new(anysavesrv.Repository), new(*anysaverepo.Repository)),
		anysavesrv.NewService,
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
	_ *anysavesrv.Service,
	_ *documentrepo.Repository,
) (*documentsrv.Service, error) {
	panic(wire.Build(
		ProvideDocumentServiceOpts,
		InitDocumentParserService,
		wire.Bind(new(documentsrv.FileRepository), new(*anysavesrv.Service)),
		wire.Bind(new(documentsrv.Repository), new(*documentrepo.Repository)),
		wire.Bind(new(documentsrv.Parser), new(*documentparser.Service)),
		documentsrv.NewService,
	))
}

func InitDocumentHandler(
	_ *documentsrv.Service,
) *documentapi.Handler {
	panic(wire.Build(
		wire.Bind(new(documentapi.Service), new(*documentsrv.Service)),
		documentapi.NewHandler,
	))
}

func InitAnysaveHandler(
	_ *documentstatussrv.Service,
	_ *anysavesrv.Service,
) *anysaveapi.Handler {
	panic(wire.Build(
		wire.Bind(new(anysaveapi.Service), new(*anysavesrv.Service)),
		wire.Bind(new(anysaveapi.StatusService), new(*documentstatussrv.Service)),
		anysaveapi.NewHandler,
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
		ProvideS3,
		ProvideElastic,
		ProvideNats,
		InitNatsJetstream,
		InitTika,

		InitShingleIndexRepository,
		InitShingleIndexService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitSemanticIndexRepository,
		InitSemanticIndexService,

		// ProvideDocumentStatusJetstreamStream,
		InitDocumentStatusRepository,
		InitDocumentStatusService,

		InitDocumentRepository,

		InitAnysaveService,
		InitDocumentService,

		InitAnalyzeService,

		InitDocumentHandler,
		InitAnalyzeHandler,
		InitAnysaveHandler,
		wire.Bind(new(serverhttp.DocumentHandler), new(*documentapi.Handler)),
		wire.Bind(new(serverhttp.AnalyzeHandler), new(*analyzeapi.Handler)),
		wire.Bind(new(serverhttp.FileHandler), new(*anysaveapi.Handler)),
		wire.FieldsOf(new(*config.Config), "Port"),
		wire.Struct(new(serverhttp.Opts), "*"),
		serverhttp.New,
	))
}

func InitFileSavedHandler(
	_ *documentsrv.Service,
	_ *anysavesrv.Service,
) (*filesavedhandler.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(filesavedhandler.FileService), new(*anysavesrv.Service)),
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

		InitDocumentRepository,

		InitAnysaveService,
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

		InitDocumentRepository,

		InitAnysaveService,
		InitDocumentService,

		InitFileSavedHandler,

		ProvideDocumentPipelineStages,
		wire.Bind(new(jetstream.ConsumerManager), new(jetstream.Stream)),
		wire.Bind(new(documentpipeline.StatusService), new(*documentstatussrv.Service)),

		documentpipeline.NewService,
	))
}
