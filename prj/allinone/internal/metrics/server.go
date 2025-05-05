package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	writeTimeout      = 5 * time.Second
	readHeaderTimeout = 5 * time.Second
)

type Server struct {
	port int
	reg  *prometheus.Registry
}

func NewMetricsServer(
	httpPort int,
	reg *prometheus.Registry,
) *Server {
	return &Server{
		port: httpPort,
		reg:  reg,
	}
}

func (server *Server) Start(context.Context) error {
	mux := http.DefaultServeMux

	http.Handle("/metrics", promhttp.HandlerFor(server.reg, promhttp.HandlerOpts{ //nolint:exhaustruct
		Registry:          server.reg,
		EnableOpenMetrics: true,
	}))

	slog.Info("metrics server started", "port", server.port)

	srv := http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf(":%d", server.port),
		WriteTimeout:      writeTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("serve metrics: %w", err)
	}

	return nil
}
