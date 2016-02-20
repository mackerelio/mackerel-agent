//go:generate go run _tools/gen_commands.go

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/Songmu/prompter"
	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
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
	if conf.Verbose {
		logging.SetLogLevel(logging.DEBUG)
	}
	return start(conf, make(chan struct{}))
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
