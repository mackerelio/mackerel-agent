package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/cmdutil"
	"github.com/pkg/errors"
)

var configLogger = logging.GetLogger("config")

// `apibase` and `agentName` are set from build flags
var apibase string

func getApibase() string {
	if apibase != "" {
		return apibase
	}
	return "https://api.mackerelio.com"
}

var agentName string

func getAgentName() string {
	if agentName != "" {
		return agentName
	}
	return "mackerel-agent"
}

// DefaultConfig stores standard settings for each environment
var DefaultConfig *Config

// CloudPlatform is an enum to represent which cloud platform the host is running on.
type CloudPlatform int

// CloudPlatform enum values
const (
	CloudPlatformAuto CloudPlatform = iota
	CloudPlatformNone
	CloudPlatformEC2
	CloudPlatformGCE
	CloudPlatformAzureVM
)

func (c CloudPlatform) String() string {
	switch c {
	case CloudPlatformAuto:
		return "auto"
	case CloudPlatformNone:
		return "none"
	case CloudPlatformEC2:
		return "ec2"
	case CloudPlatformGCE:
		return "gce"
	case CloudPlatformAzureVM:
		return "azurevm"
	}
	return ""
}

// UnmarshalText is used by toml unmarshaller
func (c *CloudPlatform) UnmarshalText(text []byte) error {
	switch string(text) {
	case "auto", "":
		*c = CloudPlatformAuto
		return nil
	case "none":
		*c = CloudPlatformNone
		return nil
	case "ec2":
		*c = CloudPlatformEC2
		return nil
	case "gce":
		*c = CloudPlatformGCE
		return nil
	case "azurevm":
		*c = CloudPlatformAzureVM
		return nil
	default:
		*c = CloudPlatformNone // Avoid panic
		return fmt.Errorf("failed to parse")
	}
}

// Config represents mackerel-agent's configuration file.
type Config struct {
	Apibase       string
	Apikey        string
	Root          string
	Pidfile       string
	Conffile      string
	Roles         []string
	Verbose       bool
	Silent        bool
	Diagnostic    bool          `toml:"diagnostic"`
	DisplayName   string        `toml:"display_name"`
	HostStatus    HostStatus    `toml:"host_status"`
	Filesystems   Filesystems   `toml:"filesystems"`
	HTTPProxy     string        `toml:"http_proxy"`
	CloudPlatform CloudPlatform `toml:"cloud_platform"`

	// This Plugin field is used to decode the toml file. After reading the
	// configuration from file, this field is set to nil.
	// Please consider using MetricPlugins and CheckPlugins.
	Plugin map[string]map[string]*PluginConfig

	Include string

	// Cannot exist in configuration files
	HostIDStorage   HostIDStorage
	MetricPlugins   map[string]*MetricPlugin
	CheckPlugins    map[string]*CheckPlugin
	MetadataPlugins map[string]*MetadataPlugin
	AutoShutdown    bool
}

// PluginConfig represents a plugin configuration.
type PluginConfig struct {
	CommandConfig
	NotificationInterval  *int32        `toml:"notification_interval"`
	CheckInterval         *int32        `toml:"check_interval"`
	ExecutionInterval     *int32        `toml:"execution_interval"`
	MaxCheckAttempts      *int32        `toml:"max_check_attempts"`
	CustomIdentifier      *string       `toml:"custom_identifier"`
	PreventAlertAutoClose bool          `toml:"prevent_alert_auto_close"`
	IncludePattern        *string       `toml:"include_pattern"`
	ExcludePattern        *string       `toml:"exclude_pattern"`
	Action                CommandConfig `toml:"action"`
	Memo                  string        `toml:"memo"`
}

// CommandConfig represents an executable command configuration.
type CommandConfig struct {
	Raw            interface{} `toml:"command"`
	User           string      `toml:"user"`
	Env            Env         `toml:"env"`
	TimeoutSeconds int64       `toml:"timeout_seconds"`
}

// Env represents environments.
type Env map[string]string

// ConvertToStrings converts to a slice of the form "key=value".
func (e Env) ConvertToStrings() ([]string, error) {
	env := make([]string, 0, len(e))
	for k, v := range e {
		if strings.Contains(k, "=") {
			return nil, fmt.Errorf("failed to parse plugin env. A key of env should not contain \"=\", but %q", k)
		}
		k = strings.Trim(k, " ")
		if k == "" {
			continue
		}
		env = append(env, k+"="+v)
	}
	return env, nil
}

// Command represents an executable command.
type Command struct {
	cmdutil.CommandOption
	Cmd  string
	Args []string
}

// Run the Command.
func (cmd *Command) Run() (stdout, stderr string, exitCode int, err error) {
	if len(cmd.Args) > 0 {
		return cmdutil.RunCommandArgs(cmd.Args, cmd.CommandOption)
	}
	return cmdutil.RunCommand(cmd.Cmd, cmd.CommandOption)
}

// RunWithEnv runs the Command with Environment.
func (cmd *Command) RunWithEnv(env []string) (stdout, stderr string, exitCode int, err error) {
	opt := cmdutil.CommandOption{
		TimeoutDuration: cmd.TimeoutDuration,
		User:            cmd.User,
		Env:             append(cmd.Env, env...),
	}
	if len(cmd.Args) > 0 {
		return cmdutil.RunCommandArgs(cmd.Args, opt)
	}
	return cmdutil.RunCommand(cmd.Cmd, opt)
}

// CommandString returns the command string for log messages
func (cmd *Command) CommandString() string {
	if len(cmd.Args) > 0 {
		return strings.Join(cmd.Args, " ")
	}
	return cmd.Cmd
}

// MetricPlugin represents the configuration of a metric plugin
// The User option is ignored on Windows
type MetricPlugin struct {
	Command          Command
	CustomIdentifier *string
	IncludePattern   *regexp.Regexp
	ExcludePattern   *regexp.Regexp
}

func (pconf *PluginConfig) buildMetricPlugin() (*MetricPlugin, error) {
	cmd, err := pconf.CommandConfig.parse()
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, fmt.Errorf("failed to parse plugin command. A configuration value of `command` should be string or string slice, but %T", pconf.Raw)
	}

	var (
		includePattern *regexp.Regexp
		excludePattern *regexp.Regexp
	)
	if pconf.IncludePattern != nil {
		includePattern, err = regexp.Compile(*pconf.IncludePattern)
		if err != nil {
			return nil, err
		}
	}
	if pconf.ExcludePattern != nil {
		excludePattern, err = regexp.Compile(*pconf.ExcludePattern)
		if err != nil {
			return nil, err
		}
	}

	return &MetricPlugin{
		Command:          *cmd,
		CustomIdentifier: pconf.CustomIdentifier,
		IncludePattern:   includePattern,
		ExcludePattern:   excludePattern,
	}, nil
}

// CheckPlugin represents the configuration of a check plugin
// The User option is ignored on Windows
type CheckPlugin struct {
	Command               Command
	CustomIdentifier      *string
	NotificationInterval  *int32
	CheckInterval         *int32
	MaxCheckAttempts      *int32
	PreventAlertAutoClose bool
	Action                *Command
	Memo                  string
}

func (pconf *PluginConfig) buildCheckPlugin(name string) (*CheckPlugin, error) {
	cmd, err := pconf.CommandConfig.parse()
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, fmt.Errorf("failed to parse plugin command. A configuration value of `command` should be string or string slice, but %T", pconf.Raw)
	}

	action, err := pconf.Action.parse()
	if err != nil {
		return nil, err
	}

	if utf8.RuneCountInString(pconf.Memo) > 250 {
		configLogger.Warningf("'plugin.checks.%s.memo' size exceeds 250 characters", name)
		str := pconf.Memo
		c := 0
		n := 0
		for len(str) > 0 && c < 250 {
			_, size := utf8.DecodeRuneInString(str)
			n += size
			c++
			str = str[size:]
		}
		pconf.Memo = pconf.Memo[:n]
	}

	plugin := CheckPlugin{
		Command:               *cmd,
		CustomIdentifier:      pconf.CustomIdentifier,
		NotificationInterval:  pconf.NotificationInterval,
		CheckInterval:         pconf.CheckInterval,
		MaxCheckAttempts:      pconf.MaxCheckAttempts,
		PreventAlertAutoClose: pconf.PreventAlertAutoClose,
		Action:                action,
		Memo:                  pconf.Memo,
	}
	if plugin.MaxCheckAttempts != nil && *plugin.MaxCheckAttempts > 1 && plugin.PreventAlertAutoClose {
		*plugin.MaxCheckAttempts = 1
		configLogger.Warningf("'plugin.checks.%s.max_check_attempts' is set to 1 (Unavailable with 'prevent_alert_auto_close')", name)
	}
	return &plugin, nil
}

// MetadataPlugin represents the configuration of a metadata plugin
// The User option is ignored on Windows
type MetadataPlugin struct {
	Command           Command
	ExecutionInterval *int32
}

func (pconf *PluginConfig) buildMetadataPlugin() (*MetadataPlugin, error) {
	cmd, err := pconf.CommandConfig.parse()
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, fmt.Errorf("failed to parse plugin command. A configuration value of `command` should be string or string slice, but %T", pconf.Raw)
	}

	return &MetadataPlugin{
		Command:           *cmd,
		ExecutionInterval: pconf.ExecutionInterval,
	}, nil
}

func (cc CommandConfig) parse() (cmd *Command, err error) {
	const errFmt = "failed to parse plugin command. A configuration value of `command` should be string or string slice, but %T"
	switch t := cc.Raw.(type) {
	case string:
		cmd = &Command{Cmd: t}
	case []interface{}:
		if len(t) > 0 {
			args := []string{}
			for _, vv := range t {
				str, ok := vv.(string)
				if !ok {
					return nil, fmt.Errorf(errFmt, cc.Raw)
				}
				args = append(args, str)
			}
			cmd = &Command{Args: args}
		} else {
			return nil, fmt.Errorf(errFmt, cc.Raw)
		}
	case []string:
		cmd = &Command{Args: t}
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf(errFmt, cc.Raw)
	}
	cmd.User = cc.User
	cmd.Env, err = cc.Env.ConvertToStrings()
	if err != nil {
		return nil, err
	}
	cmd.TimeoutDuration = time.Duration(cc.TimeoutSeconds * int64(time.Second))
	return cmd, nil
}

// PostMetricsInterval XXX
var PostMetricsInterval = 1 * time.Minute

// HostStatus configure host status on agent start/stop
type HostStatus struct {
	OnStart string `toml:"on_start"`
	OnStop  string `toml:"on_stop"`
}

// Filesystems configure filesystem related settings
type Filesystems struct {
	Ignore        Regexpwrapper `toml:"ignore"`
	UseMountpoint bool          `toml:"use_mountpoint"`
}

// Regexpwrapper is a wrapper type for marshalling string
type Regexpwrapper struct {
	*regexp.Regexp
}

// UnmarshalText for compiling regexp string while loading toml
func (r *Regexpwrapper) UnmarshalText(text []byte) error {
	var err error
	r.Regexp, err = regexp.Compile(string(text))
	return err
}

// ListCustomIdentifiers returns a list of customIdentifiers.
func (conf *Config) ListCustomIdentifiers() []string {
	var customIdentifiers []string
	for _, pconf := range conf.MetricPlugins {
		if pconf.CustomIdentifier != nil && index(customIdentifiers, *pconf.CustomIdentifier) == -1 {
			customIdentifiers = append(customIdentifiers, *pconf.CustomIdentifier)
		}
	}
	for _, cconf := range conf.CheckPlugins {
		if cconf.CustomIdentifier != nil && index(customIdentifiers, *cconf.CustomIdentifier) == -1 {
			customIdentifiers = append(customIdentifiers, *cconf.CustomIdentifier)
		}
	}
	return customIdentifiers
}

func index(xs []string, y string) int {
	for i, x := range xs {
		if x == y {
			return i
		}
	}
	return -1
}

// LoadConfig loads a Config from a file.
func LoadConfig(conffile string) (*Config, error) {
	config, err := loadConfigFile(conffile)
	if err != nil {
		return nil, err
	}

	// set default values if config does not have values
	if config.Apibase == "" {
		config.Apibase = DefaultConfig.Apibase
	}
	if config.Root == "" {
		config.Root = DefaultConfig.Root
	}
	if config.Pidfile == "" {
		config.Pidfile = DefaultConfig.Pidfile
	}
	if config.Verbose == false {
		config.Verbose = DefaultConfig.Verbose
	}
	if config.Diagnostic == false {
		config.Diagnostic = DefaultConfig.Diagnostic
	}

	return config, err
}

func (conf *Config) setEachPlugins() error {
	if pconfs, ok := conf.Plugin["metrics"]; ok {
		var err error
		for name, pconf := range pconfs {
			conf.MetricPlugins[name], err = pconf.buildMetricPlugin()
			if err != nil {
				return errors.Wrap(err, "plugin.metrics."+name)
			}
		}
	}
	if pconfs, ok := conf.Plugin["checks"]; ok {
		var err error
		for name, pconf := range pconfs {
			conf.CheckPlugins[name], err = pconf.buildCheckPlugin(name)
			if err != nil {
				return errors.Wrap(err, "plugin.checks."+name)
			}
		}
	}
	if pconfs, ok := conf.Plugin["metadata"]; ok {
		var err error
		for name, pconf := range pconfs {
			conf.MetadataPlugins[name], err = pconf.buildMetadataPlugin()
			if err != nil {
				return errors.Wrap(err, "plugin.metadata."+name)
			}
		}
	}
	// Make Plugins empty because we should not use this later.
	// Use MetricPlugins, CheckPlugins and MetadataPlugins.
	conf.Plugin = nil
	return nil
}

func loadConfigFile(file string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(file, config); err != nil {
		return config, err
	}

	config.MetricPlugins = make(map[string]*MetricPlugin)
	config.CheckPlugins = make(map[string]*CheckPlugin)
	config.MetadataPlugins = make(map[string]*MetadataPlugin)
	if err := config.setEachPlugins(); err != nil {
		return nil, err
	}

	if config.Include != "" {
		if err := includeConfigFile(config, config.Include); err != nil {
			return config, err
		}
	}

	return config, nil
}

func includeConfigFile(config *Config, include string) error {
	files, err := filepath.Glob(include)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Save current "roles" value and reset it
		// because toml.DecodeFile()-ing on a fulfilled struct
		// produces bizarre array values.
		rolesSaved := config.Roles
		config.Roles = nil

		meta, err := toml.DecodeFile(file, &config)
		if err != nil {
			return fmt.Errorf("while loading included config file %s: %s", file, err)
		}

		// If included config does not have "roles" key,
		// use the previous roles configuration value.
		if meta.IsDefined("roles") == false {
			config.Roles = rolesSaved
		}

		// Add new plugin or overwrite a plugin with the same plugin name.
		if err := config.setEachPlugins(); err != nil {
			return err
		}
	}

	return nil
}

func (conf *Config) hostIDStorage() HostIDStorage {
	if conf.HostIDStorage == nil {
		conf.HostIDStorage = &FileSystemHostIDStorage{Root: conf.Root}
	}
	return conf.HostIDStorage
}

// LoadHostID loads the previously saved host id.
func (conf *Config) LoadHostID() (string, error) {
	return conf.hostIDStorage().LoadHostID()
}

// SaveHostID saves the host id, which may be restored by LoadHostID.
func (conf *Config) SaveHostID(id string) error {
	return conf.hostIDStorage().SaveHostID(id)
}

// DeleteSavedHostID deletes the host id saved by SaveHostID.
func (conf *Config) DeleteSavedHostID() error {
	return conf.hostIDStorage().DeleteSavedHostID()
}

// HostIDStorage is an interface which maintains persistency
// of the "Host ID" for the current host where the agent is running on.
// The ID is always generated and given by Mackerel (mackerel.io).
type HostIDStorage interface {
	LoadHostID() (string, error)
	SaveHostID(id string) error
	DeleteSavedHostID() error
}

// FileSystemHostIDStorage is the default HostIDStorage
// which saves/loads the host id using an id file on the local filesystem.
// The file will be located at /var/lib/mackerel-agent/id by default on linux.
type FileSystemHostIDStorage struct {
	Root string
}

const idFileName = "id"

// HostIDFile is the location of the host id file.
func (s FileSystemHostIDStorage) HostIDFile() string {
	return filepath.Join(s.Root, idFileName)
}

// LoadHostID loads the current host ID from the mackerel-agent's id file.
func (s FileSystemHostIDStorage) LoadHostID() (string, error) {
	content, err := ioutil.ReadFile(s.HostIDFile())
	if err != nil {
		return "", err
	}
	hostID := strings.TrimRight(string(content), "\r\n")
	if hostID == "" {
		return "", fmt.Errorf("HostIDFile found, but the content is empty")
	}
	return hostID, nil
}

// SaveHostID saves the host ID to the mackerel-agent's id file.
func (s FileSystemHostIDStorage) SaveHostID(id string) error {
	err := os.MkdirAll(s.Root, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(s.HostIDFile())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(id))
	if err != nil {
		return err
	}

	return nil
}

// DeleteSavedHostID deletes the mackerel-agent's id file.
func (s FileSystemHostIDStorage) DeleteSavedHostID() error {
	return os.Remove(s.HostIDFile())
}
