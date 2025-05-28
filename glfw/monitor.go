package glfw

import "C"

// GetMonitors returns a slice of handles for all currently connected monitors.
func GetMonitors() []*Monitor {
	var length = 1
	mC := &Monitor{} // C.glfwGetMonitors((*C.int)(unsafe.Pointer(&length)))
	{
		assert(count != NULL);

		*count = 0;

		_GLFW_REQUIRE_INIT_OR_RETURN(NULL);

		*count = _glfw.monitorCount;
		return (GLFWmonitor**) _glfw.monitors;

	panicError()
	if mC == nil {
		return nil
	}
	m := make([]*Monitor, length)
	for i := 0; i < length; i++ {
		m[i] = &Monitor{} // &Monitor{C.GetMonitorAtIndex(mC, C.int(i))}
	}
	return m
}

// GetPhysicalSize returns the size, in millimetres, of the display area of the
// monitor.
//
// Note: Some operating systems do not provide accurate information, either
// because the monitor's EDID data is incorrect, or because the driver does not
// report it accurately.
func (m *Monitor) GetPhysicalSize() (width, height int) {
	var wi, h C.int
	// C.glfwGetMonitorPhysicalSize(m.data, &wi, &h)
	panicError()
	return int(wi), int(h)
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
	var cX, cY, cWidth, cHeight C.int
	// C.glfwGetMonitorWorkarea(m.data, &cX, &cY, &cWidth, &cHeight)
	x, y, width, height = int(cX), int(cY), int(cWidth), int(cHeight)
	return
}

// GetContentScale function retrieves the content scale for the specified monitor.
// The content scale is the ratio between the current DPI and the platform's
// default DPI. If you scale all pixel dimensions by this scale then your content
// should appear at an appropriate size. This is especially important for text
// and any UI elements.
//
// This function must only be called from the main thread.
func (m *Monitor) GetContentScale() (float32, float32) {
	var x, y C.float
	// C.glfwGetMonitorContentScale(m.data, &x, &y)
	return float32(x), float32(y)
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
	var xpos, ypos C.int
	// C.glfwGetMonitorPos(m.data, &xpos, &ypos)
	panicError()
	return int(xpos), int(ypos)
}
