package spec

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/Songmu/retry"
)

// If the OS is Linux, check /sys/hypervisor/uuid file first. If UUID seems to be EC2-ish, call the metadata API (up to 3 times).
// ref. http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/identify_ec2_instances.html
func isEC2() bool {
	data, err := ioutil.ReadFile("/sys/hypervisor/uuid")
	if err != nil {
		// Not EC2.
		return false
	}
	// Not EC2.
	if !strings.HasPrefix(string(data), "ec2") {
		return false
	}
	res := false
	cl := httpCli()
	err = retry.Retry(3, 2*time.Second, func() error {
		// '/ami-id` is probably an AWS specific URL
		resp, err := cl.Get(ec2BaseURL.String() + "/ami-id")
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
