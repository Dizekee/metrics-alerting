package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	models "github.com/Dizekee/metrics-alerting.git/internal/model"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ct := r.Header.Get("Content-type")
	if ct != "" && !strings.HasPrefix(ct, "text/plain") {
		http.Error(w, "Content-Type must be text/plain", http.StatusUnsupportedMediaType)
		return
	}

	metricType := r.PathValue("type")
	metricName := r.PathValue("name")
	metricValue := r.PathValue("value")

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
	}
	w.WriteHeader(http.StatusOK)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lr := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lr, r)

		log.Printf("HTTP %s %s | Status: %d | Duration: %v", r.Method, r.URL.Path, lr.statusCode, time.Since(start))
	})
}
