package mackerel

import (
	"fmt"
	"net/http"

	"github.com/mackerelio/golib/logging"
	mkr "github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mackerel-agent/checks"
)

var logger = logging.GetLogger("api")

// API is ...
type API interface {
	// implemented by mkr.Client
	CreateGraphDefs(payloads []*mkr.GraphDefsParam) error
	CreateHost(param *mkr.CreateHostParam) (string, error)
	FindHost(id string) (*mkr.Host, error)
	PostHostMetricValues(metricValues [](*mkr.HostMetricValue)) error
	PutHostMetaData(hostID, namespace string, metadata mkr.HostMetaData) error
	RetireHost(id string) error
	UpdateHost(hostID string, param *mkr.UpdateHostParam) (string, error)
	UpdateHostStatus(hostID string, status string) error
	// wrapped in APIImpl
	FindHostByCustomIdentifier(customIdentifier string) (*mkr.Host, error)
	ReportCheckMonitors(hostID string, reports []*checks.Report) error
	// accessor to Client
	SetAdditionalHeader(header http.Header)
	SetUserAgent(userAgent string)
}

// APIImpl is the main interface of Mackerel API.
type APIImpl struct {
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
func NewAPI(rawurl string, apiKey string, verbose bool) (API, error) {
	c, err := mkr.NewClientWithOptions(apiKey, rawurl, verbose)
	if err != nil {
		return nil, err
	}
	c.PrioritizedLogger = logger
	return &APIImpl{Client: c}, nil
}

// FindHostByCustomIdentifier find the host by the custom identifier
func (api *APIImpl) FindHostByCustomIdentifier(customIdentifier string) (*mkr.Host, error) {
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

// SetUserAgent sets User-Agent
func (api *APIImpl) SetUserAgent(userAgent string) {
	api.UserAgent = userAgent
}

// SetAdditionalHeader sets AdditionalHeader
func (api *APIImpl) SetAdditionalHeader(header http.Header) {
	api.AdditionalHeaders = header
}
