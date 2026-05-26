package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	EventsConsumedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "url_shortener",
			Subsystem: "analytics",
			Name:      "events_consumed_total",
			Help:      "Total number of successfully consumed link events.",
		},
		[]string{"event_type"},
	)

	EventsConsumeErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "url_shortener",
			Subsystem: "analytics",
			Name:      "events_consume_errors_total",
			Help:      "Total number of link event consume errors.",
		},
	)

	EventHandleDurationSeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "url_shortener",
			Subsystem: "analytics",
			Name:      "event_handle_duration_seconds",
			Help:      "Duration of link event handling in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
	)
)

func Register(registry *prometheus.Registry) {
	registry.MustRegister(
		EventsConsumedTotal,
		EventsConsumeErrorsTotal,
		EventHandleDurationSeconds,
	)
}
