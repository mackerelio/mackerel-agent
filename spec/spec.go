package spec

import (
	"os/exec"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/version"
)

var logger = logging.GetLogger("spec")

type Generator interface {
	Key() string
	Generate() (interface{}, error)
}

func GetHostname() (string, error) {
	out, err := exec.Command("uname", "-n").Output()

	if err != nil {
		return "", err
	}
	str := strings.TrimSpace(string(out))

	return str, nil
}

func CollectMeta(metaGenerators []Generator) map[string]interface{} {
	meta := make(map[string]interface{})
	for _, g := range metaGenerators {
		value, err := g.Generate()
		if err != nil {
			logger.Errorf("Failed to collect meta in %T (skip this spec): %s", g, err.Error())
		}
		meta[g.Key()] = value
	}
	meta["agent-version"] = version.VERSION
	meta["agent-revision"] = version.GITCOMMIT
	meta["agent-name"] = version.UserAgent()
	return meta
}
