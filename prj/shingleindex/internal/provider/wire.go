//go:build wireinject

package provider

import (
	"context"
	serverevent "shingleindex/api/nats/event"
	indexerevent "shingleindex/api/nats/event/handler/indexer"
	shingleindexmicro "shingleindex/api/nats/micro/handler/shingleindex"
	servermicro "shingleindex/api/nats/micro/server"
	"shingleindex/internal/config"
	documentrepo "shingleindex/internal/repository/document"
	shingleindexrepo "shingleindex/internal/repository/shingleindex"
	documentsrv "shingleindex/internal/service/document"
	shingleindexsrv "shingleindex/internal/service/shingleindex"

	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

func InitConfig(path string) (*config.Config, error) {
	panic(wire.Build(config.New))
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

func ProvideRedis(
	cfg *config.Config,
) (*redis.Client, error) {
	u, err := redis.ParseURL(cfg.Redis.DSN)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(u), nil
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

func InitShingleIndexRepository(
	_ *redis.Client,
	_ *config.Config,
) (*shingleindexrepo.Repository, error) {
	panic(wire.Build(
		wire.Value(shingleindexrepo.Opts{}),
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

func InitShingleHandler(
	_ *shingleindexsrv.Service,
	_ *documentsrv.Service,
) (*shingleindexmicro.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(shingleindexmicro.Service), new(*shingleindexsrv.Service)),
		wire.Bind(new(shingleindexmicro.DocumentService), new(*documentsrv.Service)),
		shingleindexmicro.NewHandler,
	))
}

func InitIndexerHandler(
	_ *shingleindexsrv.Service,
	_ *documentsrv.Service,
) (*indexerevent.Handler, error) {
	panic(wire.Build(
		wire.Bind(new(indexerevent.Service), new(*shingleindexsrv.Service)),
		wire.Bind(new(indexerevent.DocumentService), new(*documentsrv.Service)),
		indexerevent.NewHandler,
	))
}

func InitNatsMicroAPI(
	_ context.Context,
	_ *config.Config,
) (*servermicro.Server, error) {
	panic(wire.Build(
		ProvideRedis,
		InitNats,

		InitDocumentRepository,
		InitDocumentService,

		InitShingleIndexRepository,
		InitShingleIndexService,

		InitShingleHandler,

		wire.Bind(new(servermicro.FileindexHandler), new(*shingleindexmicro.Handler)),
		servermicro.NewServer,
	))
}

func InitNatsEventAPI(
	_ context.Context,
	_ *config.Config,
) (*serverevent.Server, error) {
	panic(wire.Build(
		ProvideRedis,
		InitNats,

		InitDocumentRepository,
		InitDocumentService,

		InitShingleIndexRepository,
		InitShingleIndexService,

		InitIndexerHandler,

		wire.Bind(new(serverevent.IndexerHandler), new(*indexerevent.Handler)),
		serverevent.NewServer,
	))
}
