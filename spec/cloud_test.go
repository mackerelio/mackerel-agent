package spec

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCloudKey(t *testing.T) {
	g := &CloudGenerator{}

	if g.Key() != "cloud" {
		t.Error("key should be cloud")
	}
}

func TestCloudGenerate(t *testing.T) {
	handler := func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "i-4f90d537")
	}
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler(res, req)
	}))
	defer ts.Close()

	g, err := NewCloudGenerator(ts.URL)
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}
	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}

	cloud, typeOk := value.(map[string]interface{})
	if !typeOk {
		t.Errorf("value should be map. %+v", value)
	}

	value, ok := cloud["metadata"]
	if !ok {
		t.Error("results should have metadata.")
	}

	metadata, typeOk := value.(map[string]string)
	if !typeOk {
		t.Errorf("v should be map. %+v", value)
	}

	if len(metadata["instance-id"]) == 0 {
		t.Error("instance-id should be filled")
	}
}
