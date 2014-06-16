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

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/version"
)

var logger = logging.GetLogger("api")

type CreatingMetricsValue struct {
	HostId string      `json:"hostId"`
	Name   string      `json:"name"`
	Time   float64     `json:"time"`
	Value  interface{} `json:"value"`
}

type API struct {
	BaseUrl *url.URL
	ApiKey  string
	Verbose bool
}

func NewApi(rawurl string, apiKey string, verbose bool) (*API, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &API{u, apiKey, verbose}, nil
}

func (api *API) urlFor(path string) *url.URL {
	newUrl, err := url.Parse(api.BaseUrl.String())
	if err != nil {
		panic("invalid url passed")
	}

	newUrl.Path = path

	return newUrl
}

func (api *API) Do(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add("X-Api-Key", api.ApiKey)
	req.Header.Add("X-Agent-Version", version.VERSION)
	req.Header.Add("X-Revision", version.GITCOMMIT)
	req.Header.Set("User-Agent", version.UserAgent())

	if api.Verbose {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			logger.Tracef("%s", dump)
		}
	}
	resp, err = http.DefaultClient.Do(req)
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

func (api *API) FindHost(id string) (*Host, error) {
	req, err := http.NewRequest("GET", api.urlFor(fmt.Sprintf("/api/v0/hosts/%s", id)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := api.Do(req)
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

func (api *API) CreateHost(name string, meta map[string]interface{}, interfaces []map[string]interface{}, roleFullnames []string) (string, error) {
	requestJson, err := json.Marshal(map[string]interface{}{
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
		bytes.NewReader(requestJson),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := api.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("API result failed: %s", resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.Id, nil
}

// UpdateHost updates the host information on Mackerel.
func (api *API) UpdateHost(hostId string, name string, meta map[string]interface{}, interfaces []map[string]interface{}, roleFullnames []string) error {
	url := api.urlFor("/api/v0/hosts/" + hostId)

	requestJson, err := json.Marshal(map[string]interface{}{
		"name":          name,
		"meta":          meta,
		"interfaces":    interfaces,
		"roleFullnames": roleFullnames,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url.String(), bytes.NewReader(requestJson))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := api.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (api *API) PostMetricsValues(metricsValues [](*CreatingMetricsValue)) error {
	requestJson, err := json.Marshal(metricsValues)
	if err != nil {
		return err
	}
	logger.Debugf("Metrics Post Request: %s", string(requestJson))

	req, err := http.NewRequest(
		"POST",
		api.urlFor("/api/v0/tsdb").String(),
		bytes.NewReader(requestJson),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := api.Do(req)
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
		return errors.New(fmt.Sprintf("API result failed: %s", resp.Status))
	}

	return nil
}
