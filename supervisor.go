package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type supervisor struct {
	prog string
	argv []string

	cmd      *exec.Cmd
	startAt  time.Time
	signaled bool
	hupped   bool
}

var spawnInterval = 30 * time.Second

func (sv *supervisor) launched() bool {
	return sv.cmd.Process != nil && time.Now().After(sv.startAt.Add(spawnInterval))
}

func (sv *supervisor) buildCmd() *exec.Cmd {
	argv := append(sv.argv, "-child")
	cmd := exec.Command(sv.prog, argv...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

func (sv *supervisor) start() error {
	sv.hupped = false
	sv.cmd = sv.buildCmd()
	sv.startAt = time.Now()
	return sv.cmd.Start()
}

func (sv *supervisor) stop(sig os.Signal) error {
	sv.signaled = true
	return sv.cmd.Process.Signal(sig)
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
	sv.hupped = true
	return sv.cmd.Process.Signal(syscall.SIGTERM)
}

func (sv *supervisor) wait() (err error) {
	for {
		err = sv.cmd.Wait()
		if sv.signaled || (!sv.hupped && !sv.launched()) {
			break
		}
		sv.start()
	}
	return
}

func (sv *supervisor) supervise() error {
	sv.start()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		for sig := range c {
			if sig == syscall.SIGHUP {
				err := sv.reload()
				if err != nil {
					logger.Warningf("failed to reload: %s", err.Error())
				}
			} else {
				sv.stop(sig)
			}
		}
	}()
	return sv.wait()
}
