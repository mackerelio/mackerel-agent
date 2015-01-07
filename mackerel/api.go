package mackerel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

// API XXX
type API struct {
	BaseURL *url.URL
	APIKey  string
	Verbose bool
}

// NewAPI XXX
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

// FindHost XXX
func (api *API) FindHost(id string) (*Host, error) {
	req, err := http.NewRequest("GET", api.urlFor(fmt.Sprintf("/api/v0/hosts/%s", id)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := api.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Host *Host `json:"host"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Host, err
}

// CreateHost XXX
func (api *API) CreateHost(name string, meta map[string]interface{}, interfaces []map[string]interface{}, roleFullnames []string) (string, error) {
	requestJSON, err := json.Marshal(map[string]interface{}{
		"name":          name,
		"type":          "unknown",
		"status":        "working",
		"meta":          meta,
		"interfaces":    interfaces,
		"roleFullnames": roleFullnames,
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
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API result failed: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		ID string `json:"id"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

// UpdateHost updates the host information on Mackerel.
func (api *API) UpdateHost(hostID string, name string, meta map[string]interface{}, interfaces []map[string]interface{}, roleFullnames []string) error {
	url := api.urlFor("/api/v0/hosts/" + hostID)

	requestJSON, err := json.Marshal(map[string]interface{}{
		"name":          name,
		"meta":          meta,
		"interfaces":    interfaces,
		"roleFullnames": roleFullnames,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url.String(), bytes.NewReader(requestJSON))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := api.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Debugf("Metrics Post Response: %s", string(body))
	defer resp.Body.Close()

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
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		logger.Warningf("Create graph defs response: %s", string(body))
	} else {
		logger.Debugf("Create graph defs response: %s", string(body))
	}

	return nil
}
