//go:build wireinject

package provider

import (
	"context"
	"io"
	"os"
	"simrep/api/rest/server"
	filesapi "simrep/api/rest/server/handler/files"
	"simrep/api/rest/server/handler/similarity"
	"simrep/internal/config"
	documentrepo "simrep/internal/repository/document"
	documentfilerepo "simrep/internal/repository/documentfile"
	imagerepo "simrep/internal/repository/image"
	documentsrv "simrep/internal/service/document"
	documentparsersrv "simrep/internal/service/documentparser"
	"simrep/pkg/elasticutil"
	"simrep/pkg/minioutil"

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

func InitDocumentParserService() (*documentparsersrv.Service, error) {
	panic(wire.Build(
		documentparsersrv.NewService,
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

func InitDocumentService(
	_ *documentparsersrv.Service,
	_ *imagerepo.Repository,
	_ *documentfilerepo.Repository,
	_ *documentrepo.Repository,
) (*documentsrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(documentsrv.FileParser), new(*documentparsersrv.Service)),
		wire.Bind(new(documentsrv.FileRepository), new(*documentfilerepo.Repository)),
		wire.Bind(new(documentsrv.ImageRepository), new(*imagerepo.Repository)),
		wire.Bind(new(documentsrv.Repository), new(*documentrepo.Repository)),
		documentsrv.NewService,
	))
}

func InitFilesHandler(
	_ *documentsrv.Service,
) *filesapi.Handler {
	panic(wire.Build(
		wire.Bind(new(filesapi.DocumentService), new(*documentsrv.Service)),
		filesapi.NewHandler,
	))
}

func InitSimilarityHandler() *similarity.Handler {
	panic(wire.Build(similarity.NewHandler))
}

func InitRestAPI(
	_ context.Context,
	_ *config.Config,
) (*server.Server, error) {
	panic(wire.Build(
		wire.Value("localhost:8000"),
		ProvideSpec,
		InitS3,
		InitElastic,
		InitImageRepository,
		InitDocumentFileRepository,
		InitDocumentRepository,
		InitDocumentParserService,
		InitDocumentService,
		InitFilesHandler,
		InitSimilarityHandler,
		wire.Bind(new(server.FileHandler), new(*filesapi.Handler)),
		wire.Bind(new(server.SimilarityHandler), new(*similarity.Handler)),
		wire.Struct(new(server.Opts), "*"),
		server.New,
	))
}
