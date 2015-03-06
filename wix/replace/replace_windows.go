package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 5 {
		log.Fatal("Usage: replace [in file] [out file] [old string] [new string]")
	}
	inFile := os.Args[1]
	outFile := os.Args[2]
	oldStr := os.Args[3]
	newStr := os.Args[4]

	content, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(outFile)
	outFileIsExists := err == nil
	if !(outFileIsExists) {
		err = ioutil.WriteFile(outFile, []byte(strings.Replace(string(content), oldStr, newStr, -1)), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
