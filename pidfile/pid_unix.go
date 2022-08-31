//go:build linux || freebsd || netbsd
// +build linux freebsd netbsd

package pidfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func existsPid(pid int) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/", pid))
	return err == nil
}

func getCmdName(pid int) string {
	cnt, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return ""
	}

	out := string(cnt)
	if i := strings.IndexRune(out, '\x00'); i > 0 {
		out = out[:i]
	}
	return filepath.Base(out)
}
