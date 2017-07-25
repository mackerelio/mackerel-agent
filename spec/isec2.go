// +build !linux

package spec

// For instances other than Linux, call Metadata API only once.
func isEC2() bool {
	cl := httpCli()
	// '/ami-id` is probably an AWS specific URL
	resp, err := cl.Get(ec2BaseURL.String() + "/ami-id")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}
