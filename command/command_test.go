package command

import (
	"testing"

	"github.com/mackerelio/mackerel-agent/mackerel"
)

func TestDelayByHost(t *testing.T) {
	delay1 := delayByHost(&mackerel.Host{
		Id:     "246PUVUngPo",
		Name:   "hogehoge2.host.h",
		Type:   "unknown",
		Status: "working",
	})

	delay2 := delayByHost(&mackerel.Host{
		Id:     "21GZjCE5Etb",
		Name:   "hogehoge2.host.h",
		Type:   "unknown",
		Status: "working",
	})

	if !(0 <= delay1.Seconds() && delay1.Seconds() < 60) {
		t.Errorf("delay shoud be between 0 and 60 but %v", delay1)
	}

	if delay1 == delay2 {
		t.Error("delays shoud be different")
	}
}
