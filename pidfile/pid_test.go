// +build linux freebsd darwin netbsd

package pidfile

import (
	"math"
	"os"
	"testing"
)

func TestExistsPid(t *testing.T) {
	if !ExistsPid(os.Getpid()) {
		t.Errorf("something went wrong")
	}
	if ExistsPid(math.MaxInt32) {
		t.Errorf("something went wrong")
	}
}
