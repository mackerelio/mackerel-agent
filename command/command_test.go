package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
)

func TestDelayByHost(t *testing.T) {
	delay1 := delayByHost(&mackerel.Host{
		Id:     "246PUVUngPo",
		Name:   "hogehoge2.host.h",
		Type:   "unknown",
		Status: "working",
	})

	delay2 := delayByHost(&mackerel.Host{
		Id:     "21GZjCE5Etb",
		Name:   "hogehoge2.host.h",
		Type:   "unknown",
		Status: "working",
	})

	if !(0 <= delay1.Seconds() && delay1.Seconds() < 60) {
		t.Errorf("delay shoud be between 0 and 60 but %v", delay1)
	}

	if delay1 == delay2 {
		t.Error("delays shoud be different")
	}
}

type jsonObject map[string]interface{}

// newMockAPIServer makes a dummy root directry, a mock API server, a conf.Config to using them
// and returns the Config, mock handlers map and the server.
// The mock handlers map is "<method> <path>"-to-jsonObject-generator map.
func newMockAPIServer(t *testing.T) (config.Config, map[string]func(*http.Request) jsonObject, *httptest.Server) {
	mockHandlers := map[string]func(*http.Request) jsonObject{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := req.Method + " " + req.URL.Path
		handler, ok := mockHandlers[key]
		if !ok {
			t.Fatal("Unexpected request: " + key)
		}

		data := handler(req)

		respJSON, err := json.Marshal(data)
		if err != nil {
			t.Fatal("marshalling JSON failed: ", err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(respJSON))
	}))

	root, err := ioutil.TempDir("", "mackerel-agent-test")
	if err != nil {
		t.Fatal(err)
	}

	conf := config.Config{
		Apibase: ts.URL,
		Root:    root,
	}

	return conf, mockHandlers, ts
}

func TestPrepare(t *testing.T) {
	conf, mockHandlers, ts := newMockAPIServer(t)
	defer ts.Close()

	mockHandlers["POST /api/v0/hosts"] = func(req *http.Request) jsonObject {
		return jsonObject{
			"id": "xxx1234567890",
		}
	}

	mockHandlers["GET /api/v0/hosts/xxx1234567890"] = func(req *http.Request) jsonObject {
		return jsonObject{
			"host": mackerel.Host{
				Id:     "xxx1234567890",
				Name:   "host.example.com",
				Type:   "unknown",
				Status: "standby",
			},
		}
	}

	api, host := Prepare(conf)

	if api.BaseUrl.String() != ts.URL {
		t.Errorf("Apibase mismatch: %s != %s", api.BaseUrl, ts.URL)
	}

	if host.Id != "xxx1234567890" {
		t.Error("Host ID mismatch", host)
	}

	if host.Name != "host.example.com" {
		t.Error("Host name mismatch", host)
	}
}
