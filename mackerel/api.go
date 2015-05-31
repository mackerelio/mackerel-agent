package mackerel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/version"
)

var logger = logging.GetLogger("api")

// CreatingMetricsValue XXX
type CreatingMetricsValue struct {
	HostID string      `json:"hostId"`
	Name   string      `json:"name"`
	Time   float64     `json:"time"`
	Value  interface{} `json:"value"`
}

// API is the main interface of Mackerel API.
type API struct {
	BaseURL *url.URL
	APIKey  string
	Verbose bool
}

// NewAPI creates a new instance of API.
func NewAPI(rawurl string, apiKey string, verbose bool) (*API, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &API{u, apiKey, verbose}, nil
}

func (api *API) urlFor(path string) *url.URL {
	newURL, err := url.Parse(api.BaseURL.String())
	if err != nil {
		panic("invalid url passed")
	}

	newURL.Path = path
	return newURL
}

var apiRequestTimeout = 30 * time.Second

func (api *API) do(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add("X-Api-Key", api.APIKey)
	req.Header.Add("X-Agent-Version", version.VERSION)
	req.Header.Add("X-Revision", version.GITCOMMIT)
	req.Header.Set("User-Agent", version.UserAgent())

	if api.Verbose {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			logger.Tracef("%s", dump)
		}
	}

	client := &http.Client{} // same as http.DefaultClient
	client.Timeout = apiRequestTimeout
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if api.Verbose {
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			logger.Tracef("%s", dump)
		}
	}
	return resp, nil
}

func closeResp(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}

// FindHost XXX
func (api *API) FindHost(id string) (*Host, error) {
	req, err := http.NewRequest("GET", api.urlFor(fmt.Sprintf("/api/v0/hosts/%s", id)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := api.do(req)
	defer closeResp(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	var data struct {
		Host *Host `json:"host"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Host, err
}

// CreateHost XXX
func (api *API) CreateHost(name string, meta map[string]interface{}, interfaces []map[string]interface{}, roleFullnames []string, displayName string) (string, error) {
	requestJSON, err := json.Marshal(map[string]interface{}{
		"name":          name,
		"type":          "unknown",
		"status":        "working",
		"meta":          meta,
		"interfaces":    interfaces,
		"roleFullnames": roleFullnames,
		"displayName":   displayName,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		api.urlFor("/api/v0/hosts").String(),
		bytes.NewReader(requestJSON),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := api.do(req)
	defer closeResp(resp)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API result failed: %s", resp.Status)
	}

	var data struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

// UpdateHost updates the host information on Mackerel.
func (api *API) UpdateHost(hostID string, hostSpec HostSpec) error {
	url := api.urlFor("/api/v0/hosts/" + hostID)

	requestJSON, err := json.Marshal(hostSpec)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url.String(), bytes.NewReader(requestJSON))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := api.do(req)
	defer closeResp(resp)
	if err != nil {
		return err
	}

	return nil
}

// PostMetricsValues XXX
func (api *API) PostMetricsValues(metricsValues [](*CreatingMetricsValue)) error {
	requestJSON, err := json.Marshal(metricsValues)
	if err != nil {
		return err
	}
	logger.Debugf("Metrics Post Request: %s", string(requestJSON))

	req, err := http.NewRequest(
		"POST",
		api.urlFor("/api/v0/tsdb").String(),
		bytes.NewReader(requestJSON),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := api.do(req)
	defer closeResp(resp)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API result failed: %s", resp.Status)
	}

	return nil
}

// CreateGraphDefsPayload XXX
type CreateGraphDefsPayload struct {
	Name        string                         `json:"name"`
	DisplayName string                         `json:"displayName"`
	Unit        string                         `json:"unit"`
	Metrics     []CreateGraphDefsPayloadMetric `json:"metrics"`
}

// CreateGraphDefsPayloadMetric XXX
type CreateGraphDefsPayloadMetric struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IsStacked   bool   `json:"isStacked"`
}

// CreateGraphDefs XXX
func (api *API) CreateGraphDefs(payloads []CreateGraphDefsPayload) error {
	requestJSON, err := json.Marshal(payloads)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		api.urlFor("/api/v0/graph-defs/create").String(),
		bytes.NewReader(requestJSON),
	)
	if err != nil {
		return err
	}

	logger.Debugf("Create grpah defs request: %s", string(requestJSON))

	req.Header.Add("Content-Type", "application/json")
	resp, err := api.do(req)
	defer closeResp(resp)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) requestJSON(method, path string, payload interface{}) (*http.Response, error) {
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return nil, err
	}
	logger.Debugf("%s %s %s", method, path, body.String())

	req, err := http.NewRequest(method, api.urlFor(path).String(), &body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := api.do(req)
	if err != nil {
		return resp, err
	}

	logger.Debugf("%s %s status=%q", method, path, resp.Status)
	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("request failed: [%s]", resp.Status)
	}
	return resp, nil
}

func (api *API) postJSON(path string, payload interface{}) (*http.Response, error) {
	return api.requestJSON("POST", path, payload)
}

func (api *API) putJSON(path string, payload interface{}) (*http.Response, error) {
	return api.requestJSON("PUT", path, payload)
}

// Time is a type for sending time information to Mackerel API server.
// It is encoded as an epoch seconds integer in JSON.
type Time time.Time

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Unix())
}
