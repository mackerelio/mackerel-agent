package util

import (
	"bytes"
	"os"
	"os/exec"
)

func RunCommand(command string) (string, string, error) {
	var outBuffer, errBuffer bytes.Buffer

	wd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	cmd := exec.Command("cmd", "/c", "pushd "+wd+" & "+command)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err = cmd.Run()

	if err != nil {
		return "", "", err
	}

	return string(outBuffer.Bytes()), string(errBuffer.Bytes()), nil
}
