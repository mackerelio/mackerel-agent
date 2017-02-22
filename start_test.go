package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/mackerel"
)

func respJSON(w http.ResponseWriter, data map[string]interface{}) {
	respJSON, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(respJSON))
}

func TestStart(t *testing.T) {
	hostID := "xxx1234567890"
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/hosts/"+hostID, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			respJSON(w, map[string]interface{}{
				"host": mackerel.Host{
					ID:     hostID,
					Name:   "host.example.com",
					Status: "standby",
				},
			})
		case "PUT":
			respJSON(w, map[string]interface{}{
				"result": "OK",
			})
		default:
			t.Errorf("request method should be PUT or GET but :%s", r.Method)
		}
	})
	mux.HandleFunc("/api/v0/tsdb", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			payload := []mackerel.CreatingMetricsValue{}
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				t.Errorf("decode failed: %s", err)
			}

			respJSON(w, map[string]interface{}{
				"success": true,
			})
		default:
			t.Errorf("request method should be POST but: %s", r.Method)
		}
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	root, err := ioutil.TempDir("", "mackerel-config-test")
	if err != nil {
		t.Fatalf("Could not create temporary dir for test")
	}
	defer os.RemoveAll(root)

	confFile, err := os.Create(filepath.Join(root, "mackerel-agent.conf"))
	if err != nil {
		t.Fatalf("Could not create temporary file for test")
	}
	confFile.WriteString(`apikey="DUMMYAPIKEY"` + "\n")
	confFile.Sync()
	confFile.Close()
	argv := []string{
		"-conf=" + confFile.Name(),
		"-apibase=" + ts.URL,
		"-pidfile=" + root + "/pid",
		"-root=" + root,
		"-verbose",
	}
	conf, err := resolveConfig(&flag.FlagSet{}, argv)
	if err != nil {
		t.Errorf("err should be nil, but got: %s", err)
	}
	conf.SaveHostID(hostID)
	termCh := make(chan struct{})
	done := make(chan struct{})
	go func() {
		err = start(conf, termCh)
		done <- struct{}{}
	}()
	time.Sleep(5 * time.Second)
	termCh <- struct{}{}
	<-done
	if err != nil {
		t.Errorf("err should be nil but: %s", err)
	}
}
