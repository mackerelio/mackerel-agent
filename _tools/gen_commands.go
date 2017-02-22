package main

import (
	"log"
	"os"

	"github.com/motemen/go-cli/gen"
)

func main() {
	out, err := os.Create("commands_gen.go")
	if err != nil {
		log.Fatal(err)
	}

	err = gen.Generate(out, "commands.go", nil)
	if err != nil {
		log.Fatal(err)
	}
}
