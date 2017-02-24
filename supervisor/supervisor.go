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

type supervisor struct {
	prog string
	argv []string

	cmd     *exec.Cmd
	startAt time.Time
	mu      sync.RWMutex

	signaled   bool
	signaledMu sync.RWMutex

	hupped   bool
	huppedMu sync.RWMutex
}

// Supervise starts a child mackerel-agent process and supervises it.
// 'c' can be nil and it's typically nil. When you pass signal channel to this
// method, the channel will be closed internally.
func Supervise(agentProg string, argv []string, c chan os.Signal) error {
	return (&supervisor{
		prog: agentProg,
		argv: argv,
	}).supervise(c)
}

func (sv *supervisor) setSignaled(signaled bool) {
	sv.signaledMu.Lock()
	defer sv.signaledMu.Unlock()
	sv.signaled = signaled
}

func (sv *supervisor) getSignaled() bool {
	sv.signaledMu.RLock()
	defer sv.signaledMu.RUnlock()
	return sv.signaled
}

func (sv *supervisor) setHupped(hupped bool) {
	sv.huppedMu.Lock()
	defer sv.huppedMu.Unlock()
	sv.hupped = hupped
}

func (sv *supervisor) getHupped() bool {
	sv.huppedMu.RLock()
	defer sv.huppedMu.RUnlock()
	return sv.hupped
}

func (sv *supervisor) getCmd() *exec.Cmd {
	sv.mu.RLock()
	defer sv.mu.RUnlock()
	return sv.cmd
}

func (sv *supervisor) getStartAt() time.Time {
	sv.mu.RLock()
	defer sv.mu.RUnlock()
	return sv.startAt
}

// If the child process dies within 30 seconds, it is regarded as launching failure
// and terminate the process without crash recovery
var spawnInterval = 30 * time.Second

func (sv *supervisor) launched() bool {
	return sv.getCmd().Process != nil && time.Now().After(sv.getStartAt().Add(spawnInterval))
}

func (sv *supervisor) buildCmd() *exec.Cmd {
	argv := append(sv.argv, "-child")
	cmd := exec.Command(sv.prog, argv...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

func (sv *supervisor) start() error {
	sv.mu.Lock()
	sv.setHupped(false)
	defer sv.mu.Unlock()
	sv.cmd = sv.buildCmd()
	sv.startAt = time.Now()
	return sv.cmd.Start()
}

func (sv *supervisor) stop(sig os.Signal) error {
	sv.setSignaled(true)
	return sv.getCmd().Process.Signal(sig)
}

func (sv *supervisor) configtest() error {
	argv := append([]string{"configtest"}, sv.argv...)
	cmd := exec.Command(sv.prog, argv...)
	buf := &bytes.Buffer{}
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("configtest failed: %s", buf.String())
	}
	return nil
}

func (sv *supervisor) reload() error {
	err := sv.configtest()
	if err != nil {
		return err
	}
	sv.setHupped(true)
	return sv.getCmd().Process.Signal(syscall.SIGTERM)
}

func (sv *supervisor) wait() (err error) {
	for {
		err = sv.cmd.Wait()
		if sv.getSignaled() || (!sv.getHupped() && !sv.launched()) {
			break
		}
		if err != nil {
			logger.Warningf("mackerel-agent abnormally finished with following error and try to restart it: %s", err.Error())
		}
		if err = sv.start(); err != nil {
			break
		}
	}
	return
}

func (sv *supervisor) handleSignal(ch <-chan os.Signal) {
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

func (sv *supervisor) supervise(c chan os.Signal) error {
	if err := sv.start(); err != nil {
		return err
	}
	if c == nil {
		c = make(chan os.Signal, 1)
	}
	defer close(c)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go sv.handleSignal(c)
	return sv.wait()
}
