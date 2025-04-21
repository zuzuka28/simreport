package cmd

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/similarity/internal/provider"
)

func RunNATSAPI(ctx context.Context, cfg *Config, opts ...AppOpt) error {
	appopts := new(appOpts)

	for _, opt := range opts {
		opt(appopts)
	}

	m := provider.ProvideMetrics()

	if appopts.reg != nil {
		appopts.reg.MustRegister(m.Collectors()...)
	}

	svc, err := provider.InitNatsAPI(ctx, cfg, m)
	if err != nil {
		return fmt.Errorf("provide api: %w", err)
	}

	return runService(ctx, svc, func() {})
}
