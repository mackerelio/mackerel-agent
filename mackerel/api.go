package mackerel

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mackerelio/golib/logging"
	mkr "github.com/mackerelio/mackerel-client-go"
)

var logger = logging.GetLogger("api")

// API is the main interface of Mackerel API.
type API struct {
	*mkr.Client
}

// IsNetworkError returns true if err is url.Error caused by net/http
func IsNetworkError(err error) bool {
	var e *url.Error
	return errors.As(err, &e)
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
func NewAPI(rawurl string, apiKey string, verbose bool, disableKeepAlive bool) (*API, error) {
	c, err := mkr.NewClientWithOptions(apiKey, rawurl, verbose)
	if err != nil {
		return nil, err
	}
	c.PrioritizedLogger = logger
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = disableKeepAlive
	c.HTTPClient.Transport = t

	return &API{Client: c}, nil
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
