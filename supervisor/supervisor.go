package supervisor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("supervisor")

// Supervisor supervise the mackerel-agent
type Supervisor struct {
	Prog string
	Argv []string

	cmd     *exec.Cmd
	startAt time.Time
	mu      sync.RWMutex

	signaled   bool
	signaledMu sync.RWMutex

	hupped   bool
	huppedMu sync.RWMutex
}

func (sv *Supervisor) setSignaled(signaled bool) {
	sv.signaledMu.Lock()
	defer sv.signaledMu.Unlock()
	sv.signaled = signaled
}

func (sv *Supervisor) getSignaled() bool {
	sv.signaledMu.RLock()
	defer sv.signaledMu.RUnlock()
	return sv.signaled
}

func (sv *Supervisor) setHupped(hupped bool) {
	sv.huppedMu.Lock()
	defer sv.huppedMu.Unlock()
	sv.hupped = hupped
}

func (sv *Supervisor) getHupped() bool {
	sv.huppedMu.RLock()
	defer sv.huppedMu.RUnlock()
	return sv.hupped
}

func (sv *Supervisor) getCmd() *exec.Cmd {
	sv.mu.RLock()
	defer sv.mu.RUnlock()
	return sv.cmd
}

func (sv *Supervisor) getStartAt() time.Time {
	sv.mu.RLock()
	defer sv.mu.RUnlock()
	return sv.startAt
}

// If the child process dies within 30 seconds, it is regarded as launching failure
// and terminate the process without crash recovery
var spawnInterval = 30 * time.Second

func (sv *Supervisor) launched() bool {
	return sv.getCmd().Process != nil && time.Now().After(sv.getStartAt().Add(spawnInterval))
}

func (sv *Supervisor) buildCmd() *exec.Cmd {
	argv := append(sv.Argv, "-child")
	cmd := exec.Command(sv.Prog, argv...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

func (sv *Supervisor) start() error {
	sv.mu.Lock()
	sv.setHupped(false)
	defer sv.mu.Unlock()
	sv.cmd = sv.buildCmd()
	sv.startAt = time.Now()
	return sv.cmd.Start()
}

func (sv *Supervisor) stop(sig os.Signal) error {
	sv.setSignaled(true)
	return sv.getCmd().Process.Signal(sig)
}

func (sv *Supervisor) configtest() error {
	argv := append([]string{"configtest"}, sv.Argv...)
	cmd := exec.Command(sv.Prog, argv...)
	buf := &bytes.Buffer{}
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("configtest failed: %s", buf.String())
	}
	return nil
}

func (sv *Supervisor) reload() error {
	err := sv.configtest()
	if err != nil {
		return err
	}
	sv.setHupped(true)
	return sv.getCmd().Process.Signal(syscall.SIGTERM)
}

func (sv *Supervisor) wait() (err error) {
	for {
		err = sv.cmd.Wait()
		if sv.getSignaled() || (!sv.getHupped() && !sv.launched()) {
			break
		}
		if err != nil {
			logger.Warningf("mackerel-agent abnormally finished with following error and try to restart it: %s", err.Error())
		}
		err = sv.start()
		if err != nil {
			break
		}
	}
	return
}

func (sv *Supervisor) handleSignal(ch <-chan os.Signal) {
	for sig := range ch {
		if sig == syscall.SIGHUP {
			logger.Infof("receiving HUP, spawning a new mackerel-agent")
			err := sv.reload()
			if err != nil {
				logger.Warningf("failed to reload: %s", err.Error())
			}
		} else {
			sv.stop(sig)
		}
	}
}

// Supervise the mackerel-agent
func (sv *Supervisor) Supervise(c chan os.Signal) error {
	err := sv.start()
	if err != nil {
		return err
	}
	if c == nil {
		c = make(chan os.Signal, 1)
	}
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go sv.handleSignal(c)
	return sv.wait()
}
