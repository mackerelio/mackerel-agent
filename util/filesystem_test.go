// +build linux darwin freebsd netbsd

package util

import (
	"reflect"
	"testing"
)

func TestParseDfLines(t *testing.T) {
	dfout := `Filesystem     1024-blocks     Used Available Capacity Mounted on
/dev/sda1           19734388 16868164 1863772  91% /
tmpfs                 517224        0  517224   0% /lib/init/rw
udev                  512780       96  512684   1% /dev
tmpfs                 517224        4  517220   1% /dev/shm
`
	expect := []*DfStat{
		{
			Name:      "/dev/sda1",
			Blocks:    19734388,
			Used:      16868164,
			Available: 1863772,
			Capacity:  91,
			Mounted:   "/",
		},
		{
			Name:      "tmpfs",
			Blocks:    517224,
			Used:      0,
			Available: 517224,
			Capacity:  0,
			Mounted:   "/lib/init/rw",
		},
		{
			Name:      "udev",
			Blocks:    512780,
			Used:      96,
			Available: 512684,
			Capacity:  1,
			Mounted:   "/dev",
		},
		{
			Name:      "tmpfs",
			Blocks:    517224,
			Used:      4,
			Available: 517220,
			Capacity:  1,
			Mounted:   "/dev/shm",
		},
	}
	ret := parseDfLines(dfout)
	if !reflect.DeepEqual(ret, expect) {
		t.Errorf("dfvalues are not expected: %#v", ret)
	}
}
