//go:build linux || freebsd || darwin || netbsd

package pidfile

import (
	"math"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestExistsPid(t *testing.T) {
	if !ExistsPid(os.Getpid()) {
		t.Errorf("something went wrong")
	}
	if ExistsPid(math.MaxInt32) {
		t.Errorf("something went wrong")
	}
}

func TestGetCmdName(t *testing.T) {
	ctx := t.Context()
	cmd := exec.CommandContext(ctx, "sleep", "10")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	// hack: wait for launching process
	time.Sleep(time.Millisecond)
	pid := cmd.Process.Pid

	expected := "sleep"
	if got := GetCmdName(pid); got != expected {
		t.Errorf("GetCmdName should return %q but got: %q", expected, got)
	}
}
