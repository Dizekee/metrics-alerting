package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSender_post(t *testing.T) {
	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer okServer.Close()

	badRequestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer badRequestServer.Close()

	internalErrorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer internalErrorServer.Close()

	tests := []struct {
		name string

		curl    string
		cclient *http.Client

		url     string
		wantErr bool
	}{
		{
			name: "successful request",

			curl:    okServer.URL,
			cclient: okServer.Client(),

			url:     okServer.URL + "/update/counter/PollCount/1",
			wantErr: false,
		},
		{
			name: "bad request status",

			curl:    badRequestServer.URL,
			cclient: badRequestServer.Client(),

			url:     badRequestServer.URL + "/update/gauge/Alloc/1.5",
			wantErr: true,
		},
		{
			name: "internal server error status",

			curl:    internalErrorServer.URL,
			cclient: internalErrorServer.Client(),

			url:     internalErrorServer.URL + "/update/counter/PollCount/1",
			wantErr: true,
		},
		{
			name: "invalid URL",

			curl:    okServer.URL,
			cclient: okServer.Client(),

			url:     "http://[::1]:namedport/update/counter/PollCount/1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSender(tt.curl, tt.cclient)

			err := s.post(tt.url)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("post() failed: %v", err)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("post() succeeded unexpectedly")
			}
		})
	}
}

func TestSender_SendCounter(t *testing.T) {
	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer okServer.Close()

	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errorServer.Close()

	tests := []struct {
		name string

		url    string
		client *http.Client

		metricName string
		value      int64

		wantErr bool
	}{
		{
			name: "successful counter request",

			url:    okServer.URL,
			client: okServer.Client(),

			metricName: "PollCount",
			value:      42,

			wantErr: false,
		},
		{
			name: "server returns internal error",

			url:    errorServer.URL,
			client: errorServer.Client(),

			metricName: "PollCount",
			value:      1,

			wantErr: true,
		},
		{
			name: "invalid base URL",

			url:    "http://[::1]:namedport",
			client: okServer.Client(),

			metricName: "PollCount",
			value:      1,

			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSender(tt.url, tt.client)

			err := s.SendCounter(tt.metricName, tt.value)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("SendCounter() failed: %v", err)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("SendCounter() succeeded unexpectedly")
			}
		})
	}
}

func TestSender_SendGauge(t *testing.T) {
	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer okServer.Close()

	badRequestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer badRequestServer.Close()

	internalErrorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer internalErrorServer.Close()

	tests := []struct {
		name string

		url    string
		client *http.Client

		metricName string
		value      float64

		wantErr bool
	}{
		{
			name: "successful gauge request",

			url:    okServer.URL,
			client: okServer.Client(),

			metricName: "Alloc",
			value:      123.456,

			wantErr: false,
		},
		{
			name: "zero value gauge",

			url:    okServer.URL,
			client: okServer.Client(),

			metricName: "GCCPUFraction",
			value:      0,

			wantErr: false,
		},
		{
			name: "bad request status",

			url:    badRequestServer.URL,
			client: badRequestServer.Client(),

			metricName: "Alloc",
			value:      1.5,

			wantErr: true,
		},
		{
			name: "internal server error status",

			url:    internalErrorServer.URL,
			client: internalErrorServer.Client(),

			metricName: "Alloc",
			value:      1.5,

			wantErr: true,
		},
		{
			name: "invalid base URL",

			url:    "http://[::1]:namedport",
			client: okServer.Client(),

			metricName: "Alloc",
			value:      1.5,

			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSender(tt.url, tt.client)

			err := s.SendGauge(tt.metricName, tt.value)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("SendGauge() failed: %v", err)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("SendGauge() succeeded unexpectedly")
			}
		})
	}
}
