package agent

import (
	"log"
	"math/rand/v2"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type Agent struct {
	storage        *MemStorage
	sender         *Sender
	pollInterval   time.Duration
	reportInterval time.Duration
	stop           chan struct{}
	wg             sync.WaitGroup
}

func NewAgent(serverURL string, pollInterval, reportInterval time.Duration) *Agent {
	return &Agent{
		storage:        NewMemStorage(),
		sender:         NewSender(serverURL, &http.Client{}),
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		stop:           make(chan struct{}),
	}
}

func (a *Agent) Run() {
	a.wg.Add(2)

	go func() {
		defer a.wg.Done()
		a.poolLoop()
	}()
	go func() {
		defer a.wg.Done()
		a.reportLoop()
	}()
}

func (a *Agent) Stop() {
	close(a.stop)
	a.wg.Wait()
}

func (a *Agent) poolLoop() {
	tim := time.NewTicker(a.pollInterval)
	defer tim.Stop()
	for {
		select {
		case <-tim.C:
			a.Poll()
		case <-a.stop:
			return
		}
	}
}

func (a *Agent) reportLoop() {
	tim := time.NewTicker(a.reportInterval)
	defer tim.Stop()
	for {
		select {
		case <-tim.C:
			a.Report()
		case <-a.stop:
			return
		}
	}
}

func (a *Agent) Poll() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s := a.storage
	s.SetGauge("Alloc", float64(m.Alloc))
	s.SetGauge("BuckHashSys", float64(m.BuckHashSys))
	s.SetGauge("Frees", float64(m.Frees))
	s.SetGauge("GCCPUFraction", m.GCCPUFraction)
	s.SetGauge("GCSys", float64(m.GCSys))
	s.SetGauge("HeapAlloc", float64(m.HeapAlloc))
	s.SetGauge("HeapIdle", float64(m.HeapIdle))
	s.SetGauge("HeapInuse", float64(m.HeapInuse))
	s.SetGauge("HeapObjects", float64(m.HeapObjects))
	s.SetGauge("HeapReleased", float64(m.HeapReleased))
	s.SetGauge("HeapSys", float64(m.HeapSys))
	s.SetGauge("LastGC", float64(m.LastGC))
	s.SetGauge("Lookups", float64(m.Lookups))
	s.SetGauge("MCacheInuse", float64(m.MCacheInuse))
	s.SetGauge("MCacheSys", float64(m.MCacheSys))
	s.SetGauge("MSpanInuse", float64(m.MSpanInuse))
	s.SetGauge("MSpanSys", float64(m.MSpanSys))
	s.SetGauge("Mallocs", float64(m.Mallocs))
	s.SetGauge("NextGC", float64(m.NextGC))
	s.SetGauge("NumForcedGC", float64(m.NumForcedGC))
	s.SetGauge("NumGC", float64(m.NumGC))
	s.SetGauge("OtherSys", float64(m.OtherSys))
	s.SetGauge("PauseTotalNs", float64(m.PauseTotalNs))
	s.SetGauge("StackInuse", float64(m.StackInuse))
	s.SetGauge("StackSys", float64(m.StackSys))
	s.SetGauge("Sys", float64(m.Sys))
	s.SetGauge("TotalAlloc", float64(m.TotalAlloc))
	s.SetGauge("RandomValue", rand.Float64())

	s.AddCounter("PollCount", 1)
}

func (a *Agent) Report() {
	gauges, counter := a.storage.Snapshot()
	for n, v := range gauges {
		if err := a.sender.SendGauge(n, v); err != nil {
			log.Printf("send gauge %s: %v", n, v)
		}
	}
	for n, v := range counter {
		if err := a.sender.SendCounter(n, v); err != nil {
			log.Printf("send counter %s: %v", n, v)
		}
	}
}
