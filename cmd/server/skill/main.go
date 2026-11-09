package main

import (
	"log"
	"net/http"

	"github.com/Dizekee/metrics-alerting.git/internal/handler"
	models "github.com/Dizekee/metrics-alerting.git/internal/model"
)

func main() {
	s := models.NewMemStorage()
	h := handler.NewHandler(s)

	mux := http.NewServeMux()

	mux.HandleFunc("/update/{type}/{name}/{value}", h.UpdateMetric)
	mux.HandleFunc("/", webhook)
	loggedMux := handler.LoggingMiddleware(mux)

	addr := ":8080"
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, loggedMux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
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
