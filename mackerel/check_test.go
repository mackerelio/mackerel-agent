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
		{
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

// TestReportCheckMonitorsCompat tests to be equality of before and after migration to mackerel-client.
func TestReportCheckMonitorsCompat(t *testing.T) {
	// We can't use mkr.CheckReports because mkr.CheckReports.Reports.Source is interface.
	type Report struct {
		NotificationInterval uint `json:"notificationInterval,omitempty"`
		MaxCheckAttempts     uint `json:"maxCheckAttempts,omitempty"`
	}
	type Reports struct {
		Reports []*Report `json:"reports"`
	}

	received := make(chan *Report, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := json.NewDecoder(r.Body)
		var a Reports
		if err := d.Decode(&a); err != nil {
			t.Fatalf("can't decode: %v", err)
		}
		if n := len(a.Reports); n != 1 {
			t.Fatalf("len(Reports) = %d; want 1", n)
		}
		received <- a.Reports[0]
		fmt.Fprintf(w, "OK")
	}))
	defer ts.Close()

	n32 := func(n int32) *int32 { return &n }
	tests := []struct {
		name       string
		value      *int32
		interval   uint
		maxAttemts uint
	}{
		{
			name:       "case 0",
			value:      n32(0),
			interval:   NotificationIntervalFallback,
			maxAttemts: 0,
		},
		{
			name:       "case 1",
			value:      n32(1),
			interval:   NotificationIntervalFallback,
			maxAttemts: 1,
		},
		{
			name:       "case 9",
			value:      n32(9),
			interval:   NotificationIntervalFallback,
			maxAttemts: 9,
		},
		{
			name:       "case 10",
			value:      n32(10),
			interval:   NotificationIntervalFallback,
			maxAttemts: 10,
		},
		{
			name:       "case 11",
			value:      n32(11),
			interval:   11,
			maxAttemts: 11,
		},
		{
			name:       "case -1",
			value:      n32(-1),
			interval:   NotificationIntervalFallback,
			maxAttemts: 0,
		},
		{
			name:       "case nil",
			value:      nil,
			interval:   0,
			maxAttemts: 0,
		},
	}
	api, _ := NewAPI(ts.URL, "dummy-key", false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api.ReportCheckMonitors("xxx", []*checks.Report{
				{
					NotificationInterval: tt.value,
					MaxCheckAttempts:     tt.value,
				},
			})
			r := <-received
			if r == nil {
				return
			}
			if n := r.NotificationInterval; n != tt.interval {
				t.Errorf("NotificationInterval = %d; want %d", n, tt.interval)
			}
			if n := r.MaxCheckAttempts; n != tt.maxAttemts {
				t.Errorf("MaxCheckAttempts = %d; want %d", n, tt.maxAttemts)
			}
		})
	}
}
