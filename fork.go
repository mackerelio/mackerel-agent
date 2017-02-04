package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type cmdManager struct {
	prog    string
	argv    []string
	cmd     *exec.Cmd
	startAt time.Time
}

func (cm *cmdManager) buildCmd() *exec.Cmd {
	cmd := exec.Command(cm.prog, cm.argv...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

func (cm *cmdManager) start() error {
	cm.cmd = cm.buildCmd()
	cm.startAt = time.Now()
	return cm.cmd.Start()
}

func (cm *cmdManager) wait() error {
	return cm.cmd.Wait()
}

func handleFork(prog string, argv []string) error {
	cm := &cmdManager{
		prog: prog,
		argv: argv,
	}
	cm.start()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		for sig := range c {
			if sig == syscall.SIGHUP {
				// reload agent
				cm.cmd.Process.Signal(sig)
			} else {
				cm.cmd.Process.Signal(sig)
			}
		}
	}()
	return cm.wait()
}
