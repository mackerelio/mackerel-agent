package util

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func ExistsPid(pid int) bool {
	cmd := exec.Command("/usr/sbin/lsof", "-p", fmt.Sprint(pid))
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	return err == nil
}
