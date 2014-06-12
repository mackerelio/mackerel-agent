package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
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

func TestPrepare(t *testing.T) {
	_respondJSON := func(w http.ResponseWriter, data map[string]interface{}) {
		respJson, err := json.Marshal(data)
		if err != nil {
			t.Fatal("marshalling JSON failed: ", err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(respJson))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/hosts", func(w http.ResponseWriter, req *http.Request) {
		response := map[string]interface{}{
			"id": "ThisHostId",
		}
		_respondJSON(w, response)
	})

	mux.HandleFunc("/api/v0/hosts/ThisHostId", func(w http.ResponseWriter, req *http.Request) {
		response := map[string]interface{}{
			"host": mackerel.Host{},
		}
		_respondJSON(w, response)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	root, err := ioutil.TempDir("", "mackerel-agent-test")
	if err != nil {
		t.Fatal(err)
	}

	// test preparation done

	logging.ConfigureLoggers("DEBUG")
	conf := config.Config{
		Apibase: ts.URL,
		Root:    root,
	}

	Prepare(conf)
}
