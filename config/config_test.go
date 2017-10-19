package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var sampleConfig = `
apikey = "abcde"
display_name = "fghij"
diagnostic = true

[filesystems]
ignore = "/dev/ram.*"

[connection]
post_metrics_retry_delay_seconds = 600
post_metrics_retry_max = 5

[plugin.metrics.mysql]
command = "ruby /path/to/your/plugin/mysql.rb"
user = "mysql"
custom_identifier = "app1.example.com"

[plugin.metrics.mysql2]
command = "ruby /path/to/your/plugin/mysql.rb"
include_pattern = '^mysql\.innodb\..+'
exclude_pattern = '^mysql\.innodb\.ignore'

[plugin.checks.heartbeat]
command = "heartbeat.sh"
user = "xyz"
notification_interval = 60
max_check_attempts = 3
action = { command = "cardiac_massage", user = "doctor" }

[plugin.metadata.hostinfo]
command = "hostinfo.sh"
user = "zzz"
execution_interval = 60
`

func TestLoadConfig(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfig)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if config.Apibase != "https://api.mackerelio.com" {
		t.Error("should be https://api.mackerelio.com (arg value should be used)")
	}

	if config.Apikey != "abcde" {
		t.Error("should be abcde (config value should be used)")
	}

	if config.DisplayName != "fghij" {
		t.Error("should be fghij (config value should be used)")
	}

	if config.Diagnostic != true {
		t.Error("should be true (config value should be used)")
	}

	if config.Filesystems.UseMountpoint != false {
		t.Error("should be false (default value should be used)")
	}

	if config.Connection.PostMetricsDequeueDelaySeconds != 30 {
		t.Error("should be 30 (default value should be used)")
	}

	if config.Connection.PostMetricsRetryDelaySeconds != 180 {
		t.Error("should be 180 (max retry delay seconds is 180)")
	}

	if config.Connection.PostMetricsRetryMax != 5 {
		t.Error("should be 5 (config value should be used)")
	}
}

var sampleConfigWithHostStatus = `
apikey = "abcde"
display_name = "fghij"

[host_status]
on_start = "working"
on_stop  = "poweroff"
`

func TestLoadConfigWithHostStatus(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfigWithHostStatus)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if config.Apikey != "abcde" {
		t.Error("should be abcde (config value should be used)")
	}

	if config.DisplayName != "fghij" {
		t.Error("should be fghij (config value should be used)")
	}

	if config.HostStatus.OnStart != "working" {
		t.Error(`HostStatus.OnStart should be "working"`)
	}

	if config.HostStatus.OnStop != "poweroff" {
		t.Error(`HostStatus.OnStop should be "poweroff"`)
	}
}

var sampleConfigWithMountPoint = `
apikey = "abcde"
display_name = "fghij"

[filesystems]
use_mountpoint = true
`

func TestLoadConfigWithMountPoint(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfigWithMountPoint)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if config.Filesystems.UseMountpoint != true {
		t.Error("should be true (config value should be used)")
	}
}

var sampleConfigWithInvalidIgnoreRegexp = `
apikey = "abcde"
display_name = "fghij"

[filesystems]
ignore = "**"
`

func TestLoadConfigWithInvalidIgnoreRegexp(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfigWithInvalidIgnoreRegexp)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = LoadConfig(tmpFile.Name())
	if err == nil {
		t.Errorf("should raise error: %v", err)
	}
}

func TestLoadConfigFile(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfig)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	config, err := loadConfigFile(tmpFile.Name())
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if config.Apikey != "abcde" {
		t.Error("Apikey should be abcde")
	}

	if config.DisplayName != "fghij" {
		t.Error("DisplayName should be fghij")
	}

	if config.Diagnostic != true {
		t.Error("Diagnostic should be true")
	}

	if config.Connection.PostMetricsRetryMax != 5 {
		t.Error("PostMetricsRetryMax should be 5")
	}

	if config.MetricPlugins == nil {
		t.Error("plugin should have metrics")
	}
	pluginConf := config.MetricPlugins["mysql"]
	if pluginConf.Command.Cmd != "ruby /path/to/your/plugin/mysql.rb" {
		t.Errorf("plugin conf command should be 'ruby /path/to/your/plugin/mysql.rb' but %v", pluginConf.Command.Cmd)
	}
	if pluginConf.Command.User != "mysql" {
		t.Error("plugin user_name should be 'mysql'")
	}
	if *pluginConf.CustomIdentifier != "app1.example.com" {
		t.Errorf("plugin custom_identifier should be 'app1.example.com' but got %v", *pluginConf.CustomIdentifier)
	}
	customIdentifiers := config.ListCustomIdentifiers()
	if len(customIdentifiers) != 1 {
		t.Errorf("config should have 1 custom_identifier")
	}
	if customIdentifiers[0] != "app1.example.com" {
		t.Errorf("first custom_identifier should be 'app1.example.com'")
	}
	if pluginConf.IncludePattern != nil {
		t.Errorf("plugin include_pattern should be nil but got %v", pluginConf.IncludePattern)
	}
	if pluginConf.ExcludePattern != nil {
		t.Errorf("plugin exclude_pattern should be nil but got %v", pluginConf.ExcludePattern)
	}

	pluginConf2 := config.MetricPlugins["mysql2"]
	if pluginConf2.IncludePattern.String() != regexp.MustCompile(`^mysql\.innodb\..+`).String() {
		t.Errorf("unexpected include_pattern: %v", pluginConf2.IncludePattern)
	}
	if pluginConf2.ExcludePattern.String() != regexp.MustCompile(`^mysql\.innodb\.ignore`).String() {
		t.Errorf("unexpected exclude_pattern: %v", pluginConf2.ExcludePattern)
	}

	if config.CheckPlugins == nil {
		t.Error("plugin should have checks")
	}
	checks := config.CheckPlugins["heartbeat"]
	if checks.Command.Cmd != "heartbeat.sh" {
		t.Error("check command should be 'heartbeat.sh'")
	}
	if checks.Command.User != "xyz" {
		t.Error("check user_name should be 'xyz'")
	}
	if *checks.NotificationInterval != 60 {
		t.Error("notification_interval should be 60")
	}
	if *checks.MaxCheckAttempts != 3 {
		t.Error("max_check_attempts should be 3")
	}
	if checks.Action.Cmd != "cardiac_massage" {
		t.Error("action.command should be 'cardiac_massage'")
	}
	if checks.Action.User != "doctor" {
		t.Error("action.user should be 'doctor'")
	}

	if config.MetadataPlugins == nil {
		t.Error("config should have metadata plugin list")
	}
	metadataPlugin := config.MetadataPlugins["hostinfo"]
	if metadataPlugin.Command.Cmd != "hostinfo.sh" {
		t.Errorf("command of metadata plugin should be 'hostinfo.sh' but got '%v'", metadataPlugin.Command.Cmd)
	}
	if metadataPlugin.Command.User != "zzz" {
		t.Errorf("user of metadata plugin should be 'zzz' but got '%v'", metadataPlugin.Command.User)
	}
	if *metadataPlugin.ExecutionInterval != 60 {
		t.Errorf("execution interval of metadata plugin should be 60 but got '%v'", *metadataPlugin.ExecutionInterval)
	}

	if config.Plugin != nil {
		t.Error("plugin config should be set nil, use MetricPlugins, CheckPlugins and MetadataPlugins instead")
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assert(t *testing.T, ok bool, msg string) {
	if !ok {
		t.Error(msg)
	}
}

var tomlQuotedReplacer = strings.NewReplacer(
	"\t", "\\t",
	"\n", "\\n",
	"\r", "\\r",
	"\"", "\\\"",
	"\\", "\\\\",
)

func TestLoadConfigFileInclude(t *testing.T) {
	configDir, err := ioutil.TempDir("", "mackerel-config-test")
	assertNoError(t, err)
	defer os.RemoveAll(configDir)

	includedFile, err := os.Create(filepath.Join(configDir, "sub1.conf"))
	assertNoError(t, err)

	configContent := fmt.Sprintf(`
apikey = "abcde"
pidfile = "/path/to/pidfile"
root = "/var/lib/mackerel-agent"
verbose = false

roles = [ "roles", "to be overwritten" ]

include = "%s/*.conf"

[plugin.metrics.foo1]
command = "foo1"

[plugin.metrics.bar]
command = "this will be overwritten"
`, tomlQuotedReplacer.Replace(configDir))

	configFile, err := newTempFileWithContent(configContent)
	assertNoError(t, err)
	defer os.Remove(configFile.Name())

	includedContent := `
roles = [ "Service:role" ]

[plugin.metrics.foo2]
command = "foo2"

[plugin.metrics.bar]
command = "bar"
`

	_, err = includedFile.WriteString(includedContent)
	assertNoError(t, err)
	includedFile.Close()

	config, err := loadConfigFile(configFile.Name())
	assertNoError(t, err)

	assert(t, config.Apikey == "abcde", "apikey should be kept as it is when not configured in the included file")
	assert(t, config.Pidfile == "/path/to/pidfile", "pidfile should be kept as it is when not configured in the included file")
	assert(t, config.Root == "/var/lib/mackerel-agent", "root should be kept as it is when not configured in the included file")
	assert(t, config.Verbose == false, "verbose should be kept as it is when not configured in the included file")
	assert(t, len(config.Roles) == 1, "roles should be overwritten")
	assert(t, config.Roles[0] == "Service:role", "roles should be overwritten")
	assert(t, config.MetricPlugins["foo1"].Command.Cmd == "foo1", "plugin.metrics.foo1 should exist")
	assert(t, config.MetricPlugins["foo2"].Command.Cmd == "foo2", "plugin.metrics.foo2 should exist")
	assert(t, config.MetricPlugins["bar"].Command.Cmd == "bar", "plugin.metrics.bar should be overwritten")
}

func TestLoadConfigFileIncludeOverwritten(t *testing.T) {
	configDir, err := ioutil.TempDir("", "mackerel-config-test")
	assertNoError(t, err)
	defer os.RemoveAll(configDir)

	includedFile, err := os.Create(filepath.Join(configDir, "sub2.conf"))
	assertNoError(t, err)

	configContent := fmt.Sprintf(`
apikey = "abcde"
pidfile = "/path/to/pidfile"
root = "/var/lib/mackerel-agent"
verbose = false

include = "%s/*.conf"
`, tomlQuotedReplacer.Replace(configDir))

	configFile, err := newTempFileWithContent(configContent)
	assertNoError(t, err)
	defer os.Remove(configFile.Name())

	includedContent := `
apikey = "new-api-key"
pidfile = "/path/to/pidfile2"
root = "/tmp"
verbose = true
`

	_, err = includedFile.WriteString(includedContent)
	assertNoError(t, err)
	includedFile.Close()

	config, err := loadConfigFile(configFile.Name())
	assertNoError(t, err)

	assert(t, config.Apikey == "new-api-key", "apikey should be overwritten")
	assert(t, config.Pidfile == "/path/to/pidfile2", "pidfile should be overwritten")
	assert(t, config.Root == "/tmp", "root should be overwritten")
	assert(t, config.Verbose == true, "verbose should be overwritten")
}

func TestFileSystemHostIDStorage(t *testing.T) {
	root, err := ioutil.TempDir("", "mackerel-agent-test")
	if err != nil {
		t.Fatal(err)
	}

	s := FileSystemHostIDStorage{Root: root}
	err = s.SaveHostID("test-host-id")
	assertNoError(t, err)

	hostID, err := s.LoadHostID()
	assertNoError(t, err)
	assert(t, hostID == "test-host-id", "SaveHostID and LoadHostID should preserve the host id")

	err = s.DeleteSavedHostID()
	assertNoError(t, err)

	_, err = s.LoadHostID()
	assert(t, err != nil, "LoadHostID after DeleteSavedHostID must fail")
}

func TestConfig_HostIDStorage(t *testing.T) {
	conf := Config{
		Root: "test-root",
	}

	storage, ok := conf.hostIDStorage().(*FileSystemHostIDStorage)
	assert(t, ok, "Default hostIDStorage must be *FileSystemHostIDStorage")
	assert(t, storage.Root == "test-root", "FileSystemHostIDStorage must have the same Root of Config")
}

func TestLoadConfigWithSilent(t *testing.T) {
	conff, err := newTempFileWithContent(`
apikey = "abcde"
silent = true
`)
	if err != nil {
		t.Fatalf("should not raise error: %s", err)
	}
	defer os.Remove(conff.Name())

	config, err := loadConfigFile(conff.Name())
	assertNoError(t, err)

	if !config.Silent {
		t.Error("silent should be ture")
	}
}

func TestLoadConfig_WithCommandArgs(t *testing.T) {
	conff, err := newTempFileWithContent(`
apikey = "abcde"
[plugin.metrics.hoge]
command = ["perl", "-E", "say 'Hello'"]
`)
	if err != nil {
		t.Fatalf("should not raise error: %s", err)
	}
	defer os.Remove(conff.Name())

	config, err := loadConfigFile(conff.Name())
	assertNoError(t, err)

	expected := []string{"perl", "-E", "say 'Hello'"}
	p := config.MetricPlugins["hoge"]
	output := p.Command.Args

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("command args not expected: %+v", output)
	}

	if p.Command.Cmd != "" {
		t.Errorf("p.Command should be empty but: %s", p.Command.Cmd)
	}
}

func newTempFileWithContent(content string) (*os.File, error) {
	tmpf, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		return nil, err
	}
	if _, err := tmpf.WriteString(content); err != nil {
		os.Remove(tmpf.Name())
		return nil, err
	}
	tmpf.Sync()
	tmpf.Close()
	return tmpf, nil
}
