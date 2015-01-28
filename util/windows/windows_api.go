// +build windows

package windows

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

// SystemInfo XXX
type SystemInfo struct {
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

// MemoryStatusEx XXX
type MemoryStatusEx struct {
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

// PdhFmtCountervalueDouble XXX
type PdhFmtCountervalueDouble struct {
	CStatus     uint32
	DoubleValue float64
}

// PdhFmtCountervalueItemDouble XXX
type PdhFmtCountervalueItemDouble struct {
	Name     *uint16
	FmtValue PdhFmtCountervalueDouble
}

const (
	// ErrorSuccess XXX
	ErrorSuccess      = 0
	// DriveRemovable XXX
	DriveRemovable    = 2
	// DriveFixed XXX
	DriveFixed        = 3
	// HkeyLocalMachine XXX
	HkeyLocalMachine = 0x80000002
	// RrfRtRegSz XXX
	RrfRtRegSz      = 0x00000002
	// RrdRtRegDWord XXX
	RrdRtRegDWord   = 0x00000010
	// PdhFmtDouble XXX
	PdhFmtDouble     = 0x00000200
	// PdhInvalidData XXX
	PdhInvalidData   = 0xc0000bc6
)

var (
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	modpdh      = syscall.NewLazyDLL("pdh.dll")

	// RegGetValue XXX
	RegGetValue                 = modadvapi32.NewProc("RegGetValueW")
	// GetSystemInfo XXX
	GetSystemInfo               = modkernel32.NewProc("GetSystemInfo")
	// GetTickCount XXX
	GetTickCount                = modkernel32.NewProc("GetTickCount")
	// GetDiskFreeSpaceEx XXX
	GetDiskFreeSpaceEx          = modkernel32.NewProc("GetDiskFreeSpaceExW")
	// GetLogicalDriveStrings XXX
	GetLogicalDriveStrings      = modkernel32.NewProc("GetLogicalDriveStringsW")
	// GetDriveType XXX
	GetDriveType                = modkernel32.NewProc("GetDriveTypeW")
	// QueryDosDevice XXX
	QueryDosDevice              = modkernel32.NewProc("QueryDosDeviceW")
	// GetVolumeInformationW XXX
	GetVolumeInformationW       = modkernel32.NewProc("GetVolumeInformationW")
	// GlobalMemoryStatusEx XXX
	GlobalMemoryStatusEx        = modkernel32.NewProc("GlobalMemoryStatusEx")
	// GetLastError XXX
	GetLastError				= modkernel32.NewProc("GetLastError")
	// PdhOpenQuery XXX
	PdhOpenQuery                = modpdh.NewProc("PdhOpenQuery")
	// PdhAddCounter XXX
	PdhAddCounter               = modpdh.NewProc("PdhAddCounterW")
	// PdhCollectQueryData XXX
	PdhCollectQueryData         = modpdh.NewProc("PdhCollectQueryData")
	// PdhGetFormattedCounterValue XXX
	PdhGetFormattedCounterValue = modpdh.NewProc("PdhGetFormattedCounterValue")
	// PdhCloseQuery XXX
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
		uintptr(RrdRtRegDWord),
		0,
		uintptr(unsafe.Pointer(&num)),
		uintptr(unsafe.Pointer(&numlen)))
	if ret != ErrorSuccess {
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
		uintptr(RrfRtRegSz),
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
		uintptr(RrfRtRegSz),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)))
	if ret != ErrorSuccess {
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
