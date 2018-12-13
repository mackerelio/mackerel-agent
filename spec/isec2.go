// +build !linux,!windows

package spec

import (
	"context"
	"net/http"
	"time"

	"github.com/Songmu/retry"
)

// For instances other than Linux, retry only 1 times to shorten whole process
func isEC2(ctx context.Context) bool {
	isEC2 := false
	err := retry.WithContext(ctx, 2, 2*time.Second, func() error {
		cl := httpCli()
		// '/ami-id` is probably an AWS specific URL
		req, err := http.NewRequest("GET", ec2BaseURL.String()+"/ami-id", nil)
		if err != nil {
			return err
		}
		resp, err := cl.Do(req.WithContext(ctx))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		isEC2 = resp.StatusCode == 200
		return nil
	})
	return err == nil && isEC2
}
