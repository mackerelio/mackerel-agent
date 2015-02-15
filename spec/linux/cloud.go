// +build linux

package linux

import (
	"github.com/mackerelio/mackerel-agent/logging"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// This Generator collects metadata about cloud instances.
// Currently only EC2 is supported.
// EC2: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AESDG-chapter-instancedata.html
// GCE: https://developers.google.com/compute/docs/metadata
// DigitalOcean: https://developers.digitalocean.com/metadata/

// CloudGenerator
type CloudGenerator struct {
	baseURL *url.URL
}

// Key is a root key for the generator.
func (g *CloudGenerator) Key() string {
	return "cloud"
}

var cloudLogger = logging.GetLogger("spec.cloud")

// NewCloudGenerator creates a Cloud Generator instance with specified baseurl.
func NewCloudGenerator(baseurl string) (*CloudGenerator, error) {
	if baseurl == "" {
		baseurl = "http://169.254.169.254/latest/meta-data"
	}
	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, err
	}
	return &CloudGenerator{u}, nil
}

// Generate collects metadata from cloud platform.
func (g *CloudGenerator) Generate() (interface{}, error) {

	timeout := time.Duration(100 * time.Millisecond)
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
		resp, err := client.Get(g.baseURL.String() + "/" + key)
		if err != nil {
			cloudLogger.Infof("This host may not be running on EC2. Error while reading '%s'", key)
			return nil, nil
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				cloudLogger.Errorf("Results of requesting metadata cannot be read.")
				break
			}
			metadata[key] = string(body)
			cloudLogger.Debugf("results %s:%s", key, string(body))
		}
	}

	results := make(map[string]interface{})
	results["provider"] = "ec2"
	results["metadata"] = metadata

	return results, nil
}
