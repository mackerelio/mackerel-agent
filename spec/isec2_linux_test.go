package spec

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestIsEC2UUID(t *testing.T) {
	tests := []struct {
		uuid   string
		expect bool
	}{
		{"ec2e1916-9099-7caf-fd21-01234abcdef", true},
		{"EC2E1916-9099-7CAF-FD21-01234ABCDEF", true},
		{"45e12aec-dcd1-b213-94ed-01234abcdef", true}, // little endian
		{"45E12AEC-DCD1-B213-94ED-01234ABCDEF", true}, // little endian
		{"abcd1916-9099-7caf-fd21-01234abcdef", false},
		{"ABCD1916-9099-7CAF-FD21-01234ABCDEF", false},
		{"", false},
	}

	for _, tc := range tests {
		if isEC2UUID(tc.uuid) != tc.expect {
			t.Errorf("isEC2() should be %v: %q", tc.expect, tc.uuid)
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
				res.WriteHeader(http.StatusOK)
				res.Write([]byte("aws"))
			}
			ts := httptest.NewServer(http.HandlerFunc(handler))
			defer ts.Close()
			u, _ := url.Parse(ts.URL)
			g := &EC2Generator{
				baseURL: u,
			}

			uuidFiles := make([]string, 0, 2)
			for _, exist := range tc.existsUUIDFiles {
				tf, err := ioutil.TempFile("", "")
				if err != nil {
					t.Errorf("should not raise error: %s", err)
				}

				tn := tf.Name()
				uuidFiles = append(uuidFiles, tn)

				if exist {
					defer os.Remove(tn)
				} else {
					os.Remove(tn)
					continue
				}

				tf.Write([]byte("ec2e1916-9099-7caf-fd21-012345abcdef")) // valid EC2 UUID
				tf.Close()
			}

			if g.isEC2WithSpecifiedUUIDFiles(context.Background(), uuidFiles) != tc.expect {
				t.Errorf("isEC2() should be %v: %#v", tc.expect, tc)
			}
		}()
	}
}
