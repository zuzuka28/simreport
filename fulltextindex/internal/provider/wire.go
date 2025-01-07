//go:build wireinject

package provider

import (
	"context"
	serverevent "fulltextindex/api/nats/event"
	indexerevent "fulltextindex/api/nats/event/handler/indexer"
	fulltextindexmicro "fulltextindex/api/nats/micro/handler/fulltextindex"
	servermicro "fulltextindex/api/nats/micro/server"
	"fulltextindex/internal/config"
	documentrepo "fulltextindex/internal/repository/document"
	fulltextindexrepo "fulltextindex/internal/repository/fulltextindex"
	documentsrv "fulltextindex/internal/service/document"
	fulltextindexsrv "fulltextindex/internal/service/fulltextindex"
	"fulltextindex/pkg/elasticutil"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
)

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

func InitFulltextIndexRepository(
	_ *elasticsearch.Client,
	_ *config.Config,
) (*fulltextindexrepo.Repository, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "FulltextRepo"),
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

func InitFulltextHandler(
	_ *fulltextindexsrv.Service,
	_ *documentsrv.Service,
) (*fulltextindexmicro.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(fulltextindexmicro.Service), new(*fulltextindexsrv.Service)),
		wire.Bind(new(fulltextindexmicro.DocumentService), new(*documentsrv.Service)),
		fulltextindexmicro.NewHandler,
	))
}

func InitIndexerHandler(
	_ *fulltextindexsrv.Service,
	_ *documentsrv.Service,
) (*indexerevent.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(indexerevent.Service), new(*fulltextindexsrv.Service)),
		wire.Bind(new(indexerevent.DocumentService), new(*documentsrv.Service)),
		indexerevent.NewHandler,
	))
}

func InitNatsMicroAPI(
	_ context.Context,
	_ *config.Config,
) (*servermicro.Server, error) {
	panic(wire.Build(
		InitElastic,
		InitNats,

		InitDocumentRepository,
		InitDocumentService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitFulltextHandler,

		wire.Bind(new(servermicro.FileindexHandler), new(*fulltextindexmicro.Handler)),
		servermicro.NewServer,
	))
}

func InitNatsEventAPI(
	_ context.Context,
	_ *config.Config,
) (*serverevent.Server, error) {
	panic(wire.Build(
		InitElastic,
		InitNats,

		InitDocumentRepository,
		InitDocumentService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitIndexerHandler,

		wire.Bind(new(serverevent.IndexerHandler), new(*indexerevent.Handler)),
		serverevent.NewServer,
	))
}
