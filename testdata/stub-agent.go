package main

import (
	"fmt"
	"os"
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

	for {
		time.Sleep(time.Second)
	}
}
