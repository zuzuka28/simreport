package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "simrep"
	subsystem = "semanticindex"
)

type Metrics struct {
	filestorageRequests               *prometheus.CounterVec
	filestorageRequestDurations       *prometheus.HistogramVec
	semanticIndexRequests             *prometheus.CounterVec
	semanticIndexRequestDuration      *prometheus.HistogramVec
	vectorizerRequests                *prometheus.CounterVec
	vectorizerRequestDuration         *prometheus.HistogramVec
	documentRepositoryRequests        *prometheus.CounterVec
	documentRepositoryRequestDuration *prometheus.HistogramVec
	natsMicroRequests                 *prometheus.CounterVec
	natsMicroRequestDurations         *prometheus.HistogramVec
	natsMicroSizes                    *prometheus.HistogramVec
}

func New() *Metrics {
	return &Metrics{
		filestorageRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "filestorage_requests_total",
				Help:      "Tracks requests to filestorage",
			},
			[]string{
				"op",
				"status",
			},
		),
		filestorageRequestDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "filestorage_request_duration_seconds",
				Help:      "Tracks request durations filestorage",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"op",
				"status",
			},
		),
		semanticIndexRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "semanticindex_requests_total",
				Help:      "Tracks requests to semanticindex",
			},
			[]string{
				"op",
				"status",
			},
		),
		semanticIndexRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "semanticindex_duration_seconds",
				Help:      "Tracks request durations in semanticindex",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"op",
				"status",
			},
		),
		vectorizerRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "vectorizer_requests_total",
				Help:      "Tracks requests to vectorizer",
			},
			[]string{
				"op",
				"status",
			},
		),
		vectorizerRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "vectorizer_duration_seconds",
				Help:      "Tracks request durations in vectorizer",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"op",
				"status",
			},
		),
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

func (m *Metrics) IncFilestorageRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.filestorageRequestDurations.With(labels).Observe(dur)
	m.filestorageRequests.With(labels).Inc()
}

func (m *Metrics) IncSemanticIndexRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.semanticIndexRequests.With(labels).Inc()
	m.semanticIndexRequestDuration.With(labels).Observe(dur)
}

func (m *Metrics) IncVectorizerRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.vectorizerRequests.With(labels).Inc()
	m.vectorizerRequestDuration.With(labels).Observe(dur)
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
		m.filestorageRequests,
		m.filestorageRequestDurations,
		m.semanticIndexRequests,
		m.semanticIndexRequestDuration,
		m.vectorizerRequests,
		m.vectorizerRequestDuration,
		m.documentRepositoryRequests,
		m.documentRepositoryRequestDuration,
		m.natsMicroRequests,
		m.natsMicroRequestDurations,
		m.natsMicroSizes,
	}
}
