package glfw

import (
	"errors"
	"fmt"
	"golang.design/x/clipboard"
	"log/slog"
	"syscall"
)

// Mouse buttons.
type MouseButton int

const (
	MouseButtonFirst  MouseButton = 0
	MouseButtonLeft   MouseButton = 0
	MouseButtonRight  MouseButton = 1
	MouseButtonMiddle MouseButton = 2
	MouseButtonLast   MouseButton = 2
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

func WindowHint(hint int, value int) {
	switch hint {
	case GLFW_RED_BITS:
		_glfw.hints.framebuffer.redBits = value
		return
	case GLFW_GREEN_BITS:
		_glfw.hints.framebuffer.greenBits = value
		return
	case GLFW_BLUE_BITS:
		_glfw.hints.framebuffer.blueBits = value
		return
	case GLFW_ALPHA_BITS:
		_glfw.hints.framebuffer.alphaBits = value
		return
	case GLFW_DEPTH_BITS:
		_glfw.hints.framebuffer.depthBits = value
		return
	case GLFW_STENCIL_BITS:
		_glfw.hints.framebuffer.stencilBits = value
		return
	case GLFW_ACCUM_RED_BITS:
		_glfw.hints.framebuffer.accumRedBits = value
		return
	case GLFW_ACCUM_GREEN_BITS:
		_glfw.hints.framebuffer.accumGreenBits = value
		return
	case GLFW_ACCUM_BLUE_BITS:
		_glfw.hints.framebuffer.accumBlueBits = value
		return
	case GLFW_ACCUM_ALPHA_BITS:
		_glfw.hints.framebuffer.accumAlphaBits = value
		return
	case GLFW_AUX_BUFFERS:
		_glfw.hints.framebuffer.auxBuffers = value
		return
	case GLFW_DOUBLEBUFFER:
		_glfw.hints.framebuffer.doublebuffer = value != 0
		return
	case GLFW_TRANSPARENT_FRAMEBUFFER:
		_glfw.hints.framebuffer.transparent = value != 0
		return
	case GLFW_SAMPLES:
		_glfw.hints.framebuffer.samples = value
		return
	case GLFW_SRGB_CAPABLE:
		_glfw.hints.framebuffer.sRGB = value != 0
		return
	case GLFW_RESIZABLE:
		_glfw.hints.window.resizable = value != 0
		return
	case GLFW_DECORATED:
		_glfw.hints.window.decorated = value != 0
		return
	case GLFW_FOCUSED:
		_glfw.hints.window.focused = value != 0
		return
	case GLFW_AUTO_ICONIFY:
		_glfw.hints.window.autoIconify = value != 0
		return
	case GLFW_FLOATING:
		_glfw.hints.window.floating = value != 0
		return
	case GLFW_MAXIMIZED:
		_glfw.hints.window.maximized = value != 0
		return
	case GLFW_VISIBLE:
		_glfw.hints.window.visible = value != 0
		return
	case GLFW_POSITION_X:
		_glfw.hints.window.xpos = value
		return
	case GLFW_POSITION_Y:
		_glfw.hints.window.ypos = value
		return
	case GLFW_SCALE_TO_MONITOR:
		_glfw.hints.window.scaleToMonitor = value != 0
		return
	case GLFW_SCALE_FRAMEBUFFER:
	case GLFW_COCOA_RETINA_FRAMEBUFFER:
		// _glfw.hints.window.scaleFramebuffer = value != 0
		return
	case GLFW_CENTER_CURSOR:
		_glfw.hints.window.centerCursor = value != 0
		return
	case GLFW_FOCUS_ON_SHOW:
		_glfw.hints.window.focusOnShow = value != 0
		return
	case GLFW_MOUSE_PASSTHROUGH:
		_glfw.hints.window.mousePassthrough = value != 0
		return
	case GLFW_CLIENT_API:
		_glfw.hints.context.client = value
		return
	case GLFW_CONTEXT_CREATION_API:
		_glfw.hints.context.source = value
		return
	case GLFW_CONTEXT_VERSION_MAJOR:
		_glfw.hints.context.major = value
		return
	case GLFW_CONTEXT_VERSION_MINOR:
		_glfw.hints.context.minor = value
		return
	case GLFW_CONTEXT_ROBUSTNESS:
		_glfw.hints.context.robustness = value
		return
	case GLFW_OPENGL_FORWARD_COMPAT:
		_glfw.hints.context.forward = value != 0
		return
	case GLFW_CONTEXT_DEBUG:
		_glfw.hints.context.debug = value != 0
		return
	case GLFW_CONTEXT_NO_ERROR:
		_glfw.hints.context.noerror = value != 0
		return
	case GLFW_OPENGL_PROFILE:
		_glfw.hints.context.profile = value
		return
	case GLFW_CONTEXT_RELEASE_BEHAVIOR:
		_glfw.hints.context.release = value
		return
	case GLFW_REFRESH_RATE:
		_glfw.hints.refreshRate = value
		return
	}
	slog.Error("Invalid window hint", "hint", hint, "value", value)
}

// GetClipboardString returns the contents of the system clipboard, if it
// contains or is convertible to a UTF-8 encoded string.
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
	if shape != ArrowCursor && shape != IBeamCursor && shape != CrosshairCursor &&
		shape != HandCursor && shape != HResizeCursor && shape != VResizeCursor {
		panic("Invalid standard cursor")
	}
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

func LoadCursor(cursorID uint16) HANDLE {
	h, err := LoadImage(0, uint32(cursorID), IMAGE_CURSOR, 0, 0, LR_DEFAULTSIZE|LR_SHARED)
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("LoadCursor failed, " + err.Error())
	}
	if h == 0 {
		panic("LoadCursor failed")
	}
	return HANDLE(h)
}

func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	var s *_GLFWwindow
	if share != nil {
		s = share
	}
	w, err := glfwCreateWindow(width, height, title, monitor, s)
	if err != nil {
		return nil, fmt.Errorf("glfwCreateWindow failed: %v", err)
	}
	wnd := w
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

// SetPos sets the position, in screen coordinates, of the upper-left corner of the client area of the Window.
func (w *Window) SetPos(xpos, ypos int) {
	glfwSetWindowPos(w, xpos, ypos)
}

// Init() is GLFWAPI int glfwInit(void) from init.c
func Init() error {
	var err error

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	// Repeated calls do nothing
	if _glfw.initialized {
		return nil
	}
	_glfw.hints.init = _GLFWinitconfig{}
	return glfwPlatformInit()
}

// Terminate destroys all remaining Windows, frees any allocated resources and
// sets the library to an uninitialized state.
func Terminate() {
	glfwTerminate()
}

// GetContentScale function retrieves the content scale for the specified
// Window. The content scale is the ratio between the current DPI and the
// platform's default DPI.
func (w *Window) GetContentScale() (float32, float32) {
	return glfwGetContentScale(w)
}

// GetFrameSize retrieves the size, in screen coordinates, of each edge of the frame
// of the specified Window. This size includes the title bar, if the Window has one.
func (w *Window) GetFrameSize() (left, top, right, bottom int) {
	var l, t, r, b int
	glfwGetWindowFrameSize(w, &l, &t, &r, &b)
	return int(l), int(t), int(r), int(b)
}

// GetCursorPos returns the last reported position of the cursor.
func (w *Window) GetCursorPos() (x float64, y float64) {
	var xpos, ypos int
	glfwGetCursorPos(w, &xpos, &ypos)
	return float64(xpos), float64(ypos)
}

// GetSize returns the size, in screen coordinates, of the client area of the
// specified Window.
func (w *Window) GetSize() (width int, height int) {
	var wi, h int
	glfwGetWindowSize(w, &wi, &h)
	return int(wi), int(h)
}

// Focus brings the specified Window to front and sets input focus.
func (w *Window) Focus() {
	glfwFocusWindow(w)
}

// ShouldClose reports the value of the close flag of the specified Window.
func (w *Window) ShouldClose() bool {
	return w.shouldClose
}

// Destroy destroys the specified window and its context. On calling this
// function, no further callbacks will be called for that window.
//
// This function may only be called from the main thread.
func (w *Window) Destroy() {
	// windows.remove(w.data)
	glfwDestroyWindow(w)
}

// SetSize sets the size, in screen coordinates, of the client area of the Window.
func (window *Window) SetSize(width, height int) {
	if window.monitor != nil {
		if window.monitor.window == window {
			// acquireMonitor(window)
			// fitToMonitor(window)
		}
	} else {
		glfwSetSize(window, width, height)
	}
}

// Show makes the Window visible, if it was previously hidden.
func (w *Window) Show() {
	if w.monitor != nil {
		return
	}
	glfwShowWindow(w)
	if w.focusOnShow {
		glfwFocusWindow(w)
	}
}

var previous *Window

func (w *Window) MakeContextCurrent() {
	// _GLFWWindow * Window = (_GLFWWindow *)hMonitor;
	// _GLFWWindow * previous;
	// _GLFW_REQUIRE_INIT();
	// previous := glfwPlatformGetTls(&_glfw.contextSlot);
	if previous != nil {
		previous.context.makeCurrent(nil)
	}
	previous = w
	if w == nil {
		panic("Window is nil")
	}
	w.context.makeCurrent(w)
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
