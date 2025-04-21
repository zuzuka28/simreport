package cmd

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zuzuka28/simreport/prj/fulltextindex/internal/config"
)

type Config = config.Config

type appOpts struct {
	reg prometheus.Registerer
}

type AppOpt func(*appOpts)

func WithPrometheusRegistrer(
	reg prometheus.Registerer,
) AppOpt {
	return func(ao *appOpts) {
		ao.reg = reg
	}
}
