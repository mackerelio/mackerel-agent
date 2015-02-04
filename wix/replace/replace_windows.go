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
	oldStr := os.Args[2]
	newStr := os.Args[3]

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(file, []byte(strings.Replace(string(content), oldStr, newStr, -1)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
