package allinone

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/allinone/internal/config"
	fulltextindexcmd "github.com/zuzuka28/simreport/prj/fulltextindex/pkg/cmd"
	semanticindexcmd "github.com/zuzuka28/simreport/prj/semanticindex/pkg/cmd"
	shingleindexcmd "github.com/zuzuka28/simreport/prj/shingleindex/pkg/cmd"
	similaritycmd "github.com/zuzuka28/simreport/prj/similarity/pkg/cmd"

	"github.com/prometheus/client_golang/prometheus"
	documentcmd "github.com/zuzuka28/simreport/prj/document/pkg/cmd"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	cfg config.Config
	reg prometheus.Registerer
}

func NewService(
	cfg *config.Config,
	reg prometheus.Registerer,
) *Service {
	if reg == nil {
		reg = prometheus.NewRegistry()
	}

	return &Service{
		cfg: *cfg,
		reg: reg,
	}
}

func (s *Service) Start(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "rest"},
			s.reg,
		)

		return documentcmd.RunRESTAPI(
			egCtx,
			&s.cfg.DocumentService,
			documentcmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "nats"},
			s.reg,
		)

		return documentcmd.RunAsyncProcessing(
			egCtx,
			&s.cfg.DocumentService,
			documentcmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "processing"},
			s.reg,
		)

		return documentcmd.RunAsyncProcessing(
			egCtx,
			&s.cfg.DocumentService,
			documentcmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "rest"},
			s.reg,
		)

		return similaritycmd.RunRESTAPI(
			egCtx,
			&s.cfg.SimilarityService,
			similaritycmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "nats"},
			s.reg,
		)

		return similaritycmd.RunNATSAPI(
			egCtx,
			&s.cfg.SimilarityService,
			similaritycmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "nats"},
			s.reg,
		)

		return shingleindexcmd.RunNATSAPI(
			egCtx,
			&s.cfg.ShingleIndexService,
			shingleindexcmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "processing"},
			s.reg,
		)

		return shingleindexcmd.RunAsyncProcessing(
			egCtx,
			&s.cfg.ShingleIndexService,
			shingleindexcmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "nats"},
			s.reg,
		)

		return fulltextindexcmd.RunNATSAPI(
			egCtx,
			&s.cfg.FulltextIndexService,
			fulltextindexcmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "processing"},
			s.reg,
		)

		return fulltextindexcmd.RunAsyncProcessing(
			egCtx,
			&s.cfg.FulltextIndexService,
			fulltextindexcmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "nats"},
			s.reg,
		)

		return semanticindexcmd.RunNATSAPI(
			egCtx,
			&s.cfg.SemanticIndexService,
			semanticindexcmd.WithPrometheusRegistrer(reg),
		)
	})

	eg.Go(func() error {
		reg := prometheus.WrapRegistererWith(
			prometheus.Labels{"app_part": "processing"},
			s.reg,
		)

		return semanticindexcmd.RunAsyncProcessing(
			egCtx,
			&s.cfg.SemanticIndexService,
			semanticindexcmd.WithPrometheusRegistrer(reg),
		)
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}
