package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCounter(t *testing.T) {
	tests := []struct {
		name  string
		calls []struct {
			key   string
			value int64
		}
		key  string
		want int64
	}{
		{
			name: "one add",
			calls: []struct {
				key   string
				value int64
			}{
				{"counter", 5},
			},
			key:  "counter",
			want: 5,
		},
		{
			name: "several additions",
			calls: []struct {
				key   string
				value int64
			}{
				{"test", 5},
				{"test", 10},
				{"test", 18},
			},
			key:  "test",
			want: 33,
		},
		{
			name: "negative result",
			calls: []struct {
				key   string
				value int64
			}{
				{"error", 5},
				{"error", 10},
				{"error", 18},
			},
			key:  "error",
			want: 8,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewMemStorage()
			for _, call := range test.calls {
				storage.AddCounter(call.key, call.value)
			}
			assert.Equal(t, storage.counter[test.key], test.want)
		})
	}
}

func TestSetGauge(t *testing.T) {
	tests := []struct {
		name  string
		calls []struct {
			key   string
			value float64
		}
		key  string
		want float64
	}{
		{
			name: "one set",
			calls: []struct {
				key   string
				value float64
			}{
				{"temp", 5.8},
			},
			key:  "temp",
			want: 5.8,
		},
		{
			name: "several sets",
			calls: []struct {
				key   string
				value float64
			}{
				{"requests", 3.8},
				{"requests", 68.5},
				{"requests", 43.2},
			},
			key:  "requests",
			want: 43.2,
		},
		{
			name: "negative result",
			calls: []struct {
				key   string
				value float64
			}{
				{"errors", 8.5},
				{"errors", 66.3},
				{"errors", 75.4},
			},
			key:  "errors",
			want: 35.7,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewMemStorage()
			for _, call := range test.calls {
				storage.SetGauge(call.key, call.value)
			}
			assert.Equal(t, storage.gauges[test.key], test.want)
		})
	}
}
