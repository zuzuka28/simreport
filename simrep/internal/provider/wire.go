//go:build wireinject

package provider

import (
	"context"
	"io"
	"os"
	asyncanalyzeconsumer "simrep/api/amqp/asyncanalyze/consumer"
	asyncanalyzeproducer "simrep/api/amqp/asyncanalyze/producer"
	"simrep/api/rest/server"
	analyzeapi "simrep/api/rest/server/handler/analyze"
	documentapi "simrep/api/rest/server/handler/document"
	"simrep/internal/config"
	analyzerepo "simrep/internal/repository/analyze"
	documentrepo "simrep/internal/repository/document"
	documentfilerepo "simrep/internal/repository/documentfile"
	imagerepo "simrep/internal/repository/image"
	analyzesrv "simrep/internal/service/analyze"
	documentsrv "simrep/internal/service/document"
	vectorizersrv "simrep/internal/service/vectorizer"
	"simrep/pkg/elasticutil"
	"simrep/pkg/minioutil"
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

func InitDocumentFileRepository(
	_ *minio.Client,
	_ *config.Config,
) (*documentfilerepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "DocumentFileRepo"),
		documentfilerepo.NewRepository,
	))
}

func InitImageRepository(
	_ *minio.Client,
	_ *config.Config,
) (*imagerepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "ImageRepo"),
		imagerepo.NewRepository,
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

func InitAsyncAnalyzeProducerService(
	_ *config.Config,
) (*asyncanalyzeproducer.Producer, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "AnalyzeProducer"),
		asyncanalyzeproducer.New,
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
	_ *vectorizersrv.Service,
	_ *analyzerepo.Repository,
) (*analyzesrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(analyzesrv.VectorizerService), new(*vectorizersrv.Service)),
		wire.Bind(new(analyzesrv.Repository), new(*analyzerepo.Repository)),
		analyzesrv.NewService,
	))
}

func InitDocumentService(
	_ *imagerepo.Repository,
	_ *documentfilerepo.Repository,
	_ *documentrepo.Repository,
	_ *asyncanalyzeproducer.Producer,
) (*documentsrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(documentsrv.AnalyzeService), new(*asyncanalyzeproducer.Producer)),
		wire.Bind(new(documentsrv.FileRepository), new(*documentfilerepo.Repository)),
		wire.Bind(new(documentsrv.ImageRepository), new(*imagerepo.Repository)),
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
		InitImageRepository,
		InitDocumentFileRepository,
		InitDocumentRepository,
		InitDocumentService,
		InitDocumentHandler,
		InitAnalyzedDocumentRepository,
		InitAnalyzeService,
		InitAnalyzeHandler,
		InitAsyncAnalyzeProducerService,
		wire.Bind(new(server.DocumentHandler), new(*documentapi.Handler)),
		wire.Bind(new(server.AnalyzeHandler), new(*analyzeapi.Handler)),
		wire.FieldsOf(new(*config.Config), "Port"),
		wire.Struct(new(server.Opts), "*"),
		server.New,
	))
}

func InitAsyncAnalyzeAPI(
	_ context.Context,
	_ *config.Config,
) (*asyncanalyzeconsumer.Consumer, error) {
	panic(wire.Build(
		InitS3,
		InitElastic,
		InitImageRepository,
		InitDocumentFileRepository,
		InitDocumentRepository,
		InitDocumentService,
		InitVectorizerClient,
		InitVectorizerService,
		InitAnalyzedDocumentRepository,
		InitAnalyzeService,
		InitAsyncAnalyzeProducerService,
		wire.Bind(new(asyncanalyzeconsumer.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(asyncanalyzeconsumer.AnalyzeService), new(*analyzesrv.Service)),
		wire.FieldsOf(new(*config.Config), "AnalyzeConsumer"),
		asyncanalyzeconsumer.New,
	))
}
