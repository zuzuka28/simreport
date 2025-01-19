//go:build wireinject

package provider

import (
	serverhttp "anysave/api/rest/server"
	anysaveapi "anysave/api/rest/server/handler/anysave"
	"anysave/internal/config"
	anysaverepo "anysave/internal/repository/anysave"
	anysavesrv "anysave/internal/service/anysave"
	"context"
	"io"
	"os"
	"sync"

	"github.com/zuzuka28/simreport/lib/minioutil"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
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

func InitNatsJetstream(
	_ *nats.Conn,
) (jetstream.JetStream, error) {
	panic(wire.Build(
		wire.Value([]jetstream.JetStreamOpt(nil)),
		jetstream.New,
	))
}

func InitDocumentFileRepository(
	_ *minio.Client,
	_ *config.Config,
) (*anysaverepo.Repository, error) {
	panic(wire.Build(
		anysaverepo.NewRepository,
	))
}

func ProvideAnysaveServiceOpts() anysavesrv.Opts {
	return anysavesrv.Opts{} //nolint:exhaustruct
}

func InitAnysaveService(
	_ *minio.Client,
	_ *config.Config,
) (*anysavesrv.Service, error) {
	panic(wire.Build(
		InitDocumentFileRepository,
		ProvideAnysaveServiceOpts,
		wire.Bind(new(anysavesrv.Repository), new(*anysaverepo.Repository)),
		anysavesrv.NewService,
	))
}

func InitAnysaveHandler(
	_ *anysavesrv.Service,
) *anysaveapi.Handler {
	panic(wire.Build(
		wire.Bind(new(anysaveapi.Service), new(*anysavesrv.Service)),
		anysaveapi.NewHandler,
	))
}

func InitRestAPI(
	_ context.Context,
	_ *config.Config,
) (*serverhttp.Server, error) {
	panic(wire.Build(
		ProvideSpec,
		ProvideS3,

		InitAnysaveService,

		InitAnysaveHandler,
		wire.Bind(new(serverhttp.FileHandler), new(*anysaveapi.Handler)),
		wire.FieldsOf(new(*config.Config), "Port"),
		wire.Struct(new(serverhttp.Opts), "*"),
		serverhttp.New,
	))
}
