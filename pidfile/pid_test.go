// +build linux freebsd darwin netbsd

package pidfile

import (
	"context"
	"math"
	"os"
	"os/exec"
	"testing"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "sleep", "10")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	pid := cmd.Process.Pid

	expected := "sleep"
	if got := GetCmdName(pid); got != expected {
		t.Errorf("GetCmdName should return %q bot got: %q", expected, got)
	}
}
