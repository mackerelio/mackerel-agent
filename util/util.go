// +build linux darwin
package util

import (
	"bytes"
	"os/exec"
)

// RunCommand runs command (in one string) and returns stdout, stderr strings.
func RunCommand(command string) (string, string, error) {
	var outBuffer, errBuffer bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()

	if err != nil {
		return "", "", err
	}

	return string(outBuffer.Bytes()), string(errBuffer.Bytes()), nil
}
