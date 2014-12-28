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
	"time"

	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
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
	conf, printVersion := resolveConfig()

	if printVersion {
		fmt.Printf("mackerel-agent version %s (rev %s)\n", version.VERSION, version.GITCOMMIT)
		exitWithoutPidfileCleaning(0)
	}

	if conf.Verbose {
		logging.ConfigureLoggers("DEBUG")
	} else {
		logging.ConfigureLoggers("INFO")
	}

	logger.Infof("Starting mackerel-agent version:%s, rev:%s", version.VERSION, version.GITCOMMIT)

	if conf.Apikey == "" {
		logger.Criticalf("Apikey must be specified in the command-line flag or in the config file")
		exit(1, conf)
	}

	if err := start(conf); err != nil {
		exit(1, conf)
	}
}

// resolveConfig parses command line arguments and loads config file to
// return config.Config information.
// As a special case, if `-version` flag is given it stops processing
// and return true for the second return value.
func resolveConfig() (*config.Config, bool) {
	conf := &config.Config{}

	var (
		conffile     = flag.String("conf", config.DefaultConfig.Conffile, "Config file path (Configs in this file are over-written by command line options)")
		apibase      = flag.String("apibase", config.DefaultConfig.Apibase, "API base")
		pidfile      = flag.String("pidfile", config.DefaultConfig.Pidfile, "File containing PID")
		root         = flag.String("root", config.DefaultConfig.Root, "Directory containing variable state information")
		apikey       = flag.String("apikey", "", "API key from mackerel.io web site")
		printVersion = flag.Bool("version", false, "Prints version and exit")
	)

	var verbose bool
	flag.BoolVar(&verbose, "verbose", config.DefaultConfig.Verbose, "Toggle verbosity")
	flag.BoolVar(&verbose, "v", config.DefaultConfig.Verbose, "Toggle verbosity (shorthand)")

	// The value of "role" option is internally "roll fullname",
	// but we call it "role" here for ease.
	var roleFullnames roleFullnamesFlag
	flag.Var(&roleFullnames, "role", "Set this host's roles (format: <service>:<role>)")

	flag.Parse()

	if *printVersion {
		return nil, true
	}

	conf, confErr := config.LoadConfig(*conffile)
	if confErr != nil {
		logger.Criticalf("Failed to load the config file: %s", confErr)
		exitWithoutPidfileCleaning(1)
	}

	// overwrite config from file by config from args
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "apibase":
			conf.Apibase = *apibase
		case "apikey":
			conf.Apikey = *apikey
		case "pidfile":
			conf.Pidfile = *pidfile
		case "root":
			conf.Root = *root
		case "verbose", "v":
			conf.Verbose = verbose
		case "role":
			conf.Roles = roleFullnames
		}
	})

	return conf, false
}

func createPidFile(pidfile string) error {
	if pidString, err := ioutil.ReadFile(pidfile); err == nil {
		if pid, err := strconv.Atoi(string(pidString)); err == nil {
			if _, err := os.Stat(fmt.Sprintf("/proc/%d/", pid)); err == nil {
				return fmt.Errorf("Pidfile found, try stopping another running mackerel-agent or delete %s", pidfile)
			}
			logger.Warningf("Pidfile found, but there seems no another process of mackerel-agent. Ignoring %s", pidfile)
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

func exit(exitCode int, conf *config.Config) {
	removePidFile(conf.Pidfile)
	exitWithoutPidfileCleaning(exitCode)
}

func exitWithoutPidfileCleaning(exitCode int) {
	os.Exit(exitCode)
}

const MAX_TERMINATING_INTERVAL = 30

func start(conf *config.Config) error {
	if err := createPidFile(conf.Pidfile); err != nil {
		return err
	}

	api, host, err := command.Prepare(conf)
	if err != nil {
		logger.Criticalf(err.Error())
		exit(1, conf)
	}

	termCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		received := false
		for sig := range c {
			if sig == syscall.SIGHUP {
				logger.Debugf("Received signal '%v'", sig)
				// TODO reload configuration file

				command.UpdateHostSpecs(conf, api, host)
			} else {
				if !received {
					received = true
					logger.Infof(
						"Received signal '%v', try graceful shutdown up to %d seconds. If you want force shutdown immediately, send a signal again.",
						sig,
						MAX_TERMINATING_INTERVAL)
				} else {
					logger.Infof("Received signal '%v' again, force shutdown.", sig)
				}
				termCh <- struct{}{}
				go func() {
					time.Sleep(MAX_TERMINATING_INTERVAL * time.Second)
					logger.Infof("Timed out. force shutdown.")
					termCh <- struct{}{}
				}()
			}
		}
	}()

	exitCode := command.Run(conf, api, host, termCh)
	exit(exitCode, conf)

	return nil
}
