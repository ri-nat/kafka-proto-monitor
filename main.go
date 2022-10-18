package main

import (
	"fmt"
	"os"

	monitor "github.com/ri-nat/kafka-proto-monitor/internal/app"
)

func main() {
	if err := monitor.RunMonitor(); err != nil {
		fmt.Printf("Can't start `kafka-proto-monitor`: %v\n", err)
		os.Exit(1)
	}
}
