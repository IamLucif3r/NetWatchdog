# 🕵️‍♂️ NetWatchdog

NetWatchdog is a lightweight network watchdog tool written in Go that scans your local Wi-Fi network (and optional Tailscale network) to detect newly connected devices and alert you via Discord — just like a paranoid hacker monitoring their home base 👀.

---

## 🚀 Features

- 🧠 Smart detection of new devices on your network
- 🔍 OS detection using `nmap`
- 🌐 MAC vendor lookup
- 📶 SSID identification on macOS/Linux
- 🛎️ Real-time Discord alerting
- 📜 Initial full inventory dump on startup
- 🧠 Works with Tailscale (optional)
- 💻 Minimal dependencies, no server required

---

## 🛠 Requirements

- Go 1.20+
- macOS or Linux
- `nmap` installed (`brew install nmap` or `apt install nmap`)
- A Discord Webhook URL
- (macOS only) `airport` command (comes pre-installed)

---

## 📦 Installation

```bash
git clone https://github.com/yourusername/netwatchdog
cd netwatchdog
go build -o netwatchdog
```

## ⚙️ Configuration

Create a config file:

```bash
# ~/.netwatchdog/config.yaml

interface: en1
cache_file: /Users/yourname/.netwatchdog/cache.json
discord_webhook: https://discord.com/api/webhooks/xxxx/yyyy
enable_tailscale: false

```
**Note** You can detect your Wi-Fi interface using ifconfig or networksetup -listallhardwareports.

## ▶️ Running the Tool

```bash
sudo ./netwatchdog
```
You’ll see a full inventory alert on first run and Discord notifications for any new device joining the network.

## 🔁 Run at Startup (macOS LaunchDaemon)
Move binary to global path:

```bash
sudo cp ./netwatchdog /usr/local/bin/
sudo chmod +x /usr/local/bin/netwatchdog
```
Create launch daemon file:

```bash
sudo nano /Library/LaunchDaemons/com.netwatchdog.plist
```
Paste this:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
 "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.netwatchdog</string>

  <key>ProgramArguments</key>
  <array>
    <string>/usr/local/bin/netwatchdog</string>
  </array>

  <key>RunAtLoad</key>
  <true/>

  <key>KeepAlive</key>
  <true/>

  <key>StandardOutPath</key>
  <string>/var/log/netwatchdog.log</string>

  <key>StandardErrorPath</key>
  <string>/var/log/netwatchdog.err</string>
</dict>
</plist>
```

Set permissions and load it:

```bash
sudo chown root:wheel /Library/LaunchDaemons/com.netwatchdog.plist
sudo chmod 644 /Library/LaunchDaemons/com.netwatchdog.plist

sudo launchctl load /Library/LaunchDaemons/com.netwatchdog.plist
```

Check logs:
```bash
sudo tail -f /var/log/netwatchdog.log
```

## 🧪 Example Output (Discord)

```markdown
🚨 **New Device Joined the Network!**

📍 **Source:** Wi-Fi  
📶 **SSID:** `iPhone_5G`  
📡 **IP:** `192.168.1.48`  
🔗 **MAC:** `EA:38:41:BC:AC:2A`  
💻 **Hostname:** `iPhone`  
🏷️ **Vendor:** `Raspberry Pi Foundation`  
🧠 **OS:** `Linux 5.10 (Raspbian)`  
🕒 **Seen at:** `Mon, 14 Jul 2025 22:18:31 IST`

## 🔐 Notes

- Requires root to run nmap -O (OS detection)
- Tailscale support is optional (scans connected peers)

## 🤘 Author

Built by @[iamlucif3r](github.com/iamlucif3r) — Platform Security Engineer, Automation Enthusiast.