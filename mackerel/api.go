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

// CreatingMetricsValue parameters of metric values
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
	newURL, _ := url.Parse(api.BaseURL.String())
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

// FindHost find the host
func (api *API) FindHost(id string) (*Host, error) {
	resp, err := api.get(fmt.Sprintf("/api/v0/hosts/%s", id))
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

// CreateHost register the host to mackerel
func (api *API) CreateHost(name string, meta map[string]interface{}, interfaces []map[string]interface{}, roleFullnames []string, displayName string) (string, error) {
	resp, err := api.postJSON("/api/v0/hosts", map[string]interface{}{
		"name":          name,
		"type":          "unknown",
		"status":        "working",
		"meta":          meta,
		"interfaces":    interfaces,
		"roleFullnames": roleFullnames,
		"displayName":   displayName,
	})
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
	resp, err := api.putJSON("/api/v0/hosts/"+hostID, hostSpec)
	defer closeResp(resp)
	if err != nil {
		return err
	}

	return nil
}

// UpdateHostStatus updates the status of the host
func (api *API) UpdateHostStatus(hostId string, status string) error {
	resp, err := api.postJSON(fmt.Sprintf("/api/v0/hosts/%s/status", hostId), map[string]string{
		"status": status,
	})
	defer closeResp(resp)
	if err != nil {
		return err
	}
	return nil
}

// PostMetricsValues post metrics
func (api *API) PostMetricsValues(metricsValues [](*CreatingMetricsValue)) error {
	resp, err := api.postJSON("/api/v0/tsdb", metricsValues)
	defer closeResp(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("API result failed: %s", resp.Status)
	}

	return nil
}

// CreateGraphDefsPayload payload for post graph defs
type CreateGraphDefsPayload struct {
	Name        string                         `json:"name"`
	DisplayName string                         `json:"displayName"`
	Unit        string                         `json:"unit"`
	Metrics     []CreateGraphDefsPayloadMetric `json:"metrics"`
}

// CreateGraphDefsPayloadMetric repreesnt graph defs of each metric
type CreateGraphDefsPayloadMetric struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IsStacked   bool   `json:"isStacked"`
}

// CreateGraphDefs register graph defs
func (api *API) CreateGraphDefs(payloads []CreateGraphDefsPayload) error {
	resp, err := api.postJSON("/api/v0/graph-defs/create", payloads)
	defer closeResp(resp)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", api.urlFor(path).String(), nil)
	if err != nil {
		return nil, err
	}
	return api.do(req)
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
