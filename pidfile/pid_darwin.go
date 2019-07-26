package pidfile

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

func existsPid(pid int) bool {
	cmd := exec.Command("/usr/sbin/lsof", "-p", fmt.Sprint(pid))
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	return err == nil
}

func getCmdName(pid int) string {
	cmd := exec.Command("/bin/ps", "-o", "command=", fmt.Sprint(pid))

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if err != nil {
		return ""
	}

	out := stdout.String()
	if i := strings.IndexRune(out, ' '); i > 0 {
		out = out[:i]
	}
	return filepath.Base(out)
}
