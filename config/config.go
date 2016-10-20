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
	"github.com/mackerelio/mackerel-agent/logging"
)

var configLogger = logging.GetLogger("config")

// `apibase` and `agentName` are set from build flags
var apibase string

func getApibase() string {
	if apibase != "" {
		return apibase
	}
	return "https://mackerel.io"
}

var agentName string

func getAgentName() string {
	if agentName != "" {
		return agentName
	}
	return "mackerel-agent"
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

	// Corresponds to the set of [plugin.<kind>.<name>] sections
	// the key of the map is <kind>, which should be one of "metrics" or "checks".
	Plugin map[string]PluginConfigs

	Include string

	// Cannot exist in configuration files
	HostIDStorage HostIDStorage
}

// PluginConfigs represents a set of [plugin.<kind>.<name>] sections in the configuration file
// under a specific <kind>. The key of the map is <name>, for example "mysql" of "plugin.metrics.mysql".
type PluginConfigs map[string]PluginConfig

// PluginConfig represents a section of [plugin.*].
// `MaxCheckAttempts`, `NotificationInterval` and `CheckInterval` options are used with check monitoring plugins. Custom metrics plugins ignore these options.
// `User` option is ignore in windows
type PluginConfig struct {
	Command              string
	User                 string
	NotificationInterval *int32  `toml:"notification_interval"`
	CheckInterval        *int32  `toml:"check_interval"`
	MaxCheckAttempts     *int32  `toml:"max_check_attempts"`
	CustomIdentifier     *string `toml:"custom_identifier"`
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

// CheckNames return list of plugin.checks._name_
func (conf *Config) CheckNames() []string {
	checks := []string{}
	for name := range conf.Plugin["checks"] {
		checks = append(checks, name)
	}
	return checks
}

// LoadConfig XXX
func LoadConfig(conffile string) (*Config, error) {
	config, err := loadConfigFile(conffile)

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

func loadConfigFile(file string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(file, config); err != nil {
		return config, err
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

		// Also, save plugin values for later merging
		pluginSaved := map[string]PluginConfigs{}
		for kind, plugins := range config.Plugin {
			pluginSaved[kind] = plugins
		}

		meta, err := toml.DecodeFile(file, &config)
		if err != nil {
			return fmt.Errorf("while loading included config file %s: %s", file, err)
		}

		// If included config does not have "roles" key,
		// use the previous roles configuration value.
		if meta.IsDefined("roles") == false {
			config.Roles = rolesSaved
		}

		for kind, plugins := range config.Plugin {
			for key, conf := range plugins {
				if pluginSaved[kind] == nil {
					pluginSaved[kind] = PluginConfigs{}
				}
				pluginSaved[kind][key] = conf
			}
		}

		config.Plugin = pluginSaved
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
