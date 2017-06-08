//go:generate go run _tools/gen_commands.go

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/Songmu/prompter"
	"github.com/Songmu/retry"
	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/pidfile"
	"github.com/mackerelio/mackerel-agent/supervisor"
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

Initialize mackerel-agent.conf with api key.

- The conf file doesn't exist:
    create new file and set the apikey.
- The conf file exists and apikey is unset:
    set the apikey.
- The conf file exists and apikey already set:
    skip initializing. Don't overwrite apikey and exit normally.
- The conf file exists, but the contents of it is invalid toml:
    exit with error.
*/
func doInit(fs *flag.FlagSet, argv []string) error {
	err := doInitialize(fs, argv)
	if _, ok := err.(apikeyAlreadySetError); ok {
		logger.Infof("%s", err)
		return nil
	}
	return err
}

/* +command supervise - supervisor mode

	supervise -conf mackerel-agent.conf ...

run as supervisor mode enabling configuration reloading and crash recovery
*/
func doSupervise(fs *flag.FlagSet, argv []string) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("supervise mode is not supported on windows")
	}
	copiedArgv := make([]string, len(argv))
	copy(copiedArgv, argv)
	conf, err := resolveConfig(fs, argv)
	if err != nil {
		return err
	}
	setLogLevel(conf.Silent, conf.Verbose)
	err = pidfile.Create(conf.Pidfile)
	if err != nil {
		return err
	}
	defer pidfile.Remove(conf.Pidfile)

	return supervisor.Supervise(os.Args[0], copiedArgv, nil)
}

/* +command version - display version of mackerel-agent

	version

display the version of mackerel-agent
*/
func doVersion(_ *flag.FlagSet, _ []string) error {
	fmt.Printf("mackerel-agent version %s (rev %s) [%s %s %s] \n",
		version, gitcommit, runtime.GOOS, runtime.GOARCH, runtime.Version())
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
		return fmt.Errorf("hostID file is not found")
	}

	api, err := command.NewMackerelClient(conf.Apibase, conf.Apikey, version, gitcommit, conf.Verbose)
	if err != nil {
		return fmt.Errorf("faild to create api client: %s", err)
	}

	if !force && !prompter.YN(fmt.Sprintf("retire this host? (hostID: %s)", hostID), false) {
		return fmt.Errorf("retirement is canceled")
	}

	err = retry.Retry(10, 3*time.Second, func() error {
		return api.RetireHost(hostID)
	})
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
	command.RunOnce(conf, &command.AgentMeta{
		Version:  version,
		Revision: gitcommit,
	})
	return nil
}
