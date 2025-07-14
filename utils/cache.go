package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/iamlucif3r/netwatchdog/scan"
)

type KnownDevices struct {
	Devices map[string]scan.Device `json:"devices"`
}

var mu sync.Mutex

func LoadCache(path string) KnownDevices {
	mu.Lock()
	defer mu.Unlock()

	var known KnownDevices
	known.Devices = make(map[string]scan.Device)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {

			log.Printf("⚠️ Cache file not found, initializing empty cache...")
			return known
		}
		log.Printf("❌ Error reading cache file: %v", err)
		return known
	}

	if err := json.Unmarshal(data, &known); err != nil {
		log.Printf("❌ Invalid cache file format. Starting fresh. Error: %v", err)
		known.Devices = make(map[string]scan.Device)
	}

	return known
}

func IsKnownDevice(mac string, known KnownDevices) bool {
	_, exists := known.Devices[mac]
	return exists
}

func UpdateCache(path string, known KnownDevices, newDevs []scan.Device) {
	mu.Lock()
	defer mu.Unlock()

	for _, dev := range newDevs {
		known.Devices[dev.MAC] = dev
	}

	data, err := json.MarshalIndent(known, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal updated cache: %v", err)
		return
	}

	err = os.MkdirAll("data", 0755)
	if err != nil {
		log.Printf("❌ Failed to create data directory: %v", err)
		return
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		log.Printf("❌ Failed to write cache file: %v", err)
		return
	}

	fmt.Println("💾 Cache updated.")
}
