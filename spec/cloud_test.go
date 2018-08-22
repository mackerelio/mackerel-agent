package spec

import (
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

func TestSuggestCloudGenerator(t *testing.T) {
	// All Cloud meta URLs are unreachable
	unreachableURL, _ := url.Parse("http://unreachable.localhost")
	ec2BaseURL = unreachableURL
	gceMetaURL = unreachableURL
	azureVMBaseURL = unreachableURL

	conf := config.Config{}

	cGen := SuggestCloudGenerator(&conf)
	if cGen != nil {
		t.Errorf("cGen should be nil but, %s", cGen)
	}

	func() { // ec2BaseURL is reachable but returns 404
		ts := httptest.NewServer(http.NotFoundHandler())
		defer ts.Close()
		u, _ := url.Parse(ts.URL)
		ec2BaseURL = u
		defer func() { ec2BaseURL = unreachableURL }()

		cGen = SuggestCloudGenerator(&conf)
		if cGen != nil {
			t.Errorf("cGen should be nil but, %s", cGen)
		}
	}()

	func() { // suggest GCEGenerator
		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			fmt.Fprint(res, "GCE:OK")
		}))
		defer ts.Close()
		u, _ := url.Parse(ts.URL)
		gceMetaURL = u
		defer func() { gceMetaURL = unreachableURL }()

		cGen = SuggestCloudGenerator(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		gceGen, ok := cGen.CloudMetaGenerator.(*GCEGenerator)
		if !ok {
			t.Errorf("cGen should be *GCEGenerator")
		}
		if gceGen.metaURL != gceMetaURL {
			t.Errorf("something went wrong")
		}
	}()

	func() { // suggest AzureVMGenerator
		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.Header.Get("Metadata") != "true" {
				http.NotFound(res, req)
			}
			fmt.Fprint(res, "ok")
		}))
		defer ts.Close()
		u, _ := url.Parse(ts.URL)
		azureVMBaseURL = u
		defer func() { azureVMBaseURL = unreachableURL }()

		cGen = SuggestCloudGenerator(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		gen, ok := cGen.CloudMetaGenerator.(*AzureVMGenerator)
		if !ok {
			t.Errorf("cGen should be *AzureVMGenerator")
		}
		if gen.baseURL != azureVMBaseURL {
			t.Errorf("something went wrong")
		}
	}()

	func() { // multiple generators are available, but suggest the first responded one (in this case EC2)
		// azure. ok immediately
		tsA := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			fmt.Fprint(res, "ok")
		}))
		defer tsA.Close()
		uA, _ := url.Parse(tsA.URL)
		azureVMBaseURL = uA
		// ec2 / gce. ok after 1 second
		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(2 * time.Second)
			fmt.Fprint(res, "ok")
		}))
		defer ts.Close()
		u, _ := url.Parse(ts.URL)
		ec2BaseURL = u
		gceMetaURL = u
		defer func() {
			ec2BaseURL = unreachableURL
			gceMetaURL = unreachableURL
			azureVMBaseURL = unreachableURL
		}()

		cGen = SuggestCloudGenerator(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(*AzureVMGenerator)
		if !ok {
			t.Errorf("cGen should be *AzureVMGenerator")
		}
	}()
}

func TestSuggestCloudGenerator_CloudPlatformSpecified(t *testing.T) {
	// All Cloud meta URLs are unreachable
	unreachableURL, _ := url.Parse("http://unreachable.localhost")
	ec2BaseURL = unreachableURL
	gceMetaURL = unreachableURL
	azureVMBaseURL = unreachableURL

	{
		conf := config.Config{
			CloudPlatform: config.CloudPlatformNone,
		}

		cGen := SuggestCloudGenerator(&conf)
		if cGen != nil {
			t.Errorf("cGen should be nil.")
		}
	}

	{
		conf := config.Config{
			CloudPlatform: config.CloudPlatformEC2,
		}

		cGen := SuggestCloudGenerator(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(*EC2Generator)
		if !ok {
			t.Errorf("cGen should be *EC2Generator")
		}
	}

	{
		conf := config.Config{
			CloudPlatform: config.CloudPlatformGCE,
		}

		cGen := SuggestCloudGenerator(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(*GCEGenerator)
		if !ok {
			t.Errorf("cGen should be *GCEGenerator")
		}
	}

	{
		conf := config.Config{
			CloudPlatform: config.CloudPlatformAzureVM,
		}

		cGen := SuggestCloudGenerator(&conf)
		if cGen == nil {
			t.Errorf("cGen should not be nil.")
		}

		_, ok := cGen.CloudMetaGenerator.(*AzureVMGenerator)
		if !ok {
			t.Errorf("cGen should be *AzureVMGenerator")
		}
	}
}
