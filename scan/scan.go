package scan

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

type Device struct {
	IP           string `json:"ip"`
	MAC          string `json:"mac"`
	Hostname     string `json:"hostname,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	OSGuess      string `json:"os_guess,omitempty"`
}

func ScanLocal(interfaceName string) ([]Device, error) {
	cmd := exec.Command("arp", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run arp -a: %w", err)
	}

	var devices []Device
	scanner := bufio.NewScanner(&out)

	re := regexp.MustCompile(`\(([^)]+)\) at ([0-9a-f:]+)`)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, interfaceName) {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		ip := matches[1]
		mac := strings.ToUpper(matches[2])

		devices = append(devices, Device{
			IP:       ip,
			MAC:      mac,
			Hostname: extractHostnameFromArpLine(line),
		})
	}

	return devices, nil
}

func extractHostnameFromArpLine(line string) string {
	parts := strings.Fields(line)
	if len(parts) > 0 && parts[0] != "?" {
		return parts[0]
	}
	return ""
}

func getLocalIPv4(ifaceName string) (net.IP, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ip, _, _ := net.ParseCIDR(addr.String())
		if ip.To4() != nil {
			return ip, nil
		}
	}
	return nil, errors.New("no IPv4 address found on interface")
}

type tailscaleStatus struct {
	Peer map[string]struct {
		HostName     string   `json:"HostName"`
		TailscaleIPs []string `json:"TailscaleIPs"`
		User         string   `json:"User"`
		Online       bool     `json:"Online"`
	}
}

func ScanTailscale() ([]Device, error) {
	cmd := exec.Command("tailscale", "status", "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run tailscale command: %w", err)
	}

	var status tailscaleStatus
	if err := json.Unmarshal(output, &status); err != nil {
		return nil, fmt.Errorf("failed to parse tailscale JSON: %w", err)
	}

	var devices []Device
	for peerKey, peer := range status.Peer {
		if peer.Online && len(peer.TailscaleIPs) > 0 {
			dev := Device{
				IP:       peer.TailscaleIPs[0],
				MAC:      "ts-" + peerKey[:8],
				Hostname: peer.HostName,
			}
			devices = append(devices, dev)
		}
	}

	return devices, nil
}

func EnrichDevice(dev *Device) {
	dev.Manufacturer = lookupVendor(dev.MAC)

	osGuess := runOSFingerprint(dev.IP)
	if osGuess == "" {
		dev.OSGuess = "Unknown (nmap failed)"
	} else {
		dev.OSGuess = osGuess
	}
}
func runOSFingerprint(ip string) string {
	cmd := exec.Command("nmap", "-O", "-T4", ip)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return ""
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "OS details:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "OS details:"))
		}
		if strings.HasPrefix(line, "Running:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Running:"))
		}
	}

	return ""
}

func lookupVendor(mac string) string {
	prefix := strings.ToUpper(strings.ReplaceAll(mac[:8], ":", "-"))

	vendors := map[string]string{
		"E8-6D-AA": "Raspberry Pi Foundation",
		"F0-9F-C2": "Apple Inc.",
		"3C-5A-B4": "Samsung Electronics",
		"BC-92-6B": "Intel Corp.",
		"40-B0-FA": "Xiaomi",
		"00-1A-2B": "Cisco Systems",
		"00-1B-63": "Apple Inc.",
		"00-1C-B3": "Dell Inc.",
		"00-1D-7E": "Hewlett Packard",
		"00-1E-65": "Sony Mobile",
		"00-1F-3B": "LG Electronics",
		"00-21-6A": "Microsoft",
		"00-22-48": "Nokia",
		"00-23-69": "Huawei Technologies",
		"00-24-E8": "ASUSTek Computer",
		"00-25-9C": "TP-LINK Technologies",
		"00-26-BB": "Motorola Mobility",
		"00-27-0E": "Lenovo Mobile",
		"00-28-F8": "Amazon Technologies",
		"00-30-65": "Google Inc.",
		"00-50-56": "VMware, Inc.",
		"00-90-A9": "ASRock Incorporation",
		"00-0C-29": "VMware, Inc.",
		"00-16-3E": "Xensource, Inc.",
		"00-15-5D": "Microsoft Hyper-V",
		"00-0F-4B": "Nintendo Co., Ltd.",
		"00-13-02": "Toshiba",
		"00-17-88": "Netgear",
		"00-18-4D": "D-Link Corporation",
		"00-19-E0": "Hon Hai Precision (Foxconn)",
		"00-1A-11": "Samsung Electronics",
		"00-1B-77": "Hewlett Packard",
		"00-1C-23": "Cisco Systems",
		"00-1D-0F": "Sony Corporation",
		"00-1E-8C": "Apple Inc.",
		"00-1F-16": "Dell Inc.",
		"00-21-5C": "ASUSTek Computer",
		"00-22-41": "Nokia",
		"00-23-54": "Huawei Technologies",
		"00-24-21": "TP-LINK Technologies",
		"00-25-86": "Motorola Mobility",
		"00-26-5A": "Lenovo Mobile",
		"00-27-15": "Amazon Technologies",
		"00-28-38": "Google Inc.",
		"00-50-43": "Cisco Systems",
		"00-90-27": "ASRock Incorporation",
	}

	if name, ok := vendors[prefix]; ok {
		return name
	}
	return "Unknown Vendor"
}
