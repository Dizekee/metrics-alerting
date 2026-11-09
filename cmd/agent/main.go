package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dizekee/metrics-alerting.git/internal/agent"
)

func main() {
	parseFlags()

	agent := agent.NewAgent(
		urlAddr,
		time.Duration(pollSec)*time.Second,
		time.Duration(reportSec)*time.Second,
	)

	go agent.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Println("The agent is running. Press Ctrl+C to stop.")
	<-stop

	log.Println("A stop signal has been received. We are completing the work...")
	agent.Stop()

	log.Println("The agent has been stopped.")
}
