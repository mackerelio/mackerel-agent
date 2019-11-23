package spec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mackerel-agent/config"
)

type emptyCloudMetaGenerator struct{}

func (g *emptyCloudMetaGenerator) Generate() (interface{}, error) {
	return nil, nil
}

func (g *emptyCloudMetaGenerator) SuggestCustomIdentifier() (string, error) {
	return "", nil
}

type mockAzureCloudMetaGenerator struct {
	emptyCloudMetaGenerator
	isAzureVM bool
}

func (g *mockAzureCloudMetaGenerator) IsAzureVM(ctx context.Context) bool {
	return g.isAzureVM
}

type mockEC2CloudMetaGenerator struct {
	emptyCloudMetaGenerator
	isEC2 bool
}

func (g *mockEC2CloudMetaGenerator) IsEC2(ctx context.Context) bool {
	return g.isEC2
}

type mockGCECloudMetaGenerator struct {
	emptyCloudMetaGenerator
	isGCE bool
}

func (g *mockGCECloudMetaGenerator) IsGCE(ctx context.Context) bool {
	return g.isGCE
}

func TestCloudGenerate(t *testing.T) {
	handler := func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "i-4f90d537")
	}
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler(res, req)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}
	g := &CloudGenerator{&EC2Generator{u}}

	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}

	cloud, typeOk := value.(*mackerel.Cloud)
	if !typeOk {
		t.Errorf("value should be *mackerel.Cloud. %+v", value)
	}

	metadata, typeOk := cloud.MetaData.(map[string]string)
	if !typeOk {
		t.Errorf("MetaData should be map. %+v", cloud.MetaData)
	}

	if len(metadata["instance-id"]) == 0 {
		t.Error("instance-id should be filled")
	}

	customIdentifier, err := g.SuggestCustomIdentifier()
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}

	if len(customIdentifier) == 0 {
		t.Error("customIdentifier should be retrieved")
	}
}

func TestEC2SuggestCustomIdentifier(t *testing.T) {
	i := 0
	threshold := 100
	handler := func(res http.ResponseWriter, req *http.Request) {
		if i < threshold {
			http.Error(res, "not found", 404)
		} else {
			fmt.Fprint(res, "i-4f90d537")
		}
		i++
	}
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler(res, req)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}
	g := &CloudGenerator{&EC2Generator{u}}

	// 404, 404, 404 => give up
	{
		_, err := g.SuggestCustomIdentifier()
		if err == nil {
			t.Errorf("should raise error: %s", err)
		}
	}
	i = 0
	threshold = 0
	// 200 => ok
	{
		customIdentifier, err := g.SuggestCustomIdentifier()
		if err != nil {
			t.Errorf("should not raise error: %s", err)
		}
		if customIdentifier != "i-4f90d537.ec2.amazonaws.com" {
			t.Error("customIdentifier mismatch")
		}
	}
	i = 0
	threshold = 1
	// 404, 200 => ok
	{
		customIdentifier, err := g.SuggestCustomIdentifier()
		if err != nil {
			t.Errorf("should not raise error: %s", err)
		}
		if customIdentifier != "i-4f90d537.ec2.amazonaws.com" {
			t.Error("customIdentifier mismatch")
		}
	}
	i = 0
	threshold = 3
	// 404, 404, 404(give up), 200, ...
	{
		_, err := g.SuggestCustomIdentifier()
		if err == nil {
			t.Errorf("should raise error: %s", err)
		}
	}
}

func TestGCEGenerate(t *testing.T) {
	// curl "http://metadata.google.internal./computeMetadata/v1/?recursive=true" -H "Metadata-Flavor: Google"
	sampleJSON := []byte(`{
	  "instance": {
		"attributes": {},
		"cpuPlatform": "Intel Ivy Bridge",
		"description": "",
		"disks": [
		  {
			"deviceName": "gce-1",
			"index": 0,
			"mode": "READ_WRITE",
			"type": "PERSISTENT"
		  }
		],
		"hostname": "gce-1.c.dummyproj-987.internal",
		"id": 4567890123456789123,
		"image": "",
		"machineType": "projects/1234567890123/machineTypes/g1-small",
		"maintenanceEvent": "NONE",
		"networkInterfaces": [
		  {
			"accessConfigs": [
			  {
				"externalIp": "203.0.113.1",
				"type": "ONE_TO_ONE_NAT"
			  }
			],
			"forwardedIps": [],
			"ip": "192.0.2.1",
			"network": "projects/1234567890123/networks/default"
		  }
		],
		"scheduling": {
		  "automaticRestart": "TRUE",
		  "onHostMaintenance": "MIGRATE"
		},
		"serviceAccounts": {
		  "1234567890123-compute@developer.gserviceaccount.com": {
			"aliases": [
			  "default"
			],
			"email": "1234567890123-compute@developer.gserviceaccount.com",
			"scopes": [
			  "https://www.googleapis.com/auth/devstorage.read_only",
			  "https://www.googleapis.com/auth/logging.write"
			]
		  },
		  "default": {
			"aliases": [
			  "default"
			],
			"email": "1234567890123-compute@developer.gserviceaccount.com",
			"scopes": [
			  "https://www.googleapis.com/auth/devstorage.read_only",
			  "https://www.googleapis.com/auth/logging.write"
			]
		  }
		},
		"tags": [],
		"virtualClock": {
		  "driftToken": "12345678901234567890"
		},
		"zone": "projects/1234567890123/zones/asia-east1-a"
	  },
	  "project": {
		"attributes": {
		  "google-compute-default-region": "us-central1",
		  "google-compute-default-zone": "us-central1-f",
		  "sshKeys": "dummy_user:ssh-rsa AAAhogehoge google-ssh {\"userName\":\"dummy_user@example.com\",\"expireOn\":\"2015-07-12T11:11:43+0000\"}\ndummy_user:ecdsa-sha2-nistp256 AAAhogefuga google-ssh {\"userName\":\"dummy_user@example.com\",\"expireOn\":\"2015-07-12T11:11:39+0000\"}\n"
		},
		"numericProjectId": 1234567890123,
		"projectId": "dummyprof-987"
	  }
	}`)

	var data gceMeta
	json.Unmarshal(sampleJSON, &data)

	if !reflect.DeepEqual(data.Instance, &gceInstance{
		Zone:         "projects/1234567890123/zones/asia-east1-a",
		InstanceType: "projects/1234567890123/machineTypes/g1-small",
		Hostname:     "gce-1.c.dummyproj-987.internal",
		InstanceID:   4567890123456789123,
	}) {
		t.Errorf("data.Instance should be assigned")
	}

	if !reflect.DeepEqual(data.Project, &gceProject{
		ProjectID:        "dummyprof-987",
		NumericProjectID: 1234567890123,
	}) {
		t.Errorf("data.Project should be assigned")
	}

	if d := data.toGeneratorMeta(); !reflect.DeepEqual(d, map[string]string{
		"zone":          "asia-east1-a",
		"instance-type": "g1-small",
		"hostname":      "gce-1.c.dummyproj-987.internal",
		"instance-id":   "4567890123456789123",
		"projectId":     "dummyprof-987",
	}) {
		t.Errorf("data.Project should be assigned")
	}

}

type slowCloudMetaGenerator struct {
	emptyCloudMetaGenerator
}

func (g *slowCloudMetaGenerator) IsEC2(ctx context.Context) bool {
	time.Sleep(2 * time.Second)
	return true
}

func (g *slowCloudMetaGenerator) IsGCE(ctx context.Context) bool {
	time.Sleep(2 * time.Second)
	return true
}

func TestCloudGeneratorSuggestor(t *testing.T) {
	conf := config.Config{}
	// none
	{
		suggestor := &cloudGeneratorSuggestor{
			ec2Generator:     &mockEC2CloudMetaGenerator{isEC2: false},
			gceGenerator:     &mockGCECloudMetaGenerator{isGCE: false},
			azureVMGenerator: &mockAzureCloudMetaGenerator{isAzureVM: false},
		}
		cGen := suggestor.Suggest(&conf)
		if cGen != nil {
			t.Errorf("cGen should be nil but, %s", cGen)
		}
	}

	// EC2
	{
		suggestor := &cloudGeneratorSuggestor{
			ec2Generator:     &mockEC2CloudMetaGenerator{isEC2: true},
			gceGenerator:     &mockGCECloudMetaGenerator{isGCE: false},
			azureVMGenerator: &mockAzureCloudMetaGenerator{isAzureVM: false},
		}
		cGen := suggestor.Suggest(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(ec2Generator)
		if !ok {
			t.Errorf("cGen should be ec2Generator")
		}
	}

	// GCE
	{
		suggestor := &cloudGeneratorSuggestor{
			ec2Generator:     &mockEC2CloudMetaGenerator{isEC2: false},
			gceGenerator:     &mockGCECloudMetaGenerator{isGCE: true},
			azureVMGenerator: &mockAzureCloudMetaGenerator{isAzureVM: false},
		}
		cGen := suggestor.Suggest(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(gceGenerator)
		if !ok {
			t.Errorf("cGen should be gceGenerator")
		}
	}

	// AzureVM
	{
		suggestor := &cloudGeneratorSuggestor{
			ec2Generator:     &mockEC2CloudMetaGenerator{isEC2: false},
			gceGenerator:     &mockGCECloudMetaGenerator{isGCE: false},
			azureVMGenerator: &mockAzureCloudMetaGenerator{isAzureVM: true},
		}
		cGen := suggestor.Suggest(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(azureVMGenerator)
		if !ok {
			t.Errorf("cGen should be azureVMGenerator")
		}
	}

	// multiple generators are available, but suggest the first responded one (in this case Azure)
	{
		suggestor := &cloudGeneratorSuggestor{
			ec2Generator:     &slowCloudMetaGenerator{},
			gceGenerator:     &slowCloudMetaGenerator{},
			azureVMGenerator: &mockAzureCloudMetaGenerator{isAzureVM: true},
		}
		cGen := suggestor.Suggest(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(azureVMGenerator)
		if !ok {
			t.Errorf("cGen should be azureVMGenerator")
		}
	}
}

func TestCloudGeneratorSuggestor_CloudPlatformSpecified(t *testing.T) {
	suggestor := &cloudGeneratorSuggestor{
		ec2Generator:     &mockEC2CloudMetaGenerator{isEC2: false},
		gceGenerator:     &mockGCECloudMetaGenerator{isGCE: false},
		azureVMGenerator: &mockAzureCloudMetaGenerator{isAzureVM: false},
	}
	{
		conf := config.Config{
			CloudPlatform: config.CloudPlatformNone,
		}

		cGen := suggestor.Suggest(&conf)
		if cGen != nil {
			t.Errorf("cGen should be nil.")
		}
	}

	{
		conf := config.Config{
			CloudPlatform: config.CloudPlatformEC2,
		}

		cGen := suggestor.Suggest(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(ec2Generator)
		if !ok {
			t.Errorf("cGen should be ec2Generator")
		}
	}

	{
		conf := config.Config{
			CloudPlatform: config.CloudPlatformGCE,
		}

		cGen := suggestor.Suggest(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(gceGenerator)
		if !ok {
			t.Errorf("cGen should be gceGenerator")
		}
	}

	{
		conf := config.Config{
			CloudPlatform: config.CloudPlatformAzureVM,
		}

		cGen := suggestor.Suggest(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(azureVMGenerator)
		if !ok {
			t.Errorf("cGen should be azureVMGenerator")
		}
	}
}

func TestCloudGeneratorSuggestor_Public(t *testing.T) {
	{
		gen, ok := CloudGeneratorSuggestor.ec2Generator.(*EC2Generator)
		if !ok {
			t.Error("EC2Generator should be injected as ec2Generator")
		}
		if gen.baseURL.String() != ec2BaseURL.String() {
			t.Error("real baseURL should be embedded to ec2Generator")
		}
	}
	{
		gen, ok := CloudGeneratorSuggestor.gceGenerator.(*GCEGenerator)
		if !ok {
			t.Error("GCEGenerator should be injected as gceGenerator")
		}
		if gen.metaURL.String() != gceMetaURL.String() {
			t.Error("real metaURL should be embedded to gceGenerator")
		}
	}
	{
		gen, ok := CloudGeneratorSuggestor.azureVMGenerator.(*AzureVMGenerator)
		if !ok {
			t.Error("AzureVMGenerator should be injected as azureVMGenerator")
		}
		if gen.baseURL.String() != azureVMBaseURL.String() {
			t.Error("real baseURL should be embedded to azureVMGenerator ")
		}
	}
}
