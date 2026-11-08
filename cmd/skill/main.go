package main

import (
	"fmt"
	"net/http"

	"github.com/Dizekee/metrics-alerting.git/internal/handler"
	models "github.com/Dizekee/metrics-alerting.git/internal/model"
	"github.com/go-chi/chi/v5"
)

func main() {
	s := models.NewMemStorage()
	h := handler.NewHandler(s)

	r := chi.NewRouter()

	r.Use(handler.LoggingMiddleware)

	r.Post("/update/{type}/{name}/{value}", h.UpdateMetric)
	r.Get("/value/{type}/{name}", h.GetMetric)
	r.Get("/", h.GetAllMetrics)
	r.Post("/", webhook)

	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println("Running server on", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, http.HandlerFunc(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write([]byte(`
      {
        "response": {
          "text": "Извините, я пока ничего не умею"
        },
        "version": "1.0"
      }
    `))
}
