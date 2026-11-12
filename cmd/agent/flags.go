package main

import (
	"flag"
	"time"

	"github.com/Dizekee/metrics-alerting.git/internal/agent"
)

func parseFlags() agent.AgentConfig {
	var (
		urlAddr   string
		pollSec   int
		reportSec int
	)

	flag.StringVar(&urlAddr, "u", "http://localhost:8080", "server URL")
	flag.IntVar(&pollSec, "p", 2, "pool interval")
	flag.IntVar(&reportSec, "r", 10, "report interval")

	flag.Parse()

	return agent.AgentConfig{
		Sender: agent.SenderConfig{
			UrlAddr:   urlAddr,
			ReportSec: time.Duration(reportSec) * time.Second,
		},
		Poller: agent.PollerConfig{
			PollSec: time.Duration(pollSec) * time.Second,
		},
	}
}
