package metrics

import (
	"reflect"
	"testing"
)

func TestMergeValuesCustomIdentifiers(t *testing.T) {
	var v0 = Values{
		"aa": 10,
	}
	var v1 = Values{
		"bb": 20,
		"cc": 30,
	}
	var v2 = Values{
		"dd": 40,
		"ee": 50,
	}
	var v3 = Values{
		"ff": 60,
		"gg": 70,
	}

	v := MergeValuesCustomIdentifiers([]*ValuesCustomIdentifier{
		{Values: v0},
	}, &ValuesCustomIdentifier{Values: v1})

	if !reflect.DeepEqual(v, []*ValuesCustomIdentifier{
		{
			Values: Values{
				"aa": 10,
				"bb": 20,
				"cc": 30,
			},
			CustomIdentifier: nil,
		}}) {
		t.Errorf("somthing went wrong")
	}

	customIdentifiers := "foo-bar"
	v = MergeValuesCustomIdentifiers(v, &ValuesCustomIdentifier{Values: v2, CustomIdentifier: &customIdentifiers})

	if !reflect.DeepEqual(v, []*ValuesCustomIdentifier{
		{
			Values: Values{
				"aa": 10,
				"bb": 20,
				"cc": 30,
			},
			CustomIdentifier: nil,
		},
		{
			Values: Values{
				"dd": 40,
				"ee": 50,
			},
			CustomIdentifier: &customIdentifiers,
		},
	}) {
		t.Errorf("somthing went wrong")
	}

	sameCustomIdentifiers := "foo-bar"
	v = MergeValuesCustomIdentifiers(v, &ValuesCustomIdentifier{Values: v3, CustomIdentifier: &sameCustomIdentifiers})

	if !reflect.DeepEqual(v, []*ValuesCustomIdentifier{
		{
			Values: Values{
				"aa": 10,
				"bb": 20,
				"cc": 30,
			},
			CustomIdentifier: nil,
		},
		{
			Values: Values{
				"dd": 40,
				"ee": 50,
				"ff": 60,
				"gg": 70,
			},
			CustomIdentifier: &customIdentifiers,
		},
	}) {
		t.Errorf("somthing went wrong")
	}
}
