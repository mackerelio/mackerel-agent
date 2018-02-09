package cmdutil

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

var stubcmd = "testdata/stubcmd"

func init() {
	if runtime.GOOS == "windows" {
		stubcmd = `.\testdata\stubcmd.exe`
	}
	err := exec.Command("go", "build", "-o", stubcmd, "testdata/stubcmd.go").Run()
	if err != nil {
		panic(err)
	}
}

var testCmdOpt = CommandOption{
	TimeoutDuration: 1 * time.Second,
}

func TestRunCommand(t *testing.T) {
	stdout, stderr, exitCode, err := RunCommand("echo 1", testCmdOpt)
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

func TestRunCommandWithTimeout(t *testing.T) {
	command := fmt.Sprintf("%s -sleep 2", stubcmd)

	stdout, stderr, _, err := RunCommand(command, testCmdOpt)
	if stdout != "" {
		t.Errorf("stdout shoud be empty")
	}
	if stderr != "" {
		t.Errorf("stderr shoud be empty")
	}
	if err == nil {
		t.Error("err should have error but nil")
	}
}

func TestRunCommandWithEnv(t *testing.T) {
	command := `echo $TEST_RUN_COMMAND_ENV`
	if runtime.GOOS == "windows" {
		command = `echo %TEST_RUN_COMMAND_ENV%`
	}

	opt := CommandOption{
		TimeoutDuration: testCmdOpt.TimeoutDuration,
		Env:             []string{"TEST_RUN_COMMAND_ENV=mackerel-agent"},
	}
	stdout, stderr, exitCode, err := RunCommand(command, opt)
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
