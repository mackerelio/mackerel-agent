package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type cmdManager struct {
	prog     string
	argv     []string
	cmd      *exec.Cmd
	startAt  time.Time
	signaled bool
	hupped   bool
}

var spawnInterval = 60 * time.Second

func (cm *cmdManager) launched() bool {
	return cm.cmd.Process != nil && time.Now().After(cm.startAt.Add(spawnInterval))
}

func (cm *cmdManager) buildCmd() *exec.Cmd {
	cmd := exec.Command(cm.prog, cm.argv...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

func (cm *cmdManager) start() error {
	cm.hupped = false
	cm.cmd = cm.buildCmd()
	cm.startAt = time.Now()
	return cm.cmd.Start()
}

func (cm *cmdManager) stop(sig os.Signal) error {
	cm.signaled = true
	return cm.cmd.Process.Signal(sig)
}

func (cm *cmdManager) configtest() error {
	argv := append([]string{"configtest"}, cm.argv...)
	cmd := exec.Command(cm.prog, argv...)
	return cmd.Run()
}

func (cm *cmdManager) reload() error {
	err := cm.configtest()
	if err != nil {
		return err
	}
	cm.hupped = true
	return cm.cmd.Process.Signal(syscall.SIGTERM)
}

func (cm *cmdManager) wait() (err error) {
	for {
		err = cm.cmd.Wait()
		if cm.signaled || (!cm.hupped && !cm.launched()) {
			break
		}
		cm.start()
	}
	return
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
				err := cm.reload()
				if err != nil {
					logger.Warningf("failed to reload: %s", err.Error())
				}
			} else {
				cm.stop(sig)
			}
		}
	}()
	return cm.wait()
}
