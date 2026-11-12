package agent

import "time"

type SenderConfig struct {
	UrlAddr   string
	ReportSec time.Duration
}

type PollerConfig struct {
	PollSec time.Duration
}

type AgentConfig struct {
	Sender SenderConfig
	Poller PollerConfig
}
