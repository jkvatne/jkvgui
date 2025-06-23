package glfw

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// Monitor structure
//
type Monitor struct {
	name        [128]byte
	userPointer unsafe.Pointer
	widthMM     int
	heightMM    int
	modes       *GLFWvidmode
	modeCount   int
	currentMode GLFWvidmode

	// This is defined in the window API's platform.h
	// _GLFW_PLATFORM_MONITOR_STATE;
	hMonitor HMONITOR
	hDc      HDC
	// The window whose video mode is current on this monitor
	window *_GLFWwindow
	Bounds RECT
	// This size matches the static size of DISPLAY_DEVICE.DeviceName
	adapterName       [32]uint16
	displayName       [32]uint16
	publicAdapterName string
	publicDisplayName string
	modesPruned       bool
	modeChanged       bool
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162897.aspx
type RECT struct {
	Left, Top, Right, Bottom int32
}

type HANDLE windows.Handle
type HDC HANDLE
type HMONITOR HANDLE

type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}

// GetMonitorInfo automatically sets the MONITORINFO's CbSize field.
func GetMonitorInfo(hMonitor HMONITOR, lmpi *MONITORINFO) bool {
	if lmpi != nil {
		lmpi.CbSize = uint32(unsafe.Sizeof(*lmpi))
	}
	// lmpi.CbSize = 24
	ret, _, err := _GetMonitorInfo.Call(uintptr(hMonitor), uintptr(unsafe.Pointer(lmpi)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetMonitorInfo failed, " + err.Error())
	}
	return ret != 0
}

// Use NewEnumDisplayMonitorsCallback to create the callback function.
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

func enumMonitorCallback(monitor HMONITOR, hdc HDC, bounds RECT, lParam uintptr) bool {
	m := (*Monitor)(unsafe.Pointer(lParam))
	m.hMonitor = monitor
	m.hDc = hdc
	m.Bounds = bounds
	// Monitors = append(Monitors, &m)
	var monitorInfo MONITORINFO
	GetMonitorInfo(monitor, &monitorInfo)
	// slog.Info("EnumMonitors RcMonitor", "left", monitorInfo.RcMonitor.Left, "top", monitorInfo.RcMonitor.Top, "right", monitorInfo.RcMonitor.Right, "bottom", monitorInfo.RcMonitor.Bottom)
	// slog.Info("EnumMonitors RcWork", "left", monitorInfo.RcWork.Left, "top", monitorInfo.RcWork.Top, "right", monitorInfo.RcWork.Right, "bottom", monitorInfo.RcWork.Bottom)
	// var scaleFactor int
	// r1, _, err := _GetScaleFactorForMonitor.Call(uintptr(m.hMonitor), uintptr(unsafe.Pointer(&scaleFactor)))
	// if !errors.Is(err, syscall.Errno(0)) || r1 != 0 {
	//	panic("_GetScaleFactorForMonitor failed, " + err.Error())
	// }
	// slog.Info("ScaleFactor", "hmonitor", monitor, "value", scaleFactor)
	// if monitorInfo.DwFlags != 0 {
	//	slog.Info("Primary monitor")
	// }
	return true
}

// NewEnumDisplayMonitorsCallback is used in EnumDisplayMonitors to create the callback.
func NewEnumDisplayMonitorsCallback(callback func(monitor HMONITOR, hdc HDC, bounds RECT, lParam uintptr) bool) uintptr {
	return syscall.NewCallback(
		func(monitor HMONITOR, hdc HDC, bounds *RECT, lParam uintptr) uintptr {
			var r RECT
			if bounds != nil {
				r = *bounds
			}
			if callback(monitor, hdc, r, lParam) {
				return 1
			}
			return 0
		},
	)
}

// GetMonitors returns a slice of handles for all currently connected monitors.
func GetMonitors() []*Monitor {
	return _glfw.monitors
}

// GetPhysicalSize returns the size, in millimetres, of the display area of the monitor.

func (m *Monitor) GetPhysicalSize() (width, height int) {
	return m.widthMM, m.heightMM
}

// GetWorkarea returns the position, in screen coordinates, of the upper-left
// corner of the work area of the specified monitor along with the work area
// size in screen coordinates.
func (m *Monitor) GetWorkarea() (x, y, width, height int) {
	var mi MONITORINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))
	_, _, err := _GetMonitorInfo.Call(uintptr(m.hMonitor), uintptr(unsafe.Pointer(&mi)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic(err)
	}
	x = int(mi.RcWork.Left)
	y = int(mi.RcWork.Top)
	width = int(mi.RcWork.Right - mi.RcWork.Left)
	height = int(mi.RcWork.Bottom - mi.RcWork.Top)
	return x, y, width, height
}

func GetDeviceCaps(dc HDC, flags int) int {
	r1, _, err := _GetDeviceCaps.Call(uintptr(unsafe.Pointer(dc)), uintptr(flags))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("GetDeviceCaps failed, " + err.Error())
	}
	return int(r1)
}

// GetContentScale function retrieves the content scale for the specified monitor.
// The content scale is the ratio between the current DPI and the platform's
// default DPI. If you scale all pixel dimensions by this scale then your content
// should appear at an appropriate size. This is especially important for text
// and any UI elements.
//
// This function must only be called from the main thread.
func (m *Monitor) GetContentScale() (float32, float32) {
	var dpiX, dpiY int
	if IsWindows8Point1OrGreater() {
		r1, _, err := _GetDpiForMonitor.Call(uintptr(m.hMonitor), uintptr(0), uintptr(unsafe.Pointer(&dpiX)), uintptr(unsafe.Pointer(&dpiY)))
		if !errors.Is(err, syscall.Errno(0)) || r1 != 0 {
			panic("GetContentScale failed, " + err.Error())
		}
	} else {
		dc := GetDC(0)
		dpiX = GetDeviceCaps(dc, LOGPIXELSX)
		dpiX = GetDeviceCaps(dc, LOGPIXELSY)
		ReleaseDC(0, dc)
	}
	return float32(dpiX) / USER_DEFAULT_SCREEN_DPI, float32(dpiX) / USER_DEFAULT_SCREEN_DPI
}

// GetPrimaryMonitor returns the primary monitor. This is usually the monitor
// where elements like the Windows task bar or the OS X menu bar is located.
func GetPrimaryMonitor() *Monitor {
	/* m := C.glfwGetPrimaryMonitor()
	panicError()
	if m == nil {
		return nil
	}*/
	return nil // &Monitor{m}
}

// GetPos returns the position, in screen coordinates, of the upper-left
// corner of the monitor.
func (m *Monitor) GetPos() (x, y int) {
	/*
		// C.glfwGetMonitorPos(m.Data, &xpos, &ypos)
		dm.dmSize = sizeof(dm);

		EnumDisplaySettingsExW(monitor->Win32.adapterName,
			ENUM_CURRENT_SETTINGS,
			&dm,
			EDS_ROTATEDMODE);

		*xpos = dm.dmPosition.x;
		*ypos = dm.dmPosition.y;

		panicError() */
	return int(m.Bounds.Left), int(m.Bounds.Top)
}
