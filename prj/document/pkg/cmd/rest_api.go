package cmd

import (
	"context"
	"fmt"

	"github.com/zuzuka28/simreport/prj/document/internal/provider"
)

func RunRESTAPI(ctx context.Context, cfg *Config, opts ...AppOpt) error {
	appopts := new(appOpts)

	for _, opt := range opts {
		opt(appopts)
	}

	m := provider.ProvideMetrics()

	if appopts.reg != nil {
		appopts.reg.MustRegister(m.Collectors()...)
	}

	svc, err := provider.InitRestAPI(ctx, cfg)
	if err != nil {
		return fmt.Errorf("provide api: %w", err)
	}

	return runService(ctx, svc, func() {})
}
