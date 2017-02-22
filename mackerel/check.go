package mackerel

import (
	"encoding/json"

	"github.com/mackerelio/mackerel-agent/checks"
)

type monitoringChecksPayload struct {
	Reports []*checkReport `json:"reports"`
}

type checkReport struct {
	Source               monitorTargetHost `json:"source"`
	Name                 string            `json:"name"`
	Status               checks.Status     `json:"status"`
	Message              string            `json:"message"`
	OccurredAt           Time              `json:"occurredAt"`
	NotificationInterval *int32            `json:"notificationInterval,omitempty"`
	MaxCheckAttempts     *int32            `json:"maxCheckAttempts,omitempty"`
}

type monitorTargetHost struct {
	HostID string
}

// MarshalJSON implements json.Marshaler.
func (h monitorTargetHost) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":   "host",
		"hostId": h.HostID,
	})
}

// ReportCheckMonitors sends reports of *checks.Checker() to Mackrel API server.
func (api *API) ReportCheckMonitors(hostID string, reports []*checks.Report) error {
	payload := &monitoringChecksPayload{
		Reports: make([]*checkReport, len(reports)),
	}
	const messageLengthLimit = 1024
	for i, report := range reports {
		msg := report.Message
		runes := []rune(msg)
		if len(runes) > messageLengthLimit {
			msg = string(runes[0:messageLengthLimit])
		}
		payload.Reports[i] = &checkReport{
			Source:               monitorTargetHost{HostID: hostID},
			Name:                 report.Name,
			Status:               report.Status,
			Message:              msg,
			OccurredAt:           Time(report.OccurredAt),
			NotificationInterval: report.NotificationInterval,
			MaxCheckAttempts:     report.MaxCheckAttempts,
		}
	}
	resp, err := api.postJSON("/api/v0/monitoring/checks/report", payload)
	defer closeResp(resp)
	return err
}
