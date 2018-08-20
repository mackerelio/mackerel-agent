// +build !linux

package spec

import (
	"context"
	"net/http"
)

// For instances other than Linux, call Metadata API only once.
func isEC2(ctx context.Context) bool {
	cl := httpCli()
	// '/ami-id` is probably an AWS specific URL
	req, err := http.NewRequest("GET", ec2BaseURL.String()+"/ami-id", nil)
	if err != nil {
		return false
	}
	resp, err := cl.Do(req.WithContext(ctx))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}
