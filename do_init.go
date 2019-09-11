package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mackerelio/mackerel-agent/config"
)

func doInitialize(fs *flag.FlagSet, argv []string) error {
	var (
		conffile = fs.String("conf", config.DefaultConfig.Conffile, "Config file path")
		apikey   = fs.String("apikey", "", "API key from mackerel.io web site (Required)")
	)
	fs.Parse(argv)

	if *apikey == "" {
		// Setting apikey via environment variable should be supported or not?
		return fmt.Errorf("-apikey option is required")
	}

	// make the config file directory if not exist
	root := filepath.Dir(*conffile)
	if _, err := os.Stat(root); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err := os.MkdirAll(root, 0755)
		if err != nil {
			return err
		}
	}

	_, err := os.Stat(*conffile)
	confExists := err == nil
	if confExists {
		conf, err := config.LoadConfig(*conffile)
		if err != nil {
			return fmt.Errorf("failed to load the config file: %s", err)
		}
		if conf.Apikey != "" {
			return apikeyAlreadySetError(*conffile)
		}
	}
	contents := []byte(fmt.Sprintf("apikey = %q\n", *apikey))
	if confExists {
		cBytes, err := ioutil.ReadFile(*conffile)
		if err != nil {
			return err
		}
		contents = append(contents, cBytes...)
	}
	return ioutil.WriteFile(*conffile, contents, 0644)
}

type apikeyAlreadySetError string

func (a apikeyAlreadySetError) Error() string {
	return fmt.Sprintf("apikey already set in %q. Skip initializing", string(a))
}
