//go:build wireinject

package provider

import (
	"context"
	"sync"

	"github.com/zuzuka28/simreport/lib/elasticutil"
	serverevent "github.com/zuzuka28/simreport/prj/fulltextindex/api/nats/event"
	indexerevent "github.com/zuzuka28/simreport/prj/fulltextindex/api/nats/event/handler/indexer"
	fulltextindexmicro "github.com/zuzuka28/simreport/prj/fulltextindex/api/nats/micro/handler/fulltextindex"
	servermicro "github.com/zuzuka28/simreport/prj/fulltextindex/api/nats/micro/server"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/config"
	documentrepo "github.com/zuzuka28/simreport/prj/fulltextindex/internal/repository/document"
	fulltextindexrepo "github.com/zuzuka28/simreport/prj/fulltextindex/internal/repository/fulltextindex"
	documentsrv "github.com/zuzuka28/simreport/prj/fulltextindex/internal/service/document"
	fulltextindexsrv "github.com/zuzuka28/simreport/prj/fulltextindex/internal/service/fulltextindex"

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

	return natsCli, err //nolint:wrapcheck
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
		ProvideNats,

		InitDocumentRepository,
		InitDocumentService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitFulltextHandler,

		wire.Bind(new(servermicro.Handler), new(*fulltextindexmicro.Handler)),
		servermicro.NewServer,
	))
}

func InitNatsEventAPI(
	_ context.Context,
	_ *config.Config,
) (*serverevent.Server, error) {
	panic(wire.Build(
		InitElastic,
		ProvideNats,

		InitDocumentRepository,
		InitDocumentService,

		InitFulltextIndexRepository,
		InitFulltextIndexService,

		InitIndexerHandler,

		wire.Bind(new(serverevent.IndexerHandler), new(*indexerevent.Handler)),
		serverevent.NewServer,
	))
}
