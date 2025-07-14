package main

import (
	"fmt"
	"time"

	"github.com/iamlucif3r/netwatchdog/config"
	"github.com/iamlucif3r/netwatchdog/core"
)

func main() {
	config.LoadConfig("config.json")
	fmt.Println("📡 Starting NetWatchdog...")

	for {
		core.RunScan()
		time.Sleep(1 * time.Minute)
	}
}
