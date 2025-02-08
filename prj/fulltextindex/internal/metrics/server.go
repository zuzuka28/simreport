package metrics

import (
	"fmt"
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

func (server *Server) Start() error {
	mux := http.DefaultServeMux

	http.Handle("/metrics", promhttp.HandlerFor(server.reg, promhttp.HandlerOpts{ //nolint:exhaustruct
		Registry:          server.reg,
		EnableOpenMetrics: true,
	}))

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
