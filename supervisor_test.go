// +build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

const stubAgent = "testdata/stub-agent"

func init() {
	err := exec.Command("go", "build", "-o", stubAgent, "testdata/stub-agent.go").Run()
	if err != nil {
		panic(err)
	}
}

func TestSupervisor(t *testing.T) {
	sv := &supervisor{
		prog: stubAgent,
		argv: []string{"dummy"},
	}
	sv.start()
	ch := make(chan os.Signal)
	go sv.handleSignal(ch)
	done := make(chan error)
	go func() {
		done <- sv.wait()
	}()
	pid := sv.cmd.Process.Pid
	if !existsPid(pid) {
		t.Errorf("process doesn't exist")
	}
	ch <- os.Interrupt

	err := <-done
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}
	if existsPid(pid) {
		t.Errorf("child process isn't terminated")
	}
}

func TestSupervisor_reload(t *testing.T) {
	sv := &supervisor{
		prog: stubAgent,
		argv: []string{"dummy"},
	}
	sv.start()
	ch := make(chan os.Signal)
	go sv.handleSignal(ch)
	done := make(chan error)
	go func() {
		done <- sv.wait()
	}()
	oldPid := sv.cmd.Process.Pid
	if !existsPid(oldPid) {
		t.Errorf("process doesn't exist")
	}
	ch <- syscall.SIGHUP
	time.Sleep(200 * time.Millisecond)
	newPid := sv.cmd.Process.Pid
	if oldPid == newPid {
		t.Errorf("reload failed")
	}
	if existsPid(oldPid) {
		t.Errorf("old process isn't terminated")
	}
	if !existsPid(newPid) {
		t.Errorf("new process doesn't exist")
	}
	ch <- syscall.SIGTERM
	err := <-done
	if err != nil {
		t.Errorf("something went wrong")
	}
	if newPid != sv.cmd.Process.Pid {
		t.Errorf("something went wrong")
	}
	if existsPid(newPid) {
		t.Errorf("child process isn't terminated")
	}
}

func TestSupervisor_reloadFail(t *testing.T) {
	sv := &supervisor{
		prog: stubAgent,
		argv: []string{"failed"},
	}
	sv.start()
	ch := make(chan os.Signal)
	go sv.handleSignal(ch)
	done := make(chan error)
	go func() {
		done <- sv.wait()
	}()
	oldPid := sv.cmd.Process.Pid
	if !existsPid(oldPid) {
		t.Errorf("process doesn't exist")
	}
	ch <- syscall.SIGHUP
	time.Sleep(time.Second)
	newPid := sv.cmd.Process.Pid
	if oldPid != newPid {
		t.Errorf("reload should be failed, but unintentionally reloaded")
	}

	ch <- syscall.SIGTERM
	<-done
}

func TestSupervisor_launchFailed(t *testing.T) {
	sv := &supervisor{
		prog: stubAgent,
		argv: []string{"launch failure"},
	}
	sv.start()
	ch := make(chan os.Signal)
	go sv.handleSignal(ch)
	done := make(chan error)
	go func() {
		done <- sv.wait()
	}()
	pid := sv.cmd.Process.Pid
	if !existsPid(pid) {
		t.Errorf("process doesn't exist")
	}
	err := <-done
	if err == nil {
		t.Errorf("something went wrong")
	}
	if existsPid(sv.cmd.Process.Pid) {
		t.Errorf("child process isn't terminated")
	}
}

func TestSupervisor_crashRecovery(t *testing.T) {
	origSpawnInterval := spawnInterval
	spawnInterval = 300 * time.Millisecond
	defer func() { spawnInterval = origSpawnInterval }()

	sv := &supervisor{
		prog: stubAgent,
		argv: []string{"blah blah blah"},
	}
	sv.start()
	ch := make(chan os.Signal)
	go sv.handleSignal(ch)
	done := make(chan error)
	go func() {
		done <- sv.wait()
	}()
	oldPid := sv.cmd.Process.Pid
	if !existsPid(oldPid) {
		t.Errorf("process doesn't exist")
	}
	time.Sleep(spawnInterval)

	// let it crash
	sv.cmd.Process.Signal(syscall.SIGUSR1)

	time.Sleep(spawnInterval)
	newPid := sv.cmd.Process.Pid
	if oldPid == newPid {
		t.Errorf("crash recovery failed")
	}
	if existsPid(oldPid) {
		t.Errorf("old process isn't terminated")
	}
	if !existsPid(newPid) {
		t.Errorf("new process doesn't exist")
	}
	ch <- syscall.SIGTERM
	<-done
}
