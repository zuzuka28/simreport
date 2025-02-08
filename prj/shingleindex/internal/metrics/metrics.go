package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "simrep"
	subsystem = "shingleindex"
)

type Metrics struct {
	documentRepositoryRequests        *prometheus.CounterVec
	documentRepositoryRequestDuration *prometheus.HistogramVec
	natsMicroRequests                 *prometheus.CounterVec
	natsMicroRequestDurations         *prometheus.HistogramVec
	natsMicroSizes                    *prometheus.HistogramVec
}

func New() *Metrics {
	return &Metrics{
		documentRepositoryRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "document_repository_requests_total",
				Help:      "Tracks requests to document repository",
			},
			[]string{
				"op",
				"status",
			},
		),
		documentRepositoryRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "document_repository_duration_seconds",
				Help:      "Tracks request durations in document repository",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"op",
				"status",
			},
		),
		natsMicroRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "nats_micro_requests_total",
				Help:      "Tracks requests to nats micro",
			},
			[]string{
				"op",
				"status",
			},
		),
		natsMicroRequestDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "nats_micro_request_duration_seconds",
				Help:      "Tracks request durations in nats micro",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"op",
				"status",
			},
		),
		natsMicroSizes: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "nats_micro_request_sizes_bytes",
				Help:      "Tracks request sizes in nats micro",
				Buckets:   prometheus.ExponentialBuckets(1024*256, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"op",
				"status",
			},
		),
	}
}

func (m *Metrics) IncDocumentRepositoryRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.documentRepositoryRequests.With(labels).Inc()
	m.documentRepositoryRequestDuration.With(labels).Observe(dur)
}

func (m *Metrics) IncNatsMicroRequest(op string, status string, size int, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.natsMicroRequests.With(labels).Inc()
	m.natsMicroRequestDurations.With(labels).Observe(dur)
	m.natsMicroSizes.With(labels).Observe(float64(size))
}

func (m *Metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.documentRepositoryRequests,
		m.documentRepositoryRequestDuration,
		m.natsMicroRequests,
		m.natsMicroRequestDurations,
		m.natsMicroSizes,
	}
}
