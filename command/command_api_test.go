package command

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIRequestHeader(t *testing.T) {
	ver := "1.0.0"
	rev := "1234beaf"
	apiKey := "dummy-apikey"
	handler := func(res http.ResponseWriter, req *http.Request) {
		ua := buildUA(ver, rev)
		if req.Header.Get("X-Api-Key") != apiKey {
			t.Error("X-Api-Key header should contains passed key")
		}

		if h := req.Header.Get("X-Agent-Version"); h != ver {
			t.Errorf("X-Agent-Version shoud be %s but %s", ver, h)
		}

		if h := req.Header.Get("X-Revision"); h != rev {
			t.Errorf("X-Revision shoud be %s but %s", rev, h)
		}

		if h := req.Header.Get("User-Agent"); h != ua {
			t.Errorf("User-Agent shoud be '%s' but %s", ua, h)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler(res, req)
	}))
	defer ts.Close()

	api, err := NewMackerelClient(ts.URL, apiKey, ver, rev, false)
	if err != nil {
		t.Errorf("something went wrong while creating new mackerel client: %+v", err)
	}
	api.FindHost("dummy-id")
}
