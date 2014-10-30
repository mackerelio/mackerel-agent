package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
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
`)
	confFile.Sync()
	confFile.Close()
	defer os.Remove(confFile.Name())

	os.Args = []string{"mackerel-agent", "-conf=" + confFile.Name(), "-role=My-Service:default", "-verbose"}
	// Overrides Args from go test command
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)

	mergedConfig, _ := resolveConfig()

	t.Logf("      apibase: %v", mergedConfig.Apibase)
	t.Logf("       apikey: %v", mergedConfig.Apikey)
	t.Logf("         root: %v", mergedConfig.Root)
	t.Logf("      pidfile: %v", mergedConfig.Pidfile)
	t.Logf("roleFullnames: %v", mergedConfig.Roles)
	t.Logf("      verbose: %v", mergedConfig.Verbose)

	if mergedConfig.Root != "/hoge/fuga" {
		t.Error("Root(confing from file) should be /hoge/fuga but: %v", mergedConfig.Root)
	}

	if len(mergedConfig.Roles) != 1 || mergedConfig.Roles[0] != "My-Service:default" {
		t.Error("Roles(config from command line option) should be parsed")
	}

	if mergedConfig.Verbose != true {
		t.Error("Verbose(overwritten by command line option) shoud be true")
	}
}

func TestParseFlagsPrintVersion(t *testing.T) {
	os.Args = []string{"mackerel-agent", "-version"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)

	config, printVersion := resolveConfig()

	if config != nil {
		t.Error("with -version args, config should be null")
	}

	if printVersion == false {
		t.Error("with -version args, printVersion should be true")
	}
}
