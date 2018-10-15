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
	"time"
)

var sampleConfig = `
apikey = "abcde"
display_name = "fghij"
diagnostic = true

[filesystems]
ignore = "/dev/ram.*"

[plugin.metrics.mysql]
command = "ruby /path/to/your/plugin/mysql.rb"
user = "mysql"
custom_identifier = "app1.example.com"
timeout_seconds = 60

[plugin.metrics.mysql2]
command = "ruby /path/to/your/plugin/mysql.rb"
include_pattern = '^mysql\.innodb\..+'
exclude_pattern = '^mysql\.innodb\.ignore'

[plugin.metrics.mysql3]
command = "ruby /path/to/your/plugin/mysql.rb"
env = { "MYSQL_USERNAME" = "USERNAME", "MYSQL_PASSWORD" = "PASSWORD" }

[plugin.checks.heartbeat]
command = "heartbeat.sh"
user = "xyz"
notification_interval = 60
max_check_attempts = 3
timeout_seconds = 60
action = { command = "cardiac_massage", user = "doctor" }

[plugin.checks.heartbeat2]
command = "heartbeat.sh"
env = { "ES_HOSTS" = "10.45.3.2:9220,10.45.3.1:9230" }
action = { command = "cardiac_massage", user = "doctor", env = { "NAME_1" = "VALUE_1", "NAME_2" = "VALUE_2", "NAME_3" = "VALUE_3" } }

[plugin.checks.heartbeat3]
command = "heartbeat.sh"

[plugin.metadata.hostinfo]
command = "hostinfo.sh"
user = "zzz"
execution_interval = 60
timeout_seconds = 60

[plugin.metadata.hostinfo2]
command = "hostinfo.sh"
env = { "NAME_1" = "VALUE_1" }
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

var sampleConfigWithInvalidMetricsCommand = `
apikey = "abcde"

[plugin.metrics.mysql]
command = 100
user = "mysql"
`

func TestLoadConfigWithInvalidMetricsCommand(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfigWithInvalidMetricsCommand)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = LoadConfig(tmpFile.Name())
	if err == nil {
		t.Errorf("should raise error: %v", err)
	}
	if !strings.Contains(err.Error(), "should be string or string slice, but int64") {
		t.Errorf("should raise error containing type information: %v", err)
	}
	if !strings.Contains(err.Error(), "plugin.metrics.mysql") {
		t.Errorf("should raise error containing metrics key: %v", err)
	}
}

var sampleConfigWithInvalidCheckCommand = `
apikey = "abcde"

[plugin.checks.dice]
`

func TestLoadConfigWithInvalidCheckCommand(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfigWithInvalidCheckCommand)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = LoadConfig(tmpFile.Name())
	if err == nil {
		t.Errorf("should raise error: %v", err)
	}
	if !strings.Contains(err.Error(), "should be string or string slice, but <nil>") {
		t.Errorf("should raise error containing type information: %v", err)
	}
	if !strings.Contains(err.Error(), "plugin.checks.dice") {
		t.Errorf("should raise error containing metrics key: %v", err)
	}
}

var sampleConfigWithTooLargeCheckMemo = `
apikey = "abcde"

[plugin.checks.memo]
command = "memo"
memo = "あいうえお"

[plugin.checks.toolargememo]
command = "toolargememo"
memo = "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"

[plugin.checks.toolargememo2]
command = "toolargememo"
memo = "012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678あいう"
`

func TestLoadConfigWithTooLargeCheckMemo(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfigWithTooLargeCheckMemo)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	config, err := LoadConfig(tmpFile.Name())

	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	check1 := config.CheckPlugins["memo"]
	if check1.Memo != "あいうえお" {
		t.Errorf("check command should be 'あいうえお': %v", check1.Memo)
	}

	check2 := config.CheckPlugins["toolargememo"]
	if check2.Memo != "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" {
		t.Errorf("check command should have starting 250 charcters: %v", check2.Memo)
	}

	check3 := config.CheckPlugins["toolargememo2"]
	if check3.Memo != "012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678あ" {
		t.Errorf("check command should have starting 250 charcters: %v", check3.Memo)
	}
}

var sampleConfigWithInvalidMetadataCommand = `
apikey = "abcde"

[plugin.metadata.sample]
command = [ 10, 20, 30 ]
`

func TestLoadConfigWithInvalidMetadataCommand(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfigWithInvalidMetadataCommand)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = LoadConfig(tmpFile.Name())
	if err == nil {
		t.Errorf("should raise error: %v", err)
	}
	if !strings.Contains(err.Error(), "should be string or string slice, but []interface {}") {
		t.Errorf("should raise error containing type information: %v", err)
	}
	if !strings.Contains(err.Error(), "plugin.metadata.sample") {
		t.Errorf("should raise error containing metrics key: %v", err)
	}
}

var sampleConfigWithCloudPlatformTemplate = `
apikey = "abcde"
cloud_platform = "%s"
`

var LoadConfigWithCloudPlatformTests = []struct {
	value    string
	expected CloudPlatform
}{
	{"", CloudPlatformAuto},
	{"auto", CloudPlatformAuto},
	{"none", CloudPlatformNone},
	{"ec2", CloudPlatformEC2},
	{"gce", CloudPlatformGCE},
	{"azurevm", CloudPlatformAzureVM},
}

func TestLoadConfigWithCloudPlatform(t *testing.T) {
	for _, test := range LoadConfigWithCloudPlatformTests {
		content := fmt.Sprintf(sampleConfigWithCloudPlatformTemplate, test.value)
		tmpFile, err := newTempFileWithContent(content)
		if err != nil {
			t.Errorf("should not raise error: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		config, err := LoadConfig(tmpFile.Name())
		if err != nil {
			t.Errorf("should not raise error: %v", err)
		}

		if config.CloudPlatform != test.expected {
			t.Errorf("CloudPlatform should be set to %s, but %s", test.expected, config.CloudPlatform)
		}
	}
}

var sampleConfigWithInvalidCloudPlatform = `
apikey = "abcde"
cloud_platform = "unknown"
`

func TestLoadConfigWithInvalidCloudPlatform(t *testing.T) {
	tmpFile, err := newTempFileWithContent(sampleConfigWithInvalidCloudPlatform)
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
	if pluginConf.Command.TimeoutDuration != 60*time.Second {
		t.Error("plugin timeout_seconds should be 60s")
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

	pluginConf3 := config.MetricPlugins["mysql3"]
	if pluginConf3.Command.Env == nil {
		t.Error("config should have env")
	}
	if len(pluginConf3.Command.Env) != 2 {
		t.Errorf("env should have 2 keys: %v", pluginConf3.Command.Env)
	}
	if !expectContainsString(pluginConf3.Command.Env, "MYSQL_USERNAME=USERNAME") {
		t.Errorf("Command.Env should contain 'MYSQL_USERNAME=USERNAME'")
	}
	if !expectContainsString(pluginConf3.Command.Env, "MYSQL_PASSWORD=PASSWORD") {
		t.Errorf("Command.Env should contain 'MYSQL_PASSWORD=PASSWORD'")
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
	if checks.Command.TimeoutDuration != 60*time.Second {
		t.Error("check timeout_seconds should be 60s")
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
	if expected := ""; checks.Memo != expected {
		t.Errorf("memo should be %q but got %q", expected, checks.Memo)
	}

	checks2 := config.CheckPlugins["heartbeat2"]
	if checks2.Command.Env == nil {
		t.Error("config should have env of check plugin")
	}
	if len(checks2.Command.Env) != 1 {
		t.Errorf("env of check plugin should have a key: %v", checks2.Command.Env)
	}
	if !expectContainsString(checks2.Command.Env, "ES_HOSTS=10.45.3.2:9220,10.45.3.1:9230") {
		t.Errorf("Command.Env should contain 'ES_HOSTS=10.45.3.2:9220,10.45.3.1:9230'")
	}
	if checks2.Action.Env == nil {
		t.Error("config should have action.env of check plugin")
	}
	if len(checks2.Action.Env) != 3 {
		t.Errorf("action.env of check plugin should have 3 keys: %v", checks2.Action.Env)
	}
	if !expectContainsString(checks2.Action.Env, "NAME_1=VALUE_1") {
		t.Errorf("Command.Env should contain 'NAME_1=VALUE_1'")
	}
	if !expectContainsString(checks2.Action.Env, "NAME_2=VALUE_2") {
		t.Errorf("Command.Env should contain 'NAME_2=VALUE_2'")
	}

	checks3 := config.CheckPlugins["heartbeat3"]
	if checks3.Action != nil {
		t.Error("config should not have action of check plugin")
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
	if metadataPlugin.Command.TimeoutDuration != 60*time.Second {
		t.Errorf("timeout duration of metadata plugin should be 60s, but got '%v'",
			metadataPlugin.Command.TimeoutDuration)
	}

	metadataPlugin2 := config.MetadataPlugins["hostinfo2"]
	if metadataPlugin2.Command.Env == nil {
		t.Error("config should have env of metadata plugin")
	}
	if len(metadataPlugin2.Command.Env) != 1 {
		t.Errorf("env of metadata plugin should have a key: %v", metadataPlugin2.Command.Env)
	}
	if !expectContainsString(metadataPlugin2.Command.Env, "NAME_1=VALUE_1") {
		t.Errorf("Command.Env should contain 'NAME_1=VALUE_1'")
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

	// Write an empty id to simulate a case that could not save id properly
	err = s.SaveHostID("")
	assertNoError(t, err)

	_, err = s.LoadHostID()
	assert(t, err != nil, "LoadHostID from empty HostID file must fail")
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

func TestEnv_ConvertToStrings(t *testing.T) {
	cases := []struct {
		env         Env
		expected    []string
		expectError bool
	}{
		{Env{}, []string{}, false},
		{Env{"KEY": "VALUE"}, []string{"KEY=VALUE"}, false},
		{Env{"KEY1": "VALUE1", "KEY2": "VALUE2", "KEY3": "VALUE3"}, []string{"KEY1=VALUE1", "KEY2=VALUE2", "KEY3=VALUE3"}, false},
		{Env{"KEY1": "VALUE1 VALUE2 VALUE3", "KEY2": "VALUE4 VALUE5 VALUE6"}, []string{"KEY1=VALUE1 VALUE2 VALUE3", "KEY2=VALUE4 VALUE5 VALUE6"}, false},
		{Env{"KEY": ""}, []string{"KEY="}, false},
		{Env{"   KEY   ": "   VALUE   "}, []string{"KEY=   VALUE   "}, false},
		{Env{"": ""}, []string{}, false},
		{Env{"KEY=KEY": "VALUE"}, nil, true},
	}

	for _, c := range cases {
		got, err := c.env.ConvertToStrings()
		if err != nil && c.expectError == false {
			t.Errorf("should raise error: %v", c.env)
		}
		if len(got) != len(c.expected) {
			t.Errorf("env strings should contains %d keys but: %d", len(c.expected), len(got))
		}
		for _, v := range got {
			if !expectContainsString(c.expected, v) {
				t.Errorf("env strings not expected %+v", got)
			}
		}
	}
}

func TestCommandRunWithEnv(t *testing.T) {
	tmpf, err := newTempFileWithContent(`
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print(os.Getenv("SAMPLE_KEY1"), os.Getenv("SAMPLE_KEY2"))
}`)
	if err != nil {
		os.Remove(tmpf.Name())
		t.Fatalf("should not raise error: %s", err)
	}

	gof := tmpf.Name() + ".go"

	err = os.Rename(tmpf.Name(), gof)

	if err != nil {
		os.Remove(tmpf.Name())
		t.Fatalf("should not raise error: %s", err)
	}
	defer os.Remove(gof)

	conf := fmt.Sprintf(`
apikey = "abcde"

[plugin.metrics.sample]
command = ["go", "run", '%s']
env = { "SAMPLE_KEY1" = " foo bar ", "SAMPLE_KEY2" = " baz qux " }
`, gof)

	cases := []struct {
		env      []string
		expected string
	}{
		{nil, " foo bar  baz qux "},
		{[]string{"SAMPLE_KEY1=v1 v2", "SAMPLE_KEY2= v3 v4"}, "v1 v2 v3 v4"},
	}

	for _, c := range cases {
		var stdout, stderr string
		var exitCode int
		var err error

		conff, err := newTempFileWithContent(conf)

		if err != nil {
			t.Fatalf("should not raise error: %s", err)
		}
		defer os.Remove(conff.Name())

		config, err := loadConfigFile(conff.Name())
		assertNoError(t, err)

		p := config.MetricPlugins["sample"]

		if c.env != nil {
			stdout, stderr, exitCode, err = p.Command.RunWithEnv(c.env)
		} else {
			stdout, stderr, exitCode, err = p.Command.Run()
		}

		if stdout != c.expected {
			t.Errorf("stdout not expected: %+v", stdout)
		}

		if stderr != "" {
			t.Errorf("stderr not expected: %+v", stderr)
		}

		if exitCode != 0 {
			t.Errorf("exitCode not expected: %+v", exitCode)
		}

		if err != nil {
			t.Errorf("err not expected: %+v", err)
		}
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

func expectContainsString(slice []string, contains string) bool {
	for _, v := range slice {
		if v == contains {
			return true
		}
	}
	return false
}
