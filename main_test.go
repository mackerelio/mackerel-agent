package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/pidfile"
)

func TestParseFlags(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")

	if err != nil {
		t.Fatalf("Could not create temporary config file for test")
	}
	confFile.WriteString(`verbose=false
root="/hoge/fuga"
apikey="DUMMYAPIKEY"
diagnostic=false
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())
	mergedConfig, _ := resolveConfig(&flag.FlagSet{}, []string{"-conf=" + confFile.Name(), "-role=My-Service:default,INVALID#SERVICE", "-verbose", "-diagnostic"})

	t.Logf("      apibase: %v", mergedConfig.Apibase)
	t.Logf("       apikey: %v", mergedConfig.Apikey)
	t.Logf("         root: %v", mergedConfig.Root)
	t.Logf("      pidfile: %v", mergedConfig.Pidfile)
	t.Logf("   diagnostic: %v", mergedConfig.Diagnostic)
	t.Logf("roleFullnames: %v", mergedConfig.Roles)
	t.Logf("      verbose: %v", mergedConfig.Verbose)

	if mergedConfig.Root != "/hoge/fuga" {
		t.Errorf("Root(confing from file) should be /hoge/fuga but: %v", mergedConfig.Root)
	}

	if len(mergedConfig.Roles) != 1 || mergedConfig.Roles[0] != "My-Service:default" {
		t.Error("Roles(config from command line option) should be parsed")
	}

	if mergedConfig.Verbose != true {
		t.Error("Verbose(overwritten by command line option) shoud be true")
	}

	if mergedConfig.Diagnostic != true {
		t.Error("Diagnostic(overwritten by command line option) shoud be true")
	}
}

func TestParseFlagsFallback(t *testing.T) {
	tests := []struct {
		Name           string
		Args           []string
		Dir            string
		CrunchDefaults bool // will overwrite config.DefaultConfig.Xxxfile

		Pidfile  string
		Conffile string
		Root     string
	}{
		{
			Name: "default settings",
			Args: []string{},
			Dir:  "",

			Conffile: "testdata/default/mackerel-agent.conf",
			Pidfile:  "testdata/default/pid",
			Root:     "testdata/default",
		},
		{
			Name: "overwritten by config(exist)",
			Args: []string{"-conf", "testdata/case1/mackerel-agent.conf"},
			Dir:  "",

			Conffile: "testdata/case1/mackerel-agent.conf",
			Pidfile:  "testdata/case1/pid",
			Root:     "testdata/case1",
		},
		{
			Name: "overwritten by config(not exist)",
			Args: []string{"-conf", "testdata/case2/mackerel-agent.conf"},
			Dir:  "",

			Conffile: "testdata/case2/mackerel-agent.conf",
			Pidfile:  "testdata/case2x/pid",
			Root:     "testdata/case2x",
		},
		{
			Name: "overwritten by options",
			Args: []string{
				"-conf", "testdata/case2/mackerel-agent.conf",
				"-pidfile", "testdata/case2/pid",
				"-root", "testdata/case2",
			},
			Dir: "",

			Conffile: "testdata/case2/mackerel-agent.conf",
			Pidfile:  "testdata/case2/pid",
			Root:     "testdata/case2",
		},
		{
			Name:           "overwritten by fallback",
			Args:           []string{},
			Dir:            "testdata/case3",
			CrunchDefaults: true,

			Conffile: "testdata/case3/mackerel-agent.conf",
			Pidfile:  "testdata/case3/pid",
			Root:     "testdata/case3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Dir == "" {
				os.Unsetenv("MACKEREL_CONFIG_FALLBACK")
			} else {
				os.Setenv("MACKEREL_CONFIG_FALLBACK", tt.Dir)
			}
			defaultConfig := *config.DefaultConfig
			if tt.CrunchDefaults {
				config.DefaultConfig.Conffile = "blahblah"
				config.DefaultConfig.Pidfile = "blahblah"
				config.DefaultConfig.Root = "blahblah"
			} else {
				config.DefaultConfig.Conffile = "testdata/default/mackerel-agent.conf"
				config.DefaultConfig.Pidfile = "testdata/default/pid"
				config.DefaultConfig.Root = "testdata/default"
			}
			var fs flag.FlagSet
			c, err := resolveConfig(&fs, tt.Args)
			*config.DefaultConfig = defaultConfig
			if err != nil {
				t.Fatal(err)
			}
			if s := filepath.ToSlash(c.Conffile); s != tt.Conffile {
				t.Errorf("Conffile = %s; want %s", s, tt.Conffile)
			}
			if s := filepath.ToSlash(c.Pidfile); s != tt.Pidfile {
				t.Errorf("Pidfile = %s; want %s", s, tt.Pidfile)
			}
			if s := filepath.ToSlash(c.Root); s != tt.Root {
				t.Errorf("Root = %s; want %s", s, tt.Root)
			}
		})
	}
}

func TestDetectForce(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temporary config file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	argv := []string{"-conf=" + confFile.Name()}
	conf, force, _ := resolveConfigForRetire(&flag.FlagSet{}, argv)
	if force {
		t.Errorf("force should be false")
	}
	if conf.Apikey != "DUMMYAPIKEY" {
		t.Errorf("Apikey should be 'DUMMYAPIKEY'")
	}

	argv = append(argv, "-force")
	conf, force, _ = resolveConfigForRetire(&flag.FlagSet{}, argv)
	if !force {
		t.Errorf("force should be true")
	}
	if conf.Apikey != "DUMMYAPIKEY" {
		t.Errorf("Apikey should be 'DUMMYAPIKEY'")
	}
}

func TestResolveConfigForRetire(t *testing.T) {
	confFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temporary config file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	// Allow accepting unnecessary options, pidfile, diagnostic and role.
	// Because, these options are potentially passed in initd script by using $OTHER_OPTS.
	argv := []string{
		"-conf=" + confFile.Name(),
		"-apibase=https://api.mackerelio.com",
		"-pidfile=hoge",
		"-root=hoge",
		"-verbose",
		"-diagnostic",
		"-apikey=hogege",
		"-role=hoge:fuga",
	}

	conf, force, _ := resolveConfigForRetire(&flag.FlagSet{}, argv)
	if force {
		t.Errorf("force should be false")
	}
	if conf.Apikey != "hogege" {
		t.Errorf("Apikey should be 'hogege'")
	}
}

func TestCreateAndRemovePidFile(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Errorf("failed to create tmpfile, %s", err)
	}
	fpath := file.Name()
	defer os.Remove(fpath)

	err = pidfile.Create(fpath)
	if err != nil {
		t.Errorf("pid file should be created but, %s", err)
	}

	pidfile.Remove(fpath)
	if err := pidfile.Create(fpath); err != nil {
		t.Errorf("pid file should be created but, %s", err)
	}

	pidfile.Remove(fpath)
	ioutil.WriteFile(fpath, []byte(fmt.Sprint(math.MaxInt32)), 0644)
	if err := pidfile.Create(fpath); err != nil {
		t.Errorf("old pid file should be ignored and new pid file should be created but, %s", err)
	}
}

func TestSignalHandler(t *testing.T) {
	app := &command.App{}
	termCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go signalHandler(c, app, termCh)

	resultCh := make(chan int)

	maxTerminatingInterval = 100 * time.Millisecond
	c <- os.Interrupt
	c <- os.Interrupt

	go func() {
		<-termCh
		<-termCh
		<-termCh
		<-termCh
		resultCh <- 0
	}()

	go func() {
		time.Sleep(time.Second)
		resultCh <- 1
	}()

	if r := <-resultCh; r != 0 {
		t.Errorf("Something went wrong")
	}
}

func TestNotifyUpdateFile(t *testing.T) {
	app := &command.App{}
	termCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go signalHandler(c, app, termCh)

	file := "testdata/fake-agent"
	interval := 100 * time.Millisecond
	go notifyUpdateFile(c, file, interval)
	time.Sleep(interval)
	os.Chtimes(file, time.Now(), time.Now())
	select {
	case <-termCh:
	case <-time.After(time.Second):
		t.Errorf("Interrupt signal is not received in a second")
	}
}

func TestNotifyUpdateFileDelete(t *testing.T) {
	app := &command.App{}
	termCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go signalHandler(c, app, termCh)

	f, err := ioutil.TempFile("", "mackerel-agent.test.*")
	if err != nil {
		t.Fatalf("can't create a temporary file: %v", err)
	}
	file := f.Name()
	f.Close()

	interval := 100 * time.Millisecond
	go notifyUpdateFile(c, file, interval)
	time.Sleep(interval)
	if err := os.Remove(file); err != nil {
		t.Fatalf("can't remove %v: %v", file, err)
	}
	select {
	case <-termCh:
	case <-time.After(time.Second):
		t.Errorf("Interrupt signal is not received in a second")
	}
}

func TestConfigTestOK(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temporary config file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	argv := []string{"-conf=" + confFile.Name()}
	err = doConfigtest(&flag.FlagSet{}, argv)

	if err != nil {
		t.Errorf("configtest(ok) must be return nil")
	}
}

func TestConfigTestNotFound(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temporary config file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	argv := []string{"-conf=" + confFile.Name() + "xxx"}
	err = doConfigtest(&flag.FlagSet{}, argv)

	if err == nil {
		t.Errorf("configtest(failed) must be return error")
	}
}

func TestConfigTestInvalidFormat(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temporary config file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"
invalid!!!
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	argv := []string{"-conf=" + confFile.Name()}
	err = doConfigtest(&flag.FlagSet{}, argv)

	if err == nil {
		t.Errorf("configtest(failed) must be return error")
	}
}

func TestDoOnce(t *testing.T) {
	err := doOnce(&flag.FlagSet{}, []string{})
	if err != nil {
		t.Errorf("doOnce should return nil even if argv is empty, but returns %s", err)
	}
}

func TestDoVersion(t *testing.T) {
	err := doVersion(&flag.FlagSet{}, []string{})
	if err != nil {
		t.Errorf("doVersion should return nil, but returns %s", err)
	}
}
