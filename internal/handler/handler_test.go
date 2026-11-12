package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/Dizekee/metrics-alerting.git/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAlert(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name        string
		method      string
		url         string
		contentType string
		want        want
	}{
		{
			name:        "counter ok",
			method:      http.MethodPost,
			url:         "/update/counter/testCounter/527",
			contentType: "text/plain",
			want: want{
				code:     http.StatusOK,
				response: "",
			},
		},
		{
			name:        "gauge ok",
			method:      http.MethodPost,
			url:         "/update/gauge/testGauge/3.14",
			contentType: "text/plain",
			want: want{
				code:     http.StatusOK,
				response: "",
			},
		},
		{
			name:        "method not allowed",
			method:      http.MethodGet,
			url:         "/update/counter/testCounter/527",
			contentType: "text/plain",
			want: want{
				code:     http.StatusMethodNotAllowed,
				response: "Method not allowed\n",
			},
		},
		{
			name:        "wrong content-type",
			method:      http.MethodPost,
			url:         "/update/counter/testCounter/527",
			contentType: "application/json",
			want: want{
				code:     http.StatusUnsupportedMediaType,
				response: "Content-Type must be text/plain\n",
			},
		},
		{
			name:        "invalid counter value",
			method:      http.MethodPost,
			url:         "/update/counter/testCounter/abc",
			contentType: "text/plain",
			want: want{
				code:     http.StatusBadRequest,
				response: "invalid counter type\n",
			},
		},
		{
			name:        "unknown metric type",
			method:      http.MethodPost,
			url:         "/update/unknown/testMetric/1",
			contentType: "text/plain",
			want: want{
				code:     http.StatusBadRequest,
				response: "uncnow metric type\n",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := models.NewMemStorage()
			h := NewHandler(storage)
			r := chi.NewRouter()
			r.HandleFunc("/update/{type}/{name}/{value}", h.UpdateMetric)

			req := httptest.NewRequest(test.method, test.url, nil)
			if test.contentType != "" {
				req.Header.Set("Content-Type", test.contentType)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
		})
	}
}
