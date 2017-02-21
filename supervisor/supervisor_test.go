// +build !windows

package supervisor

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/util"
)

const stubAgent = "testdata/stub-agent"

func init() {
	err := exec.Command("go", "build", "-o", stubAgent, "testdata/stub-agent.go").Run()
	if err != nil {
		panic(err)
	}
}

func TestSupervisor(t *testing.T) {
	sv := &Supervisor{
		Prog: stubAgent,
		Argv: []string{"dummy"},
	}
	ch := make(chan os.Signal, 1)
	done := make(chan error)
	go func() {
		done <- sv.Supervise(ch)
	}()
	time.Sleep(50 * time.Millisecond)
	pid := sv.getCmd().Process.Pid
	if !util.ExistsPid(pid) {
		t.Errorf("process doesn't exist")
	}
	time.Sleep(50 * time.Millisecond)
	ch <- os.Interrupt

	err := <-done
	if err != nil {
		t.Errorf("error should be nil but: %v", err)
	}
	if util.ExistsPid(pid) {
		t.Errorf("child process isn't terminated")
	}
}

func TestSupervisor_reload(t *testing.T) {
	sv := &Supervisor{
		Prog: stubAgent,
		Argv: []string{"dummy"},
	}
	ch := make(chan os.Signal, 1)
	done := make(chan error)
	go func() {
		done <- sv.Supervise(ch)
	}()
	time.Sleep(50 * time.Millisecond)
	oldPid := sv.getCmd().Process.Pid
	if !util.ExistsPid(oldPid) {
		t.Errorf("process doesn't exist")
	}
	ch <- syscall.SIGHUP
	time.Sleep(200 * time.Millisecond)
	newPid := sv.getCmd().Process.Pid
	if oldPid == newPid {
		t.Errorf("reload failed")
	}
	if util.ExistsPid(oldPid) {
		t.Errorf("old process isn't terminated")
	}
	if !util.ExistsPid(newPid) {
		t.Errorf("new process doesn't exist")
	}
	ch <- syscall.SIGTERM
	err := <-done
	if err != nil {
		t.Errorf("something went wrong")
	}
	if newPid != sv.getCmd().Process.Pid {
		t.Errorf("something went wrong")
	}
	if util.ExistsPid(newPid) {
		t.Errorf("child process isn't terminated")
	}
}

func TestSupervisor_reloadFail(t *testing.T) {
	sv := &Supervisor{
		Prog: stubAgent,
		Argv: []string{"failed"},
	}
	ch := make(chan os.Signal, 1)
	done := make(chan error)
	go func() {
		done <- sv.Supervise(ch)
	}()
	time.Sleep(50 * time.Millisecond)
	oldPid := sv.getCmd().Process.Pid
	if !util.ExistsPid(oldPid) {
		t.Errorf("process doesn't exist")
	}
	ch <- syscall.SIGHUP
	time.Sleep(time.Second)
	newPid := sv.getCmd().Process.Pid
	if oldPid != newPid {
		t.Errorf("reload should be failed, but unintentionally reloaded")
	}

	ch <- syscall.SIGTERM
	<-done
}

func TestSupervisor_launchFailed(t *testing.T) {
	sv := &Supervisor{
		Prog: stubAgent,
		Argv: []string{"launch failure"},
	}
	ch := make(chan os.Signal, 1)
	done := make(chan error)
	go func() {
		done <- sv.Supervise(ch)
	}()
	time.Sleep(50 * time.Millisecond)
	pid := sv.getCmd().Process.Pid
	if !util.ExistsPid(pid) {
		t.Errorf("process doesn't exist")
	}
	err := <-done
	if err == nil {
		t.Errorf("something went wrong")
	}
	if util.ExistsPid(sv.getCmd().Process.Pid) {
		t.Errorf("child process isn't terminated")
	}
}

func TestSupervisor_crashRecovery(t *testing.T) {
	origSpawnInterval := spawnInterval
	spawnInterval = 300 * time.Millisecond
	defer func() { spawnInterval = origSpawnInterval }()

	sv := &Supervisor{
		Prog: stubAgent,
		Argv: []string{"blah blah blah"},
	}
	ch := make(chan os.Signal, 1)
	done := make(chan error)
	go func() {
		done <- sv.Supervise(ch)
	}()
	time.Sleep(50 * time.Millisecond)
	oldPid := sv.getCmd().Process.Pid
	if !util.ExistsPid(oldPid) {
		t.Errorf("process doesn't exist")
	}
	time.Sleep(spawnInterval)

	// let it crash
	sv.getCmd().Process.Signal(syscall.SIGUSR1)

	time.Sleep(spawnInterval)
	newPid := sv.getCmd().Process.Pid
	if oldPid == newPid {
		t.Errorf("crash recovery failed")
	}
	if util.ExistsPid(oldPid) {
		t.Errorf("old process isn't terminated")
	}
	if !util.ExistsPid(newPid) {
		t.Errorf("new process doesn't exist")
	}
	ch <- syscall.SIGTERM
	<-done
}
