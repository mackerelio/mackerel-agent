package spec

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mackerelio/mackerel-agent/logging"
)

// This Generator collects metadata about cloud instances.
// Currently only EC2 is supported.
// EC2: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AESDG-chapter-instancedata.html
// GCE: https://developers.google.com/compute/docs/metadata
// DigitalOcean: https://developers.digitalocean.com/metadata/

// CloudGenerator definition
type CloudGenerator struct {
	CloudMetaGenerator
}

// CloudMetaGenerator interface of metadata generator for each cloud platform
type CloudMetaGenerator interface {
	Generate() (interface{}, error)
}

// Key is a root key for the generator.
func (g *CloudGenerator) Key() string {
	return "cloud"
}

var cloudLogger = logging.GetLogger("spec.cloud")

var ec2BaseURL, gceMetaURL, digitalOceanBaseURL *url.URL

func init() {
	ec2BaseURL, _ = url.Parse("http://169.254.169.254/latest/meta-data")
	gceMetaURL, _ = url.Parse("http://metadata.google.internal/computeMetadata/v1/?recursive=true")
	digitalOceanBaseURL, _ = url.Parse("http://169.254.169.254/metadata/v1") // has not been yet used
}

var timeout = 100 * time.Millisecond

// SuggestCloudGenerator returns suitable CloudGenerator
func SuggestCloudGenerator() *CloudGenerator {
	if isEC2() {
		return &CloudGenerator{&EC2Generator{ec2BaseURL}}
	}
	if isGCE() {
		return &CloudGenerator{&GCEGenerator{gceMetaURL}}
	}

	return nil
}

func isEC2() bool {
	cl := http.Client{
		Timeout: timeout,
	}
	// '/ami-id` is may be aws specific URL
	resp, err := cl.Get(ec2BaseURL.String() + "/ami-id")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func isGCE() bool {
	_, err := requestGCEMeta()
	return err == nil
}

func requestGCEMeta() ([]byte, error) {
	cl := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", gceMetaURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to request gce meta. response code: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

// EC2Generator meta generator for EC2
type EC2Generator struct {
	baseURL *url.URL
}

// Generate collects metadata from cloud platform.
func (g *EC2Generator) Generate() (interface{}, error) {
	client := http.Client{
		Timeout: timeout,
	}

	metadataKeys := []string{
		"instance-id",
		"instance-type",
		"placement/availability-zone",
		"security-groups",
		"ami-id",
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
			cloudLogger.Debugf("This host may not be running on EC2. Error while reading '%s'", key)
			return nil, nil
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				cloudLogger.Errorf("Results of requesting metadata cannot be read: '%s'", err)
				break
			}
			metadata[key] = string(body)
			cloudLogger.Debugf("results %s:%s", key, string(body))
		} else {
			cloudLogger.Warningf("Status code of the result of requesting metadata '%s' is '%d'", key, resp.StatusCode)
		}
	}

	results := make(map[string]interface{})
	results["provider"] = "ec2"
	results["metadata"] = metadata

	return results, nil
}

// GCEGenerator generate for GCE
type GCEGenerator struct {
	metaURL *url.URL
}

// Generate collects metadata from cloud platform.
func (g *GCEGenerator) Generate() (interface{}, error) {
	bytes, err := requestGCEMeta()
	if err != nil {
		return nil, err
	}
	var data gceMeta
	json.Unmarshal(bytes, &data)
	return data.toGeneratorResults(), nil
}

type gceInstance struct {
	Zone         string
	InstanceType string `json:"machineType"`
	Hostname     string
	InstanceID   uint64 `json:"id"`
}

type gceProject struct {
	ProjectID        string
	NumericProjectID uint64
}

type gceMeta struct {
	Instance *gceInstance
	Project  *gceProject
}

func (g gceMeta) toGeneratorMeta() map[string]string {
	meta := make(map[string]string)

	lastS := func(s string) string {
		ss := strings.Split(s, "/")
		return ss[len(ss)-1]
	}

	if ins := g.Instance; ins != nil {
		meta["hostname"] = ins.Hostname
		meta["instance-id"] = fmt.Sprint(ins.InstanceID)
		meta["instance-type"] = lastS(ins.InstanceType)
		meta["zone"] = lastS(ins.Zone)
	}

	if proj := g.Project; proj != nil {
		meta["projectId"] = proj.ProjectID
	}

	return meta
}

func (g gceMeta) toGeneratorResults() interface{} {
	results := make(map[string]interface{})
	results["provider"] = "gce"
	results["metadata"] = g.toGeneratorMeta()

	return results
}
