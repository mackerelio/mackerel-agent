package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/version"
)

// allow options like -role=... -role=...
type roleFullnamesFlag []string

var roleFullnamePattern = regexp.MustCompile(`^[\w-]+:\s*[\w-]+$`)

func (r *roleFullnamesFlag) String() string {
	return fmt.Sprint(*r)
}

func (r *roleFullnamesFlag) Set(input string) error {
	inputRoles := strings.Split(input, ",")

	for _, inputRole := range inputRoles {
		if roleFullnamePattern.MatchString(inputRole) == false {
			return fmt.Errorf("Bad format for role fullname (expecting <service>:<role>): %s", inputRole)
		}
	}

	*r = append(*r, inputRoles...)

	return nil
}

var logger = logging.GetLogger("main")

func main() {
	config := resolveConfig()

	if config.Verbose {
		logging.ConfigureLoggers("DEBUG")
	} else {
		logging.ConfigureLoggers("INFO")
	}

	logger.Infof("Starting mackerel-agent version:%s, rev:%s", version.VERSION, version.GITCOMMIT)

	if config.Apikey == "" {
		logger.Criticalf("Apikey must be specified in the command-line flag or in the config file")
		os.Exit(1)
	}

	if err := start(config); err != nil {
		os.Exit(1)
	}
}

func resolveConfig() (config mackerel.Config) {
	conffile := flag.String("conf", "/etc/mackerel-agent/mackerel-agent.conf", "Config file path (Configs in this file are over-written by command line options)")
	apibase := flag.String("apibase", mackerel.DefaultConfig.Apibase, "API base")
	pidfile := flag.String("pidfile", mackerel.DefaultConfig.Pidfile, "File containing PID")
	root := flag.String("root", mackerel.DefaultConfig.Root, "Directory containing variable state information")
	apikey := flag.String("apikey", "", "API key from mackerel.io web site")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", mackerel.DefaultConfig.Verbose, "Toggle verbosity")
	flag.BoolVar(&verbose, "v", mackerel.DefaultConfig.Verbose, "Toggle verbosity (shorthand)")

	// The value of "role" option is internally "roll fullname",
	// but we call it "role" here for ease.
	var roleFullnames roleFullnamesFlag
	flag.Var(&roleFullnames, "role", "Set this host's roles (format: <service>:<role>)")

	flag.Parse()

	config, confErr := mackerel.LoadConfig(*conffile)
	if confErr != nil {
		logger.Criticalf("Failed to load the config file: %s", confErr)
		os.Exit(1)
	}

	// overwrite config from file by config from args
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "apibase":
			config.Apibase = *apibase
		case "apikey":
			config.Apikey = *apikey
		case "pidfile":
			config.Pidfile = *pidfile
		case "root":
			config.Root = *root
		case "verbose", "v":
			config.Verbose = verbose
		case "role":
			config.Roles = roleFullnames
		}
	})

	return
}

func createPidFile(pidfile string) error {
	if pidString, err := ioutil.ReadFile(pidfile); err == nil {
		if pid, err := strconv.Atoi(string(pidString)); err == nil {
			if _, err := os.Stat(fmt.Sprintf("/proc/%d/", pid)); err == nil {
				return fmt.Errorf("Pidfile found, try stopping another running mackerel-agent or delete %s", pidfile)
			} else {
				logger.Warningf("Pidfile found, but there seems no another process of mackerel-agent. Ignoring %s", pidfile)
			}
		} else {
			logger.Warningf("Malformed pidfile found. Ignoring %s", pidfile)
		}
	}

	file, err := os.Create(pidfile)
	if err != nil {
		logger.Criticalf("Failed to create a pidfile: %s", err)
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d", os.Getpid())
	return err
}

func removePidFile(pidfile string) {
	if err := os.Remove(pidfile); err != nil {
		logger.Errorf("Failed to remove the pidfile: %s: %s", pidfile, err)
	}
}

func start(config mackerel.Config) error {
	if err := createPidFile(config.Pidfile); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		for sig := range c {
			if sig == syscall.SIGHUP {
				// nop
				// TODO reload configuration file
				logger.Debugf("Received signal '%v'", sig)
			} else {
				logger.Infof("Received signal '%v', exiting", sig)
				removePidFile(config.Pidfile)
				os.Exit(0)
			}
		}
	}()

	command.Run(config)
	return nil
}
