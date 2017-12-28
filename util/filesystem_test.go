// +build linux darwin freebsd netbsd

package util

import (
	"reflect"
	"testing"
)

func TestCollectDfValues(t *testing.T) {
	_, err := CollectDfValues()
	if err != nil {
		t.Errorf("err should be nil but: %s", err)
	}
}

func TestParseDfLines(t *testing.T) {
	dfout := `Filesystem                         1024-blocks     Used Available Capacity Mounted on
/dev/sda1                             19734388 16868164 1863772        91% /
tmpfs                                   517224        0  517224         0% /lib/init/rw
udev                                    512780       96  512684         1% /dev
tmpfs                                   517224        4  517220         1% /dev/shm
/dev/mapper/docker-000:0-000-00000    10190136   168708 9480756         2% /var/lib/docker/devicemapper/mnt/00000
/dev/dm-4                             10474496   149684 10324812        2% /var/lib/docker/devicemapper/mnt/11111
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
