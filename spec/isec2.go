// +build !linux,!windows

package spec

import (
	"context"
	"time"

	"github.com/Songmu/retry"
)

// For instances other than Linux, retry only 1 times to shorten whole process
func (g *EC2Generator) isEC2(ctx context.Context) bool {
	var res bool
	err := retry.WithContext(ctx, 2, 2*time.Second, func() error {
		res0, err := g.hasMetadataService(ctx)
		if err != nil {
			return err
		}
		res = res0
		return nil
	})
	return err == nil && res
}
