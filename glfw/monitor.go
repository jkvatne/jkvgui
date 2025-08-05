package glfw

import (
	"syscall"
	"unsafe"
)

const EDS_ROTATEDMODE = 0x00000004

// Monitor structure
//
type Monitor struct {
	name        [128]byte
	userPointer unsafe.Pointer
	widthMM     int
	heightMM    int
	modes       []_GLFWvidmode
	currentMode _GLFWvidmode

	// This is defined in the window API's platform.h _GLFW_PLATFORM_MONITOR_STATE;
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

type POINT struct {
	X, Y int32
}

// GetMonitors returns a slice of handles for all currently connected monitors.
func GetMonitors() []*Monitor {
	return _glfw.monitors
}

// GetPrimaryMonitor returns the primary monitor. This is usually the monitor
// where elements like the Windows task bar or the OS X menu bar is located.
func GetPrimaryMonitor() *Monitor {
	if len(_glfw.monitors) == 0 {
		return nil
	}
	return _glfw.monitors[0]
}

// GetPos returns the position, in screen coordinates, of the upper-left
// corner of the monitor.
func (m *Monitor) GetPos() (x, y int) {
	// This is from _glfwPlatformGetMonitorPos
	var dm DEVMODEW
	dm.dmSize = uint16(unsafe.Sizeof(dm))
	EnumDisplaySettingsEx(&m.adapterName[0],
		ENUM_CURRENT_SETTINGS,
		&dm,
		EDS_ROTATEDMODE)

	x = int(dm.dmPosition.X)
	y = int(dm.dmPosition.Y)
	return x, y
	// return int(m.Bounds.Left), int(m.Bounds.Top)
}

func enumMonitorCallback(monitor HMONITOR, hdc HDC, bounds RECT, lParam uintptr) bool {
	m := (*Monitor)(unsafe.Pointer(lParam))
	m.hMonitor = monitor
	m.hDc = hdc
	m.Bounds = bounds
	// monitorInfo := GetMonitorInfo(m.hMonitor)
	// Monitors = append(Monitors, &m)
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

// GetPhysicalSize returns the size, in millimetres, of the display area of the monitor.
func (m *Monitor) GetPhysicalSize() (width, height int) {
	return m.widthMM, m.heightMM
}

// GetWorkarea returns the position, in screen coordinates, of the upper-left
// corner of the work area of the specified monitor along with the work area
// size in screen coordinates.
func (m *Monitor) GetWorkarea() (x, y, width, height int) {
	mi := GetMonitorInfo(m.hMonitor)
	x = int(mi.RcWork.Left)
	y = int(mi.RcWork.Top)
	width = int(mi.RcWork.Right - mi.RcWork.Left)
	height = int(mi.RcWork.Bottom - mi.RcWork.Top)
	return x, y, width, height
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
		dpiX, dpiY = GetDpiForMonitor(m.hMonitor, MDT_EFFECTIVE_DPI)
	} else {
		dc := getDC(0)
		dpiX = GetDeviceCaps(dc, LOGPIXELSX)
		dpiX = GetDeviceCaps(dc, LOGPIXELSY)
		releaseDC(0, dc)
	}
	return float32(dpiX) / USER_DEFAULT_SCREEN_DPI, float32(dpiY) / USER_DEFAULT_SCREEN_DPI
}

func (m *Monitor) GetMonitorName() string {
	s := GoStr(&m.name[0])
	return s
}

func GetVideoMode(monitor *Monitor) _GLFWvidmode {
	mode := monitor.currentMode
	var dm DEVMODEW
	dm.dmSize = uint16(unsafe.Sizeof(dm))
	EnumDisplaySettingsEx(&monitor.adapterName[0], ENUM_CURRENT_SETTINGS, &dm, 0)
	mode.Width = int(dm.dmPelsWidth)
	mode.Height = int(dm.dmPelsHeight)
	mode.refreshRate = int(dm.dmDisplayFrequency)
	n := int(max(24, dm.dmBitsPerPel) / 3)
	mode.blueBits = n
	mode.redBits = n
	mode.greenBits = n
	delta := dm.dmBitsPerPel - uint32(mode.redBits*3)
	if delta >= 1 {
		mode.greenBits++
	}
	if delta == 2 {
		mode.redBits++
	}
	return mode
}
