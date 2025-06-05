package glfw

import "C"
import (
	"errors"
	"fmt"
	"github.com/jkvatne/jkvgui/gl"
	"golang.org/x/sys/windows"
	"log/slog"
	"strings"
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

func WindowHint(target Hint, hint int) {
	// C.glfwWindowHint(C.int(target), C.int(hint))
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
	_SetProcessDPIAware.Call()
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
	if ctxconfig.source != GLFW_NATIVE_CONTEXT_API &&
		ctxconfig.source != GLFW_EGL_CONTEXT_API &&
		ctxconfig.source != GLFW_OSMESA_CONTEXT_API {
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
		frameX = CW_USEDEFAULT
		frameY = CW_USEDEFAULT
		frameWidth = rect.Right - rect.Left
		frameHeight = rect.Bottom - rect.Top
	}
	// wideTitle = _glfwCreateWideStringFromUTF8Win32(wndconfig.title)
	window.Win32.handle, err = CreateWindowEx(
		WS_OVERLAPPED|WS_EX_APPWINDOW,
		_glfw.class,
		wndconfig.title,
		WS_OVERLAPPED|WS_CLIPSIBLINGS|WS_CLIPCHILDREN,
		frameX, frameY, // Window position
		int32(frameWidth), int32(frameHeight), // Window width/heigth
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
	opengl32          = windows.NewLazySystemDLL("opengl32")
	wglGetProcAddress = opengl32.NewProc("wglGetProcAddress")
)

const (
	PFD_DRAW_TO_WINDOW = 0x04
	PFD_SUPPORT_OPENGL = 0x20
	PFD_DOUBLEBUFFER   = 0x01
	PFD_TYPE_RGBA      = 0x00
)

func glfwMakeContextCurrent(window *_GLFWwindow) error {
	// _GLFWwindow* window = (_GLFWwindow*) handle;
	// previous := _glfwPlatformGetTls(&_glfw.contextSlot);
	if window != nil && window.context.client == GLFW_NO_API {
		return fmt.Errorf("Cannot make current with a window that has no OpenGL or OpenGL ES context")
	}
	// if previous!=nil && window!=nil || window->context.source != previous->context.source)
	//		previous->context.makeCurrent(NULL);
	// }

	if window != nil {
		window.context.makeCurrent(&window.context)
	}
	return nil
}

// Initialize WGL
func _glfwInitWGL() error {
	var pfd PIXELFORMATDESCRIPTOR
	var pdc uintptr
	if _glfw.wgl.instance != 0 {
		return nil
	}
	// _glfw.wgl.instance = windows.NewLazySystemDLL("opengl32.dll");
	// getProcAddress := opengl32.NewProc("wglGetProcAddress")
	_glfw.wgl.wglCreateContext = opengl32.NewProc("wglCreateContext")
	_glfw.wgl.wglDeleteContext = opengl32.NewProc("wglDeleteContext")
	_glfw.wgl.wglGetProcAddress = opengl32.NewProc("wglGetProcAddress")
	_glfw.wgl.wglGetCurrentDC = opengl32.NewProc("wglGetCurrentDC")
	_glfw.wgl.wglGetCurrentContext = opengl32.NewProc("wglGetCurrentContext")
	_glfw.wgl.wglMakeCurrent = opengl32.NewProc("wglMakeCurrent")
	_glfw.wgl.wglShareLists = opengl32.NewProc("wglShareLists")
	_glfw.wgl.wglGetString = opengl32.NewProc("wglGetString")

	// NOTE: A dummy context has to be created for opengl32.dll to load the
	//       OpenGL ICD, from which we can then query WGL extensions
	// NOTE: This code will accept the Microsoft GDI ICD; accelerated context
	//       creation failure occurs during manual pixel format enumeration

	// dc := GetDC(_glfw.win32.helperWindowHandle)

	pfd.nSize = uint16(unsafe.Sizeof(pfd))
	pfd.nVersion = 1
	pfd.dwFlags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER
	pfd.iPixelType = PFD_TYPE_RGBA
	pfd.cColorBits = 24

	// if !SetPixelFormat(dc, ChoosePixelFormat(dc, &pfd), &pfd) {
	//	return fmt.Errorf("WGL: Failed to set pixel format for dummy context")
	// }

	// rc, _, _ := _glfw.wgl.wglCreateContext.Call(uintptr(dc))
	// if rc == 0 {
	// 	return fmt.Errorf("WGL: Failed to create dummy context")
	// }

	pdc, _, _ = _glfw.wgl.wglGetCurrentDC.Call()
	prc, _, _ := _glfw.wgl.wglGetCurrentContext.Call()

	// ret, _, _ := _glfw.wgl.wglMakeCurrent.Call(dc, rc)
	// if ret == 0 {
	// 	_, _, _ = _glfw.wgl.wglMakeCurrent.Call(pdc, prc)
	// 	_, _, _ = _glfw.wgl.wglDeleteContext.Call(rc)
	// 	return fmt.Errorf("WGL: Failed to make dummy context current")
	// }

	// NOTE: Functions must be loaded first as they're needed to retrieve the
	//       extension string that tells us whether the functions are supported
	/*
		_glfw.wgl.GetExtensionsStringEXT = (PFNWGLGETEXTENSIONSSTRINGEXTPROC)
		wglGetProcAddress("wglGetExtensionsStringEXT")
		_glfw.wgl.GetExtensionsStringARB = (PFNWGLGETEXTENSIONSSTRINGARBPROC)
		wglGetProcAddress("wglGetExtensionsStringARB")
		_glfw.wgl.CreateContextAttribsARB = (PFNWGLCREATECONTEXTATTRIBSARBPROC)
		wglGetProcAddress("wglCreateContextAttribsARB")
		_glfw.wgl.SwapIntervalEXT = (PFNWGLSWAPINTERVALEXTPROC)
		wglGetProcAddress("wglSwapIntervalEXT")
		_glfw.wgl.GetPixelFormatAttribivARB = (PFNWGLGETPIXELFORMATATTRIBIVARBPROC)
		wglGetProcAddress("wglGetPixelFormatAttribivARB")
	*/

	// NOTE: WGL_ARB_extensions_string and WGL_EXT_extensions_string are not
	//       checked below as we are already using them
	/*
		_glfw.wgl.ARB_multisample =
			extensionSupportedWGL("WGL_ARB_multisample")
		_glfw.wgl.ARB_framebuffer_sRGB =
			extensionSupportedWGL("WGL_ARB_framebuffer_sRGB")
		_glfw.wgl.EXT_framebuffer_sRGB =
			extensionSupportedWGL("WGL_EXT_framebuffer_sRGB")
		_glfw.wgl.ARB_create_context =
			extensionSupportedWGL("WGL_ARB_create_context")
		_glfw.wgl.ARB_create_context_profile =
			extensionSupportedWGL("WGL_ARB_create_context_profile")
		_glfw.wgl.EXT_create_context_es2_profile =
			extensionSupportedWGL("WGL_EXT_create_context_es2_profile")
		_glfw.wgl.ARB_create_context_robustness =
			extensionSupportedWGL("WGL_ARB_create_context_robustness")
		_glfw.wgl.ARB_create_context_no_error =
			extensionSupportedWGL("WGL_ARB_create_context_no_error")
		_glfw.wgl.EXT_swap_control =
			extensionSupportedWGL("WGL_EXT_swap_control")
		_glfw.wgl.EXT_colorspace =
			extensionSupportedWGL("WGL_EXT_colorspace")
		_glfw.wgl.ARB_pixel_format =
			extensionSupportedWGL("WGL_ARB_pixel_format")
		_glfw.wgl.ARB_context_flush_control =
			extensionSupportedWGL("WGL_ARB_context_flush_control")
	*/
	_, _, _ = _glfw.wgl.wglMakeCurrent.Call(pdc, prc)
	// _, _, _ = _glfw.wgl.wglDeleteContext.Call(rc)
	return nil
}

func _glfwCreateContextWGL(window *_GLFWwindow, ctxconfig *_GLFWctxconfig, fbconfig *_GLFWfbconfig) error {
	return nil
}

func _glfwCreateContextEGL(window *_GLFWwindow, ctxconfig *_GLFWctxconfig, fbconfig *_GLFWfbconfig) error {
	return nil
}

func glfwInitEGL() error {
	return nil
}

func glfwInitOSMesa() error {
	return nil
}
func _glfwCreateContextOSMesa(window *_GLFWwindow, ctxconfig *_GLFWctxconfig, fbconfig *_GLFWfbconfig) error {
	return nil
}
func _glfwRefreshContextAttribs(window *_GLFWwindow, ctxconfig *_GLFWctxconfig) error {
	prefixes := []string{
		"OpenGL ES-CM ",
		"OpenGL ES-CL ",
		"OpenGL ES ",
	}
	window.context.source = ctxconfig.source
	window.context.client = GLFW_OPENGL_API
	// previous = _glfwPlatformGetTls(&_glfw.contextSlot);
	glfwMakeContextCurrent(window)
	// if (_glfwPlatformGetTls(&_glfw.contextSlot) != window)
	//    return GLFW_FALSE;

	// window.context.GetIntegerv = (PFNGLGETINTEGERVPROC)
	// window.context.getProcAddress("glGetIntegerv");
	// window.context.GetString = (PFNGLGETSTRINGPROC)
	// window.context.getProcAddress("glGetString");
	// if (!window.context.GetIntegerv || !window.context.GetString)
	// {
	// 	_glfwInputError(GLFW_PLATFORM_ERROR, "Entry point retrieval is broken");
	// 	glfwMakeContextCurrent((GLFWwindow*) previous);
	// 	return GLFW_FALSE;
	// }
	version := gl.GoStr(gl.GetString(gl.VERSION))
	slog.Info("GLFW got", "versions", version)
	// if (!version) {
	// 	if (ctxconfig.client == GLFW_OPENGL_API) {
	// 		_glfwInputError(GLFW_PLATFORM_ERROR, "OpenGL version string retrieval is broken");
	// 	} else {
	// 		_glfwInputError(GLFW_PLATFORM_ERROR, "OpenGL ES version string retrieval is broken");
	// 	}
	// 	glfwMakeContextCurrent((GLFWwindow*) previous);
	// 	return nil
	// }

	for _, pref := range prefixes {
		if strings.HasPrefix(pref, version) {
			window.context.client = GLFW_OPENGL_ES_API
			break
		}
	}

	window.context.major = 3
	window.context.minor = 3
	window.context.revision = 3
	if window.context.major == 0 {
		// glfwMakeContextCurrent((GLFWwindow *)
		if window.context.client == GLFW_OPENGL_API {
			return fmt.Errorf("No version found in OpenGL version string")
		} else {
			return fmt.Errorf("No version found in OpenGL ES version string")
		}
	}
	if window.context.major < ctxconfig.major || window.context.major == ctxconfig.major && window.context.minor < ctxconfig.minor {
		// The desired OpenGL version is greater than the actual version
		// This only happens if the machine lacks {GLX|WGL}_ARB_create_context
		// /and/ the user has requested an OpenGL version greater than 1.0
		// glfwMakeContextCurrent((GLFWwindow*) previous);
		if window.context.client == GLFW_OPENGL_API {
			return fmt.Errorf("Requested OpenGL version %i.%i, got version %i.%i", ctxconfig.major, ctxconfig.minor, window.context.major, window.context.minor)
		} else {
			return fmt.Errorf("Requested OpenGL ES version %i.%i, got version %i.%i", ctxconfig.major, ctxconfig.minor, window.context.major, window.context.minor)
		}
	}
	/*
		if (window.context.major >= 3) {
			// OpenGL 3.0+ uses a different function for extension string retrieval
			// We cache it here instead of in glfwExtensionSupported mostly to alert
			// users as early as possible that their build may be broken
			window.context.GetStringi = (PFNGLGETSTRINGIPROC)
			window.context.getProcAddress("glGetStringi");
			if (!window.context.GetStringi)	{
				_glfwInputError(GLFW_PLATFORM_ERROR,
					"Entry point retrieval is broken");
				glfwMakeContextCurrent((GLFWwindow*) previous);
				return GLFW_FALSE;
			}
		}

		if (window.context.client == GLFW_OPENGL_API) {
			// Read back context flags (OpenGL 3.0 and above)
			if (window.context.major >= 3) {
				// window.context.GetIntegerv(GL_CONTEXT_FLAGS, &flags);
				if (flags & GL_CONTEXT_FLAG_FORWARD_COMPATIBLE_BIT)	{
					window.context.forward = GLFW_TRUE;
					}
				if (flags & GL_CONTEXT_FLAG_DEBUG_BIT)	window.context.debug = GLFW_TRUE;
				else if (glfwExtensionSupported("GL_ARB_debug_output") && ctxconfig.debug) {
					window.context.debug = GLFW_TRUE;
				}

				if (flags & GL_CONTEXT_FLAG_NO_ERROR_BIT_KHR)
					window.context.noerror = GLFW_TRUE;
			}

			// Read back OpenGL context profile (OpenGL 3.2 and above)
			if (window.context.major >= 4 ||
				(window.context.major == 3 && window.context.minor >= 2))
			{
				GLint mask;
				window.context.GetIntegerv(GL_CONTEXT_PROFILE_MASK, &mask);

				if (mask & GL_CONTEXT_COMPATIBILITY_PROFILE_BIT)
					window.context.profile = GLFW_OPENGL_COMPAT_PROFILE;
				else if (mask & GL_CONTEXT_CORE_PROFILE_BIT)
				window.context.profile = GLFW_OPENGL_CORE_PROFILE;
				else if (glfwExtensionSupported("GL_ARB_compatibility"))
			{
				// HACK: This is a workaround for the compatibility profile bit
				//       not being set in the context flags if an OpenGL 3.2+
				//       context was created without having requested a specific
				//       version
				window.context.profile = GLFW_OPENGL_COMPAT_PROFILE;
			}
			}

			// Read back robustness strategy
			if (glfwExtensionSupported("GL_ARB_robustness"))
			{
				// NOTE: We avoid using the context flags for detection, as they are
				//       only present from 3.0 while the extension applies from 1.1

				GLint strategy;
				window.context.GetIntegerv(GL_RESET_NOTIFICATION_STRATEGY_ARB,
					&strategy);

				if (strategy == GL_LOSE_CONTEXT_ON_RESET_ARB)
					window.context.robustness = GLFW_LOSE_CONTEXT_ON_RESET;
				else if (strategy == GL_NO_RESET_NOTIFICATION_ARB)
				window.context.robustness = GLFW_NO_RESET_NOTIFICATION;
			}
		}
		else
		{
			// Read back robustness strategy
			if (glfwExtensionSupported("GL_EXT_robustness"))
			{
				// NOTE: The values of these constants match those of the OpenGL ARB
				//       one, so we can reuse them here

				GLint strategy;
				window.context.GetIntegerv(GL_RESET_NOTIFICATION_STRATEGY_ARB,
					&strategy);

				if (strategy == GL_LOSE_CONTEXT_ON_RESET_ARB)
					window.context.robustness = GLFW_LOSE_CONTEXT_ON_RESET;
				else if (strategy == GL_NO_RESET_NOTIFICATION_ARB)
				window.context.robustness = GLFW_NO_RESET_NOTIFICATION;
			}
		}
	*/
	/*
		if (glfwExtensionSupported("GL_KHR_context_flush_control"))	{
			GLint behavior;
			window.context.GetIntegerv(GL_CONTEXT_RELEASE_BEHAVIOR, &behavior);

			if (behavior == GL_NONE)
				window.context.release = GLFW_RELEASE_BEHAVIOR_NONE;
			else if (behavior == GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH)
			window.context.release = GLFW_RELEASE_BEHAVIOR_FLUSH;
		}
		// Clearing the front buffer to black to avoid garbage pixels left over from
		// previous uses of our bit of VRAM
		{
			PFNGLCLEARPROC glClear = (PFNGLCLEARPROC)
			window.context.getProcAddress("glClear");
			glClear(GL_COLOR_BUFFER_BIT);

			if (window.doublebuffer)
				window.context.swapBuffers(window);
		}
	*/

	// glfwMakeContextCurrent(previous);
	return nil
}

func glfwPlatformCreateWindow(window *_GLFWwindow, wndconfig *_GLFWwndconfig, ctxconfig *_GLFWctxconfig, fbconfig *_GLFWfbconfig) error {
	err := createNativeWindow(window, wndconfig, fbconfig)
	if err != nil {
		return err
	}
	if ctxconfig.client != GLFW_NO_API {
		if ctxconfig.source == GLFW_NATIVE_CONTEXT_API {
			if err := _glfwInitWGL(); err != nil {
				return fmt.Errorf("glglfwPlatformCreateWindowfw error")
			}
			if err := _glfwCreateContextWGL(window, ctxconfig, fbconfig); err != nil {
				return err
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
		// _glfwPlatformShowWindow(window)
		// _glfwPlatformFocusWindow(window)
		// acquireMonitor(window)
		// fitToMonitor(window)
		if wndconfig.centerCursor {
			// _glfwCenterCursorInContentArea(window)
		}
	} else if wndconfig.visible {
		// _glfwPlatformShowWindow(window)
		if wndconfig.focused {
			// _glfwPlatformFocusWindow(window)
		}
	}
	return nil
}

/*
SetProcessDPIAware()
var err error
Window.Win32.handle, err = CreateWindowEx(
WS_OVERLAPPED|WS_EX_APPWINDOW,
_glfw.class,
"",
WS_OVERLAPPED|WS_CLIPSIBLINGS|WS_CLIPCHILDREN,
CW_USEDEFAULT, CW_USEDEFAULT, // Window position
int32(width), int32(height), // Window width/heigth
0, // No parent
0, // No menu
resources.handle,
0)
return Window, err
*/

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
		return nil, fmt.Errorf("Error creating window")
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

func glfwSetWindowSize(w *_GLFWwindow, xpos, ypos int) {

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
	// TODO _glfwPlatformShowWindow(window);
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
	w.Data.context.makeCurrent(&w.Data.context)
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

func _glfwPlatformInit() error {
	return nil
}

// Init() is GLFWAPI int glfwInit(void) from init.c
func Init() error {
	var err error
	// Repeated calls do nothing
	if _glfw.initialized {
		return nil
	}
	_glfw.hints.init = _GLFWinitconfig{}

	// This is _glfwPlatformInit():
	// TODO SystemParametersInfoW(SPI_GETFOREGROUNDLOCKTIMEOUT, 0, &_glfw.Win32.foregroundLockTimeout, 0);
	// TODO SystemParametersInfoW(SPI_SETFOREGROUNDLOCKTIMEOUT, 0, UIntToPtr(0), SPIF_SENDCHANGE);
	// TODO createKeyTables()
	// TODO _glfwUpdateKeyNamesWin32()
	/*
		if(_glfwIsWindows10CreatorsUpdateOrGreaterWin32() {
			SetProcessDpiAwarenessContext(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2);
		} else if IsWindows8Point1OrGreater() {
			SetProcessDpiAwareness(PROCESS_PER_MONITOR_DPI_AWARE);
		} else if(IsWindowsVistaOrGreater() {
			SetProcessDPIAware()
		}
	*/
	if err := _glfwRegisterWindowClassWin32(); err != nil {
		return fmt.Errorf("glfw platform init failed, _glfwRegisterWindowClassWin32 failed, %v ", err.Error())
	}
	// _, _, err := _procGetModuleHandleExW.Call(GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS|GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT, uintptr(unsafe.Pointer(&_glfw)), uintptr(unsafe.Pointer(&_glfw.instance)))
	_glfw.instance, err = GetModuleHandle()
	if err != nil {
		return fmt.Errorf("glfw platform init failed %v ", err.Error())
	}

	// if !createHelperWindow() {
	//	return nil
	// }
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
	var x, y float32
	// C.glfwGetWindowContentScale(w.Data, &x, &y)
	return float32(x), float32(y)
}

// GetFrameSize retrieves the size, in screen coordinates, of each edge of the frame
// of the specified Window. This size includes the title bar, if the Window has one.
// The size of the frame may vary depending on the Window-related hints used to create it.
//
// Because this function retrieves the size of each Window frame edge and not the offset
// along a particular coordinate axis, the retrieved values will always be zero or positive.
func (w *Window) GetFrameSize() (left, top, right, bottom int) {
	var l, t, r, b int
	// C.glfwGetWindowFrameSize(w.Data, &l, &t, &r, &b)
	panicError()
	return int(l), int(t), int(r), int(b)
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

// GetSize returns the size, in screen coordinates, of the client area of the
// specified Window.
func (w *Window) GetSize() (width, height int) {
	var wi, h int
	// C.glfwGetWindowSize(w.Data, &wi, &h)
	panicError()
	return int(wi), int(h)
}
