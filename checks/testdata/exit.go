package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	n := flag.Int("code", 0, "exit code")
	m := flag.String("message", "", "message")
	flag.Parse()

	fmt.Println(*m)
	os.Exit(*n)
}
