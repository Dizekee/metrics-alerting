package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	models "github.com/Dizekee/metrics-alerting.git/internal/model"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	storage models.Storage
}

func NewHandler(storage models.Storage) *Handler {
	return &Handler{storage: storage}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lr *loggingResponseWriter) WriteHeader(code int) {
	lr.statusCode = code
	lr.ResponseWriter.WriteHeader(code)
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-type")
	if ct != "" && !strings.HasPrefix(ct, "text/plain") {
		http.Error(w, "Content-Type must be text/plain", http.StatusUnsupportedMediaType)
		return
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricName == "" {
		http.Error(w, "metric name is required", http.StatusNotFound)
		return
	}

	switch metricType {
	case models.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter type", http.StatusBadRequest)
			return
		}
		h.storage.AddCounter(metricName, value)
	case models.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "invalid counter type", http.StatusBadRequest)
			return
		}
		h.storage.SetGauge(metricName, value)
	default:
		http.Error(w, "uncnow metric type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	var valueStr string
	var found bool

	switch metricType {
	case models.Counter:
		value, found := h.storage.GetCounter(metricName)
		if found {
			valueStr = strconv.FormatInt(value, 10)
		}
	case models.Gauge:
		value, found := h.storage.GetGauge(metricName)
		if found {
			valueStr = strconv.FormatFloat(value, 'f', -1, 64)
		}
	default:
		http.Error(w, "unknow metric type", http.StatusBadRequest)
		return
	}

	if !found {
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(valueStr))
}

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	gauges, counters := h.storage.Snapshot()

	_, _ = w.Write([]byte("<html><body><h1></h1><ul>"))

	for name, value := range gauges {
		_, _ = fmt.Fprintf(w, "<li>gauge %s: %v</li>\n", name, value)
	}

	for name, value := range counters {
		_, _ = fmt.Fprintf(w, "<li>counter %s: %v</li>\n", name, value)
	}

	_, _ = w.Write([]byte("</ul></body></html>"))
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lr := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lr, r)

		log.Printf("HTTP %s %s | Status: %d | Duration: %v", r.Method, r.URL.Path, lr.statusCode, time.Since(start))
	})
}
