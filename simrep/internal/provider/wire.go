//go:build wireinject

package provider

import (
	"context"
	"io"
	"os"
	"simrep/api/rest/server"
	documentapi "simrep/api/rest/server/handler/document"
	"simrep/internal/config"
	documentrepo "simrep/internal/repository/document"
	documentfilerepo "simrep/internal/repository/documentfile"
	imagerepo "simrep/internal/repository/image"
	documentsrv "simrep/internal/service/document"
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
	_ *imagerepo.Repository,
	_ *documentfilerepo.Repository,
	_ *documentrepo.Repository,
) (*documentsrv.Service, error) {
	panic(wire.Build(
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

func InitRestAPI(
	_ context.Context,
	_ *config.Config,
) (*server.Server, error) {
	panic(wire.Build(
		ProvideSpec,
		InitS3,
		InitElastic,
		InitImageRepository,
		InitDocumentFileRepository,
		InitDocumentRepository,
		InitDocumentService,
		InitDocumentHandler,
		wire.Bind(new(server.DocumentHandler), new(*documentapi.Handler)),
		wire.FieldsOf(new(*config.Config), "Port"),
		wire.Struct(new(server.Opts), "*"),
		server.New,
	))
}
