package spec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestCloudKey(t *testing.T) {
	g := &CloudGenerator{}

	if g.Key() != "cloud" {
		t.Error("key should be cloud")
	}
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

	cloud, typeOk := value.(map[string]interface{})
	if !typeOk {
		t.Errorf("value should be map. %+v", value)
	}

	value, ok := cloud["metadata"]
	if !ok {
		t.Error("results should have metadata.")
	}

	metadata, typeOk := value.(map[string]string)
	if !typeOk {
		t.Errorf("v should be map. %+v", value)
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

func TestGCEGenerate(t *testing.T) {
	// curl "http://metadata.google.internal/computeMetadata/v1/?recursive=true" -H "Metadata-Flavor: Google"
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
	// both of ec2BaseURL and gceMetaURL are unreachable
	unreachableURL, _ := url.Parse("http://unreachable.localhost")
	ec2BaseURL = unreachableURL
	gceMetaURL = unreachableURL
	cGen := SuggestCloudGenerator()
	if cGen != nil {
		t.Errorf("cGen should be nil but, %s", cGen)
	}

	func() { // ec2BaseURL is reachable but returns 404
		ts := httptest.NewServer(http.NotFoundHandler())
		defer ts.Close()
		u, _ := url.Parse(ts.URL)
		ec2BaseURL = u

		cGen = SuggestCloudGenerator()
		if cGen != nil {
			t.Errorf("cGen should be nil but, %s", cGen)
		}
	}()

	func() { // suggest GCEGenerator
		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			fmt.Fprint(res, "GCE:OK")
		}))
		defer ts.Close()
		ec2BaseURL = unreachableURL
		u, _ := url.Parse(ts.URL)
		gceMetaURL = u

		cGen = SuggestCloudGenerator()
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
}
