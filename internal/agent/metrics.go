package agent

import "sync"

type Storage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
	Snapshot() (map[string]float64, map[string]int64)
}

type MemStorage struct {
	mu       sync.RWMutex
	counters map[string]int64
	gauges   map[string]float64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
	}
}

func (s *MemStorage) AddCounter(name string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counters[name] += value
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauges[name] = value
}

func (s *MemStorage) Snapshot() (map[string]float64, map[string]int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	gauges := make(map[string]float64, len(s.gauges))
	for n, v := range s.gauges {
		gauges[n] = v
	}

	counter := make(map[string]int64, len(s.counters))
	for n, v := range s.counters {
		counter[n] = v
	}
	return gauges, counter
}
