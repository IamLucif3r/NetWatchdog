package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type AppConfig struct {
	Interface       string `json:"interface"`
	DiscordWebhook  string `json:"discord_webhook"`
	EnableTailscale bool   `json:"enable_tailscale"`
	CacheFile       string `json:"cache_file"`
}

var Config AppConfig

func LoadConfig(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("❌ Failed to parse config: %v", err)
	}

	fmt.Println("⚙️ Config loaded successfully.")
}
