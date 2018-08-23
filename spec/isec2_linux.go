package spec

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Songmu/retry"
)

// If the OS is Linux, check /sys/hypervisor/uuid and /sys/devices/virtual/dmi/id/product_uuid files first. If UUID seems to be EC2-ish, call the metadata API (up to 3 times).
// ref. https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/identify_ec2_instances.html
func isEC2(ctx context.Context) bool {
	var uuidFiles = []string{
		"/sys/hypervisor/uuid",
		"/sys/devices/virtual/dmi/id/product_uuid",
	}

	return isEC2WithSpecifiedUUIDFiles(ctx, uuidFiles)
}

func isEC2WithSpecifiedUUIDFiles(ctx context.Context, uuidFiles []string) bool {
	looksLikeEC2 := false
	for _, u := range uuidFiles {
		data, err := ioutil.ReadFile(u)
		if err != nil {
			continue
		}
		if isEC2UUID(string(data)) {
			looksLikeEC2 = true
			break
		}
	}
	if !looksLikeEC2 {
		return false
	}

	// give up if ctx already closed
	select {
	case <-ctx.Done():
		return false
	default:
	}

	res := false
	cl := httpCli()
	err := retry.WithContext(ctx, 3, 2*time.Second, func() error {
		// '/ami-id` is probably an AWS specific URL
		req, err := http.NewRequest("GET", ec2BaseURL.String()+"/ami-id", nil)
		if err != nil {
			return nil // something wrong. give up
		}
		resp, err := cl.Do(req.WithContext(ctx))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		res = resp.StatusCode == 200
		return nil
	})

	if err == nil {
		return res
	}

	return false
}

func isEC2UUID(uuid string) bool {
	conds := func(uuid string) bool {
		if strings.HasPrefix(uuid, "ec2") || strings.HasPrefix(uuid, "EC2") {
			return true
		}
		return false
	}

	if conds(uuid) {
		return true
	}

	// Check as littele endian.
	// see. https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/identify_ec2_instances.html
	fields := strings.Split(uuid, "-")
	decoded, _ := hex.DecodeString(fields[0]) // fields[0]: UUID time_low(uint32)
	r := bytes.NewReader(decoded)
	var data uint32
	binary.Read(r, binary.LittleEndian, &data)

	return conds(fmt.Sprintf("%x", data))
}
