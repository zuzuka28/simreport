// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package provider

import (
	"context"
	document4 "document/api/nats/handler/document"
	server2 "document/api/nats/server"
	"document/api/rest/server"
	analyze2 "document/api/rest/server/handler/analyze"
	document3 "document/api/rest/server/handler/document"
	"document/internal/config"
	"document/internal/model"
	"document/internal/repository/analyzehistory"
	document2 "document/internal/repository/document"
	"document/internal/repository/documentstatus"
	"document/internal/repository/filestorage"
	"document/internal/repository/fulltextindexclient"
	fulltextindexclient2 "document/internal/repository/semanticindexclient"
	"document/internal/repository/shingleindexclient"
	"document/internal/service/analyze"
	"document/internal/service/document"
	"document/internal/service/documentparser"
	"document/internal/service/documentpipeline"
	"document/internal/service/documentpipeline/handler/filesaved"
	documentstatus2 "document/internal/service/documentstatus"
	"document/internal/service/fulltextindex"
	fulltextindex2 "document/internal/service/semanticindex"
	"document/internal/service/shingleindex"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"io"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/lib/minioutil"
	"github.com/zuzuka28/simreport/lib/tikaclient"
	"net/http"
	"os"
	"sync"
)

// Injectors from wire.go:

func InitConfig(path string) (*config.Config, error) {
	configConfig, err := config.New(path)
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

func InitShingleIndexService(repository *shingleindexclient.Repository, service *document.Service) (*shingleindex.Service, error) {
	shingleindexService := shingleindex.NewService(repository, service)
	return shingleindexService, nil
}

func InitFulltextIndexRepository(conn *nats.Conn) (*fulltextindexclient.Repository, error) {
	repository := fulltextindexclient.NewRepository(conn)
	return repository, nil
}

func InitFulltextIndexService(repository *fulltextindexclient.Repository) (*fulltextindex.Service, error) {
	service := fulltextindex.NewService(repository)
	return service, nil
}

func InitSemanticIndexRepository(conn *nats.Conn) (*fulltextindexclient2.Repository, error) {
	repository := fulltextindexclient2.NewRepository(conn)
	return repository, nil
}

func InitSemanticIndexService(repository *fulltextindexclient2.Repository) (*fulltextindex2.Service, error) {
	service := fulltextindex2.NewService(repository)
	return service, nil
}

func InitDocumentRepository(client *elasticsearch.Client, configConfig *config.Config) (*document2.Repository, error) {
	opts := configConfig.DocumentRepo
	repository, err := document2.NewRepository(opts, client)
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

func InitDocumentStatusRepository(ctx context.Context, js jetstream.JetStream) (*documentstatus.Repository, error) {
	keyValue, err := ProvideDocumentStatusJetstreamKV(ctx, js)
	if err != nil {
		return nil, err
	}
	repository := documentstatus.NewRepository(keyValue, js)
	return repository, nil
}

func InitDocumentStatusService(repository *documentstatus.Repository) (*documentstatus2.Service, error) {
	service := documentstatus2.NewService(repository)
	return service, nil
}

func InitAnalyzeService(configConfig *config.Config, service *shingleindex.Service, fulltextindexService *fulltextindex.Service, documentService *document.Service, service2 *fulltextindex2.Service, repository *analyzehistory.Repository) (*analyze.Service, error) {
	opts := ProvideAnalyzeServiceOpts()
	analyzeService := analyze.NewService(opts, documentService, service, fulltextindexService, service2, repository)
	return analyzeService, nil
}

func InitDocumentParserService(client *tikaclient.Client) (*documentparser.Service, error) {
	service := documentparser.NewService(client)
	return service, nil
}

func InitDocumentService(configConfig *config.Config, client *tikaclient.Client, repository *filestorage.Repository, documentRepository *document2.Repository) (*document.Service, error) {
	opts := ProvideDocumentServiceOpts()
	service, err := InitDocumentParserService(client)
	if err != nil {
		return nil, err
	}
	documentService := document.NewService(opts, documentRepository, repository, service)
	return documentService, nil
}

func InitDocumentHandler(service *document.Service) *document3.Handler {
	handler := document3.NewHandler(service)
	return handler
}

func InitAnalyzeHandler(service *analyze.Service) *analyze2.Handler {
	handler := analyze2.NewHandler(service)
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
	handler := InitDocumentHandler(service)
	conn, err := ProvideNats(contextContext, configConfig)
	if err != nil {
		return nil, err
	}
	shingleindexclientRepository, err := InitShingleIndexRepository(conn)
	if err != nil {
		return nil, err
	}
	shingleindexService, err := InitShingleIndexService(shingleindexclientRepository, service)
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
	repository2, err := InitSemanticIndexRepository(conn)
	if err != nil {
		return nil, err
	}
	service2, err := InitSemanticIndexService(repository2)
	if err != nil {
		return nil, err
	}
	analyzehistoryRepository, err := InitAnalyzeHistoryRepository(elasticsearchClient, configConfig)
	if err != nil {
		return nil, err
	}
	analyzeService, err := InitAnalyzeService(configConfig, shingleindexService, fulltextindexService, service, service2, analyzehistoryRepository)
	if err != nil {
		return nil, err
	}
	analyzeHandler := InitAnalyzeHandler(analyzeService)
	opts := server.Opts{
		Port:            int2,
		Spec:            v,
		DocumentHandler: handler,
		AnalyzeHandler:  analyzeHandler,
	}
	serverServer, err := server.New(opts)
	if err != nil {
		return nil, err
	}
	return serverServer, nil
}

func InitFileSavedHandler(service *document.Service, repository *filestorage.Repository) (*filesaved.Handler, error) {
	handler := filesaved.NewHandler(repository, service)
	return handler, nil
}

func InitDocumentNatsHandler(service *document.Service) *document4.Handler {
	handler := document4.NewHandler(service)
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
	handler := InitDocumentNatsHandler(service)
	serverServer := server2.NewServer(conn, handler)
	return serverServer, nil
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
	handler, err := InitFileSavedHandler(documentService, filestorageRepository)
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
	ctx context.Context,
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
		Retention: jetstream.InterestPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("new steream: %w", err)
	}

	return s, nil
}

func ProvideAnalyzeServiceOpts() analyze.Opts {
	return analyze.Opts{}
}

func ProvideDocumentServiceOpts() document.Opts {
	return document.Opts{}
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
