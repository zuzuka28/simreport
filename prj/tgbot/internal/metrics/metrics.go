package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "simrep"
	subsystem = "tgbot"
)

type Metrics struct {
	documentRepositoryRequests          *prometheus.CounterVec
	documentRepositoryRequestDuration   *prometheus.HistogramVec
	userStateRepositoryRequests         *prometheus.CounterVec
	userStateRepositoryRequestDuration  *prometheus.HistogramVec
	similarityRepositoryRequests        *prometheus.CounterVec
	similarityRepositoryRequestDuration *prometheus.HistogramVec
	botErrors                           *prometheus.CounterVec
	botRequestsByUser                   *prometheus.CounterVec
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
		userStateRepositoryRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "user_state_repository_requests_total",
				Help:      "Tracks requests to user state repository",
			},
			[]string{
				"op",
				"status",
			},
		),
		userStateRepositoryRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "user_state_repository_duration_seconds",
				Help:      "Tracks request durations in user state repository",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"op",
				"status",
			},
		),
		similarityRepositoryRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "similarity_repository_requests_total",
				Help:      "Tracks requests to similarity repository",
			},
			[]string{
				"op",
				"status",
			},
		),
		similarityRepositoryRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "similarity_repository_duration_seconds",
				Help:      "Tracks request durations in similarity repository",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), //nolint:gomnd,mnd
			},
			[]string{
				"op",
				"status",
			},
		),
		botRequestsByUser: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "tgbot_user_requests_total",
				Help:      "Tracks requests to tgbot by each user",
			},
			[]string{
				"username",
				"userid",
			},
		),
		botErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{ //nolint:exhaustruct
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "tgbot_errors_total",
				Help:      "Tracks errors by tgbot handlers",
			},
			[]string{
				"desc",
			},
		),
	}
}

// FIXME: this method makes big cardinality for metric.
func (m *Metrics) IncBotRequestsByUser(username, userid string) {
	labels := prometheus.Labels{
		"username": username,
		"userid":   userid,
	}

	m.botRequestsByUser.With(labels).Inc()
}

func (m *Metrics) IncBotErrors(desc string) {
	labels := prometheus.Labels{
		"desc": desc,
	}

	m.botErrors.With(labels).Inc()
}

func (m *Metrics) IncSimilarityRepositoryRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.similarityRepositoryRequests.With(labels).Inc()
	m.similarityRepositoryRequestDuration.With(labels).Observe(dur)
}

func (m *Metrics) IncUserStateRepositoryRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.userStateRepositoryRequests.With(labels).Inc()
	m.userStateRepositoryRequestDuration.With(labels).Observe(dur)
}

func (m *Metrics) IncDocumentRepositoryRequests(op string, status string, dur float64) {
	labels := prometheus.Labels{
		"op":     op,
		"status": status,
	}

	m.documentRepositoryRequests.With(labels).Inc()
	m.documentRepositoryRequestDuration.With(labels).Observe(dur)
}

func (m *Metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.documentRepositoryRequests,
		m.documentRepositoryRequestDuration,
		m.userStateRepositoryRequests,
		m.userStateRepositoryRequestDuration,
		m.similarityRepositoryRequests,
		m.similarityRepositoryRequestDuration,
		m.botErrors,
		m.botRequestsByUser,
	}
}
