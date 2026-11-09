package main

import (
	"flag"
)

var (
	urlAddr string
	pollSec int
	reportSec int
)

func parseFlags() {
	flag.StringVar(&urlAddr, "u", "http://localhost:8080", "server URL")
	flag.IntVar(&pollSec, "p", 2, "pool interval")
	flag.IntVar(&reportSec, "r", 10, "report interval")

	flag.Parse()
}