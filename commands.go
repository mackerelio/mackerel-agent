//go:generate go run _tools/gen_commands.go

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/Songmu/prompter"
	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/version"
)

/* +main - mackerel-agent

	mackerel-agent [options]

main process of mackerel-agent
*/
func doMain(fs *flag.FlagSet, argv []string) error {
	conf, err := resolveConfig(fs, argv)
	if err != nil {
		return fmt.Errorf("failed to load config: %s", err)
	}
	return start(conf, make(chan struct{}))
}

/* +command init - initialize mackerel-agent.conf with apikey

	init -apikey=xxxxxxxxxxx [-conf=mackerel-agent.conf]

initialize mackerel-agent.conf with api key. Set the apikey to conf file.

- The conf file doesn't exist:
    create new file and set the apikey
- The conf file exists and apikey is unset:
    set the apikey
- The conf file exists and apikey already set:
    skip initializing. Don't overwrite apikey and exit normally.
- The conf file exists, but the contents of it is invalid toml:
    exit with error.
*/
func doInit(fs *flag.FlagSet, argv []string) error {
	var (
		conffile = fs.String("conf", config.DefaultConfig.Conffile, "Config file path")
		apikey = fs.String("apikey", "", "API key from mackerel.io web site")
	)
	fs.Parse(argv)

	if *apikey == "" {
		// Setting apikey via environment variable should be supported or not?
		return fmt.Errorf("-apikey option is required")
	}
	_, err := os.Stat(*conffile)
	confExists := err == nil
	if confExists {
		conf, err := config.LoadConfig(*conffile)
		if err != nil {
			return fmt.Errorf("Failed to load the config file: %s", err)
		}
		if conf.Apikey != "" {
			logger.Infof("apikey already set. skip initializing")
			return nil
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

/* +command version - display version of mackerel-agent

	version

display the version of mackerel-agent
*/
func doVersion(_ *flag.FlagSet, _ []string) error {
	fmt.Printf("mackerel-agent version %s (rev %s) [%s %s %s] \n",
		version.VERSION, version.GITCOMMIT, runtime.GOOS, runtime.GOARCH, runtime.Version())
	return nil
}

/* +command configtest - configtest

	configtest

do configtest
*/
func doConfigtest(fs *flag.FlagSet, argv []string) error {
	conf, err := resolveConfig(fs, argv)
	if err != nil {
		return fmt.Errorf("failed to test config: %s", err)
	}
	fmt.Fprintf(os.Stderr, "%s Syntax OK\n", conf.Conffile)
	return nil
}

/* +command retire - retire the host

	retire [-force]

retire the host
*/
func doRetire(fs *flag.FlagSet, argv []string) error {
	conf, force, err := resolveConfigForRetire(fs, argv)
	if err != nil {
		return fmt.Errorf("failed to load config: %s", err)
	}

	hostID, err := conf.LoadHostID()
	if err != nil {
		return fmt.Errorf("HostID file is not found")
	}

	api, err := mackerel.NewAPI(conf.Apibase, conf.Apikey, conf.Verbose)
	if err != nil {
		return fmt.Errorf("faild to create api client: %s", err)
	}

	if !force && !prompter.YN(fmt.Sprintf("retire this host? (hostID: %s)", hostID), false) {
		return fmt.Errorf("Retirement is canceled.")
	}

	err = api.RetireHost(hostID)
	if err != nil {
		return fmt.Errorf("faild to retire the host: %s", err)
	}
	logger.Infof("This host (hostID: %s) has been retired.", hostID)
	// just to try to remove hostID file.
	err = conf.DeleteSavedHostID()
	if err != nil {
		logger.Warningf("Failed to remove HostID file: %s", err)
	}
	return nil
}

/* +command once - output onetime

	once

output metrics and meta data of the host one time.
These data are only displayed and not posted to Mackerel.
*/
func doOnce(fs *flag.FlagSet, argv []string) error {
	conf, err := resolveConfig(fs, argv)
	if err != nil {
		logger.Warningf("failed to load config (but `once` must not required conf): %s", err)
		conf = &config.Config{}
	}
	command.RunOnce(conf)
	return nil
}
