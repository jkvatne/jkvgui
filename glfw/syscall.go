package glfw

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	gdi32                = windows.NewLazySystemDLL("gdi32.dll")
	_GetDeviceCaps       = gdi32.NewProc("GetDeviceCaps")
	_CreateDC            = gdi32.NewProc("CreateDCW")
	_DeleteDC            = gdi32.NewProc("DeleteDC")
	_SwapBuffers         = gdi32.NewProc("SwapBuffers")
	_SetPixelFormat      = gdi32.NewProc("SetPixelFormat")
	_ChoosePixelFormat   = gdi32.NewProc("ChoosePixelFormat")
	_DescribePixelFormat = gdi32.NewProc("DescribePixelFormat")
)

var (
	ntdll                 = windows.NewLazySystemDLL("ntdll.dll")
	_RtlVerifyVersionInfo = ntdll.NewProc("RtlVerifyVersionInfo")
)

var (
	kernel32                 = windows.NewLazySystemDLL("kernel32.dll")
	_GetModuleHandleW        = kernel32.NewProc("GetModuleHandleW")
	_SetThreadExecutionState = kernel32.NewProc("SetThreadExecutionState")
	_TlsAlloc                = kernel32.NewProc("TlsAlloc")
	_TlsGetValue             = kernel32.NewProc("TlsGetValue")
	_TlsSetValue             = kernel32.NewProc("TlsSetValue")
)

var (
	user32                         = windows.NewLazySystemDLL("user32.dll")
	_SetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
	_EnumDisplayMonitors           = user32.NewProc("EnumDisplayMonitors")
	_EnumDisplayDevices            = user32.NewProc("EnumDisplayDevicesW")
	_EnumDisplaySettingsEx         = user32.NewProc("EnumDisplaySettingsExW")
	_GetMonitorInfo                = user32.NewProc("GetMonitorInfoW")
	_AdjustWindowRectEx            = user32.NewProc("AdjustWindowRectEx")
	_CreateWindowEx                = user32.NewProc("CreateWindowExW")
	_DefWindowProc                 = user32.NewProc("DefWindowProcW")
	_DestroyWindow                 = user32.NewProc("DestroyWindow")
	_DispatchMessage               = user32.NewProc("DispatchMessageW")
	_GetWindowRect                 = user32.NewProc("GetWindowRect")
	_GetClientRect                 = user32.NewProc("GetClientRect")
	_GetDC                         = user32.NewProc("GetDC")
	_GetDpiForWindow               = user32.NewProc("GetDpiForWindow")
	_GetKeyState                   = user32.NewProc("GetKeyState")
	_LoadCursor                    = user32.NewProc("LoadCursorW")
	_LoadImage                     = user32.NewProc("LoadImageW")
	_MonitorFromWindow             = user32.NewProc("MonitorFromWindow")
	_PeekMessage                   = user32.NewProc("PeekMessageW")
	_RegisterClassExW              = user32.NewProc("RegisterClassExW")
	_ReleaseDC                     = user32.NewProc("releaseDC")
	_ScreenToClient                = user32.NewProc("ScreenToClient")
	_ShowWindow                    = user32.NewProc("ShowWindow")
	_SetCursor                     = user32.NewProc("SetCursor")
	_SetForegroundWindow           = user32.NewProc("SetForegroundWindow")
	_SetFocus                      = user32.NewProc("SetFocus")
	_SetProcessDPIAware            = user32.NewProc("SetProcessDPIAware")
	_SetWindowPos                  = user32.NewProc("SetWindowPos")
	_TranslateMessage              = user32.NewProc("TranslateMessage")
	_UnregisterClass               = user32.NewProc("UnregisterClassW")
	_BringWindowToTop              = user32.NewProc("BringWindowToTop")
	_GetCursorPos                  = user32.NewProc("GetCursorPos")
	_SystemParametersInfoW         = user32.NewProc("SystemParametersInfoW")
	_GetWindowLongW                = user32.NewProc("GetWindowLongW")
	_SetWindowLongW                = user32.NewProc("SetWindowLongW")
	_GetActiveWindow               = user32.NewProc("GetActiveWindow")
	_GetPropA                      = user32.NewProc("GetPropA")
	_SetPropA                      = user32.NewProc("SetPropA")
)

var (
	shcore            = windows.NewLazySystemDLL("shcore.dll")
	_GetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")
)

type _OSVERSIONINFOEXW struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16
	wServicePackMajor   uint16
	wServicePackMinor   uint16
	wSuiteMask          uint16
	wProductType        uint8
	wReserved           uint8
}

const (
	VER_MAJORVERSION     = 0x0000002
	VER_MINORVERSION     = 0x0000001
	VER_BUILDNUMBER      = 0x0000004
	VER_SERVICEPACKMAJOR = 0x00000020
	WIN32_WINNT_WINBLUE  = 0x0603
)

const (
	DPI_AWARENESS_CONTEXT_UNAWARE              = 0xFFFFFFFFFFFFFFFF
	DPI_AWARENESS_CONTEXT_SYSTEM_AWARE         = 0xFFFFFFFFFFFFFFFE
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE    = 0xFFFFFFFFFFFFFFFD
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = 0xFFFFFFFFFFFFFFFC
	DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED    = 0xFFFFFFFFFFFFFFFB
	PROCESS_DPI_UNAWARE                        = 0
	PROCESS_SYSTEM_DPI_AWARE                   = 1
	PROCESS_PER_MONITOR_DPI_AWARE              = 2
)

type HDC syscall.Handle
type HMONITOR syscall.Handle
type HANDLE syscall.Handle

type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}
type RECT struct {
	Left, Top, Right, Bottom int32
}

type Msg struct {
	Hwnd     syscall.Handle
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}

type Point struct {
	X, Y int32
}

type DISPLAY_DEVICEW struct {
	cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

type WndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CnClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       syscall.Handle
}

type DEVMODEW = struct {
	mDeviceName          [32]uint16
	dmSpecVersion        uint16
	dmDriverVersion      uint16
	dmSize               uint16
	dmDriverExtra        uint16
	dmFields             uint32
	dmPosition           POINTL
	dmDisplayOrientation uint32
	dmDisplayFixedOutput uint32
	dmColor              uint16
	dmDuplex             uint16
	dmYResolution        uint16
	dmTTOption           uint16
	dmCollate            uint16
	dmFormName           [32]uint16
	dmLogPixels          uint16
	dmBitsPerPel         uint32
	dmPelsWidth          int32
	dmPelsHeight         int32
	dmDisplayFlags       uint32
	dmDisplayFrequency   uint32
	dmICMMethod          uint32
	dmICMIntent          uint32
	dmMediaType          uint32
	dmDitherType         uint32
	dmReserved1          uint32
	dmReserved2          uint32
	dmPanningWidth       uint32
	dmPanningHeight      uint32
}

type POINTL = struct {
	X, Y int32
}

func GetActiveWindow() HANDLE {
	r, _, err := _GetActiveWindow.Call()
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetActiveWindow failed, " + err.Error())
	}
	return HANDLE(r)
}

func GetProp(handle HANDLE, key string) uintptr {
	// widestr, _ := syscall.UTF16PtrFromString(key)
	cstr, _ := windows.BytePtrFromString(key)
	r, _, err := _GetPropA.Call(uintptr(handle), uintptr(unsafe.Pointer(cstr)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetProp failed, " + err.Error())
	}
	return r
}

func SetProp(handle HANDLE, key string, data uintptr) {
	// widestr, _ := syscall.UTF16PtrFromString(key)
	cstr, _ := windows.BytePtrFromString(key)
	_, _, err := _SetPropA.Call(uintptr(handle), uintptr(unsafe.Pointer(cstr)), data)
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetProp failed, " + err.Error())
	}
}

func GetKeyState(nVirtKey int) uint16 {
	c, _, _ := _GetKeyState.Call(uintptr(nVirtKey))
	return uint16(c)
}

func GetModuleHandle() (syscall.Handle, error) {
	h, _, err := _GetModuleHandleW.Call(uintptr(0))
	if h == 0 {
		return 0, fmt.Errorf("GetModuleHandleW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func RegisterClassEx(cls *WndClassEx) (uint16, error) {
	a, _, err := _RegisterClassExW.Call(uintptr(unsafe.Pointer(cls)))
	if a == 0 {
		return 0, fmt.Errorf("RegisterClassExW failed: %v", err)
	}
	return uint16(a), nil
}

func LoadImage(hInst syscall.Handle, res uint32, typ uint32, cx, cy int, fuload uint32) (syscall.Handle, error) {
	h, _, err := _LoadImage.Call(uintptr(hInst), uintptr(res), uintptr(typ), uintptr(cx), uintptr(cy), uintptr(fuload))
	if h == 0 {
		return 0, fmt.Errorf("LoadImageW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func PeekMessage(m *Msg, hwnd syscall.Handle, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) bool {
	r, _, _ := _PeekMessage.Call(uintptr(unsafe.Pointer(m)), uintptr(hwnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax), uintptr(wRemoveMsg))
	return r != 0
}

func TranslateMessage(m *Msg) {
	_TranslateMessage.Call(uintptr(unsafe.Pointer(m)))
}

func DispatchMessage(m *Msg) {
	_DispatchMessage.Call(uintptr(unsafe.Pointer(m)))
}

func CreateWindowEx(dwExStyle uint32, lpClassName uint16, lpWindowName string, dwStyle uint32, x, y, w, h int32, hWndParent, hMenu, hInstance syscall.Handle, lpParam uintptr) (syscall.Handle, error) {
	wname, _ := syscall.UTF16PtrFromString(lpWindowName)
	hwnd, _, err := _CreateWindowEx.Call(
		uintptr(dwExStyle),
		uintptr(lpClassName),
		uintptr(unsafe.Pointer(wname)),
		uintptr(dwStyle),
		uintptr(x), uintptr(y),
		uintptr(w), uintptr(h),
		uintptr(hWndParent),
		uintptr(hMenu),
		uintptr(hInstance),
		uintptr(lpParam))
	if hwnd == 0 {
		return 0, fmt.Errorf("CreateWindowEx failed: %v", err)
	}
	return syscall.Handle(hwnd), nil
}

func DestroyWindow(h syscall.Handle) {
	_, _, err := _DestroyWindow.Call(uintptr(h))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("DestroyWindow failed, " + err.Error())
	}
}

func UnregisterClass(class uint16, instance syscall.Handle) {
	_, _, err := _UnregisterClass.Call(uintptr(class), uintptr(instance))
	if !errors.Is(err, syscall.Errno(0)) {
		// TODO panic("UnregisterClass failed, " + err.Error())
	}
}

func IsWindows10Version1607OrGreater() bool {
	var osvi _OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))
	osvi.dwMajorVersion = 10
	osvi.dwMinorVersion = 0
	osvi.dwBuildNumber = 14393
	var mask uint32 = VER_MAJORVERSION | VER_MINORVERSION | VER_BUILDNUMBER
	r, _, err := _RtlVerifyVersionInfo.Call(uintptr(unsafe.Pointer(&osvi)), uintptr(mask), uintptr(0x80000000000000db))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetProcessDpiAwarenessContext failed, " + err.Error())
	}
	return r == 0
}

func IsWindows10Version1703OrGreater() bool {
	var osvi _OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))
	osvi.dwMajorVersion = 10
	osvi.dwMinorVersion = 0
	osvi.dwBuildNumber = 15063
	var mask uint32 = VER_MAJORVERSION | VER_MINORVERSION | VER_BUILDNUMBER
	r, _, err := _RtlVerifyVersionInfo.Call(uintptr(unsafe.Pointer(&osvi)), uintptr(mask), uintptr(0x80000000000000db))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("isWindows10Version1703OrGreater failed, " + err.Error())
	}
	return r == 0
}

func IsWindows8Point1OrGreater() bool {
	var osvi _OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))
	osvi.dwMajorVersion = uint32(WIN32_WINNT_WINBLUE >> 8)
	osvi.dwMinorVersion = uint32(WIN32_WINNT_WINBLUE & 0xFF)
	osvi.wServicePackMajor = 0
	var mask uint32 = VER_MAJORVERSION | VER_MINORVERSION | VER_SERVICEPACKMAJOR
	r, _, err := _RtlVerifyVersionInfo.Call(uintptr(unsafe.Pointer(&osvi)), uintptr(mask), uintptr(0x800000000001801b))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetProcessDpiAwarenessContext failed, " + err.Error())
	}
	return r == 0
}

func SetProcessDpiAwareness() {
	if IsWindows10Version1703OrGreater() {
		_, _, err := _SetProcessDpiAwarenessContext.Call(uintptr(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("SetProcessDpiAwarenessContext failed, " + err.Error())
		}
	} else if IsWindows8Point1OrGreater() {
		_, _, err := _SetProcessDpiAwarenessContext.Call(uintptr(PROCESS_PER_MONITOR_DPI_AWARE))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("SetProcessDpiAwarenessContext failed, " + err.Error())
		}
	} else if isWindowsVistaOrGreater() {
		_, _, _ = _SetProcessDPIAware.Call()
	}
}

func SetWindowPos(hwnd syscall.Handle, after syscall.Handle, x, y, w, h int32, flags uint32) {
	_, _, err := _SetWindowPos.Call(uintptr(hwnd), uintptr(after), uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(flags))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetWindowPos failed, " + err.Error())
	}
}

func GetWindowLongW(hWnd syscall.Handle, index int32) uint32 {
	r1, _, err := _GetWindowLongW.Call(uintptr(hWnd), uintptr(index))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("GetWindowLongW failed, " + err.Error())
	}
	return uint32(r1)
}

func SetWindowLongW(hWnd syscall.Handle, index int32, newValue uint32) {
	_, _, err := _SetWindowLongW.Call(uintptr(hWnd), uintptr(index), uintptr(newValue))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("GetWindowLongW failed, " + err.Error())
	}
}

func EnumDisplayDevices(device uintptr, no int, adapter *DISPLAY_DEVICEW, flags uint32) bool {
	ret, _, err := _EnumDisplayDevices.Call(device, uintptr(no), uintptr(unsafe.Pointer(adapter)), uintptr(flags))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("EnumDisplayDevices failed")
	}
	return ret == 1
}

func GetMonitorInfo(hMonitor HMONITOR) *MONITORINFO {
	lmpi := MONITORINFO{}
	lmpi.CbSize = uint32(unsafe.Sizeof(lmpi))
	_, _, err := _GetMonitorInfo.Call(uintptr(hMonitor), uintptr(unsafe.Pointer(&lmpi)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetMonitorInfo failed, " + err.Error())
	}
	return &lmpi
}

func GetDeviceCaps(dc HDC, flags int) int {
	r1, _, err := _GetDeviceCaps.Call(uintptr(unsafe.Pointer(dc)), uintptr(flags))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("GetDeviceCaps failed, " + err.Error())
	}
	return int(r1)
}

func GetDpiForMonitor(h HMONITOR, kind uint32) (dpiX int, dpiY int) {
	r1, _, err := _GetDpiForMonitor.Call(uintptr(h), uintptr(kind), uintptr(unsafe.Pointer(&dpiX)), uintptr(unsafe.Pointer(&dpiY)))
	if !errors.Is(err, syscall.Errno(0)) || r1 != 0 {
		panic("GetDpiForMonitor failed, " + err.Error())
	}
	return dpiX, dpiY
}

func EnumDisplayMonitors(hdc HDC, clip *RECT, lpfnEnum uintptr, data uintptr) error {
	ret, _, _ := _EnumDisplayMonitors.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(clip)),
		lpfnEnum,
		uintptr(unsafe.Pointer(data)),
	)
	if ret == 0 {
		return fmt.Errorf("w32.EnumDisplayMonitors returned FALSE")
	}
	return nil
}

func isWindowsVistaOrGreater() bool {
	return true
}

func AdjustWindowRectEx(rect *RECT, style uint32, menu int, exStyle uint32) {
	_, _, err := _AdjustWindowRectEx.Call(uintptr(unsafe.Pointer(rect)), uintptr(style), uintptr(menu), uintptr(exStyle))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("AdjustWindowRectEx failed, " + err.Error())
	}
}

func GetDpiForWindow(handle syscall.Handle) int {
	r, _, err := _GetDpiForWindow.Call(uintptr(handle))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetDpiForWindow failed, " + err.Error())
	}
	return int(r)
}

func AdjustWindowRectExForDpi(rect *RECT, style uint32, menu int, exStyle uint32, dpi int) {
	_, _, err := _AdjustWindowRectEx.Call(uintptr(unsafe.Pointer(rect)), uintptr(style), uintptr(menu), uintptr(exStyle), uintptr(dpi))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("AdjustWindowRectExForDpi failed, " + err.Error())
	}
}

func EnumDisplaySettingsEx(name *uint16, mode int, dm *DEVMODEW, flags int) {
	_, _, err := _EnumDisplaySettingsEx.Call(uintptr(unsafe.Pointer(name)), uintptr(mode), uintptr(unsafe.Pointer(dm)), uintptr(flags))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("EnumDisplySettingsEx failed, " + err.Error())
	}
}
func TlsSetValue(index int, value uintptr) {
	_, _, err := _TlsSetValue.Call(uintptr(index), value)
	if !errors.Is(err, syscall.Errno(0)) {
		panic("_TTlsGetValue failed, " + err.Error())
	}
}

func TlsGetValue(index int) uintptr {
	r, _, err := _TlsGetValue.Call(uintptr(index))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("_TTlsGetValue failed, " + err.Error())
	}
	return r
}

func TlsAlloc() int {
	r, _, err := _TlsAlloc.Call()
	if !errors.Is(err, syscall.Errno(0)) {
		panic("_TlsAlloc failed, " + err.Error())
	}
	return int(r)
}

func glfwPlatformSetTls(tls *_GLFWtls, value uintptr) {
	if !tls.allocated {
		panic("_glfwPlatformGetTls failed: tls not allocated")
	}
	TlsSetValue(tls.index, value)
}

func glfwPlatformGetTls(tls *_GLFWtls) uintptr {
	if !tls.allocated {
		panic("_glfwPlatformGetTls failed: tls not allocated")
	}
	return TlsGetValue(tls.index)
}

func TlsFree(index int) {
	// TODO
}
func glfwPlatformDestroyTls(tls *_GLFWtls) {
	if tls.allocated {
		TlsFree(tls.index)
	}
}

func glfwPlatformCreateTls(tls *_GLFWtls) error {
	if tls.allocated {
		return fmt.Errorf("glfwPlatformCreateTls: already allocated")
	}
	tls.index = TlsAlloc()
	if tls.index == 4294967295 { // TLS_OUT_OF_INDEXES
		return fmt.Errorf("glfwPlatformCreateTls: Failed to allocate TLS index")
	}
	tls.allocated = true
	return nil
}
