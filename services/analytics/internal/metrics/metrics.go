package metrics

import (
	"time"

	"github.com/egor200512/URL_shortener/shared/pkg/events"
	"github.com/prometheus/client_golang/prometheus"
)

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

	for _, eventType := range []string{events.LinkCreated, events.LinkFetched, events.LinkDeleted} {
		EventsConsumedTotal.WithLabelValues(eventType).Add(0)
	}
}

func RecordConsumed(eventType string, startedAt time.Time) {
	EventsConsumedTotal.WithLabelValues(eventType).Inc()
	EventHandleDurationSeconds.Observe(time.Since(startedAt).Seconds())
}

func RecordConsumeError(startedAt time.Time) {
	EventsConsumeErrorsTotal.Inc()
	EventHandleDurationSeconds.Observe(time.Since(startedAt).Seconds())
}
