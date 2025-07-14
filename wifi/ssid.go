package wifi

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
)

func GetSSID(interfaceName string) string {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -I", "-I")
	case "linux":
		cmd = exec.Command("iwgetid", "-r")
	default:
		return "UnsupportedOS"
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "Unknown"
	}

	output := out.String()
	lines := strings.Split(output, "\n")

	if runtime.GOOS == "darwin" {
		for _, line := range lines {
			if strings.Contains(line, " SSID:") {
				return strings.TrimSpace(strings.Split(line, ":")[1])
			}
		}
	} else {
		return strings.TrimSpace(output)
	}

	return "Unknown"
}
