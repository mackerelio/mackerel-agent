package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/agent"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
)

func TestBuildHostSpec(t *testing.T) {
	var interfaces []map[string]interface{}
	interfaces = append(interfaces, map[string]interface{}{
		"name":       "eth0",
		"ipAddress":  "10.0.4.7",
		"macAddress": "01:23:45:67:89:ab",
		"encap":      "Ethernet",
	})
	meta := map[string]interface{}{
		"memo": "hello",
	}
	roleFullnames := []string{"My-Service:app-default"}
	checks := []string{"heartbeat"}
	displayName := "label"

	hostSpec := buildHostSpec("Host123", meta, interfaces, roleFullnames, checks, displayName)

	if hostSpec["name"] != "Host123" {
		t.Error("name should 'Host123' but:", hostSpec["name"])
	}

	_, ok := hostSpec["interfaces"].([]map[string]interface{})
	if !ok {
		t.Error("the type of interfaces should be '[]map[string]interface{}'")
	}

	_, ok = hostSpec["meta"].(map[string]interface{})
	if !ok {
		t.Error("the type of meta should be 'map[string]interface{}'")
	}

	_, ok = hostSpec["roleFullnames"].([]string)
	if !ok {
		t.Error("the type of meta should be '[]string'")
	}

	if hostSpec["displayName"] != displayName {
		t.Error("name should 'label' but:", hostSpec["displayName"])
	}
}

func TestDelayByHost(t *testing.T) {
	delay1 := time.Duration(delayByHost(&mackerel.Host{
		ID:     "246PUVUngPo",
		Name:   "hogehoge2.host.h",
		Type:   "unknown",
		Status: "working",
	})) * time.Second

	delay2 := time.Duration(delayByHost(&mackerel.Host{
		ID:     "21GZjCE5Etb",
		Name:   "hogehoge2.host.h",
		Type:   "unknown",
		Status: "working",
	})) * time.Second

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
func newMockAPIServer(t *testing.T) (config.Config, map[string]func(*http.Request) (int, jsonObject), *httptest.Server) {
	mockHandlers := map[string]func(*http.Request) (int, jsonObject){}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := req.Method + " " + req.URL.Path
		handler, ok := mockHandlers[key]
		if !ok {
			t.Fatal("Unexpected request: " + key)
		}

		statusCode, data := handler(req)

		respJSON, err := json.Marshal(data)
		if err != nil {
			t.Fatal("marshalling JSON failed: ", err)
		}

		if statusCode != 0 {
			w.WriteHeader(statusCode)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(respJSON))
	}))

	root, err := ioutil.TempDir("", "mackerel-agent-test")
	if err != nil {
		t.Fatal(err)
	}

	conf := config.Config{
		Apibase:    ts.URL,
		Root:       root,
		Connection: config.DefaultConfig.Connection,
	}

	return conf, mockHandlers, ts
}

func TestPrepare(t *testing.T) {
	conf, mockHandlers, ts := newMockAPIServer(t)
	defer ts.Close()

	mockHandlers["POST /api/v0/hosts"] = func(req *http.Request) (int, jsonObject) {
		return 200, jsonObject{
			"id": "xxx1234567890",
		}
	}

	mockHandlers["GET /api/v0/hosts/xxx1234567890"] = func(req *http.Request) (int, jsonObject) {
		return 200, jsonObject{
			"host": mackerel.Host{
				ID:     "xxx1234567890",
				Name:   "host.example.com",
				Type:   "unknown",
				Status: "standby",
			},
		}
	}

	api, host, _ := Prepare(&conf)

	if api.BaseURL.String() != ts.URL {
		t.Errorf("Apibase mismatch: %s != %s", api.BaseURL, ts.URL)
	}

	if host.ID != "xxx1234567890" {
		t.Error("Host ID mismatch", host)
	}

	if host.Name != "host.example.com" {
		t.Error("Host name mismatch", host)
	}
}

func TestCollectHostSpecs(t *testing.T) {
	hostname, meta, _ /*interfaces*/, err := collectHostSpecs()

	if err != nil {
		t.Errorf("collectHostSpecs should not fail: %s", err)
	}

	if hostname == "" {
		t.Error("hostname should not be empty")
	}

	if _, ok := meta["cpu"]; !ok {
		t.Error("meta.cpu should exist")
	}
}

type counterGenerator struct {
	counter int
}

func (g *counterGenerator) Generate() (metrics.Values, error) {
	g.counter = g.counter + 1
	return map[string]float64{"dummy.a": float64(g.counter)}, nil
}

type byTime []mackerel.CreatingMetricsValue

func (b byTime) Len() int           { return len(b) }
func (b byTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byTime) Less(i, j int) bool { return b[i].Time < b[j].Time }

func TestLoop(t *testing.T) {
	if testing.Verbose() {
		logging.ConfigureLoggers("DEBUG")
	}

	conf, mockHandlers, ts := newMockAPIServer(t)
	defer ts.Close()

	if testing.Short() {
		// Shrink time scale
		originalPostMetricsInterval := config.PostMetricsInterval

		config.PostMetricsInterval = 10 * time.Second
		ratio := config.PostMetricsInterval.Seconds() / originalPostMetricsInterval.Seconds()

		conf.Connection.PostMetricsDequeueDelaySeconds =
			int(float64(config.DefaultConfig.Connection.PostMetricsRetryDelaySeconds) * ratio)
		conf.Connection.PostMetricsRetryDelaySeconds =
			int(float64(config.DefaultConfig.Connection.PostMetricsRetryDelaySeconds) * ratio)

		defer func() {
			config.PostMetricsInterval = originalPostMetricsInterval
		}()
	}

	/// Simulate the situation that mackerel.io is down for 3 min
	// Strategy:
	// counterGenerator generates values 1,2,3,4,...
	// when we got value 3, the server will start responding 503 for three times (inclusive)
	// so the agent should queue the generated values and retry sending.
	//
	//  status: o . o . x . x . x . o o o o o
	//    send: 1 . 2 . 3 . 3 . 3 . 3 4 5 6 7
	// collect: 1 . 2 . 3 . 4 . 5 . 6 . 7 . 8
	//           ^
	//           30s
	const (
		totalFailures = 3
		totalPosts    = 7
	)
	failureCount := 0
	receivedDataPoints := []mackerel.CreatingMetricsValue{}
	done := make(chan struct{})

	mockHandlers["POST /api/v0/tsdb"] = func(req *http.Request) (int, jsonObject) {
		payload := []mackerel.CreatingMetricsValue{}
		json.NewDecoder(req.Body).Decode(&payload)

		for _, p := range payload {
			value := p.Value.(float64)
			if value == 3 {
				failureCount++
				if failureCount <= totalFailures {
					return 503, jsonObject{
						"failure": failureCount, // just for DEBUG logging
					}
				}
			}

			if value == totalPosts {
				defer func() { done <- struct{}{} }()
			}
		}

		receivedDataPoints = append(receivedDataPoints, payload...)

		return 200, jsonObject{
			"success": true,
		}
	}

	// Prepare required objects...
	ag := &agent.Agent{
		MetricsGenerators: []metrics.Generator{
			&counterGenerator{},
		},
	}

	api, err := mackerel.NewAPI(conf.Apibase, conf.Apikey, true)
	if err != nil {
		t.Fatal(err)
	}

	host := &mackerel.Host{ID: "xyzabc12345"}

	termCh := make(chan struct{})
	exitCh := make(chan int)
	c := &context{
		ag:   ag,
		conf: &conf,
		api:  api,
		host: host,
	}
	// Start looping!
	go func() {
		exitCh <- loop(c, termCh)
	}()

	<-done

	// Verify results
	if len(receivedDataPoints) != totalPosts {
		t.Errorf("the agent should have sent %d datapoints, got: %+v", totalPosts, receivedDataPoints)
	}

	sort.Sort(byTime(receivedDataPoints))

	for i := 0; i < totalPosts; i++ {
		value := receivedDataPoints[i].Value.(float64)
		if value != float64(i+1) {
			t.Errorf("the %dth datapoint should have value %d, got: %+v", i, i+1, receivedDataPoints)
		}
	}

	termCh <- struct{}{}
	exitCode := <-exitCh
	if exitCode != 0 {
		t.Errorf("exit code should be 0, got: %d", exitCode)
	}
}
