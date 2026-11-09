package main

import (
	"log"
	"net/http"

	"github.com/Dizekee/metrics-alerting.git/internal/handler"
	models "github.com/Dizekee/metrics-alerting.git/internal/model"
)

func main() {
	storage := models.NewMemStorage()
	h := handler.NewHandler(storage)

	mux := http.NewServeMux()
	
	log.Printf("starting server on %s", ":8080")
	mux.HandleFunc("/update/{type}/{name}/{value}", h.UpdateMetric)
    if err := run(); err != nil {
        panic(err)
    }
}

func run() error {
    return http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
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
