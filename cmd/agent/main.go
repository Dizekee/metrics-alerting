package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dizekee/metrics-alerting.git/internal/agent"
)

func main() {
	config := parseFlags()

	agent := agent.NewAgent(config)

	go agent.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Println("The agent is running. Press Ctrl+C to stop.")
	<-stop

	log.Println("A stop signal has been received. We are completing the work...")
	agent.Stop()

	log.Println("The agent has been stopped.")
}
