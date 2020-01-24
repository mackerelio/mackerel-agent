package mackerel

import (
	"github.com/mackerelio/mackerel-agent/checks"
	mkr "github.com/mackerelio/mackerel-client-go"
)

const (
	// NotificationIntervalFallback is the minimum minutes of Mackerel API's
	// notification_interval.
	NotificationIntervalFallback = 10
)

// ReportCheckMonitors sends reports of *checks.Checker() to Mackrel API server.
func (api *API) ReportCheckMonitors(hostID string, reports []*checks.Report) error {
	payload := &mkr.CheckReports{
		Reports: make([]*mkr.CheckReport, len(reports)),
	}
	const messageLengthLimit = 1024
	for i, report := range reports {
		msg := report.Message
		runes := []rune(msg)
		if len(runes) > messageLengthLimit {
			msg = string(runes[0:messageLengthLimit])
		}
		payload.Reports[i] = &mkr.CheckReport{
			Source:               mkr.NewCheckSourceHost(hostID),
			Name:                 report.Name,
			Status:               mkr.CheckStatus(report.Status),
			Message:              msg,
			OccurredAt:           report.OccurredAt.Unix(),
			NotificationInterval: normalize(report.NotificationInterval, NotificationIntervalFallback),
			MaxCheckAttempts:     normalize(report.MaxCheckAttempts, 0),
		}
	}
	return api.Client.PostCheckReports(payload)
}

// normalize returns rounded valid number for Mackerel's Check API.
//
// In Mackerel, notification_interval regards nil as none,
// but regards zero as default (10 minutes).
func normalize(p *int32, min int32) uint {
	// TODO(lufia): we will be able to remove this when we can ignore v0.59 or earlier.
	if min < 0 {
		panic("min must be >=0")
	}
	switch {
	case p == nil:
		return 0
	case *p < min:
		return uint(min)
	default:
		return uint(*p)
	}
}
