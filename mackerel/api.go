package mackerel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/mackerelio/golib/logging"
	mkr "github.com/mackerelio/mackerel-client-go"
)

var logger = logging.GetLogger("api")

// API is the main interface of Mackerel API.
type API struct {
	BaseURL        *url.URL
	APIKey         string
	Verbose        bool
	UA             string
	DefaultHeaders http.Header

	c *mkr.Client
}

// Error represents API error
type Error struct {
	StatusCode int
	Message    string
}

func (aperr *Error) Error() string {
	return fmt.Sprintf("API error. status: %d, msg: %s", aperr.StatusCode, aperr.Message)
}

// IsClientError 4xx
func (aperr *Error) IsClientError() bool {
	return 400 <= aperr.StatusCode && aperr.StatusCode < 500
}

// IsClientError returns true if err is HTTP 4xx.
func IsClientError(err error) bool {
	e, ok := err.(*mkr.APIError)
	if !ok {
		return false
	}
	return 400 <= e.StatusCode && e.StatusCode < 500
}

// IsServerError 5xx
func (aperr *Error) IsServerError() bool {
	return 500 <= aperr.StatusCode && aperr.StatusCode < 600
}

// IsServerError returns true if err is HTTP 5xx.
func IsServerError(err error) bool {
	e, ok := err.(*mkr.APIError)
	if !ok {
		return false
	}
	return 500 <= e.StatusCode && e.StatusCode < 600
}

func apiError(code int, msg string) *Error {
	return &Error{
		StatusCode: code,
		Message:    msg,
	}
}

// InfoError represents Error of log level INFO
type InfoError struct {
	Message string
}

func (e *InfoError) Error() string {
	return e.Message
}

func infoError(msg string) *InfoError {
	return &InfoError{
		Message: msg,
	}
}

// NewAPI creates a new instance of API.
func NewAPI(rawurl string, apiKey string, verbose bool) (*API, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	c, err := mkr.NewClientWithOptions(apiKey, rawurl, verbose)
	if err != nil {
		return nil, err
	}
	return &API{BaseURL: u, APIKey: apiKey, Verbose: verbose, c: c}, nil
}

func (api *API) urlFor(path string, query string) *url.URL {
	newURL, _ := url.Parse(api.BaseURL.String())
	newURL.Path = path
	newURL.RawQuery = query
	return newURL
}

func (api *API) getUA() string {
	if api.UA != "" {
		return api.UA
	}
	return "mackerel-agent/0.0.0"
}

var apiRequestTimeout = 30 * time.Second

func (api *API) do(req *http.Request) (resp *http.Response, err error) {
	if api.DefaultHeaders != nil {
		for k, vs := range api.DefaultHeaders {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}
	req.Header.Add("X-Api-Key", api.APIKey)
	req.Header.Set("User-Agent", api.getUA())

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
func (api *API) FindHost(id string) (*mkr.Host, error) {
	resp, err := api.get(fmt.Sprintf("/api/v0/hosts/%s", id), "")
	defer closeResp(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, apiError(resp.StatusCode, "status code is not 200")
	}

	var data struct {
		Host *mkr.Host `json:"host"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.Host, err
}

// FindHostByCustomIdentifier find the host by the custom identifier
func (api *API) FindHostByCustomIdentifier(customIdentifier string) (*mkr.Host, error) {
	v := url.Values{}
	v.Set("customIdentifier", customIdentifier)
	for _, status := range []string{"working", "standby", "maintenance", "poweroff"} {
		v.Add("status", status)
	}
	resp, err := api.get("/api/v0/hosts", v.Encode())
	defer closeResp(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, apiError(resp.StatusCode, "status code is not 200")
	}

	var data struct {
		Hosts []*mkr.Host `json:"hosts"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if len(data.Hosts) == 0 {
		return nil, infoError(fmt.Sprintf("no host was found for the custom identifier: %s", customIdentifier))
	}
	return data.Hosts[0], err
}

// CreateHost register the host to mackerel
func (api *API) CreateHost(hostParam *mkr.CreateHostParam) (string, error) {
	return api.c.CreateHost(hostParam)
}

// UpdateHost updates the host information on Mackerel.
func (api *API) UpdateHost(hostID string, hostParam *mkr.UpdateHostParam) error {
	_, err := api.c.UpdateHost(hostID, hostParam)
	return err
}

// UpdateHostStatus updates the status of the host
func (api *API) UpdateHostStatus(hostID string, status string) error {
	return api.c.UpdateHostStatus(hostID, status)
}

// PostMetricValues post metrics
func (api *API) PostMetricValues(metricsValues [](*mkr.HostMetricValue)) error {
	resp, err := api.postJSON("/api/v0/tsdb", metricsValues)
	defer closeResp(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return apiError(resp.StatusCode, "api request failed")
	}

	return nil
}

// CreateGraphDefs register graph defs
func (api *API) CreateGraphDefs(payloads []*mkr.GraphDefsParam) error {
	resp, err := api.postJSON("/api/v0/graph-defs/create", payloads)
	defer closeResp(resp)
	if err != nil {
		return err
	}
	return nil
}

// RetireHost retires the host
func (api *API) RetireHost(hostID string) error {
	resp, err := api.postJSON(fmt.Sprintf("/api/v0/hosts/%s/retire", hostID), []byte("{}"))
	defer closeResp(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return apiError(resp.StatusCode, "api request failed")
	}
	return nil
}

func (api *API) get(path string, query string) (*http.Response, error) {
	req, err := http.NewRequest("GET", api.urlFor(path, query).String(), nil)
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

	req, err := http.NewRequest(method, api.urlFor(path, "").String(), &body)
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
		return resp, apiError(resp.StatusCode, "api request failed")
	}
	return resp, nil
}

func (api *API) postJSON(path string, payload interface{}) (*http.Response, error) {
	return api.requestJSON("POST", path, payload)
}
