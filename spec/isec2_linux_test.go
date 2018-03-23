package spec

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func setEc2BaseURL(url *url.URL) func() {
	oldEC2BaseURL := ec2BaseURL
	ec2BaseURL = url
	return func() {
		ec2BaseURL = oldEC2BaseURL // restore value
	}
}

func setUUIDFiles(files [2]string) func() {
	oldUUIDFiles := uuidFiles
	uuidFiles = files
	return func() {
		uuidFiles = oldUUIDFiles // restore value
	}
}

func TestIsEC2UUID(t *testing.T) {
	tests := []struct {
		uuid   string
		expect bool
	}{
		{"ec2e1916-9099-7caf-fd21-01234abcdef", true},
		{"EC2E1916-9099-7CAF-FD21-01234ABCDEF", true},
		{"45e12aec-dcd1-b213-94ed-01234abcdef", true}, // litte endian
		{"45E12AEC-DCD1-B213-94ED-01234ABCDEF", true}, // litte endian
		{"abcd1916-9099-7caf-fd21-01234abcdef", false},
		{"ABCD1916-9099-7CAF-FD21-01234ABCDEF", false},
		{"", false},
	}

	for _, tc := range tests {
		if isEC2UUID(tc.uuid) != tc.expect {
			t.Errorf("isEC2() should be %t: %q\n", tc.expect, tc.uuid)
		}
	}
}

func TestIsEC2(t *testing.T) {
	tests := []struct {
		existsUUIDFiles [2]bool // [0]: as `/sys/hypervisor/uuid`, [1]: as `/sys/devices/virtual/dmi/id/product_uuid`
		existsAMIId     bool
		expect          bool
	}{
		{
			existsUUIDFiles: [2]bool{
				true,
				true,
			},
			existsAMIId: true,
			expect:      true,
		},
		{
			existsUUIDFiles: [2]bool{
				false,
				true,
			},
			existsAMIId: true,
			expect:      true,
		},
		{
			existsUUIDFiles: [2]bool{
				true,
				false,
			},
			existsAMIId: true,
			expect:      true,
		},
		{
			existsUUIDFiles: [2]bool{
				false,
				false,
			},
			existsAMIId: true,
			expect:      false,
		},
		{
			existsUUIDFiles: [2]bool{
				true,
				true,
			},
			existsAMIId: false,
			expect:      false,
		},
	}

	for _, tc := range tests {
		func() {
			handler := func(res http.ResponseWriter, req *http.Request) {
				if !tc.existsAMIId {
					res.WriteHeader(http.StatusNotFound)
				}
			}
			ts := httptest.NewServer(http.HandlerFunc(handler))
			defer func() { ts.Close() }()

			u, _ := url.Parse(ts.URL)
			defer setEc2BaseURL(u)()

			var uuidFiles [2]string
			for i, exist := range tc.existsUUIDFiles {
				tf, err := ioutil.TempFile("", "")
				if err != nil {
					t.Errorf("should not raise error: %s", err)
				}

				tn := tf.Name()
				uuidFiles[i] = tn

				if exist {
					defer os.Remove(tn)
				} else {
					os.Remove(tn)
					continue
				}

				tf.Write([]byte("ec2e1916-9099-7caf-fd21-012345abcdef")) // valid EC2 UUID
				tf.Close()
			}
			defer setUUIDFiles(uuidFiles)()

			if isEC2() != tc.expect {
				t.Errorf("isEC2() should be %t: %#v\n", tc.expect, tc)
			}
		}()
	}
}
