package agent

import (
	"testing"

	models "github.com/Dizekee/metrics-alerting.git/internal/model"
)

func TestMemStorage_AddCounter(t *testing.T) {
	tests := []struct {
		name string

		metricName string
		value      int64

		want int64
	}{
		{
			name: "add positive value",

			metricName: "PollCount",
			value:      1,

			want: 1,
		},
		{
			name: "add several values",

			metricName: "PollCount",
			value:      5,

			want: 5,
		},
		{
			name: "add zero value",

			metricName: "PollCount",
			value:      0,

			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := models.NewMemStorage()

			s.AddCounter(tt.metricName, tt.value)

			got, ok := s.GetCounter(tt.metricName)
			if !ok {
				t.Fatalf("counter %q not found", tt.metricName)
			}

			if got != tt.want {
				t.Errorf("AddCounter() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	tests := []struct {
		name string

		metricName string
		value      float64

		want float64
	}{
		{
			name: "set positive value",

			metricName: "Alloc",
			value:      123.456,

			want: 123.456,
		},
		{
			name: "set zero value",

			metricName: "GCCPUFraction",
			value:      0,

			want: 0,
		},
		{
			name: "set negative value",

			metricName: "Temperature",
			value:      -12.5,

			want: -12.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := models.NewMemStorage()

			s.SetGauge(tt.metricName, tt.value)

			got, ok := s.GetGauge(tt.metricName)
			if !ok {
				t.Fatalf("gauge %q not found", tt.metricName)
			}

			if got != tt.want {
				t.Errorf("SetGauge() = %f, want %f", got, tt.want)
			}
		})
	}
}
