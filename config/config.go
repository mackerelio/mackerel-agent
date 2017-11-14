package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/util"
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

var defaultConnectionConfig = ConnectionConfig{
	PostMetricsDequeueDelaySeconds: 30,     // Check the metric values queue for every half minute
	PostMetricsRetryDelaySeconds:   60,     // Wait a minute before retrying metric value posts
	PostMetricsRetryMax:            60,     // Retry up to 60 times (30s * 60 = 30min)
	PostMetricsBufferSize:          6 * 60, // Keep metric values of 6 hours span in the queue
}

// Config represents mackerel-agent's configuration file.
type Config struct {
	Apibase     string
	Apikey      string
	Root        string
	Pidfile     string
	Conffile    string
	Roles       []string
	Verbose     bool
	Silent      bool
	Diagnostic  bool `toml:"diagnostic"`
	Connection  ConnectionConfig
	DisplayName string      `toml:"display_name"`
	HostStatus  HostStatus  `toml:"host_status"`
	Filesystems Filesystems `toml:"filesystems"`
	HTTPProxy   string      `toml:"http_proxy"`

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
}

// PluginConfig represents a plugin configuration.
type PluginConfig struct {
	CommandRaw            interface{} `toml:"command"`
	User                  string
	NotificationInterval  *int32        `toml:"notification_interval"`
	CheckInterval         *int32        `toml:"check_interval"`
	ExecutionInterval     *int32        `toml:"execution_interval"`
	MaxCheckAttempts      *int32        `toml:"max_check_attempts"`
	CustomIdentifier      *string       `toml:"custom_identifier"`
	PreventAlertAutoClose bool          `toml:"prevent_alert_auto_close"`
	IncludePattern        *string       `toml:"include_pattern"`
	ExcludePattern        *string       `toml:"exclude_pattern"`
	Action                CommandConfig `toml:"action"`
}

// CommandConfig represents an executable command configuration.
type CommandConfig struct {
	Raw  interface{} `toml:"command"`
	User string
}

// Command represents an executable command.
type Command struct {
	Cmd  string
	Args []string
	User string
}

// Run the Command.
func (cmd *Command) Run() (stdout, stderr string, exitCode int, err error) {
	if len(cmd.Args) > 0 {
		return util.RunCommandArgs(cmd.Args, cmd.User, nil)
	}
	return util.RunCommand(cmd.Cmd, cmd.User, nil)
}

// RunWithEnv runs the Command with Environment.
func (cmd *Command) RunWithEnv(env []string) (stdout, stderr string, exitCode int, err error) {
	if len(cmd.Args) > 0 {
		return util.RunCommandArgs(cmd.Args, cmd.User, env)
	}
	return util.RunCommand(cmd.Cmd, cmd.User, env)
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
	cmd, err := parseCommand(pconf.CommandRaw, pconf.User)
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, fmt.Errorf("failed to parse plugin command. A configuration value of `command` should be string or string slice, but %T", pconf.CommandRaw)
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
	NotificationInterval  *int32
	CheckInterval         *int32
	MaxCheckAttempts      *int32
	PreventAlertAutoClose bool
	Action                *Command
}

func (pconf *PluginConfig) buildCheckPlugin(name string) (*CheckPlugin, error) {
	cmd, err := parseCommand(pconf.CommandRaw, pconf.User)
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, fmt.Errorf("failed to parse plugin command. A configuration value of `command` should be string or string slice, but %T", pconf.CommandRaw)
	}
	action, err := parseCommand(pconf.Action.Raw, pconf.Action.User)
	if err != nil {
		return nil, err
	}
	plugin := CheckPlugin{
		Command:               *cmd,
		NotificationInterval:  pconf.NotificationInterval,
		CheckInterval:         pconf.CheckInterval,
		MaxCheckAttempts:      pconf.MaxCheckAttempts,
		PreventAlertAutoClose: pconf.PreventAlertAutoClose,
		Action:                action,
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
	cmd, err := parseCommand(pconf.CommandRaw, pconf.User)
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, fmt.Errorf("failed to parse plugin command. A configuration value of `command` should be string or string slice, but %T", pconf.CommandRaw)
	}
	return &MetadataPlugin{
		Command:           *cmd,
		ExecutionInterval: pconf.ExecutionInterval,
	}, nil
}

func parseCommand(commandRaw interface{}, user string) (command *Command, err error) {
	const errFmt = "failed to parse plugin command. A configuration value of `command` should be string or string slice, but %T"
	switch t := commandRaw.(type) {
	case string:
		return &Command{
			Cmd:  t,
			User: user,
		}, nil
	case []interface{}:
		if len(t) > 0 {
			args := []string{}
			for _, vv := range t {
				str, ok := vv.(string)
				if !ok {
					return nil, fmt.Errorf(errFmt, commandRaw)
				}

				args = append(args, str)
			}
			return &Command{
				Args: args,
				User: user,
			}, nil
		}
		return nil, fmt.Errorf(errFmt, commandRaw)
	case []string:
		return &Command{
			Args: t,
			User: user,
		}, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf(errFmt, commandRaw)
	}
}

const postMetricsDequeueDelaySecondsMax = 59   // max delay seconds for dequeuing from buffer queue
const postMetricsRetryDelaySecondsMax = 3 * 60 // max delay seconds for retrying a request that caused errors

// PostMetricsInterval XXX
var PostMetricsInterval = 1 * time.Minute

// ConnectionConfig XXX
type ConnectionConfig struct {
	PostMetricsDequeueDelaySeconds int `toml:"post_metrics_dequeue_delay_seconds"` // delay for dequeuing from buffer queue
	PostMetricsRetryDelaySeconds   int `toml:"post_metrics_retry_delay_seconds"`   // delay for retrying a request that caused errors
	PostMetricsRetryMax            int `toml:"post_metrics_retry_max"`             // max numbers of retries for a request that causes errors
	PostMetricsBufferSize          int `toml:"post_metrics_buffer_size"`           // max numbers of requests stored in buffer queue.
}

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

// CheckNames returns a list of name of the check plugins
func (conf *Config) CheckNames() []string {
	checks := []string{}
	for name := range conf.CheckPlugins {
		checks = append(checks, name)
	}
	return checks
}

// ListCustomIdentifiers returns a list of customIdentifiers.
func (conf *Config) ListCustomIdentifiers() []string {
	var customIdentifiers []string
	for _, pconf := range conf.MetricPlugins {
		if pconf.CustomIdentifier != nil && index(customIdentifiers, *pconf.CustomIdentifier) == -1 {
			customIdentifiers = append(customIdentifiers, *pconf.CustomIdentifier)
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
	if config.Connection.PostMetricsDequeueDelaySeconds == 0 {
		config.Connection.PostMetricsDequeueDelaySeconds = DefaultConfig.Connection.PostMetricsDequeueDelaySeconds
	}
	if config.Connection.PostMetricsDequeueDelaySeconds > postMetricsDequeueDelaySecondsMax {
		configLogger.Warningf("'post_metrics_dequese_delay_seconds' is set to %d (Maximum Value).", postMetricsDequeueDelaySecondsMax)
		config.Connection.PostMetricsDequeueDelaySeconds = postMetricsDequeueDelaySecondsMax
	}
	if config.Connection.PostMetricsRetryDelaySeconds == 0 {
		config.Connection.PostMetricsRetryDelaySeconds = DefaultConfig.Connection.PostMetricsRetryDelaySeconds
	}
	if config.Connection.PostMetricsRetryDelaySeconds > postMetricsRetryDelaySecondsMax {
		configLogger.Warningf("'post_metrics_retry_delay_seconds' is set to %d (Maximum Value).", postMetricsRetryDelaySecondsMax)
		config.Connection.PostMetricsRetryDelaySeconds = postMetricsRetryDelaySecondsMax
	}
	if config.Connection.PostMetricsRetryMax == 0 {
		config.Connection.PostMetricsRetryMax = DefaultConfig.Connection.PostMetricsRetryMax
	}
	if config.Connection.PostMetricsBufferSize == 0 {
		config.Connection.PostMetricsBufferSize = DefaultConfig.Connection.PostMetricsBufferSize
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
	return strings.TrimRight(string(content), "\r\n"), nil
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
