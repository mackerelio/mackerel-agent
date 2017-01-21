package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/version"
	"github.com/motemen/go-cli"
)

// allow options like -role=... -role=...
type roleFullnamesFlag []string

var roleFullnamePattern = regexp.MustCompile(`^[a-zA-Z0-9][-_a-zA-Z0-9]*:\s*[a-zA-Z0-9][-_a-zA-Z0-9]*$`)

func (r *roleFullnamesFlag) String() string {
	return fmt.Sprint(*r)
}

func (r *roleFullnamesFlag) Set(input string) error {
	inputRoles := strings.Split(input, ",")
	*r = append(*r, inputRoles...)
	return nil
}

var logger = logging.GetLogger("main")

func main() {
	// although the possibility is very low, mackerel-agent may panic because of
	// a race condition in multi-threaded environment on some OS/Arch.
	// So fix GOMAXPROCS to 1 just to be safe.
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(1)
	}
	// force disabling http2 for now
	godebug := os.Getenv("GODEBUG")
	if godebug != "" {
		godebug += ","
	}
	godebug += "http2client=0"
	os.Setenv("GODEBUG", godebug)
	cli.Run(os.Args[1:])
}

func printRetireUsage() {
	usage := fmt.Sprintf(`Usage of mackerel-agent retire:
  -conf string
        Config file path (Configs in this file are over-written by command line options)
        (default "%s")
  -force
        force retirement without prompting
  -apibase string
        API base (default "%s")
  -apikey string
        (DEPRECATED) API key from mackerel.io web site`,
		config.DefaultConfig.Conffile,
		config.DefaultConfig.Apibase)

	fmt.Fprintln(os.Stderr, usage)
	os.Exit(2)
}

func resolveConfigForRetire(fs *flag.FlagSet, argv []string) (*config.Config, bool, error) {
	var force = fs.Bool("force", false, "force retirement without prompting")
	fs.Usage = printRetireUsage
	conf, err := resolveConfig(fs, argv)
	return conf, *force, err
}

// resolveConfig parses command line arguments and loads config file to
// return config.Config information.
func resolveConfig(fs *flag.FlagSet, argv []string) (*config.Config, error) {
	conf := &config.Config{}

	var (
		conffile      = fs.String("conf", config.DefaultConfig.Conffile, "Config file path (Configs in this file are over-written by command line options)")
		apibase       = fs.String("apibase", config.DefaultConfig.Apibase, "API base")
		pidfile       = fs.String("pidfile", config.DefaultConfig.Pidfile, "File containing PID")
		root          = fs.String("root", config.DefaultConfig.Root, "Directory containing variable state information")
		apikey        = fs.String("apikey", "", "(DEPRECATED) API key from mackerel.io web site")
		diagnostic    = fs.Bool("diagnostic", false, "Enables diagnostic features")
		verbose       bool
		roleFullnames roleFullnamesFlag
	)
	fs.BoolVar(&verbose, "verbose", config.DefaultConfig.Verbose, "Toggle verbosity")
	fs.BoolVar(&verbose, "v", config.DefaultConfig.Verbose, "Toggle verbosity (shorthand)")

	// The value of "role" option is internally "roll fullname",
	// but we call it "role" here for ease.
	fs.Var(&roleFullnames, "role", "Set this host's roles (format: <service>:<role>)")

	fs.Parse(argv)

	conf, confErr := config.LoadConfig(*conffile)
	if confErr != nil {
		return nil, fmt.Errorf("failed to load the config file: %s", confErr)
	}
	conf.Conffile = *conffile

	// overwrite config from file by config from args
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "apibase":
			conf.Apibase = *apibase
		case "apikey":
			conf.Apikey = *apikey
		case "pidfile":
			conf.Pidfile = *pidfile
		case "root":
			conf.Root = *root
		case "diagnostic":
			conf.Diagnostic = *diagnostic
		case "verbose", "v":
			conf.Verbose = verbose
		case "role":
			conf.Roles = roleFullnames
		}
	})

	r := []string{}
	for _, roleFullName := range conf.Roles {
		if !roleFullnamePattern.MatchString(roleFullName) {
			logger.Errorf("Bad format for role fullname (expecting <service>:<role>. Alphabet, numbers, hyphens and underscores are acceptable, but the first character must not be a hyphen or an underscore.): '%s'", roleFullName)
		} else {
			r = append(r, roleFullName)
		}
	}
	conf.Roles = r

	if conf.Verbose && conf.Silent {
		logger.Warningf("both of `verbose` and `silent` option are specified. In this case, `verbose` get preference over `silent`")
	}

	if conf.Apikey == "" {
		return nil, fmt.Errorf("apikey must be specified in the config file (or by the DEPRECATED command-line flag)")
	}

	if conf.HTTPProxy != "" {
		os.Setenv("HTTP_PROXY", conf.HTTPProxy)
	}
	return conf, nil
}

func createPidFile(pidfile string) error {
	if pidString, err := ioutil.ReadFile(pidfile); err == nil {
		if pid, err := strconv.Atoi(string(pidString)); err == nil {
			if existsPid(pid) {
				return fmt.Errorf("pidfile found, try stopping another running mackerel-agent or delete %s", pidfile)
			}
			// Note mackerel-agent in windows can't remove pidfile during stoping the service
			logger.Warningf("Pidfile found, but there seems no another process of mackerel-agent. Ignoring %s", pidfile)
		} else {
			logger.Warningf("Malformed pidfile found. Ignoring %s", pidfile)
		}
	}

	err := os.MkdirAll(filepath.Dir(pidfile), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(pidfile)
	if err != nil {
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

func start(conf *config.Config, termCh chan struct{}) error {
	if conf.Silent {
		logging.SetLogLevel(logging.ERROR)
	}
	if conf.Verbose {
		logging.SetLogLevel(logging.DEBUG)
	}
	logger.Infof("Starting mackerel-agent version:%s, rev:%s, apibase:%s", version.VERSION, version.GITCOMMIT, conf.Apibase)

	if err := createPidFile(conf.Pidfile); err != nil {
		return fmt.Errorf("createPidFile(%q) failed: %s", conf.Pidfile, err)
	}
	defer removePidFile(conf.Pidfile)

	ctx, err := command.Prepare(conf)
	if err != nil {
		return fmt.Errorf("command.Prepare failed: %s", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go signalHandler(c, ctx, termCh)

	return command.Run(ctx, termCh)
}

var maxTerminatingInterval = 30 * time.Second

func signalHandler(c chan os.Signal, ctx *command.Context, termCh chan struct{}) {
	received := false
	for sig := range c {
		if sig == syscall.SIGHUP {
			logger.Debugf("Received signal '%v'", sig)
			// TODO reload configuration file

			ctx.UpdateHostSpecs()
		} else {
			if !received {
				received = true
				logger.Infof(
					"Received signal '%v', try graceful shutdown up to %f seconds. If you want force shutdown immediately, send a signal again.",
					sig,
					maxTerminatingInterval.Seconds())
			} else {
				logger.Infof("Received signal '%v' again, force shutdown.", sig)
			}
			termCh <- struct{}{}
			go func() {
				time.Sleep(maxTerminatingInterval)
				logger.Infof("Timed out. force shutdown.")
				termCh <- struct{}{}
			}()
		}
	}
}
