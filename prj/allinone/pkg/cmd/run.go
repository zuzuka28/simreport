package cmd

import (
	"context"

	"github.com/zuzuka28/simreport/prj/allinone/internal/service/allinone"
)

func RunApp(ctx context.Context, cfg *Config, opts ...AppOpt) error {
	appopts := new(appOpts)

	for _, opt := range opts {
		opt(appopts)
	}

	svc := allinone.NewService(cfg, appopts.reg)

	return runService(ctx, svc, func() {})
}
