//go:build wireinject

package provider

import (
	"context"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/zuzuka28/simreport/lib/elasticutil"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/bot"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/config"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/metrics"
	userstaterepo "github.com/zuzuka28/simreport/prj/tgbot/internal/repository/userstate"
	userstatesrv "github.com/zuzuka28/simreport/prj/tgbot/internal/service/userstate"

	documentrepo "github.com/zuzuka28/simreport/prj/tgbot/internal/repository/document"
	similarityrepo "github.com/zuzuka28/simreport/prj/tgbot/internal/repository/similarity"
	documentsrv "github.com/zuzuka28/simreport/prj/tgbot/internal/service/document"
	similaritysrv "github.com/zuzuka28/simreport/prj/tgbot/internal/service/similarity"

	"github.com/google/wire"
	"github.com/nats-io/nats.go"
)

//nolint:gochecknoglobals
var (
	metricsS    *metrics.Metrics
	metricsOnce sync.Once
)

func ProvideMetrics() *metrics.Metrics {
	metricsOnce.Do(func() {
		metricsS = metrics.New()
	})

	return metricsS
}

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
	*metrics.Metrics,
) *userstaterepo.Repository {
	panic(wire.Build(
		wire.FieldsOf(new(*config.Config), "UserStateRepo"),
		wire.Bind(new(userstaterepo.Metrics), new(*metrics.Metrics)),
		userstaterepo.NewRepository,
	))
}

func ProvideDocumentRepository(
	*nats.Conn,
	*metrics.Metrics,
) *documentrepo.Repository {
	panic(wire.Build(
		wire.Bind(new(documentrepo.Metrics), new(*metrics.Metrics)),
		documentrepo.NewRepository,
	))
}

func ProvideSimilarityRepository(
	*nats.Conn,
	*metrics.Metrics,
) *similarityrepo.Repository {
	panic(wire.Build(
		wire.Bind(new(similarityrepo.Metrics), new(*metrics.Metrics)),
		similarityrepo.NewRepository,
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

func ProvideSimilarityService(
	*similarityrepo.Repository,
) *similaritysrv.Service {
	panic(wire.Build(
		wire.Bind(new(similaritysrv.Repository), new(*similarityrepo.Repository)),
		similaritysrv.NewService,
	))
}

func InitBot(
	context.Context,
	*config.Config,
) (*bot.Bot, error) {
	panic(wire.Build(
		ProvideMetrics,
		ProvideElastic,
		ProvideNats,
		ProvideUserStateRepository,
		ProvideUserStateService,
		ProvideDocumentRepository,
		ProvideDocumentService,
		ProvideSimilarityRepository,
		ProvideSimilarityService,
		wire.Bind(new(bot.UserStateService), new(*userstatesrv.Service)),
		wire.Bind(new(bot.DocumentService), new(*documentsrv.Service)),
		wire.Bind(new(bot.SimilarityService), new(*similaritysrv.Service)),
		wire.Bind(new(bot.Metrics), new(*metrics.Metrics)),
		wire.FieldsOf(new(*config.Config), "Bot"),
		bot.New,
	))
}
