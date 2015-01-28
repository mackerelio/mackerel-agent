package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: replace [filename] [old string] [new string]")
	}
	file := os.Args[1]
	old := os.Args[2]
	new := os.Args[3]

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(file, []byte(strings.Replace(string(content), old, new, -1)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
