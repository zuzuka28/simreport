//go:build wireinject

package provider

import (
	"context"
	"fmt"
	"io"
	"os"
	"simrep/api/rest/server"
	analyzeapi "simrep/api/rest/server/handler/analyze"
	anysaveapi "simrep/api/rest/server/handler/anysave"
	documentapi "simrep/api/rest/server/handler/document"
	"simrep/internal/config"
	"simrep/internal/model"
	analyzerepo "simrep/internal/repository/analyze"
	anysaverepo "simrep/internal/repository/anysave"
	documentrepo "simrep/internal/repository/document"
	documentstatusrepo "simrep/internal/repository/documentstatus"
	vectorizerrepo "simrep/internal/repository/vectorizer"
	analyzesrv "simrep/internal/service/analyze"
	anysavesrv "simrep/internal/service/anysave"
	documentsrv "simrep/internal/service/document"
	"simrep/internal/service/documentpipeline"
	documentsavedhandler "simrep/internal/service/documentpipeline/handler/documentsaved"
	filesavedhandler "simrep/internal/service/documentpipeline/handler/filesaved"
	documentstatussrv "simrep/internal/service/documentstatus"
	"simrep/pkg/elasticutil"
	"simrep/pkg/minioutil"
	vectorizerclient "simrep/pkg/vectorizerclient"

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

func InitElastic(
	_ context.Context,
	_ *config.Config,
) (*elasticsearch.Client, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "Elastic"),
		elasticutil.NewClientWithStartup,
	))
}

func InitNats(
	_ context.Context,
	_ *config.Config,
) (*nats.Conn, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "Nats"),
		wire.Value([]nats.Option(nil)),
		nats.Connect,
	))
}

func InitS3(
	_ context.Context,
	_ *config.Config,
) (*minio.Client, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "S3"),
		minioutil.NewClientWithStartup,
	))
}

func InitNatsJetstream(
	conn *nats.Conn,
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
		Retention: jetstream.WorkQueuePolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("new steream: %w", err)
	}

	return s, nil
}

func InitVectorizerClient(
	_ *config.Config,
) (*vectorizerclient.ClientWithResponses, error) {
	panic(wire.Build(
		wire.Value([]vectorizerclient.ClientOption(nil)),
		wire.FieldsOf(new(*config.Config), "VectorizerService"),
		vectorizerclient.NewClientWithResponses,
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

func InitAnalyzedDocumentRepository(
	_ *elasticsearch.Client,
	_ *config.Config,
) (*analyzerepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "AnalyzedDocumentRepo"),
		analyzerepo.NewRepository,
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

func InitVectorizerService(
	_ *vectorizerclient.ClientWithResponses,
) (*vectorizerrepo.Repository, error) {
	panic(wire.Build(
		wire.Bind(new(vectorizerclient.ClientWithResponsesInterface), new(*vectorizerclient.ClientWithResponses)),
		vectorizerrepo.NewRepository,
	))
}

func ProvideAnalyzeServiceOpts() analyzesrv.Opts {
	return analyzesrv.Opts{} //nolint:exhaustruct
}

func InitAnalyzeService(
	_ *config.Config,
	_ *vectorizerrepo.Repository,
	_ *analyzerepo.Repository,
) (*analyzesrv.Service, error) {
	panic(wire.Build(
		ProvideAnalyzeServiceOpts,
		wire.Bind(new(analyzesrv.VectorizerService), new(*vectorizerrepo.Repository)),
		wire.Bind(new(analyzesrv.Repository), new(*analyzerepo.Repository)),
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

func InitDocumentService(
	_ *config.Config,
	_ *anysavesrv.Service,
	_ *documentrepo.Repository,
) (*documentsrv.Service, error) {
	panic(wire.Build(
		ProvideDocumentServiceOpts,
		wire.Bind(new(documentsrv.FileRepository), new(*anysavesrv.Service)),
		wire.Bind(new(documentsrv.Repository), new(*documentrepo.Repository)),
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
	_ *documentsrv.Service,
	_ *analyzesrv.Service,
) *analyzeapi.Handler {
	panic(wire.Build(
		wire.Bind(new(analyzeapi.Service), new(*analyzesrv.Service)),
		wire.Bind(new(analyzeapi.DocumentParser), new(*documentsrv.Service)),
		analyzeapi.NewHandler,
	))
}

func InitRestAPI(
	_ context.Context,
	_ *config.Config,
) (*server.Server, error) {
	panic(wire.Build(
		ProvideSpec,
		InitS3,
		InitElastic,
		InitNats,
		InitNatsJetstream,
		InitDocumentStatusRepository,
		InitDocumentStatusService,
		InitDocumentRepository,
		InitVectorizerClient,
		InitVectorizerService,
		InitAnysaveService,
		InitDocumentService,
		InitDocumentHandler,
		InitAnalyzedDocumentRepository,
		InitAnalyzeService,
		InitAnalyzeHandler,
		InitAnysaveHandler,
		wire.Bind(new(server.DocumentHandler), new(*documentapi.Handler)),
		wire.Bind(new(server.AnalyzeHandler), new(*analyzeapi.Handler)),
		wire.Bind(new(server.FileHandler), new(*anysaveapi.Handler)),
		wire.FieldsOf(new(*config.Config), "Port"),
		wire.Struct(new(server.Opts), "*"),
		server.New,
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

func InitDocumentSavedHandler(
	_ *documentsrv.Service,
	_ *analyzesrv.Service,
) (*documentsavedhandler.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(documentsavedhandler.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(documentsavedhandler.AnalyzeService), new(*analyzesrv.Service)),
		documentsavedhandler.NewHandler,
	))
}

func ProvideDocumentPipelineStages(
	dsh *documentsavedhandler.Handler,
	fsh *filesavedhandler.Handler,
) []documentpipeline.Stage {
	return []documentpipeline.Stage{
		{
			Trigger: model.DocumentProcessingStatusFileSaved,
			Action:  fsh,
			Next:    model.DocumentProcessingStatusDocumentSaved,
		},
		{
			Trigger: model.DocumentProcessingStatusDocumentSaved,
			Action:  dsh,
			Next:    model.DocumentProcessingStatusDocumentAnalyzed,
		},
	}
}

func InitDocumentPipeline(
	_ context.Context,
	_ *config.Config,
) (*documentpipeline.Service, error) {
	panic(wire.Build(
		InitS3,
		InitElastic,
		InitNats,
		InitNatsJetstream,
		InitDocumentStatusRepository,
		InitDocumentStatusService,
		ProvideDocumentStatusJetstreamStream,
		InitDocumentRepository,
		InitVectorizerClient,
		InitVectorizerService,
		InitAnysaveService,
		InitDocumentService,
		InitAnalyzedDocumentRepository,
		InitAnalyzeService,
		InitFileSavedHandler,
		InitDocumentSavedHandler,
		ProvideDocumentPipelineStages,
		wire.Bind(new(jetstream.ConsumerManager), new(jetstream.Stream)),
		wire.Bind(new(documentpipeline.StatusService), new(*documentstatussrv.Service)),
		documentpipeline.NewService,
	))
}
