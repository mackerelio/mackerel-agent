package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "configtest" {
		if len(os.Args) > 2 && os.Args[2] == "failed" {
			fmt.Fprintln(os.Stderr, "[stub] configtest failed")
			os.Exit(1)
		} else {
			fmt.Fprintln(os.Stderr, "[stub] configtest succeeded")
			os.Exit(0)
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "launch failure" {
		time.Sleep(300 * time.Millisecond)
		os.Exit(1)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGUSR1)
	go func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGHUP:
				// nop
			case syscall.SIGUSR1:
				panic("[stub] panicked by SIGUSR1")
			default:
				os.Exit(0)
			}
		}
	}()

	for {
		time.Sleep(time.Second)
	}
}
