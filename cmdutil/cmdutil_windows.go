// +build windows

package cmdutil

import (
	"bytes"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func decodeBytes(b *bytes.Buffer) string {
	if b.Len()%2 != 0 {
		return b.String()
	}
	found := false
	for _, v := range b.Bytes() {
		if v == 0x00 {
			found = true
			break
		}
	}
	if !found {
		return b.String()
	}
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	bb, _, err := transform.Bytes(enc.NewDecoder(), b.Bytes())
	if err != nil {
		return b.String()
	}
	return string(bb)
}
