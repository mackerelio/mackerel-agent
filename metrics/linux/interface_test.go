// +build linux

package linux

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/metrics"
)

func TestInterfaceGenerator(t *testing.T) {
	if _, err := os.Stat("/etc/fedora-release"); err == nil {
		t.Skip("The OS seems to be Fedora. Skipping interface test for now")
	}

	g := &InterfaceGenerator{1 * time.Second}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	metrics := []string{
		"rxBytes", "txBytes",
	}

	for _, metric := range metrics {
		if value, ok := values["interface.eth0."+metric+".delta"]; !ok {
			t.Errorf("Value for interface.eth0.%s.delta should be collected", metric)
		} else {
			t.Logf("Interface eth0 '%s' delta collected: %+v", metric, value)
		}
	}

	for _, metric := range metrics {
		if _, ok := values["interface.lo."+metric+".delta"]; ok {
			t.Errorf("Value for interface.lo.%s should NOT be collected", metric)
		} else {
			t.Logf("Interface lo '%s' NOT collected", metric)
		}
	}
}

func TestParseNetdev(t *testing.T) {
	out := []byte(`Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
eth0: 5461472598 24386569    0    2    0     0          0         0 7215710422 6079810    0    0    0     0       0          0
lo: 7779878638 1952628    0    0    0     0          0         0 7779878638 1952628    0    0    0     0       0          0
docker0: 250219988  333736    0    0    0     0          0         0 2024726607 1409929    0    0    0     0       0          0`)

	result, err := parseNetdev(out)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}

	expect := metrics.Values{
		"interface.eth0.rxBytes":         5461472598,
		"interface.eth0.rxPackets":       24386569,
		"interface.eth0.rxErrors":        0,
		"interface.eth0.rxDrops":         2,
		"interface.eth0.rxFifo":          0,
		"interface.eth0.rxFrame":         0,
		"interface.eth0.rxCompressed":    0,
		"interface.eth0.rxMulticast":     0,
		"interface.eth0.txBytes":         7215710422,
		"interface.eth0.txPackets":       6079810,
		"interface.eth0.txErrors":        0,
		"interface.eth0.txDrops":         0,
		"interface.eth0.txFifo":          0,
		"interface.eth0.txColls":         0,
		"interface.eth0.txCarrier":       0,
		"interface.eth0.txCompressed":    0,
		"interface.docker0.rxBytes":      250219988,
		"interface.docker0.rxPackets":    333736,
		"interface.docker0.rxErrors":     0,
		"interface.docker0.rxDrops":      0,
		"interface.docker0.rxFifo":       0,
		"interface.docker0.rxFrame":      0,
		"interface.docker0.rxCompressed": 0,
		"interface.docker0.rxMulticast":  0,
		"interface.docker0.txBytes":      2024726607,
		"interface.docker0.txPackets":    1409929,
		"interface.docker0.txErrors":     0,
		"interface.docker0.txDrops":      0,
		"interface.docker0.txFifo":       0,
		"interface.docker0.txColls":      0,
		"interface.docker0.txCarrier":    0,
		"interface.docker0.txCompressed": 0,
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result is not expected one: %+v", result)
	}
}
