// +build !windows

package cmdutil

import (
	"strings"
	"testing"
)

func TestRunCommandArgs(t *testing.T) {
	testCases := []struct {
		Name          string
		CommandArgs   []string
		CommandOption CommandOption

		Stdout, Stderr string
		ExitCode       int
		Err            string
	}{
		{
			Name:          "signal trapped",
			CommandArgs:   []string{stubcmd, "-trap=SIGTERM", "-trap-exit=23", "-sleep=10s"},
			CommandOption: testCmdOpt,
			Stdout:        "signal received",
			ExitCode:      23,
		},
		{
			Name:        "command not found",
			CommandArgs: []string{"notfound"},
			ExitCode:    127,
			Err:         `exec: "notfound": executable file not found`,
		},
		{
			Name:        "directory",
			CommandArgs: []string{"./testdata"},
			ExitCode:    126,
			Err:         `exit code: 126`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			stdout, stderr, exitCode, err := RunCommandArgs(tc.CommandArgs, tc.CommandOption)
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
			if tc.Err == "" {
				if err != nil {
					t.Errorf("error should be nil, but: %s", err)
				}
			} else {
				if err == nil {
					t.Error("error should be occurred but nil")
				} else if !strings.Contains(err.Error(), tc.Err) {
					t.Errorf("error should be contained string %q, but: %q", tc.Err, err)
				}
			}
		})
	}
}
