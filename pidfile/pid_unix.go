// +build linux freebsd netbsd

package pidfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func existsPid(pid int) bool {
	cnt, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return false
	}

	out := string(cnt)
	if i := strings.IndexRune(out, '\x00'); i > 0 {
		out = out[:i]
	}
	return filepath.Base(out) == filepath.Base(os.Args[0])
}
