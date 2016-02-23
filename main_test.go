package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/command"
)

func TestParseFlags(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")

	if err != nil {
		t.Fatalf("Could not create temprary config file for test")
	}
	confFile.WriteString(`verbose=false
root="/hoge/fuga"
apikey="DUMMYAPIKEY"
diagnostic=false
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())
	mergedConfig, _ := resolveConfig([]string{"-conf=" + confFile.Name(), "-role=My-Service:default,INVALID#SERVICE", "-verbose", "-diagnostic"})

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

func TestDetectForce(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temprary config file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	argv := []string{"-conf=" + confFile.Name()}
	conf, force := resolveConfigForRetire(argv)
	if force {
		t.Errorf("force should be false")
	}
	if conf.Apikey != "DUMMYAPIKEY" {
		t.Errorf("Apikey should be 'DUMMYAPIKEY'")
	}

	argv = append(argv, "-force")
	conf, force = resolveConfigForRetire(argv)
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
		t.Fatalf("Could not create temprary config file for test")
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
		"-apibase=https://mackerel.io",
		"-pidfile=hoge",
		"-root=hoge",
		"-verbose",
		"-diagnostic",
		"-apikey=hogege",
		"-role=hoge:fuga",
	}

	conf, force := resolveConfigForRetire(argv)
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

	err = createPidFile(fpath)
	if err != nil {
		t.Errorf("pid file should be created but, %s", err)
	}

	if runtime.GOOS != "windows" {
		if err := createPidFile(fpath); err == nil || !strings.HasPrefix(err.Error(), "Pidfile found, try stopping another running mackerel-agent or delete") {
			t.Errorf("creating pid file should be failed when the running process exists, %s", err)
		}
	}

	removePidFile(fpath)
	if err := createPidFile(fpath); err != nil {
		t.Errorf("pid file should be created but, %s", err)
	}

	removePidFile(fpath)
	ioutil.WriteFile(fpath, []byte(fmt.Sprint(math.MaxInt32)), 0644)
	if err := createPidFile(fpath); err != nil {
		t.Errorf("old pid file should be ignored and new pid file should be created but, %s", err)
	}
}

func TestSignalHandler(t *testing.T) {
	ctx := &command.Context{}
	termCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go signalHandler(c, ctx, termCh)

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

func TestConfigTestOK(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temprary config file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	argv := []string{"-conf=" + confFile.Name()}
	status := doConfigtest(argv)

	if status != exitStatusOK {
		t.Errorf("configtest(ok) must be return exitStatusOK")
	}
}

func TestConfigTestNotFound(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temprary config file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	argv := []string{"-conf=" + confFile.Name() + "xxx"}
	status := doConfigtest(argv)

	if status != exitStatusError {
		t.Errorf("configtest(failed) must be return existStatusError")
	}
}

func TestConfigTestInvalidFormat(t *testing.T) {
	// prepare dummy config
	confFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temprary config file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"
[plugin.checks.foo ]
command = "bar"
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	argv := []string{"-conf=" + confFile.Name()}
	status := doConfigtest(argv)

	if status != exitStatusError {
		t.Errorf("configtest(failed) must be return exitStatusError")
	}
}

func TestDoOnce(t *testing.T) {
	r := doOnce([]string{})
	if r != exitStatusOK {
		t.Errorf("doOnce should return exitStatusOK even if argv is empty, but returns %d", r)
	}
}

func TestDoVersion(t *testing.T) {
	r := doVersion([]string{})
	if r != exitStatusOK {
		t.Errorf("doVersion should return exitStatusOK, but returns %d", r)
	}
}
