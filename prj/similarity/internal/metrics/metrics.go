package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "simrep"
	subsystem = "similarity"
)

type Metrics struct {
	filestorageRequests                     *prometheus.CounterVec
	filestorageRequestDurations             *prometheus.HistogramVec
	analyzeHistoryRepositoryRequests        *prometheus.CounterVec
	analyzeHistoryRepositoryRequestDuration *prometheus.HistogramVec
	documentRepositoryRequests              *prometheus.CounterVec
	documentRepositoryRequestDuration       *prometheus.HistogramVec
	similarityIndexRequests                 *prometheus.CounterVec
	similarityIndexRequestDuration          *prometheus.HistogramVec
	natsMicroRequests                       *prometheus.CounterVec
	natsMicroRequestDurations               *prometheus.HistogramVec
	natsMicroSizes                          *prometheus.HistogramVec
	httpRequests                            *prometheus.CounterVec
	httpRequestDurations                    *prometheus.HistogramVec
	httpSizes                               *prometheus.HistogramVec
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
		analyzeHistoryRepositoryRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "analyze_history_repository_requests_total",
				Help:      "Tracks requests to analyze history repository",
			},
			[]string{
				"op",
				"status",
			},
		),
		analyzeHistoryRepositoryRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "analyze_history_repository_request_duration_seconds",
				Help:      "Tracks request durations in analyze history repository",
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
		similarityIndexRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "similarity_index_requests_total",
				Help:      "Tracks requests to similarity index",
			},
			[]string{
				"index",
				"op",
				"status",
			},
		),
		similarityIndexRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "similarity_index_duration_seconds",
				Help:      "Tracks request durations in similarity index",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"index",
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
		httpRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_server_requests_total",
				Help:      "Tracks requests to http_server",
			},
			[]string{
				"path",
				"status",
			},
		),
		httpRequestDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_server_request_duration_seconds",
				Help:      "Tracks request durations in http server",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"path",
				"status",
			},
		),
		httpSizes: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_server_request_sizes_bytes",
				Help:      "Tracks request sizes in http server",
				Buckets:   prometheus.ExponentialBuckets(1024*256, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"path",
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

func (m *Metrics) IncAnalyzeHistoryRepositoryRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.analyzeHistoryRepositoryRequests.With(labels).Inc()
	m.analyzeHistoryRepositoryRequestDuration.With(labels).Observe(dur)
}

func (m *Metrics) IncDocumentRepositoryRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.documentRepositoryRequests.With(labels).Inc()
	m.documentRepositoryRequestDuration.With(labels).Observe(dur)
}

func (m *Metrics) IncSimilarityIndexRequests(index string, op string, status string, dur float64) {
	labels := prometheus.Labels{
		"index":  index,
		"op":     op,
		"status": status,
	}

	m.similarityIndexRequests.With(labels).Inc()
	m.similarityIndexRequestDuration.With(labels).Observe(dur)
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

func (m *Metrics) IncHTTPRequest(op string, status string, size int, dur float64) {
	labels := prometheus.Labels{
		"path":   op,
		"status": status,
	}

	m.httpRequests.With(labels).Inc()
	m.httpRequestDurations.With(labels).Observe(dur)
	m.httpSizes.With(labels).Observe(float64(size))
}

func (m *Metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.filestorageRequests,
		m.filestorageRequestDurations,
		m.analyzeHistoryRepositoryRequests,
		m.analyzeHistoryRepositoryRequestDuration,
		m.documentRepositoryRequests,
		m.documentRepositoryRequestDuration,
		m.similarityIndexRequests,
		m.similarityIndexRequestDuration,
		m.natsMicroRequests,
		m.natsMicroRequestDurations,
		m.natsMicroSizes,
		m.httpRequests,
		m.httpRequestDurations,
		m.httpSizes,
	}
}
