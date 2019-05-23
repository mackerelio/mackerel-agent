package mackerel

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	mkr "github.com/mackerelio/mackerel-client-go"
)

// API is the main interface of Mackerel API.
type API struct {
	BaseURL        *url.URL
	APIKey         string
	Verbose        bool
	DefaultHeaders http.Header

	c *mkr.Client
}

// IsClientError returns true if err is HTTP 4xx.
func IsClientError(err error) bool {
	e, ok := err.(*mkr.APIError)
	if !ok {
		return false
	}
	return 400 <= e.StatusCode && e.StatusCode < 500
}

// IsServerError returns true if err is HTTP 5xx.
func IsServerError(err error) bool {
	e, ok := err.(*mkr.APIError)
	if !ok {
		return false
	}
	return 500 <= e.StatusCode && e.StatusCode < 600
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
	c.AdditionalHeaders = make(http.Header)
	// TODO(lufia): should we set a timeout explicitly?
	//c.HTTPClient.Timeout = apiRequestTimeout
	return &API{BaseURL: u, APIKey: apiKey, Verbose: verbose, c: c, DefaultHeaders: c.AdditionalHeaders}, nil
}

// SetUA is a temporary function to migrate to mackerel-client-go.
func (api *API) SetUA(s string) {
	api.c.UserAgent = s
}

var apiRequestTimeout = 30 * time.Second

// FindHost find the host
func (api *API) FindHost(id string) (*mkr.Host, error) {
	return api.c.FindHost(id)
}

// FindHostByCustomIdentifier find the host by the custom identifier
func (api *API) FindHostByCustomIdentifier(customIdentifier string) (*mkr.Host, error) {
	param := mkr.FindHostsParam{
		CustomIdentifier: customIdentifier,
		Statuses:         []string{"working", "standby", "maintenance", "poweroff"},
	}
	hosts, err := api.c.FindHosts(&param)
	if err != nil {
		return nil, err
	}
	if len(hosts) == 0 {
		return nil, infoError(fmt.Sprintf("no host was found for the custom identifier: %s", customIdentifier))
	}
	return hosts[0], nil
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
	return api.c.PostHostMetricValues(metricsValues)
}

// CreateGraphDefs register graph defs
func (api *API) CreateGraphDefs(payloads []*mkr.GraphDefsParam) error {
	return api.c.CreateGraphDefs(payloads)
}

// RetireHost retires the host
func (api *API) RetireHost(hostID string) error {
	return api.c.RetireHost(hostID)
}
