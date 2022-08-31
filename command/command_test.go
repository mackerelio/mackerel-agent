package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/agent"
	"github.com/mackerelio/mackerel-agent/checks"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
	mkr "github.com/mackerelio/mackerel-client-go"
)

func init() {
	// Shrink time scale
	originalPostMetricsInterval := config.PostMetricsInterval

	config.PostMetricsInterval = 10 * time.Second
	ratio := config.PostMetricsInterval.Seconds() / originalPostMetricsInterval.Seconds()

	postMetricsDequeueDelaySeconds =
		int(float64(postMetricsDequeueDelaySeconds) * ratio)
	postMetricsRetryDelaySeconds =
		int(float64(postMetricsRetryDelaySeconds) * ratio)

	reportCheckRetryDelaySeconds = 1
}

func TestDelayByHost(t *testing.T) {
	delay1 := time.Duration(delayByHost(&mkr.Host{
		ID:     "246PUVUngPo",
		Name:   "hogehoge2.host.h",
		Size:   "standard",
		Status: "working",
	})) * time.Second

	delay2 := time.Duration(delayByHost(&mkr.Host{
		ID:     "21GZjCE5Etb",
		Name:   "hogehoge2.host.h",
		Size:   "standard",
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
// and returns the Config, mock handlers map, the server and cleanup function which should be called finally.
// The mock handlers map is "<method> <path>"-to-jsonObject-generator map.
func newMockAPIServer(t *testing.T) (config.Config, map[string]func(*http.Request) (int, jsonObject), *httptest.Server, func()) {
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

	root, err := os.MkdirTemp("", "mackerel-agent-test")
	if err != nil {
		t.Fatal(err)
	}

	conf := config.Config{
		Apibase: ts.URL,
		Root:    root,
	}

	return conf, mockHandlers, ts, func() {
		ts.Close()
		os.RemoveAll(root)
	}
}

func TestPrepareWithCreate(t *testing.T) {
	conf, mockHandlers, ts, deferFunc := newMockAPIServer(t)
	defer deferFunc()

	// To prevent unexpected behavior in the TravisCI environment
	conf.CloudPlatform = config.CloudPlatformNone

	mockHandlers["POST /api/v0/hosts"] = func(req *http.Request) (int, jsonObject) {
		return 200, jsonObject{
			"id": "xxx1234567890",
		}
	}

	mockHandlers["GET /api/v0/hosts/xxx1234567890"] = func(req *http.Request) (int, jsonObject) {
		return 200, jsonObject{
			"host": mkr.Host{
				ID:     "xxx1234567890",
				Name:   "host.example.com",
				Size:   "standard",
				Status: "standby",
			},
		}
	}

	c, _ := Prepare(&conf, &AgentMeta{})
	api := c.API
	host := c.Host

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

func TestPrepareWithCreateWithFail(t *testing.T) {
	conf, mockHandlers, _, deferFunc := newMockAPIServer(t)
	defer deferFunc()

	// To prevent unexpected behavior in the TravisCI environment
	conf.CloudPlatform = config.CloudPlatformNone

	mockHandlers["POST /api/v0/hosts"] = func(req *http.Request) (int, jsonObject) {
		return 403, jsonObject{
			"result": "error",
		}
	}

	origRetryNum := retryNum
	retryNum = 1
	defer func() {
		retryNum = origRetryNum
	}()

	_, err := Prepare(&conf, &AgentMeta{})

	if err == nil {
		t.Errorf("error should be occurred")
	}
}

func TestPrepareWithUpdate(t *testing.T) {
	conf, mockHandlers, ts, deferFunc := newMockAPIServer(t)
	defer deferFunc()
	tempDir, _ := os.MkdirTemp("", "")
	conf.Root = tempDir
	conf.SaveHostID("xxx12345678901")

	mockHandlers["PUT /api/v0/hosts/xxx12345678901"] = func(req *http.Request) (int, jsonObject) {
		return 200, jsonObject{
			"result": "OK",
		}
	}

	mockHandlers["GET /api/v0/hosts/xxx12345678901"] = func(req *http.Request) (int, jsonObject) {
		return 200, jsonObject{
			"host": mkr.Host{
				ID:     "xxx12345678901",
				Name:   "host.example.com",
				Size:   "standard",
				Status: "standby",
			},
		}
	}

	c, _ := Prepare(&conf, &AgentMeta{})
	api := c.API
	host := c.Host

	if api.BaseURL.String() != ts.URL {
		t.Errorf("Apibase mismatch: %s != %s", api.BaseURL, ts.URL)
	}

	if host.ID != "xxx12345678901" {
		t.Error("Host ID mismatch", host)
	}

	if host.Name != "host.example.com" {
		t.Error("Host name mismatch", host)
	}
}

func TestCollectHostParam(t *testing.T) {
	conf := config.Config{}
	hostParam, err := collectHostParam(&conf, &AgentMeta{})

	if err != nil {
		t.Errorf("collectHostParam should not fail: %s", err)
	}

	if hostParam.Name == "" {
		t.Error("hostname should not be empty")
	}

	if len(hostParam.Meta.CPU) == 0 {
		t.Error("meta.cpu should exist")
	}

	if len(hostParam.Meta.Memory) == 0 {
		t.Error("meta.memory should exist")
	}

	if len(hostParam.Meta.Filesystem) == 0 {
		t.Error("meta.filesystem should exist")
	}

	if len(hostParam.Meta.Kernel) == 0 {
		t.Error("meta.kernel should exist")
	}
}

func TestCollectHostParamWithChecks(t *testing.T) {
	customIdentifier := "app.example.com"
	conf := config.Config{
		CheckPlugins: map[string]*config.CheckPlugin{
			"chk1": &config.CheckPlugin{
				CustomIdentifier: nil,
			},
			"chk2": &config.CheckPlugin{
				CustomIdentifier: &customIdentifier,
			},
		},
	}
	hostParam, err := collectHostParam(&conf, &AgentMeta{})

	if err != nil {
		t.Errorf("collectHostParam should not fail: %s", err)
	}

	if len(hostParam.Checks) != 1 {
		t.Error("only checks without customIdentifier should be included in param")
	}
	if hostParam.Checks[0].Name != "chk1" {
		t.Error("only checks without customIdentifier should be included in param")
	}
}

type counterGenerator struct {
	counter int
	sync.Mutex
}

func (g *counterGenerator) Generate() (metrics.Values, error) {
	g.Lock()
	defer g.Unlock()
	g.counter = g.counter + 1
	return map[string]float64{"dummy.a": float64(g.counter)}, nil
}

type byTime []mkr.HostMetricValue

func (b byTime) Len() int           { return len(b) }
func (b byTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byTime) Less(i, j int) bool { return b[i].Time < b[j].Time }

func TestLoop(t *testing.T) {
	if testing.Verbose() {
		logging.SetLogLevel(logging.DEBUG)
	}

	conf, mockHandlers, _, deferFunc := newMockAPIServer(t)
	defer deferFunc()

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
	receivedDataPoints := []mkr.HostMetricValue{}
	done := make(chan struct{})

	mockHandlers["POST /api/v0/tsdb"] = func(req *http.Request) (int, jsonObject) {
		payload := []mkr.HostMetricValue{}
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
	mockHandlers["PUT /api/v0/hosts/xyzabc12345"] = func(req *http.Request) (int, jsonObject) {
		return 200, jsonObject{
			"result": "OK",
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

	host := &mkr.Host{ID: "xyzabc12345"}

	termCh := make(chan struct{})
	exitCh := make(chan error)
	app := &App{
		Agent:     ag,
		Config:    &conf,
		API:       api,
		Host:      host,
		AgentMeta: &AgentMeta{},
	}
	// Start looping!
	go func() {
		exitCh <- loop(app, termCh)
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
	exitErr := <-exitCh
	if exitErr != nil {
		t.Errorf("exitErr should be nil, got: %s", exitErr)
	}
}

func TestLoop_NetworkError(t *testing.T) {
	if testing.Verbose() {
		logging.SetLogLevel(logging.DEBUG)
	}

	conf, mockHandlers, _, deferFunc := newMockAPIServer(t)
	defer deferFunc()

	// 1(network error) =(retry immediately)=> 2(network error) =(wait while retry)=> 3(success)
	const (
		totalFailures = 2
		totalPosts    = 3
	)
	postCount := 0
	done := make(chan struct{})

	mockHandlers["POST /api/v0/tsdb"] = func(req *http.Request) (int, jsonObject) {
		payload := []mkr.HostMetricValue{}
		json.NewDecoder(req.Body).Decode(&payload)

		for _, p := range payload {
			value := p.Value.(float64)
			if value == 1 {
				postCount++
				if postCount <= totalFailures {
					// returning StatusSeeOther without Location header causes url.Error
					return http.StatusSeeOther, jsonObject{}
				}

				if postCount == totalPosts {
					defer func() { done <- struct{}{} }()
				}
			}
		}

		return 200, jsonObject{
			"success": true,
		}
	}
	mockHandlers["PUT /api/v0/hosts/xyzabc12345"] = func(req *http.Request) (int, jsonObject) {
		return 200, jsonObject{
			"result": "OK",
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

	host := &mkr.Host{ID: "xyzabc12345"}

	termCh := make(chan struct{})
	exitCh := make(chan error)
	app := &App{
		Agent:     ag,
		Config:    &conf,
		API:       api,
		Host:      host,
		AgentMeta: &AgentMeta{},
	}
	// Start looping!
	go func() {
		exitCh <- loop(app, termCh)
	}()

	<-done
	termCh <- struct{}{}
	exitErr := <-exitCh
	if exitErr != nil {
		t.Errorf("exitErr should be nil, got: %s", exitErr)
	}
}

func TestLoop_NetworkErrorWithRecovery(t *testing.T) {
	if testing.Verbose() {
		logging.SetLogLevel(logging.DEBUG)
	}

	conf, mockHandlers, _, deferFunc := newMockAPIServer(t)
	defer deferFunc()

	// expected scenario.
	// datapoint 1: 1(network error) =(retry immediately)=> 2(success)
	const (
		totalFailures = 1
		totalPosts    = 2
	)
	postCount := 0
	done := make(chan struct{})

	mockHandlers["POST /api/v0/tsdb"] = func(req *http.Request) (int, jsonObject) {
		payload := []mkr.HostMetricValue{}
		json.NewDecoder(req.Body).Decode(&payload)

		for _, p := range payload {
			value := p.Value.(float64)
			if value == 1 {
				postCount++
				if postCount <= totalFailures {
					// returning StatusSeeOther without Location header causes url.Error
					return http.StatusSeeOther, jsonObject{}
				}

				if postCount == totalPosts {
					defer func() { done <- struct{}{} }()
				}
			}
		}

		return 200, jsonObject{
			"success": true,
		}
	}
	mockHandlers["PUT /api/v0/hosts/xyzabc12345"] = func(req *http.Request) (int, jsonObject) {
		return 200, jsonObject{
			"result": "OK",
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

	host := &mkr.Host{ID: "xyzabc12345"}

	termCh := make(chan struct{})
	exitCh := make(chan error)
	app := &App{
		Agent:     ag,
		Config:    &conf,
		API:       api,
		Host:      host,
		AgentMeta: &AgentMeta{},
	}
	// Start looping!
	go func() {
		exitCh <- loop(app, termCh)
	}()

	<-done
	termCh <- struct{}{}
	exitErr := <-exitCh
	if exitErr != nil {
		t.Errorf("exitErr should be nil, got: %s", exitErr)
	}
}

func TestReportCheckMonitors(t *testing.T) {
	if testing.Verbose() {
		logging.SetLogLevel(logging.DEBUG)
	}

	cases := []struct {
		Status      int
		expectRetry bool
	}{
		{http.StatusOK, false},
		{http.StatusBadRequest, false},
		{http.StatusInternalServerError, true},
	}

	for _, tc := range cases {
		conf, mockHandlers, _, deferFunc := newMockAPIServer(t)
		defer deferFunc()

		postCount := 0
		retried := false
		mu := &sync.Mutex{}

		mockHandlers["POST /api/v0/monitoring/checks/report"] = func(req *http.Request) (int, jsonObject) {
			mu.Lock()
			defer mu.Unlock()
			postCount++
			if postCount > 1 {
				retried = true
			}
			return tc.Status, jsonObject{}
		}

		api, err := mackerel.NewAPI(conf.Apibase, conf.Apikey, true)
		if err != nil {
			t.Fatal(err)
		}

		host := &mkr.Host{ID: "xyzabc12345"}

		app := &App{
			Agent:     &agent.Agent{},
			Config:    &conf,
			API:       api,
			Host:      host,
			AgentMeta: &AgentMeta{},
		}

		go func() {
			reportCheckMonitors(app, "", []*checks.Report{})
		}()

		time.Sleep(time.Duration(reportCheckRetryDelaySeconds) * 3 * time.Second)

		mu.Lock()
		defer mu.Unlock()

		if retried != tc.expectRetry {
			text := http.StatusText(tc.Status)
			if tc.expectRetry {
				t.Errorf("the agent should have resend reports when got %d %q", tc.Status, text)
			} else {
				t.Errorf("the agent should not have resend reports when got %d %q", tc.Status, text)
			}
		}
	}
}

func TestReportCheckMonitors_NetworkError(t *testing.T) {
	if testing.Verbose() {
		logging.SetLogLevel(logging.DEBUG)
	}

	conf, mockHandlers, _, deferFunc := newMockAPIServer(t)
	defer deferFunc()

	postCount := 0
	mu := &sync.Mutex{}

	mockHandlers["POST /api/v0/monitoring/checks/report"] = func(req *http.Request) (int, jsonObject) {
		mu.Lock()
		defer mu.Unlock()
		postCount++
		// returning StatusSeeOther without Location header causes url.Error
		return http.StatusSeeOther, jsonObject{}
	}

	api, err := mackerel.NewAPI(conf.Apibase, conf.Apikey, true)
	if err != nil {
		t.Fatal(err)
	}

	host := &mkr.Host{ID: "xyzabc12345"}

	app := &App{
		Agent:     &agent.Agent{},
		Config:    &conf,
		API:       api,
		Host:      host,
		AgentMeta: &AgentMeta{},
	}

	expectedPostCount := 0
	checkPostCount := func() {
		mu.Lock()
		defer mu.Unlock()
		if postCount != expectedPostCount {
			if postCount != expectedPostCount {
				t.Errorf("the agent should try sending reports %d times, but %d", expectedPostCount, postCount)
			}
		}
	}

	go func() {
		reportCheckMonitors(app, "", []*checks.Report{})
	}()

	// wait 1 cycle
	time.Sleep(1 * time.Second)
	expectedPostCount += 2
	checkPostCount()

	// wait another cycle
	time.Sleep(time.Duration(reportCheckRetryDelaySeconds) * time.Second)
	expectedPostCount += 2
	checkPostCount()
}

func TestReportCheckMonitors_NetworkErrorWithRecovery(t *testing.T) {
	if testing.Verbose() {
		logging.SetLogLevel(logging.DEBUG)
	}

	conf, mockHandlers, _, deferFunc := newMockAPIServer(t)
	defer deferFunc()

	postCount := 0
	mu := &sync.Mutex{}

	// returning StatusSeeOther without Location header causes url.Error
	// therefore url.Error happens on first request, and no error happens afterwards
	mockHandlers["POST /api/v0/monitoring/checks/report"] = func(req *http.Request) (int, jsonObject) {
		mu.Lock()
		defer mu.Unlock()
		postCount++
		if postCount == 1 {
			return http.StatusSeeOther, jsonObject{}
		}
		return http.StatusOK, jsonObject{}
	}

	api, err := mackerel.NewAPI(conf.Apibase, conf.Apikey, true)
	if err != nil {
		t.Fatal(err)
	}

	host := &mkr.Host{ID: "xyzabc12345"}

	app := &App{
		Agent:     &agent.Agent{},
		Config:    &conf,
		API:       api,
		Host:      host,
		AgentMeta: &AgentMeta{},
	}

	reportCheckMonitors(app, "", []*checks.Report{})

	expectedPostCount := 2
	if postCount != expectedPostCount {
		if postCount != expectedPostCount {
			t.Errorf("the agent should try sending reports %d times, but %d", expectedPostCount, postCount)
		}
	}

}
