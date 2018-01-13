package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	os.Exit(run())
}

func run() int {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	select {
	case <-time.After(10 * time.Second):
		return 0
	case <-c:
		fmt.Println("signal received")
		return 23
	}
}
