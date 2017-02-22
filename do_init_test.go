package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoInitialize(t *testing.T) {
	{
		err := doInitialize(&flag.FlagSet{}, []string{})
		if err == nil || err.Error() != "-apikey option is required" {
			t.Errorf("-apikey option is required")
		}
	}

	root, err := ioutil.TempDir("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temporary dir for test")
	}
	defer os.RemoveAll(root)

	{
		conffile := filepath.Join(root, "mackerel-agent.conf")
		argv := []string{"-conf", conffile, "-apikey", "hoge"}
		err := doInit(&flag.FlagSet{}, argv)
		if err != nil {
			t.Errorf("err should be nil but: %s", err)
		}
		content, _ := ioutil.ReadFile(conffile)
		if string(content) != `apikey = "hoge"`+"\n" {
			t.Errorf("somthing went wrong: %s", string(content))
		}
		err = doInitialize(&flag.FlagSet{}, argv)
		if _, ok := err.(apikeyAlreadySetError); !ok {
			t.Errorf("err should be `apikeyAlreadySetError` but :%s", err)
		}
	}

	{
		f := filepath.Join(root, "mackerel-agent2.conf")
		confFile, err := os.Create(f)
		if err != nil {
			t.Fatalf("Could not create temporary file for test")
		}
		confFile.Close()
		argv := []string{"-conf", f, "-apikey", "hoge"}
		err = doInit(&flag.FlagSet{}, argv)
		if err != nil {
			t.Errorf("err should be nil but: %s", err)
		}
	}

	{
		f := filepath.Join(root, "mackerel-agent-invalid.conf")
		conffile, err := os.Create(f)
		conffile.WriteString(`dummy = "`)
		conffile.Sync()
		conffile.Close()
		err = doInit(&flag.FlagSet{}, []string{"-conf", f, "-apikey=hoge"})

		if err == nil || !strings.HasPrefix(err.Error(), "failed to load the config") {
			t.Errorf("should return load config error but: %s", err)
		}
	}
}
