package cmd

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/semanticindex/internal/provider"
)

func RunAsyncProcessing(ctx context.Context, cfg *Config, opts ...AppOpt) error {
	appopts := new(appOpts)

	for _, opt := range opts {
		opt(appopts)
	}

	m := provider.ProvideMetrics()

	if appopts.reg != nil {
		appopts.reg.MustRegister(m.Collectors()...)
	}

	svc, err := provider.InitNatsEventAPI(ctx, cfg, m)
	if err != nil {
		return fmt.Errorf("provide document pipeline: %w", err)
	}

	return runService(ctx, svc, func() {})
}
