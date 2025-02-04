// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package provider

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/lib/minioutil"
	"github.com/zuzuka28/simreport/lib/tikaclient"
	attribute4 "github.com/zuzuka28/simreport/prj/document/api/nats/handler/attribute"
	document4 "github.com/zuzuka28/simreport/prj/document/api/nats/handler/document"
	similarity3 "github.com/zuzuka28/simreport/prj/document/api/nats/handler/similarity"
	server2 "github.com/zuzuka28/simreport/prj/document/api/nats/server"
	"github.com/zuzuka28/simreport/prj/document/api/rest/server"
	attribute3 "github.com/zuzuka28/simreport/prj/document/api/rest/server/handler/attribute"
	document3 "github.com/zuzuka28/simreport/prj/document/api/rest/server/handler/document"
	similarity2 "github.com/zuzuka28/simreport/prj/document/api/rest/server/handler/similarity"
	"github.com/zuzuka28/simreport/prj/document/internal/config"
	"github.com/zuzuka28/simreport/prj/document/internal/model"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/analyzehistory"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/attribute"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/document"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/documentstatus"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/filestorage"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/fulltextindexclient"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/semanticindexclient"
	"github.com/zuzuka28/simreport/prj/document/internal/repository/shingleindexclient"
	attribute2 "github.com/zuzuka28/simreport/prj/document/internal/service/attribute"
	document2 "github.com/zuzuka28/simreport/prj/document/internal/service/document"
	"github.com/zuzuka28/simreport/prj/document/internal/service/documentparser"
	"github.com/zuzuka28/simreport/prj/document/internal/service/documentpipeline"
	"github.com/zuzuka28/simreport/prj/document/internal/service/documentpipeline/handler/filesaved"
	documentstatus2 "github.com/zuzuka28/simreport/prj/document/internal/service/documentstatus"
	"github.com/zuzuka28/simreport/prj/document/internal/service/fulltextindex"
	"github.com/zuzuka28/simreport/prj/document/internal/service/semanticindex"
	"github.com/zuzuka28/simreport/prj/document/internal/service/shingleindex"
	"github.com/zuzuka28/simreport/prj/document/internal/service/similarity"
	"io"
	"net/http"
	"os"
	"sync"
)

// Injectors from wire.go:

func InitConfig(string2 string) (*config.Config, error) {
	configConfig, err := config.New(string2)
	if err != nil {
		return nil, err
	}
	return configConfig, nil
}

func InitTika(contextContext context.Context, configConfig *config.Config) (*tikaclient.Client, error) {
	client := _wireClientValue
	string2 := configConfig.Tika
	tikaclientClient := tikaclient.New(client, string2)
	return tikaclientClient, nil
}

var (
	_wireClientValue = http.DefaultClient
)

func InitNatsJetstream(conn *nats.Conn) (jetstream.JetStream, error) {
	v := _wireValue
	jetStream, err := jetstream.New(conn, v...)
	if err != nil {
		return nil, err
	}
	return jetStream, nil
}

var (
	_wireValue = []jetstream.JetStreamOpt(nil)
)

func InitFilestorageRepository(client *minio.Client, configConfig *config.Config) (*filestorage.Repository, error) {
	repository := filestorage.NewRepository(client)
	return repository, nil
}

func InitShingleIndexRepository(conn *nats.Conn) (*shingleindexclient.Repository, error) {
	repository := shingleindexclient.NewRepository(conn)
	return repository, nil
}

func InitShingleIndexService(repository *shingleindexclient.Repository) (*shingleindex.Service, error) {
	service := shingleindex.NewService(repository)
	return service, nil
}

func InitFulltextIndexRepository(conn *nats.Conn) (*fulltextindexclient.Repository, error) {
	repository := fulltextindexclient.NewRepository(conn)
	return repository, nil
}

func InitFulltextIndexService(repository *fulltextindexclient.Repository) (*fulltextindex.Service, error) {
	service := fulltextindex.NewService(repository)
	return service, nil
}

func InitSemanticIndexRepository(conn *nats.Conn) (*semanticindexclient.Repository, error) {
	repository := semanticindexclient.NewRepository(conn)
	return repository, nil
}

func InitSemanticIndexService(repository *semanticindexclient.Repository) (*semanticindex.Service, error) {
	service := semanticindex.NewService(repository)
	return service, nil
}

func InitDocumentRepository(client *elasticsearch.Client, configConfig *config.Config) (*document.Repository, error) {
	opts := configConfig.DocumentRepo
	repository, err := document.NewRepository(opts, client)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func InitAttributeRepository(client *elasticsearch.Client, configConfig *config.Config) (*attribute.Repository, error) {
	opts := configConfig.AttributeRepo
	repository, err := attribute.NewRepository(opts, client)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func InitAnalyzeHistoryRepository(client *elasticsearch.Client, configConfig *config.Config) (*analyzehistory.Repository, error) {
	opts := configConfig.AnalyzeHistoryRepo
	repository, err := analyzehistory.NewRepository(opts, client)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func InitDocumentStatusRepository(contextContext context.Context, jetStream jetstream.JetStream) (*documentstatus.Repository, error) {
	keyValue, err := ProvideDocumentStatusJetstreamKV(contextContext, jetStream)
	if err != nil {
		return nil, err
	}
	repository := documentstatus.NewRepository(keyValue, jetStream)
	return repository, nil
}

func InitAttributeService(repository *attribute.Repository) (*attribute2.Service, error) {
	service := attribute2.NewService(repository)
	return service, nil
}

func InitDocumentStatusService(repository *documentstatus.Repository) (*documentstatus2.Service, error) {
	service := documentstatus2.NewService(repository)
	return service, nil
}

func InitAnalyzeService(configConfig *config.Config, service *shingleindex.Service, fulltextindexService *fulltextindex.Service, documentService *document2.Service, semanticindexService *semanticindex.Service, repository *analyzehistory.Repository) (*similarity.Service, error) {
	opts := ProvideAnalyzeServiceOpts()
	similarityService := similarity.NewService(opts, documentService, service, fulltextindexService, semanticindexService, repository)
	return similarityService, nil
}

func InitDocumentParserService(client *tikaclient.Client) (*documentparser.Service, error) {
	service := documentparser.NewService(client)
	return service, nil
}

func InitDocumentService(configConfig *config.Config, client *tikaclient.Client, repository *filestorage.Repository, documentRepository *document.Repository) (*document2.Service, error) {
	opts := ProvideDocumentServiceOpts()
	service, err := InitDocumentParserService(client)
	if err != nil {
		return nil, err
	}
	documentService := document2.NewService(opts, documentRepository, repository, service)
	return documentService, nil
}

func InitDocumentHandler(service *document2.Service, documentstatusService *documentstatus2.Service) *document3.Handler {
	handler := document3.NewHandler(service, documentstatusService)
	return handler
}

func InitAnalyzeHandler(service *similarity.Service) *similarity2.Handler {
	handler := similarity2.NewHandler(service)
	return handler
}

func InitAttributeHandler(service *attribute2.Service) *attribute3.Handler {
	handler := attribute3.NewHandler(service)
	return handler
}

func InitRestAPI(contextContext context.Context, configConfig *config.Config) (*server.Server, error) {
	int2 := configConfig.Port
	v, err := ProvideSpec()
	if err != nil {
		return nil, err
	}
	client, err := InitTika(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	minioClient, err := ProvideS3(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	repository, err := InitFilestorageRepository(minioClient, configConfig)
	if err != nil {
		return nil, err
	}
	elasticsearchClient, err := ProvideElastic(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	documentRepository, err := InitDocumentRepository(elasticsearchClient, configConfig)
	if err != nil {
		return nil, err
	}
	service, err := InitDocumentService(configConfig, client, repository, documentRepository)
	if err != nil {
		return nil, err
	}
	conn, err := ProvideNats(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	jetStream, err := InitNatsJetstream(conn)
	if err != nil {
		return nil, err
	}
	documentstatusRepository, err := InitDocumentStatusRepository(contextContext, jetStream)
	if err != nil {
		return nil, err
	}
	documentstatusService, err := InitDocumentStatusService(documentstatusRepository)
	if err != nil {
		return nil, err
	}
	handler := InitDocumentHandler(service, documentstatusService)
	shingleindexclientRepository, err := InitShingleIndexRepository(conn)
	if err != nil {
		return nil, err
	}
	shingleindexService, err := InitShingleIndexService(shingleindexclientRepository)
	if err != nil {
		return nil, err
	}
	fulltextindexclientRepository, err := InitFulltextIndexRepository(conn)
	if err != nil {
		return nil, err
	}
	fulltextindexService, err := InitFulltextIndexService(fulltextindexclientRepository)
	if err != nil {
		return nil, err
	}
	semanticindexclientRepository, err := InitSemanticIndexRepository(conn)
	if err != nil {
		return nil, err
	}
	semanticindexService, err := InitSemanticIndexService(semanticindexclientRepository)
	if err != nil {
		return nil, err
	}
	analyzehistoryRepository, err := InitAnalyzeHistoryRepository(elasticsearchClient, configConfig)
	if err != nil {
		return nil, err
	}
	similarityService, err := InitAnalyzeService(configConfig, shingleindexService, fulltextindexService, service, semanticindexService, analyzehistoryRepository)
	if err != nil {
		return nil, err
	}
	similarityHandler := InitAnalyzeHandler(similarityService)
	attributeRepository, err := InitAttributeRepository(elasticsearchClient, configConfig)
	if err != nil {
		return nil, err
	}
	attributeService, err := InitAttributeService(attributeRepository)
	if err != nil {
		return nil, err
	}
	attributeHandler := InitAttributeHandler(attributeService)
	opts := server.Opts{
		Port:             int2,
		Spec:             v,
		DocumentHandler:  handler,
		AnalyzeHandler:   similarityHandler,
		AttributeHandler: attributeHandler,
	}
	serverServer, err := server.New(opts)
	if err != nil {
		return nil, err
	}
	return serverServer, nil
}

func InitDocumentNatsHandler(service *document2.Service, documentstatusService *documentstatus2.Service) *document4.Handler {
	handler := document4.NewHandler(service, documentstatusService)
	return handler
}

func InitAnalyzeNatsHandler(service *similarity.Service) *similarity3.Handler {
	handler := similarity3.NewHandler(service)
	return handler
}

func InitAttributeNatsHandler(service *attribute2.Service) *attribute4.Handler {
	handler := attribute4.NewHandler(service)
	return handler
}

func InitNatsAPI(contextContext context.Context, configConfig *config.Config) (*server2.Server, error) {
	conn, err := ProvideNats(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	client, err := InitTika(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	minioClient, err := ProvideS3(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	repository, err := InitFilestorageRepository(minioClient, configConfig)
	if err != nil {
		return nil, err
	}
	elasticsearchClient, err := ProvideElastic(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	documentRepository, err := InitDocumentRepository(elasticsearchClient, configConfig)
	if err != nil {
		return nil, err
	}
	service, err := InitDocumentService(configConfig, client, repository, documentRepository)
	if err != nil {
		return nil, err
	}
	jetStream, err := InitNatsJetstream(conn)
	if err != nil {
		return nil, err
	}
	documentstatusRepository, err := InitDocumentStatusRepository(contextContext, jetStream)
	if err != nil {
		return nil, err
	}
	documentstatusService, err := InitDocumentStatusService(documentstatusRepository)
	if err != nil {
		return nil, err
	}
	handler := InitDocumentNatsHandler(service, documentstatusService)
	attributeRepository, err := InitAttributeRepository(elasticsearchClient, configConfig)
	if err != nil {
		return nil, err
	}
	attributeService, err := InitAttributeService(attributeRepository)
	if err != nil {
		return nil, err
	}
	attributeHandler := InitAttributeNatsHandler(attributeService)
	shingleindexclientRepository, err := InitShingleIndexRepository(conn)
	if err != nil {
		return nil, err
	}
	shingleindexService, err := InitShingleIndexService(shingleindexclientRepository)
	if err != nil {
		return nil, err
	}
	fulltextindexclientRepository, err := InitFulltextIndexRepository(conn)
	if err != nil {
		return nil, err
	}
	fulltextindexService, err := InitFulltextIndexService(fulltextindexclientRepository)
	if err != nil {
		return nil, err
	}
	semanticindexclientRepository, err := InitSemanticIndexRepository(conn)
	if err != nil {
		return nil, err
	}
	semanticindexService, err := InitSemanticIndexService(semanticindexclientRepository)
	if err != nil {
		return nil, err
	}
	analyzehistoryRepository, err := InitAnalyzeHistoryRepository(elasticsearchClient, configConfig)
	if err != nil {
		return nil, err
	}
	similarityService, err := InitAnalyzeService(configConfig, shingleindexService, fulltextindexService, service, semanticindexService, analyzehistoryRepository)
	if err != nil {
		return nil, err
	}
	similarityHandler := InitAnalyzeNatsHandler(similarityService)
	serverServer := server2.NewServer(conn, handler, attributeHandler, similarityHandler)
	return serverServer, nil
}

func InitFileSavedHandler(service *document2.Service) (*filesaved.Handler, error) {
	handler := filesaved.NewHandler(service)
	return handler, nil
}

func InitDocumentPipeline(contextContext context.Context, configConfig *config.Config) (*documentpipeline.Service, error) {
	conn, err := ProvideNats(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	jetStream, err := InitNatsJetstream(conn)
	if err != nil {
		return nil, err
	}
	stream, err := ProvideDocumentStatusJetstreamStream(contextContext, jetStream)
	if err != nil {
		return nil, err
	}
	repository, err := InitDocumentStatusRepository(contextContext, jetStream)
	if err != nil {
		return nil, err
	}
	service, err := InitDocumentStatusService(repository)
	if err != nil {
		return nil, err
	}
	client, err := InitTika(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	minioClient, err := ProvideS3(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	filestorageRepository, err := InitFilestorageRepository(minioClient, configConfig)
	if err != nil {
		return nil, err
	}
	elasticsearchClient, err := ProvideElastic(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	documentRepository, err := InitDocumentRepository(elasticsearchClient, configConfig)
	if err != nil {
		return nil, err
	}
	documentService, err := InitDocumentService(configConfig, client, filestorageRepository, documentRepository)
	if err != nil {
		return nil, err
	}
	handler, err := InitFileSavedHandler(documentService)
	if err != nil {
		return nil, err
	}
	v := ProvideDocumentPipelineStages(handler)
	documentpipelineService, err := documentpipeline.NewService(contextContext, stream, service, v)
	if err != nil {
		return nil, err
	}
	return documentpipelineService, nil
}

// wire.go:

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
	_ context.Context,
	cfg *config.Config,
) (*nats.Conn, error) {
	var err error

	natsCliOnce.Do(func() {
		natsCli, err = nats.Connect(cfg.Nats)
	})

	return natsCli, err
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

func ProvideDocumentStatusJetstreamKV(
	ctx context.Context,
	js jetstream.JetStream,
) (jetstream.KeyValue, error) {
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
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
	s, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      "documentstatus",
		Subjects:  []string{"documentstatus.>"},
		Retention: jetstream.WorkQueuePolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("new steream: %w", err)
	}

	return s, nil
}

func ProvideAnalyzeServiceOpts() similarity.Opts {
	return similarity.Opts{}
}

func ProvideDocumentServiceOpts() document2.Opts {
	return document2.Opts{}
}

func ProvideDocumentPipelineStages(
	fsh *filesaved.Handler,
) []documentpipeline.Stage {
	return []documentpipeline.Stage{
		{
			Trigger: model.DocumentProcessingStatusFileSaved,
			Action:  fsh,
			Next:    model.DocumentProcessingStatusDocumentSaved,
		},
	}
}
