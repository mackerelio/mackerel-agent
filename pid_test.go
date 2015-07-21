// +build linux freebsd darwin

package main

import (
	"math"
	"os"
	"testing"
)

func TestExistsPid(t *testing.T) {
	if !existsPid(os.Getpid()) {
		t.Errorf("something went wrong")
	}
	if existsPid(math.MaxInt32) {
		t.Errorf("something went wrong")
	}
}
