package glfw

import "C"
import (
	"fmt"
	"golang.design/x/clipboard"
)

// MouseButton definitions
type MouseButton int

const (
	MouseButtonFirst  MouseButton = 0
	MouseButtonLeft   MouseButton = 0
	MouseButtonRight  MouseButton = 1
	MouseButtonMiddle MouseButton = 2
	MouseButtonLast   MouseButton = 2
)

// Exported cursor types
const (
	ArrowCursor     = 0x00036001
	IBeamCursor     = 0x00036002
	CrosshairCursor = 0x00036003
	HandCursor      = 0x00036004
	HResizeCursor   = 0x00036005
	VResizeCursor   = 0x00036006
)

// Exported hints
const (
	True                    = 1
	False                   = 0
	OpenGLForwardCompatible = GLFW_OPENGL_FORWARD_COMPAT
	Focused                 = GLFW_FOCUSED
	Resizable               = GLFW_RESIZABLE
	Visible                 = GLFW_VISIBLE
	Decorated               = GLFW_DECORATED
	AutoIconify             = GLFW_AUTO_ICONIFY
	Floating                = GLFW_FLOATING
	Maximized               = GLFW_OPENGL_PROFILE
	Samples                 = GLFW_SAMPLES
	ContextVersionMajor     = GLFW_CONTEXT_VERSION_MAJOR
	ForwardCompatible       = GLFW_OPENGL_FORWARD_COMPAT
	OpenGLProfile           = GLFW_OPENGL_PROFILE
	OpenGLCoreProfile       = GLFW_OPENGL_CORE_PROFILE
	ContextVersionMinor     = GLFW_CONTEXT_VERSION_MINOR
)

type Action int

type StandardCursor uint16

type Hint uint32

// Window represents a Window.
type Window = _GLFWwindow

type Cursor struct {
	next   *Cursor
	handle HANDLE
}

// PollEvents processes only those events that have already been received and
// then returns immediately. Processing events will cause the Window and input
// callbacks associated with those events to be called.
func PollEvents() {
	glfwPollEvents()
}

func WindowHint(hint int, value int) error {
	switch hint {
	case GLFW_RED_BITS:
		_glfw.hints.framebuffer.redBits = value
	case GLFW_GREEN_BITS:
		_glfw.hints.framebuffer.greenBits = value
	case GLFW_BLUE_BITS:
		_glfw.hints.framebuffer.blueBits = value
	case GLFW_ALPHA_BITS:
		_glfw.hints.framebuffer.alphaBits = value
	case GLFW_DEPTH_BITS:
		_glfw.hints.framebuffer.depthBits = value
	case GLFW_STENCIL_BITS:
	case GLFW_ACCUM_RED_BITS:
		_glfw.hints.framebuffer.accumRedBits = value
	case GLFW_ACCUM_GREEN_BITS:
		_glfw.hints.framebuffer.accumGreenBits = value
	case GLFW_ACCUM_BLUE_BITS:
		_glfw.hints.framebuffer.accumBlueBits = value
	case GLFW_ACCUM_ALPHA_BITS:
		_glfw.hints.framebuffer.accumAlphaBits = value
	case GLFW_AUX_BUFFERS:
		_glfw.hints.framebuffer.auxBuffers = value
	case GLFW_DOUBLEBUFFER:
		_glfw.hints.framebuffer.doublebuffer = value != 0
	case GLFW_TRANSPARENT_FRAMEBUFFER:
		_glfw.hints.framebuffer.transparent = value != 0
	case GLFW_SAMPLES:
		_glfw.hints.framebuffer.samples = value
	case GLFW_SRGB_CAPABLE:
		_glfw.hints.framebuffer.sRGB = value != 0
	case GLFW_RESIZABLE:
		_glfw.hints.window.resizable = value != 0
	case GLFW_DECORATED:
		_glfw.hints.window.decorated = value != 0
	case GLFW_FOCUSED:
		_glfw.hints.window.focused = value != 0
	case GLFW_AUTO_ICONIFY:
		_glfw.hints.window.autoIconify = value != 0
	case GLFW_FLOATING:
		_glfw.hints.window.floating = value != 0
	case GLFW_MAXIMIZED:
		_glfw.hints.window.maximized = value != 0
	case GLFW_VISIBLE:
		_glfw.hints.window.visible = value != 0
	case GLFW_POSITION_X:
		_glfw.hints.window.xpos = value
	case GLFW_POSITION_Y:
		_glfw.hints.window.ypos = value
	case GLFW_SCALE_TO_MONITOR:
		_glfw.hints.window.scaleToMonitor = value != 0
	case GLFW_SCALE_FRAMEBUFFER:
	case GLFW_COCOA_RETINA_FRAMEBUFFER:
		_glfw.hints.window.scaleFramebuffer = value != 0
	case GLFW_CENTER_CURSOR:
		_glfw.hints.window.centerCursor = value != 0
	case GLFW_FOCUS_ON_SHOW:
		_glfw.hints.window.focusOnShow = value != 0
	case GLFW_MOUSE_PASSTHROUGH:
		_glfw.hints.window.mousePassthrough = value != 0
	case GLFW_CLIENT_API:
		_glfw.hints.context.client = value
	case GLFW_CONTEXT_CREATION_API:
		_glfw.hints.context.source = value
	case GLFW_CONTEXT_VERSION_MAJOR:
		_glfw.hints.context.major = value
	case GLFW_CONTEXT_VERSION_MINOR:
		_glfw.hints.context.minor = value
	case GLFW_CONTEXT_ROBUSTNESS:
		_glfw.hints.context.robustness = value
	case GLFW_OPENGL_FORWARD_COMPAT:
		_glfw.hints.context.forward = value != 0
	case GLFW_CONTEXT_DEBUG:
		_glfw.hints.context.debug = value != 0
	case GLFW_CONTEXT_NO_ERROR:
		_glfw.hints.context.noerror = value != 0
	case GLFW_OPENGL_PROFILE:
		_glfw.hints.context.profile = value
	case GLFW_CONTEXT_RELEASE_BEHAVIOR:
		_glfw.hints.context.release = value
	case GLFW_REFRESH_RATE:
		_glfw.hints.refreshRate = value
	default:
		return fmt.Errorf("Invalid window hint %d with value %d", hint, value)
	}
	return nil
}

// GetClipboardString returns the contents of the system clipboard
// if it contains or is convertible to a UTF-8 encoded string.
// This function may only be called from the main thread.
func GetClipboardString() string {
	return glfwGetClipboardString()
}

// SetClipboardString sets the system clipboard to the specified UTF-8 encoded string.
// This function may only be called from the main thread.
func SetClipboardString(str string) {
	glfwSetClipboardString(str)
}

// CreateStandardCursor returns a cursor with a standard shape,
// that can be set for a Window with SetCursor.
func CreateStandardCursor(shape int) *Cursor {
	var cursor = Cursor{}
	cursor.next = _glfw.cursorListHead
	_glfw.cursorListHead = &cursor
	var id uint16
	switch shape {
	case ArrowCursor:
		id = IDC_ARROW
	case IBeamCursor:
		id = IDC_IBEAM
	case CrosshairCursor:
		id = IDC_CROSS
	case HResizeCursor:
		id = IDC_SIZEWE
	case VResizeCursor:
		id = IDC_SIZENS
	case HandCursor:
		id = IDC_HAND
	default:
		panic("Win32: Unknown or unsupported standard cursor")
	}
	cursor.handle = LoadCursor(id)
	if cursor.handle == 0 {
		panic("Win32: Failed to create standard cursor")
	}
	return &cursor
}

func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	wnd, err := glfwCreateWindow(width, height, title, monitor, share)
	if err != nil {
		return nil, fmt.Errorf("glfwCreateWindow failed: %v", err)
	}
	return wnd, nil
}

// SwapBuffers swaps the front and back buffers of the Window.
func (w *Window) SwapBuffers() {
	glfwSwapBuffers(w)
}

// SetCursor sets the cursor image to be used when the cursor is over the client area
func (w *Window) SetCursor(c *Cursor) {
	glfwSetCursor(w, c)
}

// SetPos sets the position, in screen coordinates, of the Window's upper-left corner
func (w *Window) SetPos(xPos, yPos int) {
	rect := RECT{Left: int32(xPos), Top: int32(yPos), Right: int32(xPos), Bottom: int32(yPos)}
	AdjustWindowRect(&rect, getWindowStyle(w), 0, getWindowExStyle(w), GetDpiForWindow(w.Win32.handle), "glfwSetWindowPos")
	SetWindowPos(w.Win32.handle, 0, rect.Left, rect.Top, 0, 0, SWP_NOACTIVATE|SWP_NOZORDER|SWP_NOSIZE)
}

// SetSize sets the size, in screen coordinates, of the client area of the Window.
func (w *Window) SetSize(width, height int) {
	if w.monitor != nil {
		if w.monitor.window == w {
			acquireMonitor(w)
			fitToMonitor(w)
		}
	} else {
		rect := RECT{0, 0, int32(width), int32(height)}
		AdjustWindowRect(&rect, getWindowStyle(w), 0, getWindowExStyle(w), GetDpiForWindow(w.Win32.handle), "glfwSetWindowSize")
		SetWindowPos(w.Win32.handle, 0, 0, 0, int32(width), int32(height), SWP_NOACTIVATE|SWP_NOOWNERZORDER|SWP_NOMOVE|SWP_NOZORDER)
	}
}

// SetMonitor sets the monitor that the window uses for full screen mode or,
// if the monitor is NULL, makes it windowed mode.
func (w *Window) SetMonitor(monitor *Monitor, xpos, ypos, width, height, refreshRate int) {
	glfwSetWindowMonitor(w, monitor, xpos, ypos, width, height, refreshRate)
}

// GetMonitor returns the handle of the monitor that the window is in fullscreen on.
// Returns nil if the window is in windowed mode.
func (w *Window) GetMonitor() *Monitor {
	return glfwGetWindowMonitor(w)
}

// GetContentScale function retrieves the content scale for the specified
// Window. The content scale is the ratio between the current DPI and the
// platform's default DPI.
func (w *Window) GetContentScale() (float32, float32) {
	return glfwGetContentScale(w)
}

// GetFrameSize retrieves the size, in screen coordinates, of each edge of the frame
// This size includes the title bar if the Window has one.
func (w *Window) GetFrameSize() (left, top, right, bottom int) {
	var l, t, r, b int
	glfwGetWindowFrameSize(w, &l, &t, &r, &b)
	return l, t, r, b
}

// GetCursorPos returns the last reported position of the cursor.
func (w *Window) GetCursorPos() (x float64, y float64) {
	var xPos, yPos int
	glfwGetCursorPos(w, &xPos, &yPos)
	return float64(xPos), float64(yPos)
}

// GetSize returns the size, in screen coordinates, of the client area of the
// specified Window.
func (w *Window) GetSize() (width int, height int) {
	var wi, h int
	glfwGetWindowSize(w, &wi, &h)
	return wi, h
}

// Focus brings the specified Window to front and sets input focus.
func (w *Window) Focus() {
	glfwFocusWindow(w)
}

// ShouldClose reports the close flag value for the specified Window.
func (w *Window) ShouldClose() bool {
	return w.shouldClose
}

// Destroy destroys the specified window and its context.
func (w *Window) Destroy() {
	// windows.remove(w.data)
	glfwDestroyWindow(w)
}

// Show makes the Window visible if it was previously hidden.
func (w *Window) Show() {
	if w.monitor != nil {
		return
	}
	glfwShowWindow(w)
	if w.focusOnShow {
		glfwFocusWindow(w)
	}
}

func (w *Window) MakeContextCurrent() {
	_ = w.context.makeCurrent(w)
}

func (w *Window) Iconify() {
	w.Win32.maximized = false
	w.Win32.iconified = true
	glfwShowWindow(w)
}

func (w *Window) Maximize() {
	w.Win32.iconified = false
	w.Win32.maximized = true
	glfwShowWindow(w)
}

// Terminate destroys all remaining Windows, frees any allocated resources and
// sets the library to an uninitialized state.
func Terminate() {
	glfwTerminate()
}

// Init is glfwInit(void)
func Init() error {
	// Repeated calls do nothing
	if _glfw.initialized {
		return nil
	}
	_glfw.hints.init = _GLFWinitconfig{}
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	if err = glfwPlatformInit(); err != nil {
		return err
	}
	err = glfwPlatformCreateTls(&_glfw.errorSlot)
	if err != nil {
		return err
	}
	err = glfwPlatformCreateTls(&_glfw.contextSlot)
	if err != nil {
		return err
	}
	return nil
}
