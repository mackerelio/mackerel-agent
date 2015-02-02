// +build windows

package windows

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
	"fmt"
	"strconv"
	"strings"
	"os/exec"
)

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

// PDH_FMT_COUNTERVALUE_DOUBLE XXX
type PDH_FMT_COUNTERVALUE_DOUBLE struct {
	CStatus     uint32
	DoubleValue float64
}

// PDH_FMT_COUNTERVALUE_ITEM_DOUBLE XXX
type PDH_FMT_COUNTERVALUE_ITEM_DOUBLE struct {
	Name     *uint16
	FmtValue PDH_FMT_COUNTERVALUE_DOUBLE
}

// windows system const
const (
	ERROR_SUCCESS      = 0
	DRIVE_REMOVABLE    = 2
	DRIVE_FIXED        = 3
	HKEY_LOCAL_MACHINE = 0x80000002
	RRF_RT_REG_SZ      = 0x00000002
	RRF_RT_REG_DWORD   = 0x00000010
	PDH_FMT_DOUBLE     = 0x00000200
	PDH_INVALID_DATA   = 0xc0000bc6
)

// windows procs
var (
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	modpdh      = syscall.NewLazyDLL("pdh.dll")

	RegGetValue                 = modadvapi32.NewProc("RegGetValueW")
	GetSystemInfo               = modkernel32.NewProc("GetSystemInfo")
	GetTickCount                = modkernel32.NewProc("GetTickCount")
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
func RegGetInt(hKey uint32, subKey string, value string) (uint32, error) {
	var num, numlen uint32
	numlen = 4
	ret, _, err := RegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
		uintptr(RRF_RT_REG_DWORD),
		0,
		uintptr(unsafe.Pointer(&num)),
		uintptr(unsafe.Pointer(&numlen)))
	if ret != ERROR_SUCCESS {
		return 0, err
	}

	return num, nil
}

// RegGetString XXX
func RegGetString(hKey uint32, subKey string, value string) (string, error) {
	var bufLen uint32
	RegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
		uintptr(RRF_RT_REG_SZ),
		0,
		0,
		uintptr(unsafe.Pointer(&bufLen)))
	if bufLen == 0 {
		return "", errors.New("Can't get size of registry value")
	}

	buf := make([]uint16, bufLen)
	ret, _, err := RegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
		uintptr(RRF_RT_REG_SZ),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)))
	if ret != ERROR_SUCCESS {
		return "", err
	}

	return syscall.UTF16ToString(buf), nil
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
	r, _, err := PdhAddCounter.Call(
		uintptr(query),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(v))),
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

// GetAdapterList XXX
func GetAdapterList() (*syscall.IpAdapterInfo, error) {
	b := make([]byte, 1000)
	l := uint32(len(b))
	a := (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
	err := syscall.GetAdaptersInfo(a, &l)
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b = make([]byte, l)
		a = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
		err = syscall.GetAdaptersInfo(a, &l)
	}
	if err != nil {
		return nil, os.NewSyscallError("GetAdaptersInfo", err)
	}
	return a, nil
}

// BytePtrToString XXX
func BytePtrToString(p *uint8) string {
	a := (*[10000]uint8)(unsafe.Pointer(p))
	i := 0
	for a[i] != 0 {
		i++
	}
	return string(a[:i])
}

// FilesystemInfo XXX
type FilesystemInfo struct {
	Percent_used string
	Kb_used      float64
	Kb_size      float64
	Kb_available float64
	Mount        string
	Label        string
	Volume_name  string
	Fs_type      string
}

// GetWmic XXX
func GetWmic(target string, query string) (string, error) {
	cpuGet, err := exec.Command("wmic", target, "get", query).Output()
	if err != nil {
		return "", err
	}

	percentages := string(cpuGet)

	lines := strings.Split(percentages, "\r\r\n")

	if len(lines) <= 2 {
		return "", fmt.Errorf("wmic result malformed: [%q]", lines)
	}

	return strings.Trim(lines[1], " "), nil
}

// GetWmicToFloat XXX
func GetWmicToFloat(target string, query string) (float64, error) {
	value, err := GetWmic(target, query)
	if err != nil {
		return 0, err
	}

	ret, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return ret, nil
}
