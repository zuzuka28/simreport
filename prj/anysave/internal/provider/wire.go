//go:build wireinject

package provider

import (
	serverhttp "anysave/api/rest/server"
	anysaveapi "anysave/api/rest/server/handler/anysave"
	"anysave/internal/config"
	anysaverepo "anysave/internal/repository/anysave"
	documentstatusrepo "anysave/internal/repository/documentstatus"
	anysavesrv "anysave/internal/service/anysave"
	documentstatussrv "anysave/internal/service/documentstatus"
	"github.com/zuzuka28/simreport/lib/minioutil"
	"context"
	"fmt"
	"io"
	"os"
	"sync"

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

func ProvideDocumentStatusJetstreamKV(
	ctx context.Context,
	js jetstream.JetStream,
) (jetstream.KeyValue, error) {
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{ //nolint:exhaustruct
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
	s, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{ //nolint:exhaustruct
		Name:      "documentstatus",
		Subjects:  []string{"documentstatus.>"},
		Retention: jetstream.InterestPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("new steream: %w", err)
	}

	return s, nil
}

func InitDocumentFileRepository(
	_ *minio.Client,
	_ *config.Config,
) (*anysaverepo.Repository, error) {
	panic(wire.Build(
		anysaverepo.NewRepository,
	))
}

func InitDocumentStatusRepository(
	ctx context.Context,
	js jetstream.JetStream,
) (*documentstatusrepo.Repository, error) {
	panic(wire.Build(
		ProvideDocumentStatusJetstreamKV,
		wire.Bind(new(jetstream.Publisher), new(jetstream.JetStream)),
		documentstatusrepo.NewRepository,
	))
}

func InitDocumentStatusService(
	_ *documentstatusrepo.Repository,
) (*documentstatussrv.Service, error) {
	panic(wire.Build(
		wire.Bind(new(documentstatussrv.Repository), new(*documentstatusrepo.Repository)),
		documentstatussrv.NewService,
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
	_ *documentstatussrv.Service,
	_ *anysavesrv.Service,
) *anysaveapi.Handler {
	panic(wire.Build(
		wire.Bind(new(anysaveapi.Service), new(*anysavesrv.Service)),
		wire.Bind(new(anysaveapi.StatusService), new(*documentstatussrv.Service)),
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
		ProvideNats,
		InitNatsJetstream,

		InitDocumentStatusRepository,
		InitDocumentStatusService,

		InitAnysaveService,

		InitAnysaveHandler,
		wire.Bind(new(serverhttp.FileHandler), new(*anysaveapi.Handler)),
		wire.FieldsOf(new(*config.Config), "Port"),
		wire.Struct(new(serverhttp.Opts), "*"),
		serverhttp.New,
	))
}
