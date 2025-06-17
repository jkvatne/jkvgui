package glfw

import "C"
import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

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
	lmpi.CbSize = 24
	ret, _, _ := getMonitorInfo.Call(
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(lmpi)),
	)
	return ret != 0
}

// Use NewEnumDisplayMonitorsCallback to create the callback function.
func EnumDisplayMonitors(hdc HDC, clip *RECT, lpfnEnum, data uintptr) error {
	ret, _, _ := enumDisplayMonitors.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(clip)),
		lpfnEnum,
		data,
	)
	if ret == 0 {
		return fmt.Errorf("w32.EnumDisplayMonitors returned FALSE")
	}
	return nil
}

type Monitor struct {
	hMonitor HMONITOR
	hDc      HDC
	Bounds   RECT
}

var Monitors []Monitor
var initialized = true

func enumMonitorCallback(monitor HMONITOR, hdc HDC, bounds RECT, lParam uintptr) bool {
	m := Monitor{}
	m.hMonitor = monitor
	m.hDc = hdc
	m.Bounds = bounds
	Monitors = append(Monitors, m)
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
func GetMonitors() *[]Monitor {
	if !initialized {
		panic("GLFW not initialized")
	}
	err := EnumDisplayMonitors(0, nil, NewEnumDisplayMonitorsCallback(enumMonitorCallback), 0)
	if err != nil {
		panic(err)
	}
	return &Monitors
}

// GetPhysicalSize returns the size, in millimetres, of the display area of the
// monitor.
//
// Note: Some operating systems do not provide accurate information, either
// because the monitor's EDID Data is incorrect, or because the driver does not
// report it accurately.
func (m *Monitor) GetPhysicalSize() (width, height int) {
	// glfwGetMonitorPhysicalSize(m.Data, &wi, &h)
	panicError()
	return width, height
}

// GetWorkarea returns the position, in screen coordinates, of the upper-left
// corner of the work area of the specified monitor along with the work area
// size in screen coordinates. The work area is defined as the area of the
// monitor not occluded by the operating system task bar where present. If no
// task bar exists then the work area is the monitor resolution in screen
// coordinates.
//
// This function must only be called from the main thread.
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

func IsWindows8Point1OrGreater() bool {
	return true
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
		_, _, err := _GetDpiForMonitor.Call(uintptr(m.hMonitor), uintptr(0), uintptr(unsafe.Pointer(&dpiX)), uintptr(unsafe.Pointer(&dpiY)))
		if !errors.Is(err, syscall.Errno(0)) {
			panic(err)
		}
	} else {
		/*const HDC dc = GetDC(NULL)
		dpiX = GetDeviceCaps(dc, LOGPIXELSX)
		dpiX = GetDeviceCaps(dc, LOGPIXELSY)
		ReleaseDC(NULL, dc)
		*/
	}
	return float32(dpiX) / 72.0, float32(dpiX) / 72.0
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
