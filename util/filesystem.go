// +build linux darwin freebsd

package util

// ref. https://github.com/opscode/ohai/blob/master/lib/ohai/plugins/linux/filesystem.rb

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
)

// `df -P` sample:
//  Filesystem     1024-blocks     Used Available Capacity Mounted on
//  /dev/sda1           19734388 16868164 1863772  91% /
//  tmpfs                 517224        0  517224   0% /lib/init/rw
//  udev                  512780       96  512684   1% /dev
//  tmpfs                 517224        4  517220   1% /dev/shm

var dfHeaderPattern = regexp.MustCompile(
	`^Filesystem\s+1024-block`,
)

// DfColumnSpec XXX
type DfColumnSpec struct {
	Name  string
	IsInt bool // type of collected data  true: int64, false: string
}

var dfColumnsPattern = regexp.MustCompile(
	`^(.+?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+%)\s+(.+)$`,
)

var logger = logging.GetLogger("util.filesystem")

var dfOpt []string

func init() {
	switch runtime.GOOS {
	case "darwin":
		dfOpt = []string{"-Pkl"}
	case "freebsd":
		dfOpt = []string{"-Pkt", "noprocfs,devfs,fdescfs,nfs,cd9660"}
	default:
		dfOpt = []string{"-P"}
	}
}

// CollectDfValues XXX
func CollectDfValues(dfColumnSpecs []DfColumnSpec) (map[string]map[string]interface{}, error) {
	cmd := exec.Command("df", dfOpt...)
	cmd.Env = append(cmd.Env, "LANG=C")

	// Ignores exit status in case that df returns exit status 1
	// when the agent does not have permission to access file system info.
	out, err := cmd.Output()
	if err != nil {
		logger.Warningf("'df %s' command exited with a non-zero status: '%s'", strings.Join(dfOpt, " "), err)
	}

	lineScanner := bufio.NewScanner(bytes.NewReader(out))
	filesystems := make(map[string]map[string]interface{})

DF_LINES:
	for lineScanner.Scan() {
		line := lineScanner.Text()

		if dfHeaderPattern.MatchString(line) {
			continue
		} else if matches := dfColumnsPattern.FindStringSubmatch(line); matches != nil {
			name := matches[1]
			entry := make(map[string]interface{})

			for i, colSpec := range dfColumnSpecs {
				stringValue := matches[2+i]

				var (
					value interface{}
					err   error
				)

				if colSpec.IsInt {
					// parse as int64 to allow large size disks
					value, err = strconv.ParseInt(stringValue, 0, 64)
				} else {
					value = stringValue
				}

				if err != nil {
					logger.Warningf("Failed to parse value: [%s]", stringValue)
					continue DF_LINES
				}

				entry[colSpec.Name] = value
			}

			filesystems[name] = entry
		} else {
			logger.Warningf("Failed to parse line: [%s]", line)
		}
	}

	return filesystems, nil
}
