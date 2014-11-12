package util

import (
	"bytes"
	"os/exec"
)

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
