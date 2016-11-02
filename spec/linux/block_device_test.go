// +build linux

package linux

import (
	"regexp"
	"testing"
)

func TestBlockDeviceGenerator(t *testing.T) {
	g := &BlockDeviceGenerator{}

	if g.Key() != "block_device" {
		t.Error("key should be block_device")
	}
}

func hasValidBlockDeviceValueForKey(t *testing.T, deviceInfo map[string]interface{}, key string) {
	if value, ok := deviceInfo[key]; !ok {
		t.Errorf("value of %s should be retrieved but none", key)
	} else if regexp.MustCompile(`\n$`).MatchString(value.(string)) {
		t.Errorf("value of %s should not be end with newline", key)
	}
}

func TestBlockDeviceGenerate(t *testing.T) {
	g := &BlockDeviceGenerator{}

	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	blockDevice, typeOk := value.(map[string]map[string]interface{})
	if !typeOk {
		t.Errorf("value should be slice of map. %+v", value)
	}

	sda, ok := blockDevice["sda"]
	if !ok {
		t.Skip("should have map for sda")
	}

	hasValidBlockDeviceValueForKey(t, sda, "size")
	hasValidBlockDeviceValueForKey(t, sda, "removable")
	hasValidBlockDeviceValueForKey(t, sda, "model")
	hasValidBlockDeviceValueForKey(t, sda, "rev")
	hasValidBlockDeviceValueForKey(t, sda, "state")
	hasValidBlockDeviceValueForKey(t, sda, "timeout")
	hasValidBlockDeviceValueForKey(t, sda, "vendor")
}
