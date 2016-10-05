// +build linux darwin freebsd netbsd

package util

import (
	"testing"
	"time"
)

func init() {
	TimeoutDuration = 1 * time.Second
}

func TestRunCommand(t *testing.T) {
	stdout, stderr, exitCode, err := RunCommand("echo 1", "")
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
	stdout, stderr, _, err := RunCommand("sleep 2", "")
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

func TestSanitizeMetricKey(t *testing.T) {
	if SanitizeMetricKey("Hoge-123_") != "Hoge-123_" {
		t.Errorf("characters matching [A-Za-z0-9_-] should be kept as is")
	}
	if SanitizeMetricKey(" /pあ'*") != "__p___" {
		t.Errorf("dangerous characters should be sanitized")
	}
	if SanitizeMetricKey("p.q.r") != "p_q_r" {
		t.Errorf(". (dot) should be sanitized")
	}
}
