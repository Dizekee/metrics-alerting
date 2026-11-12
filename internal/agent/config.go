package agent

import "time"

type SenderConfig struct {
	URLAddr        string
	ReportInterval time.Duration
}

type PollerConfig struct {
	PollInterval time.Duration
}

type AgentConfig struct {
	Sender SenderConfig
	Poller PollerConfig
}
