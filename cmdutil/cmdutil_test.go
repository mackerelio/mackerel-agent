package cmdutil

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
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
	TimeoutDuration: 300 * time.Millisecond,
}

func TestRunCommand(t *testing.T) {
	testCases := []struct {
		Name          string
		Command       string
		CommandOption CommandOption

		Stdout, Stderr string
		ExitCode       int
		Err            error
	}{
		{
			Name:          "echo-1",
			Command:       "echo 1",
			CommandOption: testCmdOpt,
			Stdout:        "1",
		},
		{
			Name:          "Timeout",
			Command:       fmt.Sprintf("%s -sleep 11s", stubcmd),
			CommandOption: testCmdOpt,
			ExitCode: func() int {
				if runtime.GOOS == "windows" {
					return 1
				}
				return 128 + int(syscall.SIGTERM)
			}(),
			Err: errTimedOut,
		},
		{
			Name: "withEnv",
			Command: func() string {
				command := `echo $TEST_RUN_COMMAND_ENV`
				if runtime.GOOS == "windows" {
					command = `echo %TEST_RUN_COMMAND_ENV%`
				}
				return command
			}(),
			CommandOption: CommandOption{
				TimeoutDuration: testCmdOpt.TimeoutDuration,
				Env:             []string{"TEST_RUN_COMMAND_ENV=mackerel-agent"},
			},
			Stdout: "mackerel-agent",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			stdout, stderr, exitCode, err := RunCommand(tc.Command, tc.CommandOption)
			stdout = strings.TrimSpace(stdout)
			stderr = strings.TrimSpace(stderr)
			if stdout != tc.Stdout {
				t.Errorf("invalid stdout. out=%q, expect=%q", stdout, tc.Stdout)
			}
			if stderr != tc.Stderr {
				t.Errorf("invalid stderr. out=%q, expect=%q", stderr, tc.Stderr)
			}
			if exitCode != tc.ExitCode {
				t.Errorf("exitCode should be %d, but: %d", tc.ExitCode, exitCode)
			}
			if err != tc.Err {
				t.Errorf("err should be %v but: %v", tc.Err, err)
			}
		})
	}
}
