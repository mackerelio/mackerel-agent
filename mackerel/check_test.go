package mackerel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/checks"
)

func TestReportCheckMonitors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/monitoring/checks/report" {
			t.Error("request URL should be /api/v0/monitoring/checks/report but :", req.URL.Path)
		}
		if req.Method != "POST" {
			t.Error("request method should be POST but :", req.Method)
		}

		body, _ := ioutil.ReadAll(req.Body)
		content := string(body)
		type testrepo struct {
			Source     map[string]string `json:"source"`
			Name       string            `json:"name"`
			Status     string            `json:"status"`
			Message    string            `json:"message"`
			OccurredAt float64           `json:"occurredAt"`
		}
		var data struct {
			Reports []testrepo `json:"reports"`
		}

		err := json.Unmarshal(body, &data)
		if err != nil {
			t.Fatal("request content should be decoded as json", content)
		}

		if reflect.DeepEqual(data.Reports[0], testrepo{
			Source: map[string]string{
				"hostId": "9rxGOHfVF8F",
				"type":   "host",
			},
			Name:       "sabasaba",
			Status:     "OK",
			Message:    "mesmes",
			OccurredAt: 0,
		}) != true {
			t.Error("report format invalid: ", data.Reports[0])
		}

		respJSON, _ := json.Marshal(map[string]string{
			"result": "OK",
		})
		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	api, _ := NewAPI(ts.URL, "dummy-key", false)

	err := api.ReportCheckMonitors("9rxGOHfVF8F", []*checks.Report{
		&checks.Report{
			Name:       "sabasaba",
			Status:     checks.StatusOK,
			Message:    "mesmes",
			OccurredAt: time.Unix(0, 0),
		},
	})

	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}
}
