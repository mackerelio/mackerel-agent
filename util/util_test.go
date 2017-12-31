// +build linux darwin freebsd netbsd windows

package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func init() {
	TimeoutDuration = 1 * time.Second
}

func TestRunCommand(t *testing.T) {
	stdout, stderr, exitCode, err := RunCommand("echo 1", "", nil)
	if runtime.GOOS == "windows" {
		stdout = strings.Replace(stdout, "\r\n", "\n", -1)
		stderr = strings.Replace(stderr, "\r\n", "\n", -1)
	}
	if stdout != "1\n" {
		t.Errorf("stdout shoud be 1")
	}
	if stderr != "" {
		t.Errorf("stderr shoud be empty")
	}
	if exitCode != 0 {
		t.Errorf("exitCode should be zero")
	}
	if err != nil {
		t.Error("err should be nil but:", err)
	}
}

func makeSleep(t *testing.T) (string, string) {
	tmpdir, err := ioutil.TempDir("", "mackerel-agent")
	if err != nil {
		t.Fatal(err)
	}

	f := filepath.Join(tmpdir, "sleep.go")
	err = ioutil.WriteFile(f, []byte(`package main;import "time";func main(){time.Sleep(time.Second*2)}`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(f, " ") {
		f = `"` + f + `"`
	}
	return fmt.Sprintf(`go run %s 2`, f), tmpdir
}

func TestRunCommandWithTimeout(t *testing.T) {
	command := "sleep 2"

	var tmpdir string
	if runtime.GOOS == "windows" {
		command, tmpdir = makeSleep(t)
	}
	stdout, stderr, _, err := RunCommand(command, "", nil)
	if stdout != "" {
		t.Errorf("stdout shoud be empty")
	}
	if stderr != "" {
		t.Errorf("stderr shoud be empty")
	}
	if err == nil {
		t.Error("err should have error but nil")
	}
	if tmpdir != "" {
		os.RemoveAll(tmpdir)
	}
}

func TestRunCommandWithEnv(t *testing.T) {
	command := `echo $TEST_RUN_COMMAND_ENV`
	if runtime.GOOS == "windows" {
		command = `echo %TEST_RUN_COMMAND_ENV%`
	}

	stdout, stderr, exitCode, err := RunCommand(command, "", []string{"TEST_RUN_COMMAND_ENV=mackerel-agent"})
	if runtime.GOOS == "windows" {
		stdout = strings.Replace(stdout, "\r\n", "\n", -1)
		stderr = strings.Replace(stderr, "\r\n", "\n", -1)
	}
	if stdout != "mackerel-agent\n" {
		t.Errorf("stdout shoud be 1")
	}
	if stderr != "" {
		t.Errorf("stderr shoud be empty")
	}
	if exitCode != 0 {
		t.Errorf("exitCode should be zero")
	}
	if err != nil {
		t.Error("err should be nil but:", err)
	}
}
