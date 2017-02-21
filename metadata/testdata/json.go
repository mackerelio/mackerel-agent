package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	n := flag.Int("exit-code", 0, "exit code")
	m := flag.String("metadata", "", "metadata")
	flag.Parse()

	fmt.Print(*m)
	os.Exit(*n)
}
