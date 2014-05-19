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

func Collect(specGenerators []Generator) map[string]interface{} {
	specs := make(map[string]interface{})
	for _, g := range specGenerators {
		value, err := g.Generate()
		if err != nil {
			logger.Errorf("Failed to collect meta in %T (skip this spec): %s", g, err.Error())
		}
		specs[g.Key()] = value
	}
	specs["agent-version"] = version.VERSION
	specs["agent-revision"] = version.GITCOMMIT
	specs["agent-name"] = version.UserAgent()
	return specs
}
