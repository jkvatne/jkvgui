package glfw

import (
	"errors"
	"fmt"
	"golang.design/x/clipboard"
	"log/slog"
	"syscall"
)

type Action int

type StandardCursor uint16

type Hint uint32

// Window represents a Window.

type Window = _GLFWwindow

type Cursor struct {
	next   *Cursor
	handle syscall.Handle
}

// Cursor modes
const (
	GLFW_CURSOR_NORMAL   = 0x00034001
	GLFW_CURSOR_CAPTURED = 0x00034004
	GLFW_CURSOR_HIDDEN   = 0x00034002
	GLFW_CURSOR_DISABLED = 0x00034003
)

// Hints
const (
	GLFW_RED_BITS                 = 0x00021001
	GLFW_GREEN_BITS               = 0x00021002
	GLFW_BLUE_BITS                = 0x00021003
	GLFW_ALPHA_BITS               = 0x00021004
	GLFW_DEPTH_BITS               = 0x00021005
	GLFW_STENCIL_BITS             = 0x00021006
	GLFW_ACCUM_RED_BITS           = 0x00021007
	GLFW_ACCUM_GREEN_BITS         = 0x00021008
	GLFW_ACCUM_BLUE_BITS          = 0x00021009
	GLFW_ACCUM_ALPHA_BITS         = 0x0002100A
	GLFW_AUX_BUFFERS              = 0x0002100B
	GLFW_SAMPLES                  = 0x0002100D
	GLFW_SRGB_CAPABLE             = 0x0002100E
	GLFW_REFRESH_RATE             = 0x0002100F
	GLFW_DOUBLEBUFFER             = 0x00021010
	GLFW_CLIENT_API               = 0x00022001
	GLFW_CONTEXT_VERSION_MAJOR    = 0x00022002
	GLFW_CONTEXT_VERSION_MINOR    = 0x00022003
	GLFW_CONTEXT_ROBUSTNESS       = 0x00022005
	GLFW_OPENGL_FORWARD_COMPAT    = 0x00022006
	GLFW_CONTEXT_DEBUG            = 0x00022007
	GLFW_OPENGL_PROFILE           = 0x00022008
	GLFW_CONTEXT_RELEASE_BEHAVIOR = 0x00022009
	GLFW_CONTEXT_NO_ERROR         = 0x0002200A
	GLFW_CONTEXT_CREATION_API     = 0x0002200B
	GLFW_SCALE_TO_MONITOR         = 0x0002200C
	GLFW_SCALE_FRAMEBUFFER        = 0x0002200D
	GLFW_WIN32_KEYBOARD_MENU      = 0x00025001
	GLFW_WIN32_SHOWDEFAULT        = 0x00025002
	GLFW_TRANSPARENT_FRAMEBUFFER  = 0x0002000A
	GLFW_RESIZABLE                = 0x00020003
	GLFW_DECORATED                = 0x00020005
	GLFW_AUTO_ICONIFY             = 0x00020006
	GLFW_FLOATING                 = 0x00020007
	GLFW_MAXIMIZED                = 0x00020008
	GLFW_POSITION_X               = 0x0002000E
	GLFW_POSITION_Y               = 0x0002000F
	GLFW_FOCUSED                  = 0x00020001
	GLFW_VISIBLE                  = 0x00020004
)

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
		// _glfw.hints.window.scaleFramebuffer = value != 0
		return
		/*
			case GLFW_CENTER_CURSOR:
				_glfw.hints.window.centerCursor = value!=0
				return;
			case GLFW_FOCUS_ON_SHOW:
				_glfw.hints.window.focusOnShow = value!=0
				return;
			case GLFW_MOUSE_PASSTHROUGH:
				_glfw.hints.window.mousePassthrough = value!=0
				return;*/
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
	// return fmt.Errorf("Invalid window hint");
	slog.Error("Invalid window hint")
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
	var cursor = Cursor{}
	if shape != ArrowCursor && shape != IbeamCursor && shape != CrosshairCursor &&
		shape != HandCursor && shape != HResizeCursor && shape != VResizeCursor {
		panic("Invalid standard cursor")
	}
	cursor.next = _glfw.cursorListHead
	_glfw.cursorListHead = &cursor
	glfwCreateStandardCursorWin32(&cursor, shape)
	return &cursor
}

func LoadCursor(cursorID uint16) syscall.Handle {
	h, err := LoadImage(0, uint32(cursorID), IMAGE_CURSOR, 0, 0, LR_DEFAULTSIZE|LR_SHARED)
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("LoadCursor failed, " + err.Error())
	}
	if h == 0 {
		panic("LoadCursor failed")
	}
	return syscall.Handle(h)
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
	if c == nil {
		glfwSetCursor(w, nil)
	} else {
		glfwSetCursor(w, c)
	}
}

// SetPos sets the position, in screen coordinates, of the upper-left corner of the client area of the Window.
func (w *Window) SetPos(xpos, ypos int) {
	glfwSetWindowPos(w, xpos, ypos)
}

// CursorPosCallback the cursor position callback.
type CursorPosCallback func(w *Window, xpos float64, ypos float64)

// SetCursorPosCallback sets the cursor position callback which is called
// when the cursor is moved. The callback is provided with the position relative
// to the upper-left corner of the client area of the Window.
func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	w.cursorPosCallback = cbfun
	return nil
}

// KeyCallback is the key callback.
type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)

// SetKeyCallback sets the key callback which is called when a key is pressed, repeated or released.
func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	w.keyCallback = cbfun
	return nil
}

// CharCallback is the character callback.
type CharCallback func(w *Window, char rune)

// SetCharCallback sets the character callback which is called when a Unicode character is input.
func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	w.charCallback = cbfun
	return nil
}

// DropCallback is the drop callback which is called when an object is dropped over the Window.
type DropCallback func(w *Window, names []string)

// SetDropCallback sets the drop callback
func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	w.dropCallback = cbfun
	return nil
}

// ContentScaleCallback is the function signature for Window content scale
// callback functions.
type ContentScaleCallback func(w *Window, x float32, y float32)

// SetContentScaleCallback function sets the Window content scale callback of
// the specified Window, which is called when the content scale of the specified Window changes.
func (w *Window) SetContentScaleCallback(cbfun ContentScaleCallback) ContentScaleCallback {
	w.contentScaleCallback = cbfun
	return nil
}

// RefreshCallback is the Window refresh callback.
type RefreshCallback func(w *Window)

// SetRefreshCallback sets the refresh callback of the Window, which
// is called when the client area of the Window needs to be redrawn,
func (w *Window) SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback) {
	w.refreshCallback = cbfun
	return nil
}

// FocusCallback is the Window focus callback.
type FocusCallback func(w *Window, focused bool)

// SetFocusCallback sets the focus callback of the Window, which is called when
// the Window gains or loses focus.
//
// After the focus callback is called for a Window that lost focus, synthetic key
// and mouse button release events will be generated for all such that had been
// pressed. For more information, see SetKeyCallback and SetMouseButtonCallback.
func (w *Window) SetFocusCallback(cbfun FocusCallback) (previous FocusCallback) {
	w.focusCallback = cbfun
	return nil
}

// SizeCallback is the Window size callback.
type SizeCallback func(w *Window, width int, height int)

// SetSizeCallback sets the size callback of the Window, which is called when
// the Window is resized. The callback is provided with the size, in screen
// coordinates, of the client area of the Window.
func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	w.sizeCallback = cbfun
	return nil
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

func (w *Window) MakeContextCurrent() {
	// _GLFWWindow * Window = (_GLFWWindow *)hMonitor;
	// _GLFWWindow * previous;
	// _GLFW_REQUIRE_INIT();
	// previous := glfwPlatformGetTls(&_glfw.contextSlot);
	if w == nil {
		panic("Window is nil")
	}
	w.context.makeCurrent(w)
	w.Focus()
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
