// +build !windows

package pidfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestCreate(t *testing.T) {
	err := Create("")
	if err != nil {
		t.Errorf("err should be nil but: %v", err)
	}

	dir, err := ioutil.TempDir("", "mackerel-agent-test-pidfile")
	if err != nil {
		t.Fatalf("failed to create tempdir")
	}
	defer os.RemoveAll(dir)
	pidfile := filepath.Join(dir, "pidfile")
	err = Create(pidfile)
	if err != nil {
		t.Errorf("err should be nil but: %v", err)
	}
	pidString, _ := ioutil.ReadFile(pidfile)
	pid, _ := strconv.Atoi(string(pidString))
	if pid != os.Getpid() {
		t.Errorf("contents of pidfile does not match pid. content: %d, pid: %d", pid, os.Getpid())
	}

	err = Create(pidfile)
	if err != nil {
		t.Errorf("When the content of pidfile is the same as own pid, the error should be nil, but: %s", err.Error())
	}
}

func TestRemove(t *testing.T) {
	err := Remove("")
	if err != nil {
		t.Errorf("err should be nil but: %v", err)
	}
	dir, err := ioutil.TempDir("", "mackerel-agent-test-pidfile")
	if err != nil {
		t.Fatalf("failed to create tempdir")
	}
	defer os.RemoveAll(dir)
	pidfile := filepath.Join(dir, "pidfile")
	err = Create(pidfile)
	if err != nil {
		t.Errorf("err should be nil but: %v", err)
	}

	err = Remove(pidfile)
	if err != nil {
		t.Errorf("err should be nil but: %v", err)
	}
}
