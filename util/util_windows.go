package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// RunCommand XXX
func RunCommand(command string) (string, string, error) {
	var outBuffer, errBuffer bytes.Buffer

	cmd := exec.Command("cmd", "/c", command)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()

	if err != nil {
		return "", "", err
	}

	return string(outBuffer.Bytes()), string(errBuffer.Bytes()), nil
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
