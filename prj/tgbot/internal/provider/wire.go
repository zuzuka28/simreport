//go:build wireinject
package provider

import (
	"context"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/bot"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/config"
	userstaterepo "github.com/zuzuka28/simreport/prj/tgbot/internal/repository/userstate"
	userstatesrv "github.com/zuzuka28/simreport/prj/tgbot/internal/service/userstate"

	documentrepo "github.com/zuzuka28/simreport/prj/tgbot/internal/repository/document"
	documentsrv "github.com/zuzuka28/simreport/prj/tgbot/internal/service/document"

	"github.com/google/wire"
	"github.com/nats-io/nats.go"
)

func InitConfig(path string) (*config.Config, error) {
	panic(wire.Build(config.New))
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

	return elasticCli, err //nolint:wrapcheck
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

func ProvideUserStateRepository(
	*config.Config,
	*elasticsearch.Client,
) *userstaterepo.Repository {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "UserStateRepo"),
		userstaterepo.NewRepository,
	))
}

func ProvideDocumentRepository(
	*nats.Conn,
) *documentrepo.Repository {
	panic(wire.Build(
		documentrepo.NewRepository,
	))
}

func ProvideUserStateService(
	*userstaterepo.Repository,
) *userstatesrv.Service {
	panic(wire.Build(
		wire.Bind(new(userstatesrv.Repository), new(*userstaterepo.Repository)),
		userstatesrv.NewService,
	))
}

func ProvideDocumentService(
	*documentrepo.Repository,
) *documentsrv.Service {
	panic(wire.Build(
		wire.Bind(new(documentsrv.Repository), new(*documentrepo.Repository)),
		documentsrv.NewService,
	))
}

func InitBot(
	context.Context,
	*config.Config,
) (*bot.Bot, error) {
	panic(wire.Build(
		ProvideElastic,
		ProvideNats,
		ProvideUserStateRepository,
		ProvideUserStateService,
		ProvideDocumentRepository,
		ProvideDocumentService,
		wire.Bind(new(bot.UserStateService), new(*userstatesrv.Service)),
		wire.Bind(new(bot.DocumentService), new(*documentsrv.Service)),
		wire.FieldsOf(new(*config.Config), "Bot"),
		bot.New,
	))
}
