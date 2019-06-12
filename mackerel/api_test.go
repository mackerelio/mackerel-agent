package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mkr "github.com/mackerelio/mackerel-client-go"
)

func TestNewAPI(t *testing.T) {
	api, err := NewAPI(
		"http://example.com",
		"dummy-key",
		true,
	)

	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if api.BaseURL.String() != "http://example.com" {
		t.Error("should return URL")
	}

	if api.APIKey != "dummy-key" {
		t.Error("should return api key")
	}

	if api.Verbose != true {
		t.Error("should return verbose value")
	}
}

func TestFindHostByCustomIdentifier(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/hosts" {
			t.Error("request URL should be /api/v0/hosts but :", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but :", req.Method)
		}

		var hosts []map[string]interface{}

		customIdentifier := req.URL.Query().Get("customIdentifier")

		if customIdentifier == "foo-bar" {
			hosts = []map[string]interface{}{
				{
					"id":               "9rxGOHfVF8F",
					"CustomIdentifier": "foo-bar",
					"name":             "mydb001",
					"status":           "working",
					"memo":             "memo",
					"roles":            map[string][]string{"My-Service": {"db-master", "db-slave"}},
				},
			}
		}

		respJSON, _ := json.Marshal(map[string]interface{}{"hosts": hosts})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	api, _ := NewAPI(ts.URL, "dummy-key", false)

	var tests = []struct {
		customIdentifier string
		host             *mkr.Host
		returnInfoError  bool
	}{
		{
			customIdentifier: "foo-bar",
			host: &mkr.Host{
				ID:               "9rxGOHfVF8F",
				Name:             "mydb001",
				Type:             "",
				Status:           "working",
				CustomIdentifier: "foo-bar",
				Memo:             "memo",
				Roles: mkr.Roles{
					"My-Service": []string{"db-master", "db-slave"},
				},
			},
			returnInfoError: false,
		},
		{
			customIdentifier: "unregistered-custom_identifier",
			host:             nil,
			returnInfoError:  true,
		},
		{
			customIdentifier: "",
			host:             nil,
			returnInfoError:  true,
		},
	}

	for _, tc := range tests {
		host, err := api.FindHostByCustomIdentifier(tc.customIdentifier)
		if tc.returnInfoError {
			if _, ok := err.(*InfoError); !ok {
				t.Error("err shoud be type of *InfoError but: ", reflect.TypeOf(err))
			}
		} else {
			if err != nil {
				t.Error("err shoud be nil but: ", err)
			}
		}
		if reflect.DeepEqual(host, tc.host) != true {
			t.Error("request sends json including memo but: ", host)
		}
	}
}

func TestClientError(t *testing.T) {
	tests := []struct {
		Err  error
		Want bool
	}{
		{
			Err:  &mkr.APIError{StatusCode: 400, Message: "400"},
			Want: true,
		},
		{
			Err:  &mkr.APIError{StatusCode: 499, Message: "499"},
			Want: true,
		},
		{
			Err:  &mkr.APIError{StatusCode: 500, Message: "500"},
			Want: false,
		},
		{
			Err:  fmt.Errorf("err"),
			Want: false,
		},
	}
	for _, tt := range tests {
		v := IsClientError(tt.Err)
		if v != tt.Want {
			t.Errorf("IsClientError(%v) = %v; want %v", tt.Err, v, tt.Want)
		}
	}
}

func TestServerError(t *testing.T) {
	tests := []struct {
		Err  error
		Want bool
	}{
		{
			Err:  &mkr.APIError{StatusCode: 500, Message: "500"},
			Want: true,
		},
		{
			Err:  &mkr.APIError{StatusCode: 599, Message: "599"},
			Want: true,
		},
		{
			Err:  &mkr.APIError{StatusCode: 400, Message: "400"},
			Want: false,
		},
		{
			Err:  fmt.Errorf("err"),
			Want: false,
		},
	}
	for _, tt := range tests {
		v := IsServerError(tt.Err)
		if v != tt.Want {
			t.Errorf("IsServerError(%v) = %v; want %v", tt.Err, v, tt.Want)
		}
	}
}
