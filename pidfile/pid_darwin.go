package pidfile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func existsPid(pid int) (b bool) {
	cmd := exec.Command("/bin/ps", "-o", "command=", fmt.Sprint(pid))

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return false
	}

	out := stdout.String()
	if i := strings.IndexRune(out, ' '); i > 0 {
		out = out[:i]
	}
	return filepath.Base(out) == filepath.Base(os.Args[0])
}
