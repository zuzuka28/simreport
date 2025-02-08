package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "simrep"
	subsystem = "document"
)

type Metrics struct {
	attributeRepositoryRequests         *prometheus.CounterVec
	attributeRepositoryRequestDurations *prometheus.HistogramVec
	documentRepositoryRequests          *prometheus.CounterVec
	documentRepositoryRequestDurations  *prometheus.HistogramVec
	documentStatusRepositoryUpdates     *prometheus.CounterVec
	filestorageRequests                 *prometheus.CounterVec
	filestorageRequestDurations         *prometheus.HistogramVec
	filestorageUploads                  *prometheus.CounterVec
	natsMicroRequests                   *prometheus.CounterVec
	natsMicroRequestDurations           *prometheus.HistogramVec
	natsMicroSizes                      *prometheus.HistogramVec
	httpRequests                        *prometheus.CounterVec
	httpRequestDurations                *prometheus.HistogramVec
	httpSizes                           *prometheus.HistogramVec
}

func New() *Metrics {
	return &Metrics{
		attributeRepositoryRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "attribute_repository_requests_total",
				Help:      "Tracks document status updates",
			},
			[]string{
				"op",
				"status",
			},
		),
		attributeRepositoryRequestDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "attribute_repository_request_duration_seconds",
				Help:      "Tracks request durations in attribute repository",
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
				Help:      "Tracks document status updates",
			},
			[]string{
				"op",
				"status",
			},
		),
		documentRepositoryRequestDurations: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "document_repository_request_duration_seconds",
				Help:      "Tracks request durations in document repository",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"op",
				"status",
			},
		),
		documentStatusRepositoryUpdates: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "document_status_repository_updates_total",
				Help:      "Tracks document status updates",
			},
			[]string{
				"status",
				"result",
			},
		),
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
		filestorageUploads: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "filestorage_updates_total",
				Help:      "Tracks filestorage updates",
			},
			[]string{
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

func (m *Metrics) IncAttributeRepositoryRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.attributeRepositoryRequestDurations.With(labels).Observe(dur)
	m.attributeRepositoryRequests.With(labels).Inc()
}

func (m *Metrics) IncDocumentRepositoryRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.documentRepositoryRequestDurations.With(labels).Observe(dur)
	m.documentRepositoryRequests.With(labels).Inc()
}

func (m *Metrics) IncDocumentStatusRepositoryUpdates(status string, result string) {
	m.documentStatusRepositoryUpdates.With(prometheus.Labels{
		"status": status,
		"result": result,
	}).Inc()
}

func (m *Metrics) IncFilestorageRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.filestorageRequestDurations.With(labels).Observe(dur)
	m.filestorageRequests.With(labels).Inc()
}

func (m *Metrics) IncFilestorageUploads(status string) {
	m.filestorageUploads.With(prometheus.Labels{
		"status": status,
	}).Inc()
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
		m.attributeRepositoryRequests,
		m.attributeRepositoryRequestDurations,
		m.documentRepositoryRequests,
		m.documentRepositoryRequestDurations,
		m.documentStatusRepositoryUpdates,
		m.filestorageRequests,
		m.filestorageRequestDurations,
		m.filestorageUploads,
		m.natsMicroRequests,
		m.natsMicroRequestDurations,
		m.natsMicroSizes,
		m.httpRequests,
		m.httpRequestDurations,
		m.httpSizes,
	}
}
