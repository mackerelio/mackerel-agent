package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/mackerelio/mackerel-agent/logging"
)

var utilLogger = logging.GetLogger("util")

// RunCommand XXX
func RunCommand(command, user string) (string, string, int, error) {
	var outBuffer, errBuffer bytes.Buffer

	wd, err := os.Getwd()
	if err != nil {
		return "", "", -1, err
	}
	cmd := exec.Command("cmd", "/c", "pushd "+wd+" & "+command)
	if user != "" {
		utilLogger.Warningf("RunCommand ignore option: user = %q", user)
	}
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err = cmd.Run()

	stdout := outBuffer.String()
	stderr := errBuffer.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode := waitStatus.ExitStatus()
				return stdout, stderr, exitCode, nil
			}
		}
		return stdout, stderr, -1, err
	}

	return stdout, stderr, 0, nil
}

// GetWmic XXX
func GetWmic(target string, query string) (string, error) {
	cpuGet, err := exec.Command("wmic", target, "get", query).Output()
	if err != nil {
		return "", err
	}

	percentages := string(cpuGet)

	lines := strings.Split(percentages, "\r\r\n")

	if len(lines) <= 2 {
		return "", fmt.Errorf("wmic result malformed: [%q]", lines)
	}

	return strings.Trim(lines[1], " "), nil
}

// GetWmicToFloat XXX
func GetWmicToFloat(target string, query string) (float64, error) {
	value, err := GetWmic(target, query)
	if err != nil {
		return 0, err
	}

	ret, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

var metricKeySanitizerReg = regexp.MustCompile(`[^A-Za-z0-9_-]`)

// SanitizeMetricKey sanitizes key string
func SanitizeMetricKey(rawKey string) string {
	return metricKeySanitizerReg.ReplaceAllString(rawKey, "_")
}
