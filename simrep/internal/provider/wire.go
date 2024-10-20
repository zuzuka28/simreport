//go:build wireinject

package provider

import (
	"context"
	"io"
	"os"
	notifyprocessing "simrep/api/amqp/asyncnotify/consumer"
	documentsavedapi "simrep/api/amqp/asyncnotify/handler/documentsaved"
	filesavedapi "simrep/api/amqp/asyncnotify/handler/filesaved"
	notify "simrep/api/amqp/asyncnotify/producer"
	"simrep/api/rest/server"
	analyzeapi "simrep/api/rest/server/handler/analyze"
	documentapi "simrep/api/rest/server/handler/document"
	fileapi "simrep/api/rest/server/handler/file"
	"simrep/internal/config"
	analyzerepo "simrep/internal/repository/analyze"
	documentrepo "simrep/internal/repository/document"
	filerepo "simrep/internal/repository/file"
	analyzesrv "simrep/internal/service/analyze"
	documentsrv "simrep/internal/service/document"
	documentfilesrv "simrep/internal/service/documentfile"
	imagefilesrv "simrep/internal/service/imagefile"
	vectorizersrv "simrep/internal/service/vectorizer"
	"simrep/pkg/elasticutil"
	"simrep/pkg/minioutil"
	"simrep/pkg/rabbitmq"
	vectorizerclient "simrep/pkg/vectorizerclient"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
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

func InitS3(
	_ context.Context,
	_ *config.Config,
) (*minio.Client, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "S3"),
		minioutil.NewClientWithStartup,
	))
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

func InitRabbitNotifyFileSavedPublisher(
	_ *config.Config,
) (*rabbitmq.Producer, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "NotifyFileSavedProducer"),
		rabbitmq.NewProducer,
	))
}

func InitRabbitNotifyDocumentSavedPublisher(
	_ *config.Config,
) (*rabbitmq.Producer, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "NotifyDocumentSavedProducer"),
		rabbitmq.NewProducer,
	))
}

func InitRabbitNotifyDocumentAnalyzedPublisher(
	_ *config.Config,
) (*rabbitmq.Producer, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "NotifyDocumentAnalyzedProducer"),
		rabbitmq.NewProducer,
	))
}

func InitRabbitNotifyFileSavedConsumer(
	_ *config.Config,
) (*rabbitmq.Consumer, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "NotifyFileSavedConsumer"),
		rabbitmq.NewConsumer,
	))
}

func InitRabbitNotifyDocumentSavedConsumer(
	_ *config.Config,
) (*rabbitmq.Consumer, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "NotifyDocumentSavedConsumer"),
		rabbitmq.NewConsumer,
	))
}

func InitRabbitNotifyDocumentAnalyzedConsumer(
	_ *config.Config,
) (*rabbitmq.Consumer, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "NotifyDocumentAnalyzedConsumer"),
		rabbitmq.NewConsumer,
	))
}

func InitNotifyFileSavedProducer(
	_ *config.Config,
) (*notify.Producer, error) {
	panic(wire.Build(
		InitRabbitNotifyFileSavedPublisher,
		wire.Bind(new(notify.Publisher), new(*rabbitmq.Producer)),
		notify.New,
	))
}

func InitNotifyDocumentSavedProducer(
	_ *config.Config,
) (*notify.Producer, error) {
	panic(wire.Build(
		InitRabbitNotifyDocumentSavedPublisher,
		wire.Bind(new(notify.Publisher), new(*rabbitmq.Producer)),
		notify.New,
	))
}

func InitNotifyDocumentAnalyzedProducer(
	_ *config.Config,
) (*notify.Producer, error) {
	panic(wire.Build(
		InitRabbitNotifyDocumentAnalyzedPublisher,
		wire.Bind(new(notify.Publisher), new(*rabbitmq.Producer)),
		notify.New,
	))
}

func InitDocumentFileRepository(
	_ *minio.Client,
	_ *config.Config,
) (*filerepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "DocumentFileRepo"),
		filerepo.NewRepository,
	))
}

func InitImageFileRepository(
	_ *minio.Client,
	_ *config.Config,
) (*filerepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "ImageRepo"),
		filerepo.NewRepository,
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

func InitVectorizerService(
	_ *vectorizerclient.ClientWithResponses,
) (*vectorizersrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(vectorizerclient.ClientWithResponsesInterface), new(*vectorizerclient.ClientWithResponses)),
		vectorizersrv.NewService,
	))
}

func InitAnalyzeService(
	_ *config.Config,
	_ *vectorizersrv.Service,
	_ *analyzerepo.Repository,
) (*analyzesrv.Service, error) {
	panic(wire.Build(
		InitNotifyDocumentAnalyzedProducer,
		wire.Bind(new(analyzesrv.Notify), new(*notify.Producer)),
		wire.Bind(new(analyzesrv.VectorizerService), new(*vectorizersrv.Service)),
		wire.Bind(new(analyzesrv.Repository), new(*analyzerepo.Repository)),
		analyzesrv.NewService,
	))
}

func InitDocumentFileService(
	_ *minio.Client,
	_ *config.Config,
) (*documentfilesrv.Service, error) {
	panic(wire.Build(
		InitNotifyFileSavedProducer,
		InitDocumentFileRepository,
		wire.Bind(new(documentfilesrv.Repository), new(*filerepo.Repository)),
		wire.Bind(new(documentfilesrv.Notify), new(*notify.Producer)),
		documentfilesrv.NewService,
	))
}

func InitImageFileService(
	_ *minio.Client,
	_ *config.Config,
) (*imagefilesrv.Service, error) {
	panic(wire.Build(
		InitImageFileRepository,
		wire.Bind(new(imagefilesrv.Repository), new(*filerepo.Repository)),
		imagefilesrv.NewService,
	))
}

func InitDocumentService(
	_ *config.Config,
	_ *imagefilesrv.Service,
	_ *documentfilesrv.Service,
	_ *documentrepo.Repository,
) (*documentsrv.Service, error) {
	panic(wire.Build(
		InitNotifyDocumentSavedProducer,
		wire.Bind(new(documentsrv.Notify), new(*notify.Producer)),
		wire.Bind(new(documentsrv.FileRepository), new(*documentfilesrv.Service)),
		wire.Bind(new(documentsrv.ImageRepository), new(*imagefilesrv.Service)),
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

func InitFileHandler(
	_ *documentfilesrv.Service,
) *fileapi.Handler {
	panic(wire.Build(
		wire.Bind(new(fileapi.Service), new(*documentfilesrv.Service)),
		fileapi.NewHandler,
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
		InitVectorizerClient,
		InitVectorizerService,
		InitDocumentFileService,
		InitImageFileService,
		InitDocumentRepository,
		InitDocumentService,
		InitDocumentHandler,
		InitAnalyzedDocumentRepository,
		InitAnalyzeService,
		InitAnalyzeHandler,
		InitFileHandler,
		wire.Bind(new(server.DocumentHandler), new(*documentapi.Handler)),
		wire.Bind(new(server.AnalyzeHandler), new(*analyzeapi.Handler)),
		wire.Bind(new(server.FileHandler), new(*fileapi.Handler)),
		wire.FieldsOf(new(*config.Config), "Port"),
		wire.Struct(new(server.Opts), "*"),
		server.New,
	))
}

func InitAsyncFileSavedHandler(
	_ *documentfilesrv.Service,
	_ *documentsrv.Service,
) *filesavedapi.Handler {
	panic(wire.Build(
		wire.Bind(new(filesavedapi.FileService), new(*documentfilesrv.Service)),
		wire.Bind(new(filesavedapi.DocumentService), new(*documentsrv.Service)),
		filesavedapi.NewHandler,
	))
}

func InitAsyncDocumentSavedHandler(
	_ *documentsrv.Service,
	_ *analyzesrv.Service,
) *documentsavedapi.Handler {
	panic(wire.Build(
		wire.Bind(new(documentsavedapi.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(documentsavedapi.AnalyzeService), new(*analyzesrv.Service)),
		documentsavedapi.NewHandler,
	))
}

func InitAsyncDocumentParsing(
	_ context.Context,
	_ *config.Config,
) (*notifyprocessing.Consumer, error) {
	panic(wire.Build(
		InitS3,
		InitElastic,
		InitRabbitNotifyFileSavedConsumer,
		InitDocumentFileService,
		InitImageFileService,
		InitDocumentRepository,
		InitDocumentService,
		InitAsyncFileSavedHandler,
		wire.Bind(new(notifyprocessing.Handler), new(*filesavedapi.Handler)),
		wire.Bind(new(notifyprocessing.RMQConsumer), new(*rabbitmq.Consumer)),
		notifyprocessing.New,
	))
}

func InitAsyncDocumentAnalysis(
	_ context.Context,
	_ *config.Config,
) (*notifyprocessing.Consumer, error) {
	panic(wire.Build(
		InitS3,
		InitElastic,
		InitRabbitNotifyDocumentSavedConsumer,
		InitDocumentFileService,
		InitImageFileService,
		InitDocumentRepository,
		InitDocumentService,
		InitVectorizerClient,
		InitVectorizerService,
		InitAnalyzedDocumentRepository,
		InitAnalyzeService,
		InitAsyncDocumentSavedHandler,
		wire.Bind(new(notifyprocessing.Handler), new(*documentsavedapi.Handler)),
		wire.Bind(new(notifyprocessing.RMQConsumer), new(*rabbitmq.Consumer)),
		notifyprocessing.New,
	))
}
