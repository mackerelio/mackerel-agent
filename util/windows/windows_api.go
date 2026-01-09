//go:build windows

package windows

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ref. https://github.com/mackerelio/mackerel-agent/pull/134

/*
//#include <pdh.h>
typedef unsigned long DWORD;
// Union specialization for double values
typedef struct _PDH_FMT_COUNTERVALUE_DOUBLE {
	DWORD  CStatus;
	double DoubleValue;
} PDH_FMT_COUNTERVALUE_DOUBLE;
*/
import "C"

// SYSTEM_INFO XXX
type SYSTEM_INFO struct {
	ProcessorArchitecture     uint16
	PageSize                  uint32
	MinimumApplicationAddress *byte
	MaximumApplicationAddress *byte
	ActiveProcessorMask       *byte
	NumberOfProcessors        uint32
	ProcessorType             uint32
	AllocationGranularity     uint32
	ProcessorLevel            uint16
	ProcessorRevision         uint16
}

// MEMORY_STATUS_EX XXX
type MEMORY_STATUS_EX struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

// windows system const
const (
	ERROR_SUCCESS        = 0
	ERROR_FILE_NOT_FOUND = 2
	DRIVE_REMOVABLE      = 2
	DRIVE_FIXED          = 3
	HKEY_LOCAL_MACHINE   = 0x80000002
	RRF_RT_REG_SZ        = 0x00000002
	RRF_RT_REG_DWORD     = 0x00000010
	PDH_FMT_DOUBLE       = 0x00000200
	PDH_INVALID_DATA     = 0xc0000bc6
	PDH_INVALID_HANDLE   = 0xC0000bbc
	PDH_NO_DATA          = 0x800007d5
)

// windows procs
var (
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	modpdh      = syscall.NewLazyDLL("pdh.dll")

	RegGetValue                 = modadvapi32.NewProc("RegGetValueW")
	GetSystemInfo               = modkernel32.NewProc("GetSystemInfo")
	GetDiskFreeSpaceEx          = modkernel32.NewProc("GetDiskFreeSpaceExW")
	GetLogicalDriveStrings      = modkernel32.NewProc("GetLogicalDriveStringsW")
	GetDriveType                = modkernel32.NewProc("GetDriveTypeW")
	QueryDosDevice              = modkernel32.NewProc("QueryDosDeviceW")
	GetVolumeInformationW       = modkernel32.NewProc("GetVolumeInformationW")
	GlobalMemoryStatusEx        = modkernel32.NewProc("GlobalMemoryStatusEx")
	GetLastError                = modkernel32.NewProc("GetLastError")
	PdhOpenQuery                = modpdh.NewProc("PdhOpenQuery")
	PdhAddCounter               = modpdh.NewProc("PdhAddCounterW")
	PdhCollectQueryData         = modpdh.NewProc("PdhCollectQueryData")
	PdhGetFormattedCounterValue = modpdh.NewProc("PdhGetFormattedCounterValue")
	PdhCloseQuery               = modpdh.NewProc("PdhCloseQuery")
)

// RegGetInt XXX
func RegGetInt(hKey uint32, subKey string, value string) (uint32, uintptr, error) {
	var num, numlen uint32
	numlen = 4
	s, err := syscall.UTF16PtrFromString(subKey)
	if err != nil {
		return 0, 0, err
	}
	v, err := syscall.UTF16PtrFromString(value)
	if err != nil {
		return 0, 0, err
	}
	ret, _, err := RegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(s)),
		uintptr(unsafe.Pointer(v)),
		uintptr(RRF_RT_REG_DWORD),
		0,
		uintptr(unsafe.Pointer(&num)),
		uintptr(unsafe.Pointer(&numlen)))
	if ret != ERROR_SUCCESS {
		return 0, ret, err
	}

	return num, ret, nil
}

// RegGetString XXX
func RegGetString(hKey uint32, subKey string, value string) (string, uintptr, error) {
	var bufLen uint32
	s, err := syscall.UTF16PtrFromString(subKey)
	if err != nil {
		return "", 0, err
	}
	v, err := syscall.UTF16PtrFromString(value)
	if err != nil {
		return "", 0, err
	}
	ret, _, err := RegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(s)),
		uintptr(unsafe.Pointer(v)),
		uintptr(RRF_RT_REG_SZ),
		0,
		0,
		uintptr(unsafe.Pointer(&bufLen)))
	if ret != ERROR_SUCCESS {
		return "", ret, err
	}
	if bufLen == 0 {
		return "", ret, nil
	}

	buf := make([]uint16, bufLen)
	s, err = syscall.UTF16PtrFromString(subKey)
	if err != nil {
		return "", 0, err
	}
	v, err = syscall.UTF16PtrFromString(value)
	if err != nil {
		return "", 0, err
	}
	ret, _, err = RegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(s)),
		uintptr(unsafe.Pointer(v)),
		uintptr(RRF_RT_REG_SZ),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)))
	if ret != ERROR_SUCCESS {
		return "", ret, err
	}

	return syscall.UTF16ToString(buf), ret, nil
}

// CounterInfo XXX
type CounterInfo struct {
	PostName    string
	CounterName string
	Counter     syscall.Handle
}

// CreateQuery XXX
func CreateQuery() (syscall.Handle, error) {
	var query syscall.Handle
	r, _, err := PdhOpenQuery.Call(0, 0, uintptr(unsafe.Pointer(&query)))
	if r != 0 {
		return 0, err
	}
	return query, nil
}

// CreateCounter XXX
func CreateCounter(query syscall.Handle, k, v string) (*CounterInfo, error) {
	var counter syscall.Handle

	v2, err := syscall.UTF16PtrFromString(v)
	if err != nil {
		return nil, err
	}
	r, _, err := PdhAddCounter.Call(
		uintptr(query),
		uintptr(unsafe.Pointer(v2)),
		0,
		uintptr(unsafe.Pointer(&counter)))
	if r != 0 {
		return nil, err
	}
	return &CounterInfo{
		PostName:    k,
		CounterName: v,
		Counter:     counter,
	}, nil
}

// GetCounterValue get counter value from handle
func GetCounterValue(counter syscall.Handle) (float64, error) {
	var value C.PDH_FMT_COUNTERVALUE_DOUBLE
	r, _, err := PdhGetFormattedCounterValue.Call(uintptr(counter), PDH_FMT_DOUBLE, uintptr(0), uintptr(unsafe.Pointer(&value)))
	if r != 0 && r != PDH_INVALID_DATA {
		return 0.0, err
	}
	return float64(value.DoubleValue), nil
}

type Adapter struct {
	Index int
	Name  string
}

func GetAdapterList() ([]Adapter, error) {
	// Reference from src/net/interface_windows.go at Go 1.25
	var b []byte
	l := uint32(15000)
	for {
		b = make([]byte, l)
		const flags = windows.GAA_FLAG_INCLUDE_PREFIX | windows.GAA_FLAG_INCLUDE_GATEWAYS
		err := windows.GetAdaptersAddresses(syscall.AF_UNSPEC, flags, 0, (*windows.IpAdapterAddresses)(unsafe.Pointer(&b[0])), &l)

		if err == nil {
			if l == 0 {
				return nil, nil
			}
			break
		}
		if err.(syscall.Errno) != syscall.ERROR_BUFFER_OVERFLOW {
			return nil, os.NewSyscallError("getadaptersaddresses", err)
		}
		if l <= uint32(len(b)) {
			return nil, os.NewSyscallError("getadaptersaddresses", err)
		}
	}
	var ads []Adapter
	for aa := (*windows.IpAdapterAddresses)(unsafe.Pointer(&b[0])); aa != nil; aa = aa.Next {
		index := aa.IfIndex
		if index == 0 { // ipv6IfIndex is a substitute for ifIndex
			index = aa.Ipv6IfIndex
		}
		ads = append(ads, Adapter{
			Index: int(index),
			Name:  windows.UTF16PtrToString(aa.Description),
		})
	}
	return ads, nil
}
