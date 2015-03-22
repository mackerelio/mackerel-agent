// +build linux darwin freebsd

package util

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/Songmu/timeout"
)

var TimeoutDuration = 30 * time.Second
var TimeoutKillAfter = 10 * time.Second

// RunCommand runs command (in one string) and returns stdout, stderr strings.
func RunCommand(command string) (string, string, error) {
	tio := &timeout.Timeout{
		Cmd:       exec.Command("/bin/sh", "-c", command),
		Duration:  TimeoutDuration,
		KillAfter: TimeoutKillAfter,
	}
	exitStatus, stdout, stderr, err := tio.Run()

	if err == nil && exitStatus.IsTimedOut() {
		err = fmt.Errorf("command timed out")
	}
	return stdout, stderr, err
}
