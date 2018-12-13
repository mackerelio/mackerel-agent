// +build !windows

package cmdutil

import (
	"bytes"
)

func decodeBytes(b *bytes.Buffer) string {
	return b.String()
}
