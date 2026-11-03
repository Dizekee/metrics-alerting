package handler

import (
	"net/http"
	"strconv"
	"strings"

	models "github.com/Dizekee/metrics-alerting.git/internal/model"
)

type Handler struct {
	storage models.Storage
}

func NewHandler(storage models.Storage) *Handler{
	return &Handler{storage: storage}
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