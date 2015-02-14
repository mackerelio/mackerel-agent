// +build linux

package linux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewInstanceGenerator(t *testing.T) {
	g, err := NewInstanceGenerator(
		"http://example.com",
	)

	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}


	if g.baseURL.String() != "http://example.com" {
		t.Error("should return URL")
	}
}

func TestInstanceKey(t *testing.T) {
	g := &InstanceGenerator{}

	if g.Key() != "instance" {
		t.Error("key should be instance")
	}
}

func TestInstanceGenerate(t *testing.T) {
	handler := func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "i-4f90d537")
	}
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler(res, req)
	}))
	defer ts.Close()

	g, err := NewInstanceGenerator(ts.URL)
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}
	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}

	instance, typeOk := value.(map[string]interface{})
	if !typeOk {
		t.Errorf("value should be map. %+v", value)
	}

	value, ok := instance["metadata"]
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
