package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

	if isExists(outFile) {
		return
	}
	if dir := os.Getenv("MACKEREL_CONFIG_FALLBACK"); dir != "" {
		file := filepath.Join(dir, "mackerel-agent.conf")
		if isExists(file) {
			return
		}
	}

	content, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(outFile, []byte(strings.Replace(string(content), oldStr, newStr, -1)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func isExists(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return true // file is existed; but can't read it.
	}
	f.Close()
	return true
}
