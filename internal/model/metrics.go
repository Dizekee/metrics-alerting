package models

import "sync"

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

type Storage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
}

type MemStorage struct {
	mu      sync.RWMutex
	counter map[string]int64
	gauges  map[string]float64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counter: make(map[string]int64),
		gauges:  make(map[string]float64),
	}
}

func (s *MemStorage) AddCounter(name string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter[name] += value
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauges[name] = value
}
