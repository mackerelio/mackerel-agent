package mackerel

import (
	"fmt"
	"net/http"
	"time"

	mkr "github.com/mackerelio/mackerel-client-go"
)

// API is the main interface of Mackerel API.
type API struct {
	*mkr.Client
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
	c, err := mkr.NewClientWithOptions(apiKey, rawurl, verbose)
	if err != nil {
		return nil, err
	}
	c.AdditionalHeaders = make(http.Header)
	// TODO(lufia): should we set a timeout explicitly?
	//c.HTTPClient.Timeout = apiRequestTimeout
	return &API{Client: c}, nil
}

var apiRequestTimeout = 30 * time.Second

// FindHost find the host
func (api *API) FindHost(id string) (*mkr.Host, error) {
	return api.Client.FindHost(id)
}

// FindHostByCustomIdentifier find the host by the custom identifier
func (api *API) FindHostByCustomIdentifier(customIdentifier string) (*mkr.Host, error) {
	param := mkr.FindHostsParam{
		CustomIdentifier: customIdentifier,
		Statuses:         []string{"working", "standby", "maintenance", "poweroff"},
	}
	hosts, err := api.Client.FindHosts(&param)
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
	return api.Client.CreateHost(hostParam)
}

// UpdateHost updates the host information on Mackerel.
func (api *API) UpdateHost(hostID string, hostParam *mkr.UpdateHostParam) error {
	_, err := api.Client.UpdateHost(hostID, hostParam)
	return err
}

// UpdateHostStatus updates the status of the host
func (api *API) UpdateHostStatus(hostID string, status string) error {
	return api.Client.UpdateHostStatus(hostID, status)
}

// PostMetricValues post metrics
func (api *API) PostMetricValues(metricsValues [](*mkr.HostMetricValue)) error {
	return api.Client.PostHostMetricValues(metricsValues)
}

// CreateGraphDefs register graph defs
func (api *API) CreateGraphDefs(payloads []*mkr.GraphDefsParam) error {
	return api.Client.CreateGraphDefs(payloads)
}

// RetireHost retires the host
func (api *API) RetireHost(hostID string) error {
	return api.Client.RetireHost(hostID)
}
