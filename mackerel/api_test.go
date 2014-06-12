package mackerel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mackerelio/mackerel-agent/version"
)

func TestNewApi(t *testing.T) {
	api, err := NewApi(
		"http://example.com",
		"dummy-key",
		true,
	)

	if err != nil {
		t.Error("should not raise error")
	}

	if api.BaseUrl.String() != "http://example.com" {
		t.Error("should return URL")
	}

	if api.ApiKey != "dummy-key" {
		t.Error("should return api key")
	}

	if api.Verbose != true {
		t.Error("should return verbose value")
	}
}

func TestUrlFor(t *testing.T) {
	api, _ := NewApi(
		"http://example.com",
		"dummy-key",
		true,
	)

	if api.urlFor("/").String() != "http://example.com/" {
		t.Error("should return http://example.com/")
	}

	if api.urlFor("/path/to/api").String() != "http://example.com/path/to/api" {
		t.Error("should return http://example.com/path/to/api")
	}
}

func TestDo(t *testing.T) {
	version.VERSION = "1.0.0"
	version.GITCOMMIT = "1234beaf"
	handler := func(res http.ResponseWriter, req *http.Request) {
		user_agent := "mackerel-agent/1.0.0 (Revision 1234beaf)"
		if req.Header.Get("X-Api-Key") != "dummy-key" {
			t.Error("X-Api-Key header should contains passed key")
		}

		if h := req.Header.Get("X-Agent-Version"); h != version.VERSION {
			t.Errorf("X-Agent-Version shoud be %s but %s", version.VERSION, h)
		}

		if h := req.Header.Get("X-Revision"); h != version.GITCOMMIT {
			t.Errorf("X-Revision shoud be %s but %s", version.GITCOMMIT, h)
		}

		if h := req.Header.Get("User-Agent"); h != user_agent {
			t.Errorf("User-Agent shoud be '%s' but %s", user_agent, h)
		}

		version.GITCOMMIT = ""
		version.VERSION = ""
	}
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler(res, req)
	}))
	defer ts.Close()

	api, _ := NewApi(
		ts.URL,
		"dummy-key",
		false,
	)

	req, _ := http.NewRequest("GET", api.urlFor("/").String(), nil)
	api.Do(req)
}

func TestCreateHost(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		called = true
		if req.URL.Path != "/api/v0/hosts" {
			t.Error("request URL should be /api/v0/hosts but :", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but :", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)
		content := string(body)

		var data struct {
			Name          string              `json:"name"`
			Tame          string              `json:"type"`
			Status        string              `json:"status"`
			Meta          map[string]string   `json:"meta"`
			Interfaces    []map[string]string `json:"interfaces"`
			RoleFullnames []string            `json:"roleFullnames"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request content should be decoded as json", content)
		}

		if data.Meta["memo"] != "hello" {
			t.Error("request sends json including memo but: ", data)
		}

		if len(data.Interfaces) == 0 {
			t.Error("request sends json including interfaces but: ", data)
		}
		iface := data.Interfaces[0]
		if iface["name"] != "eth0" || iface["ipAddress"] != "10.0.4.7" {
			t.Error("interface name and ipAddress should be eth0 and 10.0.4.7, respectively, but: ", data)
		}

		if len(data.RoleFullnames) != 1 {
			t.Errorf("roleFullnames must have size 1: %v", data.RoleFullnames)
		}

		if data.RoleFullnames[0] != "My-Service:app-default" {
			t.Errorf("Wrong data for roleFullnames: %v", data.RoleFullnames)
		}

		respJson, _ := json.Marshal(map[string]interface{}{
			"id": "ABCD123",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer ts.Close()

	api, _ := NewApi(ts.URL, "dummy-key", false)

	interfaces := make([]map[string]interface{}, 0)
	interfaces = append(interfaces, map[string]interface{}{
		"name":       "eth0",
		"ipAddress":  "10.0.4.7",
		"macAddress": "01:23:45:67:89:ab",
		"encap":      "Ethernet",
	})
	hostId, err := api.CreateHost("dummy", map[string]interface{}{
		"memo": "hello",
	}, interfaces, []string{"My-Service:app-default"})

	if err != nil {
		t.Error("should not raise error: ", err)
	}

	if !called {
		t.Error("should http-request")
	}

	if hostId != "ABCD123" {
		t.Error("should returns ABCD123 but:", hostId)
	}
}

func TestCreateHostWithNilArgs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts" {
			t.Error("request URL should be /api/v0/hosts but :", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but :", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)
		content := string(body)

		var data struct {
			Name          string              `json:"name"`
			Tame          string              `json:"type"`
			Status        string              `json:"status"`
			Meta          map[string]string   `json:"meta"`
			Interfaces    []map[string]string `json:"interfaces"`
			RoleFullnames []string            `json:"roleFullnames"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request content should be decoded as json", content)
		}

		respJson, _ := json.Marshal(map[string]interface{}{
			"id": "ABCD123",
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer ts.Close()

	api, _ := NewApi(ts.URL, "dummy-key", false)

	// with nil args
	hostId, err := api.CreateHost("nilsome", nil, nil, nil)
	if err != nil {
		t.Error("should not return error but got: ", err)
	}

	if hostId != "ABCD123" {
		t.Error("should returns ABCD123 but:", hostId)
	}
}

func TestUpdateHost(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		called = true
		if req.URL.Path != "/api/v0/hosts/ABCD123" {
			t.Error("request URL should be /api/v0/hosts/ABCD123 but :", req.URL.Path)
		}

		if req.Method != "PUT" {
			t.Error("request method should be PUT but :", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)
		content := string(body)

		var data struct {
			Name          string              `json:"name"`
			Tame          string              `json:"type"`
			Status        string              `json:"status"`
			Meta          map[string]string   `json:"meta"`
			Interfaces    []map[string]string `json:"interfaces"`
			RoleFullnames []string            `json:"roleFullnames"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request content should be decoded as json", content)
		}

		if data.Meta["memo"] != "hello" {
			t.Error("request sends json including memo but: ", data)
		}

		if len(data.Interfaces) == 0 {
			t.Error("request sends json including interfaces but: ", data)
		}
		iface := data.Interfaces[0]
		if iface["name"] != "eth0" || iface["ipAddress"] != "10.0.4.7" {
			t.Error("interface name and ipAddress should be eth0 and 10.0.4.7, respectively, but: ", data)
		}

		if len(data.RoleFullnames) != 1 {
			t.Errorf("roleFullnames must have size 1: %v", data.RoleFullnames)
		}

		if data.RoleFullnames[0] != "My-Service:app-default" {
			t.Errorf("Wrong data for roleFullnames: %v", data.RoleFullnames)
		}
	}))
	defer ts.Close()

	api, _ := NewApi(ts.URL, "dummy-key", false)

	interfaces := make([]map[string]interface{}, 0)
	interfaces = append(interfaces, map[string]interface{}{
		"name":       "eth0",
		"ipAddress":  "10.0.4.7",
		"macAddress": "01:23:45:67:89:ab",
		"encap":      "Ethernet",
	})

	err := api.UpdateHost("ABCD123", "dummy", map[string]interface{}{
		"memo": "hello",
	}, interfaces, []string{"My-Service:app-default"})

	if err != nil {
		t.Error("should not raise error: ", err)
	}

	if !called {
		t.Error("should http-request")
	}
}
