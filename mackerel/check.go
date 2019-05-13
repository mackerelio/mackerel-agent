package mackerel

import (
	"github.com/mackerelio/mackerel-agent/checks"
	mkr "github.com/mackerelio/mackerel-client-go"
)

type monitoringChecksPayload struct {
	Reports []*checkReport `json:"reports"`
}

type checkReport struct {
	Source               mkr.CheckSource `json:"source"`
	Name                 string          `json:"name"`
	Status               mkr.CheckStatus `json:"status"`
	Message              string          `json:"message"`
	OccurredAt           int64           `json:"occurredAt"`
	NotificationInterval uint            `json:"notificationInterval,omitempty"`
	MaxCheckAttempts     uint            `json:"maxCheckAttempts,omitempty"`
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
			Source:               mkr.NewCheckSourceHost(hostID),
			Name:                 report.Name,
			Status:               mkr.CheckStatus(report.Status),
			Message:              msg,
			OccurredAt:           report.OccurredAt.Unix(),
			NotificationInterval: int32ptrToUint(report.NotificationInterval),
			MaxCheckAttempts:     int32ptrToUint(report.MaxCheckAttempts),
		}
	}
	resp, err := api.postJSON("/api/v0/monitoring/checks/report", payload)
	defer closeResp(resp)
	return err
}

func int32ptrToUint(p *int32) uint {
	if p == nil || *p < 0 {
		return 0
	}
	return uint(*p)
}
