package metrics

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	var v = Values(map[string]float64{
		"aa": 10,
	})
	var vv = Values(map[string]float64{
		"bb": 20,
		"cc": 30,
	})

	v.Merge(vv)

	if !reflect.DeepEqual(v, Values(map[string]float64{
		"aa": 10,
		"bb": 20,
		"cc": 30,
	})) {
		t.Errorf("somthing went wrong")
	}
}
