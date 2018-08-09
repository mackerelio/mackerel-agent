// +build linux

package linux

import (
	"bytes"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func TestCPUGenerate(t *testing.T) {
	g := &CPUGenerator{}
	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	cpu, typeOk := value.(mackerel.CPU)
	if !typeOk {
		t.Errorf("value should be mackerel.CPU. %+v", value)
	}

	if len(cpu) == 0 {
		t.Fatal("should have at least 1 cpu")
	}

	cpu1 := cpu[0]
	if _, ok := cpu1["vendor_id"]; !ok {
		t.Error("cpu should have vendor_id")
	}
	if _, ok := cpu1["family"]; !ok {
		t.Error("cpu should have family")
	}
	if _, ok := cpu1["model"]; !ok {
		t.Error("cpu should have model")
	}
	if _, ok := cpu1["stepping"]; !ok {
		t.Error("cpu should have stepping")
	}
	if _, ok := cpu1["physical_id"]; !ok {
		// fails on some environments
		// t.Error("cpu should have physical_id")
	}
	if _, ok := cpu1["core_id"]; !ok {
		// fails on some environments
		// t.Error("cpu should have core_id")
	}
	if _, ok := cpu1["cores"]; !ok {
		// fails on some environments
		// t.Error("cpu should have cores")
	}
	if _, ok := cpu1["model_name"]; !ok {
		t.Error("cpu should have model_name")
	}
	if _, ok := cpu1["mhz"]; !ok {
		t.Error("cpu should have mhz")
	}
	if _, ok := cpu1["cache_size"]; !ok {
		t.Error("cpu should have cache_size")
	}
}

func TestCPUgenerate_linux4_0_amd64(t *testing.T) {
	cpuinfo := `processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 2
model name	: QEMU Virtual CPU version 1.1.2
stepping	: 3
microcode	: 0x1
cpu MHz		: 3392.292
cache size	: 4096 KB
physical id	: 0
siblings	: 1
core id		: 0
cpu cores	: 1
apicid		: 0
initial apicid	: 0
fpu		: yes
fpu_exception	: yes
cpuid level	: 4
wp		: yes
flags		: fpu de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pse36 clflush mmx fxsr sse sse2 syscall nx lm rep_good nopl pni cx16 popcnt hypervisor lahf_lm
bugs		:
bogomips	: 6784.58
clflush size	: 64
cache_alignment	: 64
address sizes	: 40 bits physical, 48 bits virtual
power management:

processor	: 1
vendor_id	: GenuineIntel
cpu family	: 6
model		: 2
model name	: QEMU Virtual CPU version 1.1.2
stepping	: 3
microcode	: 0x1
cpu MHz		: 3392.292
cache size	: 4096 KB
physical id	: 1
siblings	: 1
core id		: 0
cpu cores	: 1
apicid		: 1
initial apicid	: 1
fpu		: yes
fpu_exception	: yes
cpuid level	: 4
wp		: yes
flags		: fpu de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pse36 clflush mmx fxsr sse sse2 syscall nx lm rep_good nopl pni cx16 popcnt hypervisor lahf_lm
bugs		:
bogomips	: 6784.58
clflush size	: 64
cache_alignment	: 64
address sizes	: 40 bits physical, 48 bits virtual
power management:`

	g := &CPUGenerator{}
	value, err := g.generate(bytes.NewBufferString(cpuinfo))
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	cpus, typeOk := value.(mackerel.CPU)
	if !typeOk {
		t.Errorf("value should be mackerel.CPU. %+v", value)
	}

	if len(cpus) != 2 {
		t.Fatal("should have exactly 2 cpus")
	}

	for _, cpu := range cpus {
		modelName, ok := cpu["model_name"]
		if !ok {
			t.Error("cpu should have model_name")
		}
		if modelName != "QEMU Virtual CPU version 1.1.2" {
			t.Error("cpu should have correct model_name")
		}

		mhz, ok := cpu["mhz"]
		if !ok {
			t.Error("cpu should have mhz")
		}
		if mhz != "3392.292" {
			t.Error("cpu should have correct mhz")
		}
	}
}

func TestCPUgenerate_linux3_18_arm(t *testing.T) {
	cpuinfo := `processor       : 0
model name      : ARMv7 Processor rev 5 (v7l)
BogoMIPS        : 38.40
Features        : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm
CPU implementer : 0x41
CPU architecture: 7
CPU variant     : 0x0
CPU part        : 0xc07
CPU revision    : 5

processor       : 1
model name      : ARMv7 Processor rev 5 (v7l)
BogoMIPS        : 38.40
Features        : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm
CPU implementer : 0x41
CPU architecture: 7
CPU variant     : 0x0
CPU part        : 0xc07
CPU revision    : 5

processor       : 2
model name      : ARMv7 Processor rev 5 (v7l)
BogoMIPS        : 38.40
Features        : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm
CPU implementer : 0x41
CPU architecture: 7
CPU variant     : 0x0
CPU part        : 0xc07
CPU revision    : 5

processor       : 3
model name      : ARMv7 Processor rev 5 (v7l)
BogoMIPS        : 38.40
Features        : half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm
CPU implementer : 0x41
CPU architecture: 7
CPU variant     : 0x0
CPU part        : 0xc07
CPU revision    : 5

Hardware        : BCM2709
Revision        : 1a01040`

	g := &CPUGenerator{}
	value, err := g.generate(bytes.NewBufferString(cpuinfo))
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	cpus, typeOk := value.(mackerel.CPU)
	if !typeOk {
		t.Errorf("value should be mackerel.CPU. %+v", value)
	}

	if len(cpus) != 4 {
		t.Fatal("should have exactly 4 cpus")
	}

	for _, cpu := range cpus {
		modelName, ok := cpu["model_name"]
		if !ok {
			t.Error("cpu should have model_name")
		}
		if modelName != "ARMv7 Processor rev 5 (v7l)" {
			t.Error("cpu should have correct model_name")
		}
	}
}

func TestCPUgenerate_linux3_4_smp_arm(t *testing.T) {
	cpuinfo := `Processor       : ARMv7 Processor rev 0 (v7l)
processor       : 0
BogoMIPS        : 38.40

processor       : 1
BogoMIPS        : 38.40

processor       : 2
BogoMIPS        : 38.40

processor       : 3
BogoMIPS        : 38.40

Features        : swp half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt
CPU implementer : 0x51
CPU architecture: 7
CPU variant     : 0x2
CPU part        : 0x06f
CPU revision    : 0

Hardware        : Qualcomm MSM 8974 (Flattened Device Tree)
Revision        : 0000
Serial          : 0000000000000000`

	g := &CPUGenerator{}
	value, err := g.generate(bytes.NewBufferString(cpuinfo))
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	cpus, typeOk := value.(mackerel.CPU)
	if !typeOk {
		t.Errorf("value should be mackerel.CPU. %+v", value)
	}

	if len(cpus) != 4 {
		t.Fatal("should have exactly 4 cpus")
	}

	for _, cpu := range cpus {
		modelName, ok := cpu["model_name"]
		if !ok {
			t.Error("cpu should have model_name")
		}
		if modelName != "ARMv7 Processor rev 0 (v7l)" {
			t.Error("cpu should have correct model_name")
		}
	}
}

func TestCPUgenerate_linux3_0_nosmp_arm(t *testing.T) {
	cpuinfo := `Processor       : Marvell PJ4Bv7 Processor rev 1 (v7l)
BogoMIPS        : 1196.85
Features        : swp half thumb fastmult vfp edsp vfpv3 vfpv3d16
CPU implementer : 0x56
CPU architecture: 7
CPU variant     : 0x1
CPU part        : 0x581
CPU revision    : 1

Hardware        : Marvell Armada-370
Revision        : 0000
Serial          : 0000000000000000`

	g := &CPUGenerator{}
	value, err := g.generate(bytes.NewBufferString(cpuinfo))
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	cpus, typeOk := value.(mackerel.CPU)
	if !typeOk {
		t.Errorf("value should be mackerel.CPU. %+v", value)
	}

	if len(cpus) != 1 {
		t.Fatal("should have exactly 1 cpu")
	}

	cpu1 := cpus[0]
	modelName, ok := cpu1["model_name"]
	if !ok {
		t.Error("cpu should have model_name")
	}
	if modelName != "Marvell PJ4Bv7 Processor rev 1 (v7l)" {
		t.Error("cpu should have correct model_name")
	}
}
