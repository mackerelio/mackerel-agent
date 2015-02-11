// +build linux

package linux

import (
	"github.com/mackerelio/mackerel-agent/logging"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// EC2: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AESDG-chapter-instancedata.html
// GCE: https://developers.google.com/compute/docs/metadata
// DigitalOcean: https://developers.digitalocean.com/metadata/

// InstanceGenerator XXX
type InstanceGenerator struct {
	BaseURL *url.URL
}

// Key XXX
func (g *InstanceGenerator) Key() string {
	return "instance"
}

var instanceLogger = logging.GetLogger("spec.instance")

// NewInstanceGenerator XXX
func NewInstanceGenerator(baseurl string) (*InstanceGenerator, error) {
	if baseurl == "" {
		baseurl = "http://169.254.169.254/latest/meta-data"
	}
	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, err
	}
	return &InstanceGenerator{u}, nil
}

// Generate XXX
func (g *InstanceGenerator) Generate() (interface{}, error) {

	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	metadataKeys := []string{
		"instance-id",
		"instance-type",
		"placement/availability-zone",
		"security-groups",
		"ami_id",
		"hostname",
		"local-hostname",
		"public-hostname",
		"local-ipv4",
		"public-keys",
		"public-ipv4",
		"reservation-id",
	}

	metadata := make(map[string]string)

	for _, key := range metadataKeys {
		resp, err := client.Get(g.BaseURL.String() + "/" + key)
		instanceLogger.Infof("reading '%s/%s'", g.BaseURL.String(), key)
		if err != nil {
			instanceLogger.Infof("This host may not be running on EC2. Error while reading '%s'", key)
			return nil, nil
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				instanceLogger.Errorf("Results of requesting metadata cannot be read.")
				break
			}
			metadata[key] = string(body)
			instanceLogger.Infof("results %s:%s", key, string(body))
		}
	}

	return metadata, nil
}
