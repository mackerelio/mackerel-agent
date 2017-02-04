package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func handleFork(prog string, argv []string) error {
	var cmdBuilder = func() *exec.Cmd {
		cmd := exec.Command(prog, argv...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd
	}
	cmd := cmdBuilder()
	startAt := time.Now()
	_ = startAt
	cmd.Start()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		for sig := range c {
			if sig == syscall.SIGHUP {
				// reload agent
				cmd.Process.Signal(sig)
			} else {
				cmd.Process.Signal(sig)
			}
		}
	}()
	return cmd.Wait()
}
