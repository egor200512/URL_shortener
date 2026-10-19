package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func TestRegisterExposesDashboardMetrics(t *testing.T) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	Register(registry)

	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}

	got := make(map[string]bool, len(metricFamilies))
	for _, family := range metricFamilies {
		got[family.GetName()] = true
	}

	for _, name := range []string{
		"go_gc_duration_seconds",
		"go_goroutines",
		"go_memstats_heap_alloc_bytes",
		"process_cpu_seconds_total",
		"process_resident_memory_bytes",
		"process_start_time_seconds",
		"url_shortener_analytics_events_consumed_total",
		"url_shortener_analytics_events_consume_errors_total",
		"url_shortener_analytics_event_handle_duration_seconds",
	} {
		if !got[name] {
			t.Errorf("metric %q is not registered", name)
		}
	}
}
