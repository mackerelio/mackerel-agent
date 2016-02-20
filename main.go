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

	"github.com/Songmu/prompter"
	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/version"
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

const (
	exitStatusOK = iota
	exitStatusError
)

func main() {
	os.Exit(dispatch(os.Args[1:]))
}

// empty string is dealt with the key of main process
const mainProcess = ""

// subcommands and processes of the mackerel-agent
func commands() map[string](func([]string) int) {
	return map[string](func([]string) int){
		mainProcess:  doMain,
		"version":    doVersion,
		"retire":     doRetire,
		"configtest": doConfigtest,
		"once":       doOnce,
	}
}

func doVersion(_ []string) int {
	fmt.Printf("mackerel-agent version %s (rev %s) [%s %s %s] \n",
		version.VERSION, version.GITCOMMIT, runtime.GOOS, runtime.GOARCH, runtime.Version())
	return exitStatusOK
}

func doConfigtest(argv []string) int {
	conf := resolveConfig(argv)
	if conf == nil {
		return exitStatusError
	}
	fmt.Fprintf(os.Stderr, "%s Syntax OK\n", conf.Conffile)
	return exitStatusOK
}

func doMain(argv []string) int {
	conf := resolveConfig(argv)
	if conf == nil {
		return exitStatusError
	}
	if conf.Verbose {
		logging.SetLogLevel(logging.DEBUG)
	}
	logger.Infof("Starting mackerel-agent version:%s, rev:%s, apibase:%s", version.VERSION, version.GITCOMMIT, conf.Apibase)
	return start(conf)
}

func doRetire(argv []string) int {
	conf, force, err := resolveConfigForRetire(argv)
	if err != nil {
		return exitStatusError
	}

	hostID, err := conf.LoadHostID()
	if err != nil {
		logger.Warningf("HostID file is not found")
		return exitStatusError
	}

	api, err := mackerel.NewAPI(conf.Apibase, conf.Apikey, conf.Verbose)
	if err != nil {
		logger.Errorf("failed to create api client: %s", err)
		return exitStatusError
	}

	if !force && !prompter.YN(fmt.Sprintf("retire this host? (hostID: %s)", hostID), false) {
		logger.Infof("Retirement is canceled.")
		return exitStatusError
	}

	err = api.RetireHost(hostID)
	if err != nil {
		logger.Errorf("failed to retire the host: %s", err)
		return exitStatusError
	}
	logger.Infof("This host (hostID: %s) has been retired.", hostID)
	// just to try to remove hostID file.
	err = conf.DeleteSavedHostID()
	if err != nil {
		logger.Warningf("Failed to remove HostID file: %s", err)
	}
	return exitStatusOK
}

func doOnce(argv []string) int {
	// dirty hack `resolveConfig` required apikey so fill up
	argvOpt := append(argv, "-apikey=dummy")
	conf := resolveConfig(argvOpt)
	if conf == nil {
		return exitStatusError
	}
	command.RunOnce(conf)
	return exitStatusOK
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
        API key from mackerel.io web site`,
		config.DefaultConfig.Conffile,
		config.DefaultConfig.Apibase)

	fmt.Fprintln(os.Stderr, usage)
	os.Exit(2)
}

var helpReg = regexp.MustCompile(`^--?h(?:elp)?$`)
var forceReg = regexp.MustCompile(`^--?force$`)

func resolveConfigForRetire(argv []string) (*config.Config, bool, error) {
	optArgs := []string{}
	isForce := false
	for _, v := range argv {
		if helpReg.MatchString(v) {
			printRetireUsage()
		}
		if forceReg.MatchString(v) {
			isForce = true
			continue
		}
		optArgs = append(optArgs, v)
	}
	conf := resolveConfig(optArgs)
	if conf == nil {
		printRetireUsage()
	}
	return conf, isForce, nil
}

func printSubCommands() {
	for c, _ := range commands() {
		if c != mainProcess {
			fmt.Fprintf(os.Stderr, "  %s\n", c)
		}
	}
}

// resolveConfig parses command line arguments and loads config file to
// return config.Config information.
func resolveConfig(argv []string) *config.Config {
	conf := &config.Config{}

	fs := flag.NewFlagSet("mackerel-agent", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of mackerel-agent:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSUB COMMANDS:\n")
		printSubCommands()
	}
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
	conf.Conffile = *conffile
	if confErr != nil {
		logger.Criticalf("Failed to load the config file: %s", confErr)
		return nil
	}

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

	if conf.Apikey == "" {
		logger.Criticalf("Apikey must be specified in the command-line flag or in the config file")
		return nil
	}
	return conf
}

func createPidFile(pidfile string) error {
	if pidString, err := ioutil.ReadFile(pidfile); err == nil {
		if pid, err := strconv.Atoi(string(pidString)); err == nil {
			if existsPid(pid) {
				return fmt.Errorf("Pidfile found, try stopping another running mackerel-agent or delete %s", pidfile)
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

func start(conf *config.Config) int {
	if err := createPidFile(conf.Pidfile); err != nil {
		return exitStatusError
	}
	defer removePidFile(conf.Pidfile)

	ctx, err := command.Prepare(conf)
	if err != nil {
		logger.Criticalf(err.Error())
		return exitStatusError
	}

	termCh := make(chan struct{})
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
