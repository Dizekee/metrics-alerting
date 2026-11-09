package agent

import (
	"testing"
	"time"
)

func TestAgent_Poll(t *testing.T) {
	tests := []struct {
		name           string
		serverURL      string
		pollInterval   time.Duration
		reportInterval time.Duration
	}{
		{
			name:           "successful poll collects all metrics",
			serverURL:      "http://localhost:8080",
			pollInterval:   2 * time.Second,
			reportInterval: 10 * time.Second,
		},
	}

	// Список всех ожидаемых gauge-метрик
	expectedGauges := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAgent(tt.serverURL, tt.pollInterval, tt.reportInterval)

			a.Poll()
			a.Poll()

			gauges, counters := a.storage.Snapshot()

			pollCount, ok := counters["PollCount"]
			if !ok {
				t.Fatalf("PollCount metric is missing")
			}
			if pollCount != 2 {
				t.Errorf("expected PollCount to be 2, got %d", pollCount)
			}

			for _, name := range expectedGauges {
				if _, ok := gauges[name]; !ok {
					t.Errorf("gauge metric %q is missing", name)
				}
			}

			randomValue := gauges["RandomValue"]
			if randomValue < 0.0 || randomValue >= 1.0 {
				t.Errorf("RandomValue out of range [0, 1): got %f", randomValue)
			}

			if gauges["Alloc"] < 0 {
				t.Errorf("Alloc cannot be negative: got %f", gauges["Alloc"])
			}
		})
	}
}
