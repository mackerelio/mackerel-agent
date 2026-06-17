package metrics

import (
	"reflect"
	"testing"
)

func TestMergeValuesCustomIdentifiers(t *testing.T) {
	var v0 = Values{
		"aa": NewValueAttribute(10),
	}
	var v1 = Values{
		"bb": NewValueAttribute(20),
		"cc": NewValueAttribute(30),
	}
	var v2 = Values{
		"dd": NewValueAttribute(40),
		"ee": NewValueAttribute(50),
	}
	var v3 = Values{
		"ff": NewValueAttribute(60),
		"gg": NewValueAttribute(70),
	}

	v := MergeValuesCustomIdentifiers([]*ValuesCustomIdentifier{
		{Values: v0},
	}, &ValuesCustomIdentifier{Values: v1})

	if !reflect.DeepEqual(v, []*ValuesCustomIdentifier{
		{
			Values: Values{
				"aa": NewValueAttribute(10),
				"bb": NewValueAttribute(20),
				"cc": NewValueAttribute(30),
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
				"aa": NewValueAttribute(10),
				"bb": NewValueAttribute(20),
				"cc": NewValueAttribute(30),
			},
			CustomIdentifier: nil,
		},
		{
			Values: Values{
				"dd": NewValueAttribute(40),
				"ee": NewValueAttribute(50),
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
				"aa": NewValueAttribute(10),
				"bb": NewValueAttribute(20),
				"cc": NewValueAttribute(30),
			},
			CustomIdentifier: nil,
		},
		{
			Values: Values{
				"dd": NewValueAttribute(40),
				"ee": NewValueAttribute(50),
				"ff": NewValueAttribute(60),
				"gg": NewValueAttribute(70),
			},
			CustomIdentifier: &customIdentifiers,
		},
	}) {
		t.Errorf("somthing went wrong")
	}
}
