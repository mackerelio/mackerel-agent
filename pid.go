// +build linux freebsd

package main

import (
	"fmt"
	"os"
)

func existsPid(pid int) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/", pid))
	return err == nil
}
