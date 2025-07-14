package core

import (
	"fmt"
	"log"
	"time"

	"github.com/iamlucif3r/netwatchdog/alert"
	"github.com/iamlucif3r/netwatchdog/config"
	"github.com/iamlucif3r/netwatchdog/scan"
	"github.com/iamlucif3r/netwatchdog/utils"
	"github.com/iamlucif3r/netwatchdog/wifi"
)

var firstRunDone = false

func RunScan() {
	ssid := wifi.GetSSID(config.Config.Interface)
	if ssid == "" || ssid == "Unknown" {
		ssid = "N/A"
	}

	known := utils.LoadCache(config.Config.CacheFile)

	devices, err := scan.ScanLocal(config.Config.Interface)
	if err != nil {
		log.Printf("❌ Failed to scan local network: %v", err)
		return
	}

	var newDevices []scan.Device
	var inventoryReport string

	for _, dev := range devices {
		scan.EnrichDevice(&dev)

		if !firstRunDone {
			inventoryReport += alert.FormatDeviceDetails(dev, "Wi-Fi", ssid) + "\n\n"
		}

		if !utils.IsKnownDevice(dev.MAC, known) {
			newDevices = append(newDevices, dev)

			if firstRunDone {
				if err := alert.SendDiscordAlert(config.Config.DiscordWebhook, dev, "Wi-Fi", ssid); err != nil {
					log.Printf("⚠️ Failed to alert for new device %s: %v", dev.MAC, err)
				}
			}
		}

		known.Devices[dev.MAC] = dev
	}

	if !firstRunDone {
		header := "🧾 **Initial Device Inventory (Wi-Fi)**\n\n"
		footer := fmt.Sprintf("🕒 Captured at: `%s`", time.Now().Format(time.RFC1123))
		fullMessage := header + inventoryReport + footer

		err := alert.SendMarkdownToDiscord(config.Config.DiscordWebhook, fullMessage)
		if err != nil {
			log.Printf("❌ Failed to send initial inventory report: %v", err)
		}
		firstRunDone = true
	}

	utils.UpdateCache(config.Config.CacheFile, known, []scan.Device{})

	log.Printf("✅ Scan complete — total: %d, new: %d", len(devices), len(newDevices))
}
