package config

import (
	"flag"
)

var flagRunAddr = flag.String("a", "localhost:8080", "address and port to run server")

func ParseFlags() string {

	flag.Parse()
	return *flagRunAddr
}
