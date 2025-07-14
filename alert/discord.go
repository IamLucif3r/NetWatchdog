package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/iamlucif3r/netwatchdog/scan"
)

const maxDiscordLength = 2000

// Discord payload format
type discordPayload struct {
	Content string `json:"content"`
}

// 📢 Format a device's details as a Markdown block
func FormatDeviceDetails(dev scan.Device, sourceLabel, ssid string) string {
	if dev.Hostname == "" {
		dev.Hostname = "Unknown"
	}
	if dev.Manufacturer == "" {
		dev.Manufacturer = "Unknown"
	}
	if dev.OSGuess == "" {
		dev.OSGuess = "Unknown"
	}
	if ssid == "" || ssid == "Unknown" {
		ssid = "N/A"
	}

	return fmt.Sprintf(
		"📡 **Device Found**\n"+
			"🔗 **MAC:** `%s`\n"+
			"📍 **IP:** `%s`\n"+
			"💻 **Hostname:** `%s`\n"+
			"🏷️ **Vendor:** `%s`\n"+
			"🧠 **OS:** `%s`\n"+
			"📶 **Network:** `%s` on `%s`",
		dev.MAC,
		dev.IP,
		dev.Hostname,
		dev.Manufacturer,
		dev.OSGuess,
		ssid,
		sourceLabel,
	)
}

// 🚨 Alert for newly joined device
func SendDiscordAlert(webhookURL string, dev scan.Device, sourceLabel string, ssid string) error {
	now := time.Now().Format(time.RFC1123)

	// Handle missing values
	hostname := dev.Hostname
	if hostname == "" {
		hostname = "Unknown"
	}
	osGuess := dev.OSGuess
	if osGuess == "" {
		osGuess = "Unknown"
	}
	vendor := dev.Manufacturer
	if vendor == "" {
		vendor = "Unknown"
	}
	if ssid == "" || ssid == "Unknown" {
		ssid = "N/A"
	}

	// Markdown-formatted alert
	message := fmt.Sprintf(
		"🚨 **New Device Joined the Network!**\n\n"+
			"📍 **Source:** %s\n"+
			"📶 **SSID:** `%s`\n"+
			"📡 **IP:** `%s`\n"+
			"🔗 **MAC:** `%s`\n"+
			"💻 **Hostname:** `%s`\n"+
			"🏷️ **Vendor:** `%s`\n"+
			"🧠 **OS:** `%s`\n"+
			"🕒 **Seen at:** `%s`",
		sourceLabel,
		ssid,
		dev.IP,
		dev.MAC,
		hostname,
		vendor,
		osGuess,
		now,
	)

	return sendToDiscord(webhookURL, message)
}

// 📄 Send markdown message (like inventory dump)
func SendMarkdownToDiscord(webhookURL string, content string) error {
	return sendToDiscord(webhookURL, content)
}

// 🔧 Core sending logic with error debug
func sendToDiscord(webhookURL, content string) error {
	if len(content) > maxDiscordLength {
		content = content[:maxDiscordLength-5] + "…"
	}

	payload := discordPayload{Content: content}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord payload: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create Discord request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Discord: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Discord returned non-2xx status: %s\nResponse: %s", resp.Status, string(body))
	}

	return nil
}
