package main

import (
	"log"
	"os"

	"github.com/motemen/go-cli/gen"
)

func main() {
	out, err := os.Create("commands.go")
	if err != nil {
		log.Fatal(err)
	}

	err = gen.Generate(out, "main.go", nil)
	if err != nil {
		log.Fatal(err)
	}
}
