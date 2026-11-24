package agent

import (
	"fmt"
	"net/http"
)

type Sender struct {
	url    string
	client *http.Client
}

func NewSender(url string, client *http.Client) *Sender {
	return &Sender{url: url, client: client}
}

func (s *Sender) SendCounter(name string, value int64) error {
	url := fmt.Sprintf("%s/update/counter/%s/%d", s.url, name, value)
	return s.post(url)
}

func (s *Sender) SendGauge(name string, value float64) error {
	url := fmt.Sprintf("%s/update/gauge/%s/%f", s.url, name, value)
	return s.post(url)
}

func (s *Sender) post(url string) error {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}
	return nil
}
