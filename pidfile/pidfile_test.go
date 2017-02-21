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

	dir, err := ioutil.TempDir("", "")
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
	if err == nil {
		t.Errorf("Successfully overwriting the pidfile unintentionally")
	}
}

func TestRemove(t *testing.T) {
	err := Remove("")
	if err != nil {
		t.Errorf("err should be nil but: %v", err)
	}
	dir, err := ioutil.TempDir("", "")
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
