// +build linux freebsd netbsd

package pidfile

import (
	"fmt"
	"os"
)

func existsPid(pid int) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/", pid))
	return err == nil
}
