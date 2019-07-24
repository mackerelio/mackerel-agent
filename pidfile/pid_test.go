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

func TestExistsPid_CheckProcessName(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "sleep", "10")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	pid := cmd.Process.Pid

	origArg0 := os.Args[0]
	defer func() { os.Args[0] = origArg0 }()

	os.Args[0] = "/bin/go-test"
	if existsPid(pid) {
		t.Errorf("existsPid should return false when os.Args[0] is %q and pid: %v", os.Args[0], pid)
	}

	os.Args[0] = "/bin/sleep" // note that only the base name is checked
	if !existsPid(pid) {
		t.Errorf("existsPid should return true when os.Args[0] is %q and pid: %v", os.Args[0], pid)
	}
}
