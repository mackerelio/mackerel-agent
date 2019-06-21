// +build linux darwin freebsd netbsd

package util

// ref. https://github.com/opscode/ohai/blob/master/lib/ohai/plugins/linux/filesystem.rb

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/timeout"
	"github.com/mackerelio/golib/logging"
)

// DfStat is disk free statistics from df command.
// Field names are taken from column names of `df -P`
type DfStat struct {
	Name      string
	Blocks    uint64
	Used      uint64
	Available uint64
	Capacity  uint8
	Mounted   string
}

// `df -P` sample:
//  Filesystem     1024-blocks     Used Available Capacity Mounted on
//  /dev/sda1           19734388 16868164 1863772  91% /
//  tmpfs                 517224        0  517224   0% /lib/init/rw
//  udev                  512780       96  512684   1% /dev
//  tmpfs                 517224        4  517220   1% /dev/shm

var dfHeaderPattern = regexp.MustCompile(
	// 1024-blocks or 1k-blocks
	`^Filesystem\s+(?:1024|1[Kk])-block`,
)

var logger = logging.GetLogger("util.filesystem")

var dfOpt = "-Pkl"

func init() {
	// Some `df` command such as busybox does not have `-P` or `-l` option.
	tio := &timeout.Timeout{
		Cmd:       exec.Command("df", dfOpt),
		Duration:  3 * time.Second,
		KillAfter: 1 * time.Second,
	}
	exitSt, _, stderr, err := tio.Run()
	if err == nil && exitSt.Code != 0 && (strings.Contains(stderr, "df: invalid option -- ") || strings.Contains(stderr, "df: unrecognized option: ")) {
		dfOpt = "-k"
	}
}

// CollectDfValues collects disk free statistics from df command
func CollectDfValues() ([]*DfStat, error) {
	cmd := exec.Command("df", dfOpt)
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  15 * time.Second,
		KillAfter: 5 * time.Second,
	}
	exitSt, stdout, stderr, err := tio.Run()
	if err != nil {
		logger.Warningf("failed to invoke 'df %s' command: %q", dfOpt, err)
		return nil, nil
	}
	// Ignores exit status in case that df returns exit status 1
	// when the agent does not have permission to access file system info.
	if exitSt.Code != 0 {
		logger.Warningf("'df %s' command exited with a non-zero status: %d: %q", dfOpt, exitSt.Code, stderr)
		return nil, nil
	}
	return parseDfLines(stdout), nil
}

func parseDfLines(out string) []*DfStat {
	lineScanner := bufio.NewScanner(strings.NewReader(out))
	var filesystems []*DfStat
	for lineScanner.Scan() {
		line := lineScanner.Text()
		if dfHeaderPattern.MatchString(line) {
			continue
		}
		dfstat, err := parseDfLine(line)
		if err != nil {
			logger.Warningf(err.Error())
			continue
		}
		// https://github.com/docker/docker/blob/v1.5.0/daemon/graphdriver/devmapper/deviceset.go#L981
		if strings.HasPrefix(dfstat.Name, "/dev/mapper/docker-") {
			continue
		}
		// https://debbugs.gnu.org/cgi/bugreport.cgi?bug=10363
		// http://git.savannah.gnu.org/gitweb/?p=coreutils.git;a=commit;h=1e18d8416f9ef43bf08982cabe54220587061a08
		// coreutils >= 8.15
		if strings.HasPrefix(dfstat.Name, "/dev/dm-") && strings.Contains(dfstat.Mounted, "devicemapper/mnt") {
			continue
		}

		if runtime.GOOS == "darwin" {
			if strings.HasPrefix(dfstat.Mounted, "/Volumes/") {
				continue
			}
			// Skip APFS vm partition, add its usage to the root filesystem.
			if dfstat.Mounted == "/private/var/vm" {
				for _, fs := range filesystems {
					if fs.Mounted == "/" {
						fs.Used += dfstat.Used
						break
					}
				}
				continue
			}
		}

		filesystems = append(filesystems, dfstat)
	}
	return filesystems
}

var dfColumnsPattern = regexp.MustCompile(`^(.+?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)%\s+(.+)$`)

func parseDfLine(line string) (*DfStat, error) {
	matches := dfColumnsPattern.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("failed to parse line: [%s]", line)
	}
	name := matches[1]
	blocks, _ := strconv.ParseUint(matches[2], 0, 64)
	used, _ := strconv.ParseUint(matches[3], 0, 64)
	available, _ := strconv.ParseUint(matches[4], 0, 64)
	capacity, _ := strconv.ParseUint(matches[5], 0, 8)
	mounted := matches[6]

	return &DfStat{
		Name:      name,
		Blocks:    blocks,
		Used:      used,
		Available: available,
		Capacity:  uint8(capacity),
		Mounted:   mounted,
	}, nil
}
