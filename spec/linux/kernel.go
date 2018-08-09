// +build linux

package linux

import (
	"os/exec"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/shirou/gopsutil/host"
)

// KernelGenerator XXX
type KernelGenerator struct {
}

var kernelLogger = logging.GetLogger("spec.kernel")

// Generate XXX
func (g *KernelGenerator) Generate() (interface{}, error) {
	commands := map[string][]string{
		"name":    {"uname", "-s"},
		"release": {"uname", "-r"},
		"version": {"uname", "-v"},
		"machine": {"uname", "-m"},
		"os":      {"uname", "-o"},
	}

	results := make(mackerel.Kernel)
	for key, command := range commands {
		out, err := exec.Command(command[0], command[1]).Output()
		if err != nil {
			kernelLogger.Errorf("Failed to run %s %s (skip this spec): %s", command[0], command[1], err)
			return nil, err
		}
		str := strings.TrimSpace(string(out))

		results[key] = str
	}

	platform, _, version, err := host.PlatformInformation()
	if err != nil {
		kernelLogger.Errorf("Failed to get platform information: %s", err)
		return results, nil
	}

	if platformName := normalizePlatform(platform); platformName != "" {
		results["platform_name"] = platformName
	}

	if version != "" {
		results["platform_version"] = version
	}

	return results, nil
}

func normalizePlatform(platform string) string {
	var normalized string

	switch platform {
	case "debian":
		normalized = "Debian"
	case "ubuntu":
		normalized = "Ubuntu"
	case "linuxmint":
		normalized = "Linux Mint"
	case "raspbian":
		normalized = "Raspbian"
	case "fedora":
		normalized = "Fedora"
	case "oracle":
		normalized = "Oracle Linux"
	case "enterpriseenterprise":
		normalized = "Oracle Enterprise Linux"
	case "centos":
		normalized = "CentOS"
	case "redhat":
		normalized = "Red Hat Enterprise Linux"
	case "scientific":
		normalized = "Scientific Linux"
	case "amazon":
		normalized = "Amazon Linux"
	case "xenserver":
		normalized = "XenServer"
	case "cloudlinux":
		normalized = "CloudLinux"
	case "ibm_powerkvm":
		normalized = "IBM PowerKVM"
	case "suse":
		normalized = "SUSE Linux Enterprise Server"
	case "opensuse":
		normalized = "openSUSE"
	case "gentoo":
		normalized = "Gentoo Linux"
	case "slackware":
		normalized = "Slackware"
	case "arch":
		normalized = "Arch Linux"
	case "exherbo":
		normalized = "Exherbo"
	case "alpine":
		normalized = "Alpine Linux"
	case "coreos":
		normalized = "CoreOS"
	default:
		normalized = platform
	}

	return normalized
}
