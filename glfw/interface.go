package glfw

import "C"
import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"log/slog"
	"syscall"
	"unsafe"
)

const (
	GL_VERSION                 = 0x1F02
	GLFW_OPENGL_ANY_PROFILE    = 0
	GLFW_OPENGL_CORE_PROFILE   = 0x00032001
	GLFW_OPENGL_COMPAT_PROFILE = 0x00032002
	_GLFW_WNDCLASSNAME         = "GLFW30"
	GLFW_DONT_CARE             = -1
	OpenGLProfile              = 0x00022008
	OpenGLCoreProfile          = 0x00032001
	OpenGLForwardCompatible    = 0x00032002
	True                       = 1
	False                      = 0
	Resizable                  = 0x00020003
	Focused                    = 0x00020001
	Iconified                  = 0x00020002
	Resizeable                 = 0x00020003
	Visible                    = 0x00020004
	Decorated                  = 0x00020005
	AutoIconify                = 0x00020006
	Floating                   = 0x00020007
	Maximized                  = 0x00020008
	ContextVersionMajor        = 0x00022002
	ContextVersionMinor        = 0x00022003
	Samples                    = 0x0002100D
	ArrowCursor                = 0x00036001
	IbeamCursor                = 0x00036002
	CrosshairCursor            = 0x00036003
	HandCursor                 = 0x00036004
	HResizeCursor              = 0x00036005
	VResizeCursor              = 0x00036006
	LR_CREATEDIBSECTION        = 0x00002000
	LR_DEFAULTCOLOR            = 0x00000000
	LR_DEFAULTSIZE             = 0x00000040
	LR_LOADFROMFILE            = 0x00000010
	LR_LOADMAP3DCOLORS         = 0x00001000
	LR_LOADTRANSPARENT         = 0x00000020
	LR_MONOCHROME              = 0x00000001
	LR_SHARED                  = 0x00008000
	LR_VGACOLOR                = 0x00000080
	IMAGE_ICON                 = 1
	CS_HREDRAW                 = 0x0002
	CS_INSERTCHAR              = 0x2000
	CS_NOMOVECARET             = 0x4000
	CS_VREDRAW                 = 0x0001
	CS_OWNDC                   = 0x0020
	KF_EXTENDED                = 0x100
	GLFW_RELEASE               = 0
	GLFW_PRESS                 = 1
	GLFW_REPEAT                = 2
	GLFW_CURSOR_NORMAL         = 0x00034001
	GLFW_CURSOR_HIDDEN         = 0x00034002
	GLFW_CURSOR_DISABLED       = 0x00034003
	GLFW_OPENGL_API            = 0x00030001
	GLFW_NATIVE_CONTEXT_API    = 0x00036001
	GLFW_OPENGL_ES_API         = 0x00030002
	GLFW_EGL_CONTEXT_API       = 0x00036002
	GLFW_OSMESA_CONTEXT_API    = 0x00036003
	GLFW_NO_API                = 0
)

type Action int

type StandardCursor uint32

type Hint uint32

//
type GLFWvidmode struct {
	width       int
	height      int
	redBits     int
	greenBits   int
	blueBits    int
	refreshRate int
}

// Window represents a Window.
type Window struct {
	Data                 *_GLFWwindow
	charCallback         CharCallback
	focusCallback        FocusCallback
	keyCallback          KeyCallback
	mouseButtonCallback  MouseButtonCallback
	cursorPosCallback    CursorPosCallback
	scrollCallback       ScrollCallback
	refreshCallback      RefreshCallback
	sizeCallback         SizeCallback
	dropCallback         DropCallback
	contentScaleCallback func(w *Window, x float32, y float32)

	fPosHolder             func(w *Window, xpos int, ypos int)
	fSizeHolder            func(w *Window, width int, height int)
	fFramebufferSizeHolder func(w *Window, width int, height int)
	fCloseHolder           func(w *Window)
	fMaximizeHolder        func(w *Window, maximized bool)
	fContentScaleHolder    func(w *Window, x float32, y float32)
	fRefreshHolder         func(w *Window)
	fFocusHolder           func(w *Window, focused bool)
	fIconifyHolder         func(w *Window, iconified bool)
	fMouseButtonHolder     func(w *Window, button MouseButton, action Action, mod ModifierKey)
	fCursorPosHolder       func(w *Window, xpos float64, ypos float64)
	fCursorEnterHolder     func(w *Window, entered bool)
	fScrollHolder          func(w *Window, xoff float64, yoff float64)
	fKeyHolder             func(w *Window, key Key, scancode int, action Action, mods ModifierKey)
	fCharHolder            func(w *Window, char rune)
	fCharModsHolder        func(w *Window, char rune, mods ModifierKey)
	fDropHolder            func(w *Window, names []string)
}

type Cursor struct {
	data *_GLFWcursor
}

// iconID is the ID of the icon in the resource file.
const iconID = 1

var resources struct {
	handle syscall.Handle
	class  uint16
	cursor syscall.Handle
}

func panicError() {
	// err := acceptError()
	// if err != nil {
	//	panic(err)
	// }
}

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
	GLFW_STEREO                   = 0x0002100C
	GLFW_SAMPLES                  = 0x0002100D
	GLFW_SRGB_CAPABLE             = 0x0002100E
	GLFW_REFRESH_RATE             = 0x0002100F
	GLFW_DOUBLEBUFFER             = 0x00021010
	GLFW_CLIENT_API               = 0x00022001
	GLFW_CONTEXT_VERSION_MAJOR    = 0x00022002
	GLFW_CONTEXT_VERSION_MINOR    = 0x00022003
	GLFW_CONTEXT_REVISION         = 0x00022004
	GLFW_CONTEXT_ROBUSTNESS       = 0x00022005
	GLFW_OPENGL_FORWARD_COMPAT    = 0x00022006
	GLFW_CONTEXT_DEBUG            = 0x00022007
	GLFW_OPENGL_DEBUG_CONTEXT     = GLFW_CONTEXT_DEBUG
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

func glfwWindowHint(hint int, value int) {
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
	case GLFW_STEREO:
		_glfw.hints.framebuffer.stereo = value != 0
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
	case GLFW_WIN32_KEYBOARD_MENU:
		// _glfw.hints.window.win32.keymenu = value != 0
		return
	case GLFW_WIN32_SHOWDEFAULT:
		// _glfw.hints.window.win32.showDefault = value != 0
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

func WindowHint(target int, hint int) {
	glfwWindowHint(target, hint)
	panicError()
}

// GetClipboardString returns the contents of the system clipboard, if it
// contains or is convertible to a UTF-8 encoded string.
// This function may only be called from the main thread.
func GetClipboardString() string {
	// return C.glfwGetClipboardString(nil)
	return ""
}

// SetClipboardString sets the system clipboard to the specified UTF-8 encoded string.
// This function may only be called from the main thread.
func SetClipboardString(str string) {
	// cp := C.CString(str)
	// defer C.free(unsafe.Pointer(cp))
	// C.glfwSetClipboardString(nil, cp)
	panicError()
}

// CreateStandardCursor returns a cursor with a standard shape,
// that can be set for a Window with SetCursor.
func CreateStandardCursor(shape StandardCursor) *Cursor {
	// c := C.glfwCreateStandardCursor(C.int(shape))
	panicError()
	return nil // &Cursor{c}
}

func SetProcessDPIAware() {
	_, _, _ = _SetProcessDPIAware.Call()
}

func LoadCursor(curID uint16) (syscall.Handle, error) {
	h, _, err := _LoadCursor.Call(0, uintptr(curID))
	if h == 0 {
		return 0, fmt.Errorf("LoadCursorW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	var s *_GLFWwindow
	if share != nil {
		s = share.Data
	}
	w, err := glfwCreateWindow(width, height, title, monitor, s)
	if err != nil {
		return nil, fmt.Errorf("glfwCreateWindow failed: %v", err)
	}
	wnd := &Window{Data: w}
	windowMap.put(wnd)
	return wnd, nil
}

func glfwIsValidContextConfig(ctxconfig *_GLFWctxconfig) error {
	if ctxconfig.source != GLFW_NATIVE_CONTEXT_API && ctxconfig.source != GLFW_EGL_CONTEXT_API && ctxconfig.source != GLFW_OSMESA_CONTEXT_API {
		return fmt.Errorf("Invalid context creation API")
	}
	if ctxconfig.client != GLFW_NO_API && ctxconfig.client != GLFW_OPENGL_API && ctxconfig.client != GLFW_OPENGL_ES_API {
		return fmt.Errorf("Invalid client API")
	}
	if ctxconfig.share != nil {
		if ctxconfig.client == GLFW_NO_API || ctxconfig.share.context.client == GLFW_NO_API {
			return fmt.Errorf("Invalid share")
		}
		if ctxconfig.source != ctxconfig.share.context.source {
			return fmt.Errorf("Invalid share")
		}
	}

	if ctxconfig.client == GLFW_OPENGL_API {
		if (ctxconfig.major < 1 || ctxconfig.minor < 0) ||
			(ctxconfig.major == 1 && ctxconfig.minor > 5) ||
			(ctxconfig.major == 2 && ctxconfig.minor > 1) ||
			(ctxconfig.major == 3 && ctxconfig.minor > 3) {
			// OpenGL 1.0 is the smallest valid version
			// OpenGL 1.x series ended with version 1.5
			// OpenGL 2.x series ended with version 2.1
			// OpenGL 3.x series ended with version 3.3
			// For now, let everything else through
			return fmt.Errorf("Invalid OpenGL version %i.%i", ctxconfig.major, ctxconfig.minor)
		}

		if ctxconfig.profile != 0 {
			if ctxconfig.profile != GLFW_OPENGL_CORE_PROFILE && ctxconfig.profile != GLFW_OPENGL_COMPAT_PROFILE {
				return fmt.Errorf("Invalid OpenGL profile 0x%08X", ctxconfig.profile)
			}
			if ctxconfig.major <= 2 || (ctxconfig.major == 3 && ctxconfig.minor < 2) {
				// Desktop OpenGL context profiles are only defined for version 3.2 and above
				return fmt.Errorf("Context profiles are only defined for OpenGL version 3.2 and above")
			}
		}
		if ctxconfig.forward && ctxconfig.major <= 2 {
			// Forward-compatible contexts are only defined for OpenGL version 3.0 and above
			return fmt.Errorf("Forward-compatibility is only defined for OpenGL version 3.0 and above")
		}
	} else if ctxconfig.client == GLFW_OPENGL_ES_API {
		if ctxconfig.major < 1 || ctxconfig.minor < 0 || (ctxconfig.major == 1 && ctxconfig.minor > 1) || (ctxconfig.major == 2 && ctxconfig.minor > 0) {
			// OpenGL ES 1.0 is the smallest valid version
			// OpenGL ES 1.x series ended with version 1.1
			// OpenGL ES 2.x series ended with version 2.0
			// For now, let everything else through
			return fmt.Errorf("Invalid OpenGL ES version %i.%i", ctxconfig.major, ctxconfig.minor)
		}
	}
	// if ctxconfig.robustness > 0 && ctxconfig.robustness != GLFW_NO_RESET_NOTIFICATION && ctxconfig.robustness != GLFW_LOSE_CONTEXT_ON_RESET {
	//	return fmt.Errorf("Invalid context robustness mode 0x%08X", ctxconfig.robustness)
	// }

	// if ctxconfig.release > 0 && ctxconfig.release != GLFW_RELEASE_BEHAVIOR_NONE && ctxconfig.release != GLFW_RELEASE_BEHAVIOR_FLUSH {
	//	return fmt.Errorf("Invalid context release behavior 0x%08X", ctxconfig.release)
	// }
	return nil
}

func createNativeWindow(window *_GLFWwindow, wndconfig *_GLFWwndconfig, fbconfig *_GLFWfbconfig) error {
	var err error
	var frameX, frameY, frameWidth, frameHeight int32
	SetProcessDPIAware()
	// style := getWindowStyle(window)
	// exStyle := getWindowExStyle(window)
	if window.monitor != nil {
		var mi MONITORINFO
		mi.CbSize = uint32(unsafe.Sizeof(mi))
		_, _, err := _GetMonitorInfo.Call(uintptr(window.monitor.hMonitor), uintptr(unsafe.Pointer(&mi)))
		if errors.Is(err, syscall.Errno(0)) {
			return err
		}
		// NOTE: This window placement is temporary and approximate, as the
		//       correct position and size cannot be known until the monitor
		//       video mode has been picked in _glfwSetVideoModeWin32
		frameX = mi.RcMonitor.Left
		frameY = mi.RcMonitor.Top
		frameWidth = mi.RcMonitor.Right - mi.RcMonitor.Left
		frameHeight = mi.RcMonitor.Bottom - mi.RcMonitor.Top
	} else {
		rect := RECT{0, 0, int32(wndconfig.width), int32(wndconfig.height)}
		window.Win32.maximized = wndconfig.maximized
		if wndconfig.maximized {
			// style |= WS_MAXIMIZE
		}
		// TODO AdjustWindowRectEx(&rect, style, FALSE, exStyle);
		frameX = 100 // CW_USEDEFAULT
		frameY = 100 //  CW_USEDEFAULT
		frameWidth = rect.Right - rect.Left
		frameHeight = rect.Bottom - rect.Top
	}

	window.Win32.handle, err = CreateWindowEx(
		WS_OVERLAPPED|WS_EX_APPWINDOW,
		_glfw.class,
		wndconfig.title,
		WS_OVERLAPPED|WS_CLIPSIBLINGS|WS_CLIPCHILDREN,
		frameX, frameY, // Window position
		frameWidth, frameHeight, // Window width/heigth
		0, // No parent
		0, // No menu
		resources.handle,
		0)
	return err
}

type PIXELFORMATDESCRIPTOR = struct {
	nSize           uint16
	nVersion        uint16
	dwFlags         uint32
	iPixelType      uint8
	cColorBits      uint8
	cRedBits        uint8
	cRedShift       uint8
	cGreenBits      uint8
	cGreenShift     uint8
	cBlueBits       uint8
	cBlueShift      uint8
	cAlphaBits      uint8
	cAlphaShift     uint8
	cAccumBits      uint8
	cAccumRedBits   uint8
	cAccumGreenBits uint8
	cAccumBlueBits  uint8
	cAccumAlphaBits uint8
	cDepthBits      uint8
	cStencilBits    uint8
	cAuxBuffers     uint8
	iLayerType      uint8
	bReserved       uint8
	dwLayerMask     uint32
	dwVisibleMask   uint32
	dwDamageMask    uint32
}

var (
	gdi32 = windows.NewLazySystemDLL("gdi32.dll")
)

const (
	PFD_DRAW_TO_WINDOW = 0x04
	PFD_SUPPORT_OPENGL = 0x20
	PFD_DOUBLEBUFFER   = 0x01
	PFD_TYPE_RGBA      = 0x00
)

func glfwPlatformCreateWindow(window *_GLFWwindow, wndconfig *_GLFWwndconfig, ctxconfig *_GLFWctxconfig, fbconfig *_GLFWfbconfig) error {
	err := createNativeWindow(window, wndconfig, fbconfig)
	if err != nil {
		return err
	}
	if ctxconfig.client != GLFW_NO_API {
		if ctxconfig.source == GLFW_NATIVE_CONTEXT_API {
			if err := _glfwInitWGL(); err != nil {
				return fmt.Errorf("_glfwInitWGL error " + err.Error())
			}
			if err := _glfwCreateContextWGL(window, ctxconfig, fbconfig); err != nil {
				return fmt.Errorf("_glfwCreateContextWGL error " + err.Error())
			}
		} else if ctxconfig.source == GLFW_EGL_CONTEXT_API {
			if err := glfwInitEGL(); err != nil {
				return err
			}
			if err := _glfwCreateContextEGL(window, ctxconfig, fbconfig); err != nil {
				return err
			}
		} else if ctxconfig.source == GLFW_OSMESA_CONTEXT_API {
			if err := glfwInitOSMesa(); err != nil {
				return err
			}
			if err := _glfwCreateContextOSMesa(window, ctxconfig, fbconfig); err != nil {
				return err
			}
		}
		if err := _glfwRefreshContextAttribs(window, ctxconfig); err != nil {
			return err
		}
	}
	if window.monitor != nil {
		_glfwPlatformShowWindow(window)
		// _glfwPlatformFocusWindow(window)
		// acquireMonitor(window)
		// fitToMonitor(window)
		if wndconfig.centerCursor {
			// _glfwCenterCursorInContentArea(window)
		}
	} else if wndconfig.visible {
		_glfwPlatformShowWindow(window)
		if wndconfig.focused {
			// _glfwPlatformFocusWindow(window)
		}
	}
	return nil
}

func glfwCreateWindow(width, height int, title string, monitor *Monitor, share *_GLFWwindow) (*_GLFWwindow, error) {

	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid width/heigth")
	}
	// End of _glfwPlatformCreateWindow
	fbconfig := _glfw.hints.framebuffer
	ctxconfig := _glfw.hints.context
	wndconfig := _glfw.hints.window
	wndconfig.width = width
	wndconfig.height = height

	wndconfig.title = title
	ctxconfig.share = share
	if glfwIsValidContextConfig(&ctxconfig) != nil {
		return nil, fmt.Errorf("glfw context config is invalid: %v", ctxconfig)
	}

	window := &_GLFWwindow{}
	window.next = _glfw.windowListHead
	_glfw.windowListHead = window

	window.videoMode.width = width
	window.videoMode.height = height
	window.videoMode.redBits = fbconfig.redBits
	window.videoMode.greenBits = fbconfig.greenBits
	window.videoMode.blueBits = fbconfig.blueBits
	window.videoMode.refreshRate = _glfw.hints.refreshRate

	window.monitor = monitor
	window.resizable = wndconfig.resizable
	window.decorated = wndconfig.decorated
	window.autoIconify = wndconfig.autoIconify
	window.floating = wndconfig.floating
	window.focusOnShow = wndconfig.focusOnShow
	window.cursorMode = GLFW_CURSOR_NORMAL
	window.doublebuffer = fbconfig.doublebuffer
	window.minwidth = GLFW_DONT_CARE
	window.minheight = GLFW_DONT_CARE
	window.maxwidth = GLFW_DONT_CARE
	window.maxheight = GLFW_DONT_CARE
	window.numer = GLFW_DONT_CARE
	window.denom = GLFW_DONT_CARE

	if err := glfwPlatformCreateWindow(window, &wndconfig, &ctxconfig, &fbconfig); err != nil {
		// glfwDestroyWindow(window)
		return nil, fmt.Errorf("Error creating window, " + err.Error())
	}
	return window, nil
}

// SwapBuffers swaps the front and back buffers of the Window.
func (w *Window) SwapBuffers() {
	glfwSwapBuffers(w.Data)
	panicError()
}

// SetCursor sets the cursor image to be used when the cursor is over the client area
func (w *Window) SetCursor(c *Cursor) {
	if c == nil {
		glfwSetCursor(w.Data, nil)
	} else {
		// TODO glfwSetCursor(w.Data, c)
	}
	panicError()
}

func glfwSetWindowPos(w *_GLFWwindow, xpos, ypos int) {

}

// SetPos sets the position, in screen coordinates, of the upper-left corner of the client area of the Window.
func (w *Window) SetPos(xpos, ypos int) {
	glfwSetWindowPos(w.Data, xpos, ypos)
	panicError()
}

func SetWindowPos(hWnd HANDLE, after HANDLE, x, y, cx, cy, flags int) {

}

const (
	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
	SWP_NOREDRAW       = 0x0008
	SWP_NOACTIVATE     = 0x0010
	SWP_FRAMECHANGED   = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOCOPYBITS     = 0x0100
	SWP_NOOWNERZORDER  = 0x0200
	SWP_NOSENDCHANGING = 0x0400
)

func glfwSetWindowSize(window *_GLFWwindow, width, height int) {
	if window.monitor != nil {
		if window.monitor.window == window {
			// acquireMonitor(window)
			// fitToMonitor(window)
		}
	} else {
		if true { // (_glfwIsWindows10Version1607OrGreaterWin32()) {
			// AdjustWindowRectExForDpi(&rect, getWindowStyle(window),	FALSE, getWindowExStyle(window), GetDpiForWindow(window.win32.handle));
		} else {
			// AdjustWindowRectEx(&rect, getWindowStyle(window), FALSE, getWindowExStyle(window));
		}
		SetWindowPos(window.Win32.handle, 0, 0, 0, width, height, SWP_NOACTIVATE|SWP_NOOWNERZORDER|SWP_NOMOVE|SWP_NOZORDER)
	}
}

// SetSize sets the size, in screen coordinates, of the client area of the Window.
func (w *Window) SetSize(width, height int) {
	glfwSetWindowSize(w.Data, width, height)
	panicError()
}

func glfwShowWindow(w *_GLFWwindow) {
	if w.monitor != nil {
		return
	}
	_ = _glfwPlatformShowWindow(w)
	if w.focusOnShow {
		// TODO _glfwPlatformFocusWindow(window)
	}
}

// Show makes the Window visible, if it was previously hidden.
func (w *Window) Show() {
	glfwShowWindow(w.Data)
	panicError()
}

func (w *Window) MakeContextCurrent() {
	// _GLFWWindow * Window = (_GLFWWindow *)handle;
	// _GLFWWindow * previous;
	// _GLFW_REQUIRE_INIT();
	// previous := glfwPlatformGetTls(&_glfw.contextSlot);
	if w == nil {
		panic("Window is nil")
	}
	if w != nil && w.Data.context.client == 0 {
		panic("Cannot make current with a Window that has no OpenGL or OpenGL ES context")
	}
	w.Data.context.makeCurrent(w.Data)
	panicError()
}

// Focus brings the specified Window to front and sets input focus.
func (w *Window) Focus() {
	// TODO glfwFocusWindow(w.Data)
}

// ShouldClose reports the value of the close flag of the specified Window.
func (w *Window) ShouldClose() bool {
	return w.Data.shouldClose
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

// PollEvents processes only those events that have already been received and
// then returns immediately. Processing events will cause the Window and input
// callbacks associated with those events to be called.
func PollEvents() {
	glfwPollEvents()
	panicError()
}

func _glfwRegisterWindowClassWin32() error {
	/*var wc WNDCLASSEXW
	wc.cbSize        = sizeof(wc);
	wc.style         = CS_HREDRAW | CS_VREDRAW | CS_OWNDC;
	wc.lpfnWndProc   = windowProc;
	wc.hInstance     = _glfw.Win32.instance;
	wc.hCursor       = LoadCursorW(NULL, IDC_ARROW);
	wc.lpszClassName = _GLFW_WNDCLASSNAME;
	// Load user-provided icon if available
	//wc.hIcon = LoadImageW(GetModuleHandleW(NULL),"GLFW_ICON", IMAGE_ICON, 0, 0, LR_DEFAULTSIZE | LR_SHARED);
	//if (!wc.hIcon) {
		// No user-provided icon found, load default icon
		//wc.hIcon = LoadImageW(NULL,	IDI_APPLICATION, IMAGE_ICON, 0, 0, LR_DEFAULTSIZE | LR_SHARED);
	//}*/
	icon := syscall.Handle(0)
	wcls := WndClassEx{
		CbSize:        uint32(unsafe.Sizeof(WndClassEx{})),
		Style:         CS_HREDRAW | CS_VREDRAW | CS_OWNDC,
		LpfnWndProc:   syscall.NewCallback(windowProc),
		HInstance:     _glfw.instance,
		HIcon:         icon,
		LpszClassName: syscall.StringToUTF16Ptr("GLFW"),
	}
	var err error
	_glfw.class, err = RegisterClassEx(&wcls)
	return err
}

// Flags used for GetModuleHandleEx
const (
	GET_MODULE_HANDLE_EX_FLAG_PIN                = 1
	GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT = 2
	GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS       = 4
)

func glfwDefaultWindowHints() {
	_glfw.hints.context.client = GLFW_OPENGL_API
	_glfw.hints.context.source = GLFW_NATIVE_CONTEXT_API
	_glfw.hints.context.major = 1
	_glfw.hints.context.minor = 0
	// The default is a focused, visible, resizable window with decorations
	_glfw.hints.window.resizable = true
	_glfw.hints.window.visible = true
	_glfw.hints.window.decorated = true
	_glfw.hints.window.focused = true
	_glfw.hints.window.autoIconify = true
	_glfw.hints.window.centerCursor = true
	_glfw.hints.window.focusOnShow = true
	// The default is 24 bits of color, 24 bits of depth and 8 bits of stencil, double buffered
	_glfw.hints.framebuffer.redBits = 8
	_glfw.hints.framebuffer.greenBits = 8
	_glfw.hints.framebuffer.blueBits = 8
	_glfw.hints.framebuffer.alphaBits = 8
	_glfw.hints.framebuffer.depthBits = 24
	_glfw.hints.framebuffer.stencilBits = 8
	_glfw.hints.framebuffer.doublebuffer = true
	// The default is to select the highest available refresh rate
	_glfw.hints.refreshRate = GLFW_DONT_CARE
	// The default is to use full Retina resolution framebuffers
	_glfw.hints.window.ns.retina = true
}

func helperWindowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	/*	switch msg	{
		//case WM_DISPLAYCHANGE:
		//    _glfwPollMonitorsWin32();

		case WM_DEVICECHANGE:
		if (!_glfw.joysticksInitialized)
				break;

		if (wParam == DBT_DEVICEARRIVAL)
		{
		DEV_BROADCAST_HDR* dbh = (DEV_BROADCAST_HDR*) lParam;
		if (dbh && dbh->dbch_devicetype == DBT_DEVTYP_DEVICEINTERFACE)
		_glfwDetectJoystickConnectionWin32();
		}
		else if (wParam == DBT_DEVICEREMOVECOMPLETE)
		{
		DEV_BROADCAST_HDR* dbh = (DEV_BROADCAST_HDR*) lParam;
		if (dbh && dbh->dbch_devicetype == DBT_DEVTYP_DEVICEINTERFACE)
		_glfwDetectJoystickDisconnectionWin32();
		}

		break;
		}
		}
	*/
	r1, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return r1
}

func _glfwPlatformShowWindow(w *_GLFWwindow) error {
	_, _, err := _ShowWindow.Call(uintptr(w.Win32.handle), windows.SW_NORMAL)
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		return err
	}
	return nil
}

func createHelperWindow() error {
	var err error
	var wc WndClassEx
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.Style = CS_OWNDC
	wc.LpfnWndProc = syscall.NewCallback(helperWindowProc)
	wc.HInstance = _glfw.instance
	wc.LpszClassName = syscall.StringToUTF16Ptr("GLFW3 Helper")

	_glfw.win32.helperWindowClass, err = RegisterClassEx(&wc)
	if _glfw.win32.helperWindowClass == 0 || err != nil {
		panic("Win32: Failed to register helper window class")
	}
	_glfw.win32.helperWindowHandle, err =
		CreateWindowEx(WS_OVERLAPPED,
			_glfw.win32.helperWindowClass,
			"Helper window",
			WS_OVERLAPPED|WS_CLIPSIBLINGS|WS_CLIPCHILDREN,
			0, 0, 500, 500,
			0, 0,
			resources.handle,
			0)

	if _glfw.win32.helperWindowHandle == 0 || err != nil {
		panic("Win32: Failed to create helper window")
	}

	// HACK: The command to the first ShowWindow call is ignored if the parent
	//       process passed along a STARTUPINFO, so clear that with a no-op call
	_, _, err = _ShowWindow.Call(uintptr(_glfw.win32.helperWindowHandle), windows.SW_NORMAL) // OBS, should be SW_HIDE
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		return err
	}

	// TODO Register for HID device notifications
	/*
		{
			dbi DEV_BROADCAST_DEVICEINTERFACE_W
			ZeroMemory(&dbi, sizeof(dbi));
			dbi.dbcc_size = sizeof(dbi);
			dbi.dbcc_devicetype = DBT_DEVTYP_DEVICEINTERFACE;
			dbi.dbcc_classguid = GUID_DEVINTERFACE_HID;

			_glfw.win32.deviceNotificationHandle =
				RegisterDeviceNotificationW(_glfw.win32.helperWindowHandle,
					(DEV_BROADCAST_HDR*) &dbi,
					DEVICE_NOTIFY_WINDOW_HANDLE);
		}

		while (PeekMessageW(&msg, _glfw.win32.helperWindowHandle, 0, 0, PM_REMOVE))
		{
			TranslateMessage(&msg);
			DispatchMessageW(&msg);
		}
	*/
	return nil
}

func _glfwInitWin32() error {
	/*err := loadLibraries()
	if err != nil {
		return err
	}*/
	// createKeyTables();
	// _glfwUpdateKeyNamesWin32();
	/*
		if (_glfwIsWindows10Version1703OrGreaterWin32()) {
			SetProcessDpiAwarenessContext(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2);
		} else if (IsWindows8Point1OrGreater()) {
			SetProcessDpiAwareness(PROCESS_PER_MONITOR_DPI_AWARE);
		} else if (IsWindowsVistaOrGreater()) {
			SetProcessDPIAware();
		} */
	SetProcessDPIAware()
	err := createHelperWindow()
	if err != nil {
		return err
	}
	// _glfwPollMonitorsWin32();
	return nil
}

func _glfwConnectWin32() {

}

// Init() is GLFWAPI int glfwInit(void) from init.c
func Init() error {
	var err error
	// Repeated calls do nothing
	if _glfw.initialized {
		return nil
	}
	_glfw.hints.init = _GLFWinitconfig{}

	// This is _glfwPlatformInit()/glfwInitWIn32()
	// TODO createKeyTables()
	// TODO _glfwUpdateKeyNamesWin32()
	_glfwConnectWin32()
	/*
		// Set dpi aware
		if(_glfwIsWindows10CreatorsUpdateOrGreaterWin32() {
			SetProcessDpiAwarenessContext(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2);
		} else if IsWindows8Point1OrGreater() {
			SetProcessDpiAwareness(PROCESS_PER_MONITOR_DPI_AWARE);
		} else if(IsWindowsVistaOrGreater() {
			SetProcessDPIAware()
		}
	*/
	SetProcessDPIAware()
	if err := _glfwRegisterWindowClassWin32(); err != nil {
		return fmt.Errorf("glfw platform init failed, _glfwRegisterWindowClassWin32 failed, %v ", err.Error())
	}
	// _, _, err := _procGetModuleHandleExW.Call(GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS|GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT, uintptr(unsafe.Pointer(&_glfw)), uintptr(unsafe.Pointer(&_glfw.instance)))
	_glfw.instance, err = GetModuleHandle()
	if err != nil {
		return fmt.Errorf("glfw platform init failed %v ", err.Error())
	}

	err = createHelperWindow()
	if err != nil {
		return err
	}
	// _glfwInitTimerWin32();
	// _glfwInitJoysticksWin32();
	// _glfwPollMonitorsWin32();
	// End of _glfwPlatformInit():

	// _glfwPlatformSetTls(&_glfw.errorSlot, &_glfwMainThreadError)
	// _glfwInitGamepadMappings()
	// _glfw.timer.offset = _glfwPlatformGetTimerValue()
	glfwDefaultWindowHints()
	_glfw.initialized = true
	return nil
}

// Terminate destroys all remaining Windows, frees any allocated resources and
// sets the library to an uninitialized state.
func Terminate() {
	/*
		if (_glfw.Win32.deviceNotificationHandle) {
			UnregisterDeviceNotification(_glfw.Win32.deviceNotificationHandle);
		}

		if (_glfw.Win32.helperWindowHandle)  {
			DestroyWindow(_glfw.Win32.helperWindowHandle);
		}
		_glfwUnregisterWindowClassWin32();
		// Restore previous foreground lock timeout system setting
		SystemParametersInfoW(SPI_SETFOREGROUNDLOCKTIMEOUT, 0, IntToPtr(_glfw.Win32.foregroundLockTimeout),	SPIF_SENDCHANGE);
		free(_glfw.Win32.clipboardString);
		free(_glfw.Win32.rawInput);
		_glfwTerminateWGL();
		_glfwTerminateEGL();
		_glfwTerminateOSMesa();
		_glfwTerminateJoysticksWin32();
		freeLibraries();
	*/
}

// GetContentScale function retrieves the content scale for the specified
// Window. The content scale is the ratio between the current DPI and the
// platform's default DPI. If you scale all pixel dimensions by this scale then
// your content should appear at an appropriate size. This is especially
// important for text and any UI elements.
//
// This function may only be called from the main thread.
func (w *Window) GetContentScale() (float32, float32) {
	// TODO
	// var x, y float32
	// C.glfwGetWindowContentScale(w.Data, &x, &y)
	// _glfwGetWindowContentScaleWin32
	// return float32(x), float32(y)
	return 1.5, 1.5
}

// GetFrameSize retrieves the size, in screen coordinates, of each edge of the frame
// of the specified Window. This size includes the title bar, if the Window has one.
// The size of the frame may vary depending on the Window-related hints used to create it.
//
// Because this function retrieves the size of each Window frame edge and not the offset
// along a particular coordinate axis, the retrieved values will always be zero or positive.
func (w *Window) GetFrameSize() (left, top, right, bottom int) {
	var l, t, r, b int
	glfwGetWindowFrameSizeWin32(w.Data, &l, &t, &r, &b)
	panicError()
	return int(l), int(t), int(r), int(b)
}

func glfwGetWindowFrameSizeWin32(window *_GLFWwindow, left, top, right, bottom *int) {
	var rect RECT
	var width, height int
	_glfwGetWindowSizeWin32(window, &width, &height)
	rect.Right = int32(width)
	rect.Bottom = int32(height)
	/*	if (_glfwIsWindows10Version1607OrGreaterWin32()) {
			AdjustWindowRectExForDpi(&rect, getWindowStyle(window),	FALSE, getWindowExStyle(window),GetDpiForWindow(window->win32.handle));
		} else {
			AdjustWindowRectEx(&rect, getWindowStyle(window),FALSE, getWindowExStyle(window));
		} */
	*left = int(-rect.Left)
	*top = int(-rect.Top)
	*right = int(rect.Right) - width
	*bottom = int(rect.Bottom) - height
}

// SwapInterval sets the swap interval for the current context, i.e. the number
// of screen updates to wait before swapping the buffers of a Window and
// returning from SwapBuffers. This is sometimes called
// 'vertical synchronization', 'vertical retrace synchronization' or 'vsync'.
//
// Contexts that support either of the WGL_EXT_swap_control_tear and
// GLX_EXT_swap_control_tear extensions also accept negative swap intervals,
// which allow the driver to swap even if a frame arrives a little bit late.
// You can check for the presence of these extensions using
// ExtensionSupported. For more information about swap tearing,
// see the extension specifications.
//
// Some GPU drivers do not honor the requested swap interval, either because of
// user settings that override the request or due to bugs in the driver.
func SwapInterval(interval int) {
	// C.glfwSwapInterval(C.int(interval))
	panicError()
}

// GetCursorPos returns the last reported position of the cursor.
//
// If the cursor is disabled (with CursorDisabled) then the cursor position is
// unbounded and limited only by the minimum and maximum values of a double.
//
// The coordinate can be converted to their integer equivalents with the floor
// function. Casting directly to an integer type works for positive coordinates,
// but fails for negative ones.
func (w *Window) GetCursorPos() (x, y float64) {
	var xpos, ypos float32
	// C.glfwGetCursorPos(w.Data, &xpos, &ypos)
	panicError()
	return float64(xpos), float64(ypos)
}

func _glfwGetWindowSizeWin32(window *_GLFWwindow, width *int, height *int) {
	var area RECT
	_, _, err := _GetClientRect.Call(uintptr(unsafe.Pointer(window.Win32.handle)), uintptr(unsafe.Pointer(&area)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic(err)
	}
	// GetClientRect(window->win32.handle, &area);
	*width = int(area.Right)
	*height = int(area.Bottom)
}

// GetSize returns the size, in screen coordinates, of the client area of the
// specified Window.
func (w *Window) GetSize() (width int, height int) {
	var wi, h int
	_glfwGetWindowSizeWin32(w.Data, &wi, &h)
	panicError()
	return int(wi), int(h)
}
