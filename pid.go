// +build linux freebsd

package main

import (
	"fmt"
	"os"
)

func existsPid(int pid) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/", pid))
	return err == nil
}
