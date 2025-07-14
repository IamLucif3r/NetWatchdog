# 🛡️ NetWatchdog
Notifies you when a new device joins your local or Tailscale network.

## 📝 Overview

**NetWatchdog** is a network monitoring tool written in Go that detects and notifies you when new devices join your local network or your Tailscale mesh. It helps you maintain awareness of your network environment and enhances security by alerting you to unexpected connections.

## ✨ Features

- Monitors local and Tailscale networks for new device connections
- Sends real-time notifications when a new device is detected
- Lightweight and easy to deploy
- Configurable notification channels (email, Slack, etc.)
- Cross-platform support

## ⚙️ Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/NetWatchdog.git
    cd NetWatchdog
    ```
2. Build the project:
    ```bash
    go build -o netwatchdog
    ```
3. Configure your notification settings in `config.yaml`.

## Usage

Run NetWatchdog with:
```bash
sudo ./netwatchdog
```

You can customize monitoring intervals and notification preferences in the configuration file.

## Configuration

Here is an example of config.json

```json

{
  "interface": "en1",
  "discord_webhook": "https://discord.com/api/webhooks/1234...",
  "enable_tailscale": false,
  "cache_file": "data/known.json"
}

```

## Contributing

Contributions are welcome! Please open issues or submit pull requests for new features, bug fixes, or improvements.


## Contact

For questions or support, please open an issue on GitHub.