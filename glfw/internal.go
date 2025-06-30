package glfw

import (
	"errors"
	"golang.org/x/sys/windows"
	"log/slog"
	"sync"
	"syscall"
	"unicode"
	"unsafe"
)

const (
	GLFW_MOD_CAPS_LOCK         = 0x0010
	GLFW_MOD_NUM_LOCK          = 0x0020
	GLFW_CONNECTED             = 0x00040001
	GLFW_DISCONNECTED          = 0x00040002
	_GLFW_STICK                = 3
	_GLFW_INSERT_FIRST         = 0
	_GLFW_INSERT_LAST          = 1
	GL_VERSION                 = 0x1F02
	GLFW_OPENGL_ANY_PROFILE    = 0
	GLFW_OPENGL_CORE_PROFILE   = 0x00032001
	GLFW_OPENGL_COMPAT_PROFILE = 0x00032002
	_GLFW_WNDCLASSNAME         = "GLFW30"
	GLFW_DONT_CARE             = -1
	OpenGLCoreProfile          = 0x00032001
	OpenGLForwardCompatible    = GLFW_OPENGL_FORWARD_COMPAT
	True                       = 1
	False                      = 0

	OpenGLProfile       = 0x00022008
	ForwardCompatible   = 0x00022006
	Resizable           = 0x00020003
	Focused             = 0x00020001
	Iconified           = 0x00020002
	Resizeable          = 0x00020003
	Visible             = 0x00020004
	Decorated           = 0x00020005
	AutoIconify         = 0x00020006
	Floating            = 0x00020007
	Maximized           = 0x00020008
	ContextVersionMajor = 0x00022002
	ContextVersionMinor = GLFW_CONTEXT_VERSION_MINOR
	Samples             = 0x0002100D

	ArrowCursor     = 0x00036001
	IBeamCursor     = 0x00036002
	CrosshairCursor = 0x00036003
	HandCursor      = 0x00036004
	HResizeCursor   = 0x00036005
	VResizeCursor   = 0x00036006

	IMAGE_BITMAP      = 0
	IMAGE_ICON        = 1
	IMAGE_CURSOR      = 2
	IMAGE_ENHMETAFILE = 3

	LR_CREATEDIBSECTION = 0x00002000
	LR_DEFAULTCOLOR     = 0x00000000
	LR_DEFAULTSIZE      = 0x00000040
	LR_LOADFROMFILE     = 0x00000010
	LR_LOADMAP3DCOLORS  = 0x00001000
	LR_LOADTRANSPARENT  = 0x00000020
	LR_MONOCHROME       = 0x00000001
	LR_SHARED           = 0x00008000
	LR_VGACOLOR         = 0x00000080

	CS_HREDRAW               = 0x0002
	CS_INSERTCHAR            = 0x2000
	CS_NOMOVECARET           = 0x4000
	CS_VREDRAW               = 0x0001
	CS_OWNDC                 = 0x0020
	KF_EXTENDED              = 0x100
	GLFW_RELEASE             = 0
	GLFW_PRESS               = 1
	GLFW_REPEAT              = 2
	GLFW_OPENGL_API          = 0x00030001
	GLFW_NATIVE_CONTEXT_API  = 0x00036001
	GLFW_OPENGL_ES_API       = 0x00030002
	GLFW_EGL_CONTEXT_API     = 0x00036002
	GLFW_OSMESA_CONTEXT_API  = 0x00036003
	GLFW_NO_API              = 0
	MONITOR_DEFAULTTONULL    = 0x00000000
	MONITOR_DEFAULTTOPRIMARY = 0x00000001
	MONITOR_DEFAULTTONEAREST = 0x00000002
	USER_DEFAULT_SCREEN_DPI  = 96
	LOGPIXELSX               = 88
	LOGPIXELSY               = 90
	SIZE_RESTORED            = 0
	SIZE_MINIMIZED           = 1
	SIZE_MAXIMIZED           = 2
	SIZE_MAXSHOW             = 3
	SIZE_MAXHIDE             = 4
)

type _GLFWvidmode struct {
	width       int
	height      int
	redBits     int
	greenBits   int
	blueBits    int
	refreshRate int
}

type (
	_GLFWmakecontextcurrentfun = func(w *Window) error
	_GLFWswapbuffersfun        = func(w *Window)
	_GLFWextensionsupportedfun = func(x byte) bool
	_GLFWgetprocaddressfun     = func()
	_GLFWdestroycontextfun     = func(w *Window)
)

// Context structure
//
type _GLFWcontext struct {
	major, minor, revision  int
	forward, debug, noerror bool
	profile                 int
	robustness              int
	release                 int
	makeCurrent             _GLFWmakecontextcurrentfun
	swapBuffers             _GLFWswapbuffersfun
	getProcAddress          _GLFWgetprocaddressfun
	destroy                 _GLFWdestroycontextfun
	wgl                     struct {
		dc       HDC
		handle   HANDLE
		interval int
	}
}

type _GLFWwindow struct {
	next *_GLFWwindow
	// Window settings and state
	resizable              bool
	decorated              bool
	autoIconify            bool
	floating               bool
	focusOnShow            bool
	shouldClose            bool
	userPointer            unsafe.Pointer
	doublebuffer           bool
	videoMode              _GLFWvidmode
	monitor                *Monitor
	cursor                 *Cursor
	minwidth               int
	minheight              int
	maxwidth               int
	maxheight              int
	numer                  int
	denom                  int
	stickyKeys             bool
	stickyMouseButtons     bool
	lockKeyMods            bool
	cursorMode             int
	mouseButtons           [MouseButtonLast + 1]byte
	keys                   [KeyLast + 1]byte
	virtualCursorPosX      float64 // Virtual cursor position when cursor is disabled
	virtualCursorPosY      float64 // Virtual cursor position when cursor is disabled
	context                _GLFWcontext
	lastCursorPosX         float64 // The last received cursor position, regardless of source
	lastCursorPosY         float64 // The last received cursor position, regardless of source
	charCallback           CharCallback
	focusCallback          FocusCallback
	keyCallback            KeyCallback
	mouseButtonCallback    MouseButtonCallback
	cursorPosCallback      CursorPosCallback
	scrollCallback         ScrollCallback
	refreshCallback        RefreshCallback
	sizeCallback           SizeCallback
	dropCallback           DropCallback
	contentScaleCallback   ContentScaleCallback
	fFramebufferSizeHolder func(w *_GLFWwindow, width int, height int)
	fCloseHolder           func(w *_GLFWwindow)
	fMaximizeHolder        func(w *_GLFWwindow, maximized bool)
	fIconifyHolder         func(w *_GLFWwindow, iconified bool)
	fCursorEnterHolder     func(w *_GLFWwindow, entered bool)
	fCharModsHolder        func(w *_GLFWwindow, char rune, mods ModifierKey)
	Win32                  _GLFWwindowWin32
}

type _GLFWwindowWin32 = struct {
	handle         syscall.Handle
	bigIcon        syscall.Handle
	smallIcon      syscall.Handle
	cursorTracked  bool
	frameAction    bool
	iconified      bool
	maximized      bool
	transparent    bool // Whether to enable framebuffer transparency on DWM
	scaleToMonitor bool
	width          int    // Cached size used to filter out duplicate events
	height         int    // Cached size used to filter out duplicate events
	highSurrogate  uint16 // The last recevied high surrogate when decoding pairs of UTF-16 messages
}

type _GLFWinitconfig = struct {
	hatButtons bool
	ns         struct {
		menubar bool
		chdir   bool
	}
	wl struct {
		libdecorMode int
	}
}
type _GLFWwndconfig = struct {
	xpos             int
	ypos             int
	width            int
	height           int
	title            string
	resizable        bool
	visible          bool
	decorated        bool
	focused          bool
	autoIconify      bool
	floating         bool
	maximized        bool
	centerCursor     bool
	focusOnShow      bool
	mousePassthrough bool
	scaleToMonitor   bool
	ns               struct {
		retina    bool
		frameName string
	}
}

type _GLFWctxconfig = struct {
	client     int
	source     int
	major      int
	minor      int
	forward    bool
	debug      bool
	noerror    bool
	profile    int
	robustness int
	release    int
	share      *_GLFWwindow
	nsgl       struct {
		offline bool
	}
}

type hints = struct {
	init        _GLFWinitconfig
	framebuffer _GLFWfbconfig
	window      _GLFWwndconfig
	context     _GLFWctxconfig
	refreshRate int
}

type _GLFWfbconfig = struct {
	redBits        int
	greenBits      int
	blueBits       int
	alphaBits      int
	depthBits      int
	stencilBits    int
	accumRedBits   int
	accumGreenBits int
	accumBlueBits  int
	accumAlphaBits int
	auxBuffers     int
	samples        int
	sRGB           bool
	doublebuffer   bool
	transparent    bool
	handle         uintptr
}

type _GLFWerror struct {
	next        *_GLFWerror
	code        int
	description string
}

type _GLFWtls = struct {
	allocated bool
	index     int
}

// Library global Data
var _glfw struct {
	hints
	class           uint16
	available       bool
	instance        syscall.Handle
	initialized     bool
	errorListHead   *_GLFWerror
	cursorListHead  *Cursor
	windowListHead  *_GLFWwindow
	monitors        []*Monitor
	monitorCallback func(w *Monitor, action int)
	monitorCount    int
	errorSlot       _GLFWtls
	contextSlot     _GLFWtls
	errorLock       sync.Mutex
	win32           struct {
		helperWindowHandle syscall.Handle
		helperWindowClass  uint16
		mainWindowClass    uint16
		blankCursor        HANDLE
		keycodes           [512]Key
		scancodes          [512]int16
		instance           syscall.Handle
	}
	wgl struct {
		dc                         HDC
		handle                     syscall.Handle
		interval                   int
		instance                   *windows.LazyDLL
		wglCreateContextAttribsARB *windows.LazyProc
		wglDeleteContext           *windows.LazyProc
		wglGetProcAddress          *windows.LazyProc
		wglGetCurrentDC            *windows.LazyProc
		wglGetCurrentContext       *windows.LazyProc
		wglMakeCurrent             *windows.LazyProc
		wglSwapBuffers             *windows.LazyProc
		wglCreateContext           *windows.LazyProc
		wglSetPixelFormat          *windows.LazyProc
		wglChoosePixelFormat       *windows.LazyProc
		wglDescribePixelFormat     *windows.LazyProc
		getProcAddress             *windows.LazyProc
		GetExtensionsStringEXT     *windows.LazyProc
		GetExtensionsStringARB     *windows.LazyProc
		GetPixelFormatAttribivARB  *windows.LazyProc
		GetDeviceCaps              *windows.LazyProc
		GetString                  *windows.LazyProc
		ARB_pixel_format           int
		ARB_multisample            bool
		ARB_framebuffer_sRGB       bool
		EXT_framebuffer_sRGB       bool
		EXT_colorspace             bool
	}
}

func glfwInputKey(window *_GLFWwindow, key Key, scancode int, action int, mods ModifierKey) {
	var repeated bool
	if key >= 0 && key <= KeyLast {
		repeated = false

		if action == GLFW_RELEASE && window.keys[key] == GLFW_RELEASE {
			return
		}

		if action == GLFW_PRESS && window.keys[key] == GLFW_PRESS {
			repeated = true
		}

		if action == GLFW_RELEASE && window.stickyKeys {
			window.keys[key] = _GLFW_STICK
		} else {
			window.keys[key] = uint8(action)
		}
		if repeated {
			action = GLFW_REPEAT
		}
	}
	if !window.lockKeyMods {
		mods &= ^(GLFW_MOD_CAPS_LOCK | GLFW_MOD_NUM_LOCK)
	}

	if window.keyCallback != nil {
		window.keyCallback(window, key, scancode, Action(action), mods)
	}
}

func glfwInputMouseClick(window *_GLFWwindow, button MouseButton, action Action, mods ModifierKey) {
	// TODO if (!window.lockKeyMods)	mods &= ~(GLFW_MOD_CAPS_LOCK | GLFW_MOD_NUM_LOCK);
	// TODO if (action == GLFW_RELEASE && window.stickyMouseButtons) window.mouseButtons[button] = _GLFW_STICK; else window.mouseButtons[button] = (char) action;
	if window.mouseButtonCallback != nil {
		window.mouseButtonCallback(window, button, action, mods)
	}
}

// Notifies shared code that a window has lost or received input focus
func glfwInputWindowFocus(window *_GLFWwindow, focused bool) {
	if window == nil {
		return
	}
	if window.focusCallback != nil {
		window.focusCallback(window, focused)
	}
	if !focused {
		// Force release of buttons
		/* TODO
		for k := Key(0);  k <= KeyLast;  k++ {
			if (window.keys[k] == GLFW_PRESS) {
				scancode := glfwPlatformGetKeyScancode(k);
				glfwInputKey(window, k, scancode, GLFW_RELEASE, 0);
			}
		}*/
		for button := MouseButton(0); button <= MouseButtonLast; button++ {
			if window.mouseButtons[button] == GLFW_PRESS {
				glfwInputMouseClick(window, button, GLFW_RELEASE, 0)
			}
		}
	}
}

func glfwInputCursorPos(window *_GLFWwindow, xpos, ypos float64) {
	if window.virtualCursorPosX == xpos && window.virtualCursorPosY == ypos {
		return
	}
	window.virtualCursorPosX = xpos
	window.virtualCursorPosY = ypos
	if window.cursorPosCallback != nil {
		window.cursorPosCallback(window, xpos, ypos)
	}
}

func glfwInputScroll(window *_GLFWwindow, xoffset, yoffset float64) {
	if window.scrollCallback != nil {
		window.scrollCallback(window, xoffset, yoffset)
	}
}

func glfwInputWindowDamage(window *_GLFWwindow) {
	if window.refreshCallback != nil {
		window.refreshCallback(window)
	}
}

func glfwInputWindowCloseRequest(window *_GLFWwindow) {
	slog.Error("Got CloseRequest")
}

func getKeyMods() ModifierKey {
	var mods ModifierKey
	if GetKeyState(VK_SHIFT)&0x8000 != 0 {
		mods |= ModShift
	}
	if GetKeyState(VK_CONTROL)&0x8000 != 0 {
		mods |= ModControl
	}
	if GetKeyState(VK_MENU)&0x8000 != 0 {
		mods |= ModAlt
	}
	if (GetKeyState(VK_LWIN)|GetKeyState(VK_RWIN))&0x8000 != 0 {
		mods |= ModSuper
	}
	if (GetKeyState(VK_CAPITAL) & 1) != 0 {
		mods |= ModCapsLock
	}
	if (GetKeyState(VK_NUMLOCK) & 1) != 0 {
		mods |= ModNumLock
	}
	return mods
}

var winMap map[syscall.Handle]*_GLFWwindow

func getProp(hwnd syscall.Handle) *_GLFWwindow {
	return winMap[hwnd]
}

func setProp(hwnd syscall.Handle, prop *_GLFWwindow) {
	if winMap == nil {
		winMap = make(map[syscall.Handle]*_GLFWwindow)
	}
	winMap[hwnd] = prop
}

func windowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	window := getProp(hwnd)
	if window == nil {
		r1, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
		return r1
	}

	switch msg {
	case WM_CLOSE:
		window.shouldClose = true
	case WM_UNICHAR:
		if wParam == UNICODE_NOCHAR {
			// Tell the system that we accept WM_UNICHAR messages.
			return TRUE
		}
		fallthrough
	case WM_CHAR, WM_SYSCHAR:
		if r := rune(wParam); unicode.IsPrint(r) {
			window.charCallback(nil, r)
		}
		return TRUE
	case WM_DPICHANGED:
		// Let Windows know we're prepared for runtime DPI changes.
		return TRUE
	case WM_ERASEBKGND:
		// Avoid flickering between GPU content and background color.
		return TRUE
	case WM_KEYDOWN, WM_KEYUP, WM_SYSKEYDOWN, WM_SYSKEYUP:
		var key Key
		action := GLFW_PRESS
		if (lParam>>16)&0x8000 != 0 {
			action = GLFW_RELEASE
		}
		mods := getKeyMods()
		scancode := int((lParam >> 16) & 0x1ff)
		switch scancode {
		case 0: // scancode = MapVirtualKeyW((UINT) wParam, MAPVK_VK_TO_VSC);
		case 0x54:
			scancode = 0x137 // Alt+PrtSc
		case 0x146:
			scancode = 0x45 // Ctrl+Pause
		case 0x136:
			scancode = 0x36 // CJK IME sets the extended bit for right Shift
		}

		key = _glfw.win32.keycodes[scancode]
		if wParam == VK_CONTROL {
			if lParam>>16&KF_EXTENDED != 0 {
				// Right side keys have the extended key bit set
				key = KeyRightControl
			} else {
				/*
					// NOTE: Alt Gr sends Left Ctrl followed by Right Alt
					// HACK: We only want one event for Alt Gr, so if we detect
					//       this sequence we discard this Left Ctrl message now
					//       and later report Right Alt normally
					MSG next;
					const DWORD time = GetMessageTime();

					if (PeekMessageW(&next, NULL, 0, 0, PM_NOREMOVE)) {
						if (next.message == WM_KEYDOWN ||
							next.message == WM_SYSKEYDOWN ||
							next.message == WM_KEYUP ||
							next.message == WM_SYSKEYUP)
						{
							if (next.wParam == VK_MENU &&
								(HIWORD(next.lParam) & KF_EXTENDED) &&
								next.time == time)
							{
								// Next message is Right Alt down so discard this
								break;
							}
						}
					}
				*/
				// This is a regular Left Ctrl message
				key = KeyLeftControl
			}
		}

		if action == GLFW_RELEASE && wParam == VK_SHIFT {
			// HACK: Release both Shift keys on Shift up event, as when both
			//       are pressed the first release does not emit any event
			// NOTE: The other half of this is in _glfwPlatformPollEvents
			glfwInputKey(window, KeyLeftShift, scancode, action, mods)
			glfwInputKey(window, KeyRightShift, scancode, action, mods)
		} else if wParam == VK_SNAPSHOT {
			// HACK: Key down is not reported for the Print Screen key
			glfwInputKey(window, key, scancode, GLFW_PRESS, mods)
			glfwInputKey(window, key, scancode, GLFW_RELEASE, mods)
		} else {
			glfwInputKey(window, key, scancode, action, mods)
		}
		break

	case WM_LBUTTONDOWN, WM_LBUTTONUP, WM_RBUTTONDOWN, WM_RBUTTONUP, WM_MBUTTONDOWN, WM_MBUTTONUP:
		var button MouseButton
		if msg == WM_LBUTTONDOWN || msg == WM_LBUTTONUP {
			button = MouseButtonLeft
		} else if msg == WM_RBUTTONDOWN || msg == WM_RBUTTONUP {
			button = MouseButtonRight
		} else if msg == WM_MBUTTONDOWN || msg == WM_MBUTTONUP {
			button = MouseButtonMiddle
		}
		var action Action
		if msg == WM_LBUTTONDOWN || msg == WM_RBUTTONDOWN || msg == WM_MBUTTONDOWN {
			action = GLFW_PRESS
		} else {
			action = GLFW_RELEASE
		}
		var i MouseButton
		for i = MouseButtonFirst; i <= MouseButtonLast; i++ {
			if window.mouseButtons[i] == GLFW_PRESS {
				break
			}
		}
		// if i > MouseButtonLast {
		// TODO SetCapture(hWnd);
		// }

		glfwInputMouseClick(window, button, action, getKeyMods())
		for i = MouseButtonFirst; i <= MouseButtonLast; i++ {
			if window.mouseButtons[i] == GLFW_PRESS {
				break
			}
		}
		// if (i > MouseButtonLast)
		// TODO ReleaseCapture();
		// }

		return 0

	// TODO case WM_CANCELMODE:

	case WM_SETFOCUS:
		glfwInputWindowFocus(window, true)
		// HACK: Do not disable cursor while the user is interacting with a caption button
		// TODO if (window.Win32.frameAction) break;
		// TODO if (window.cursorMode == GLFW_CURSOR_DISABLED)	disableCursor(window);
		return 0
	case WM_KILLFOCUS:
		// TODO if (window.cursorMode == GLFW_CURSOR_DISABLED) enableCursor(window);
		// TODO if (window.monitor && window.autoIconify) _glfwPlatformIconifyWindow(window);
		glfwInputWindowFocus(window, false)
		return 0

	case WM_MOUSEMOVE:
		x := float64(int(lParam & 0xffff))
		y := float64(int((lParam >> 16) & 0xffff))
		if !window.Win32.cursorTracked {
			// tme.dwFlags = TME_LEAVE;
			// tme.hwndTrack = window.hMonitor;
			// TrackMouseEvent(&tme);
			// window.cursorTracked = true;
			// glfwInputCursorEnter(window, GLFW_TRUE);
		}

		if window.cursorMode == GLFW_CURSOR_DISABLED {
			dx := float64(x) - window.lastCursorPosX
			dy := float64(y) - window.lastCursorPosY
			// TODO if _glfw.Win32.disabledCursorWindow != window {			break			}
			glfwInputCursorPos(window, window.virtualCursorPosX+dx, window.virtualCursorPosY+dy)
		} else {
			glfwInputCursorPos(window, x, y)
		}
		window.lastCursorPosX = x
		window.lastCursorPosY = y
		return 0

	case WM_MOUSEWHEEL:
		glfwInputScroll(window, 0.0, float64(int16(wParam>>16))/120.0)
		return 0

	case WM_MOUSEHWHEEL:
		glfwInputScroll(window, -float64(int16(wParam>>16))/120.0, 0.0)
		return 0

	case WM_PAINT:
		glfwInputWindowDamage(window)

	case WM_SIZE:
		width := int(lParam & 0xFFFF)
		height := int(lParam >> 16)
		iconified := wParam == SIZE_MINIMIZED
		maximized := wParam == SIZE_MAXIMIZED || (window.Win32.maximized && wParam != SIZE_RESTORED)
		// if (_glfw.win32.capturedCursorWindow == window) {
		//	captureCursor(window)
		// }

		if window.Win32.iconified != iconified {
			// TODO _glfwInputWindowIconify(window, iconified)
		}

		if window.Win32.maximized != maximized {
			// TODO _glfwInputWindowMaximize(window, maximized);
		}

		if width != window.Win32.width || height != window.Win32.height {
			window.Win32.width = width
			window.Win32.height = height
			// TODO _glfwInputFramebufferSize(window, width, height);
			if window.sizeCallback != nil {
				window.sizeCallback(window, width, height)
			}
		}
		if window.monitor != nil && window.Win32.iconified != iconified {
			if iconified {
				// TODO releaseMonitor(window);
			} else {
				// TODO acquireMonitor(window);
				// TODO fitToMonitor(window);
			}
		}
		window.Win32.iconified = iconified
		window.Win32.maximized = maximized
		return 0

	case WM_GETMINMAXINFO:
		// TODO

	case WM_SETCURSOR:
		// TODO
	}

	r1, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return r1
}

func glfwPlatformPollEvents() {
	var msg Msg
	for PeekMessage(&msg, 0, 0, 0, PM_REMOVE) {
		if msg.Message == WM_QUIT {
			// NOTE: While GLFW does not itself post WM_QUIT, other processes may post it to this one, for example Task Manager
			// HACK: Treat WM_QUIT as a close on all windows
			window := _glfw.windowListHead
			for window != nil {
				glfwInputWindowCloseRequest(window)
				window = window.next
			}
		} else {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}
	}

	// HACK: Release modifier keys that the system did not emit KEYUP for
	// NOTE: Shift keys on Windows tend to "stick" when both are pressed as no key up message is generated by the first key release
	// NOTE: Windows key is not reported as released by the Win+V hotkey. Other Win hotkeys are handled implicitly by _glfwInputWindowFocus
	//       because they change the input focus
	// NOTE: The other half of this is in the WM_*KEY* handler in windowProc
	/* TODO
	hMonitor = GetActiveWindow();
	if (hMonitor!=nil) {
		window := 74W(hMonitor, "GLFW");
		if window != nil {
			//const keys[4][2] = {{ VK_LSHIFT, GLFW_KEY_LEFT_SHIFT },    { VK_RSHIFT, GLFW_KEY_RIGHT_SHIFT },{ VK_LWIN, GLFW_KEY_LEFT_SUPER },{ VK_RWIN, GLFW_KEY_RIGHT_SUPER }}
			for i := 0; i < 4; i++ {
				vk := keys[i][0];
				key := keys[i][1];
				// scancode := Win32.scancodes[key];
				if GetKeyState(vk) & 0x8000 || window.keys[key] != GLFW_PRESS {
					continue;
				}
				_glfwInputKey(window, key, scancode, GLFW_RELEASE, getKeyMods());
			}
		}
	}
	window := _glfw.Win32.disabledCursorWindow;
	if window!=nil {
		var width, height int
		glfwPlatformGetWindowSize(window, &width, &height);
		// NOTE: Re-center the cursor only if it has moved since the last call, to avoid breaking glfwWaitEvents with WM_MOUSEMOVE
		if window.lastCursorPosX != width / 2 || window.lastCursorPosY != height / 2 {
			glfwPlatformSetCursorPos(window, width / 2, height / 2);
		}
	}
	*/
}

func glfwSwapBuffers(window *_GLFWwindow) {
	window.context.swapBuffers(window)
}

func cursorInContentArea(w *_GLFWwindow) bool {
	var x, y, width, height int
	glfwGetCursorPos(w, &x, &y)
	glfwGetWindowSize(w, &width, &height)
	return x >= 0 && y >= 0 && x < width && y < height // PtInRect(&area, pos);
}

func SetCursorWin32(handle HANDLE) {
	_, _, err := _SetCursor.Call(uintptr(handle))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("_SetCursor failed, " + err.Error())
	}
}

func updateCursorImage(window *_GLFWwindow) {

}

func glfwSetCursor(window *_GLFWwindow, cursor *Cursor) {
	window.cursor = cursor
	if cursorInContentArea(window) {
		if window.cursorMode == GLFW_CURSOR_NORMAL || window.cursorMode == GLFW_CURSOR_CAPTURED {
			if window.cursor != nil {
				SetCursorWin32(window.cursor.handle)
			} else {
				SetCursorWin32(LoadCursor(IDC_ARROW))
			}
		} else {
			// NOTE: Via Remote Desktop, setting the cursor to NULL does not hide it.
			// HACK: When running locally, it is set to NULL, but when connected via Remote
			//       Desktop, this is a transparent cursor.
			SetCursorWin32(_glfw.win32.blankCursor)
		}
	}
}

func SetFocus(window *_GLFWwindow) {
	r1, _, err := _SetFocus.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if r1 == 0 || err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetFocus failed, " + err.Error())
	}
	if r1 == 0 {
		panic("SetFocus failed")
	}
}

func BringWindowToTop(window *_GLFWwindow) {
	r1, _, err := _BringWindowToTop.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if r1 == 0 || err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("BringWindowToTop failed, " + err.Error())
	}
	if r1 == 0 {
		panic("BringWindowToTop failed")
	}
}

func SetForegroundWindow(window *_GLFWwindow) {
	r1, _, err := _SetForegroundWindow.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if r1 == 0 || err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetForegroundWindow failed, " + err.Error())
	}
	if r1 == 0 {
		panic("SetForegroundWindow failed")
	}
}

func glfwFocusWindow(window *_GLFWwindow) {
	BringWindowToTop(window)
	SetForegroundWindow(window)
	SetFocus(window)
}

type POINTL = struct {
	X, Y int32
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

const (
	ENUM_CURRENT_SETTINGS      = -1
	HORZSIZE                   = 4
	VERTSIZE                   = 6
	HORZRES                    = 8
	VERTRES                    = 10
	DISPLAY_DEVICE_MODESPRUNED = 0x08000000
	DISPLAY_DEVICE_REMOTE      = 0x04000000
	DISPLAY_DEVICE_DISCONNECT  = 0x02000000
)

func EnumDisplaySettings(name *uint16, mode int, dm *DEVMODEW) {
	ret, _, err := _EnumDisplaySettings.Call(uintptr(unsafe.Pointer(name)), uintptr(mode), uintptr(unsafe.Pointer(dm)))
	if ret == 0 || !errors.Is(err, syscall.Errno(0)) {
		panic("EnumDisplySetting failed, " + err.Error())
	}
}

func createMonitor(adapter *DISPLAY_DEVICEW, display *DISPLAY_DEVICEW) *Monitor {
	var monitor Monitor
	var widthMM, heightMM int
	var rect RECT
	var dm DEVMODEW

	dm.dmSize = uint16(unsafe.Sizeof(dm))
	EnumDisplaySettings(&adapter.DeviceName[0], ENUM_CURRENT_SETTINGS, &dm)
	pName, _ := syscall.UTF16PtrFromString("DISPLAY")
	ret, _, err := _CreateDC.Call(uintptr(unsafe.Pointer(pName)), uintptr(unsafe.Pointer(&adapter.DeviceName)), 0, 0)
	if !errors.Is(err, syscall.Errno(0)) {
		panic("CreateDC failed, " + err.Error())
	}
	dc := HDC(ret)
	if isWindows8Point1OrGreater() {
		widthMM = GetDeviceCaps(dc, HORZSIZE)
		heightMM = GetDeviceCaps(dc, VERTSIZE)
	} else {
		widthMM = int(float64(dm.dmPelsWidth) * 25.4 / float64(GetDeviceCaps(dc, LOGPIXELSX)))
		heightMM = int(float64(dm.dmPelsHeight) * 25.4 / float64(GetDeviceCaps(dc, LOGPIXELSY)))
	}
	ret, _, err = _DeleteDC.Call(uintptr(dc))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("CreateDC failed, " + err.Error())
	}
	monitor.heightMM = heightMM
	monitor.widthMM = widthMM

	if adapter.StateFlags&DISPLAY_DEVICE_MODESPRUNED != 0 {
		monitor.modesPruned = true
	}
	// copy(monitor.adapterName, adapter.DeviceName)
	// WideCharToMultiByte(CP_UTF8, 0, adapter.DeviceName, -1, monitor.win32.publicAdapterName, sizeof(monitor.win32.publicAdapterName), NULL, NULL)
	// if display != nil {
	//	wcscpy(monitor.win32.displayName, display.DeviceName)
	//	WideCharToMultiByte(CP_UTF8, 0, display.DeviceName, -1, monitor.win32.publicDisplayName, sizeof(monitor.win32.publicDisplayName), NULL, NULL)
	// }
	rect.Left = dm.dmPosition.X
	rect.Top = dm.dmPosition.Y
	rect.Right = dm.dmPosition.X + dm.dmPelsWidth
	rect.Bottom = dm.dmPosition.Y + dm.dmPelsHeight
	EnumDisplayMonitors(0, &rect, NewEnumDisplayMonitorsCallback(enumMonitorCallback), uintptr(unsafe.Pointer(&monitor)))
	return &monitor
}

type DISPLAY_DEVICEW struct {
	cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

func EnumDisplayDevices(device uintptr, no int, adapter *DISPLAY_DEVICEW, flags uint32) bool {
	ret, _, err := _EnumDisplayDevices.Call(device, uintptr(no), uintptr(unsafe.Pointer(adapter)), uintptr(flags))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("EnumDisplayDevices failed")
	}
	return ret == 1
}

const DISPLAY_DEVICE_ACTIVE = 0x00000001
const DISPLAY_DEVICE_ATTACHED = 0x00000002
const DISPLAY_DEVICE_PRIMARY_DEVICE = 0x00000004

func glfwPollMonitorsWin32() {

	/* disconnectedCount := _glfw.monitorCount;
	if (disconnectedCount) {
		disconnected = _glfw_calloc(_glfw.monitorCount, sizeof(Monitor*));
		memcpy(disconnected, _glfw.monitors, _glfw.monitorCount * sizeof(Monitor*));
	} */
	// var disconnected []*Monitor = _glfw.monitors

	for adapterIndex := 0; adapterIndex < 1000; adapterIndex++ {
		var adapter DISPLAY_DEVICEW
		adapterType := _GLFW_INSERT_LAST
		adapter.cb = uint32(unsafe.Sizeof(adapter))
		EnumDisplayDevices(0, adapterIndex, &adapter, 0)

		if (adapter.StateFlags & DISPLAY_DEVICE_ACTIVE) == 0 {
			continue
		}

		if (adapter.StateFlags & DISPLAY_DEVICE_PRIMARY_DEVICE) != 0 {
			adapterType = _GLFW_INSERT_FIRST
		}
		for displayIndex := 0; ; displayIndex++ {
			var display DISPLAY_DEVICEW
			display.cb = uint32(unsafe.Sizeof(display))
			if !EnumDisplayDevices(uintptr(unsafe.Pointer(&adapter.DeviceName)), displayIndex, &display, 0) {
				break
			}

			if (display.StateFlags & DISPLAY_DEVICE_ACTIVE) == 0 {
				continue
			}
			monitor := createMonitor(&adapter, &display)
			if monitor == nil {
				return
			}

			_glfwInputMonitor(monitor, GLFW_CONNECTED, adapterType)
			adapterType = _GLFW_INSERT_LAST

			// HACK: If an active adapter does not have any display devices
			//       (as sometimes happens), add it directly as a monitor
			/*
				if displayIndex == 0 {
					for i := 0; i < disconnectedCount; i++ {
						if disconnected[i] && wcscmp(disconnected[i].win32.adapterName, adapter.DeviceName) == 0 {
							disconnected[i] = NULL
							break
						}
					}
				}
				if i < disconnectedCount {
					continue
				}

				monitor = createMonitor(&adapter, NULL)
				if monitor == nil {
					_glfw_free(disconnected)
					return
				}
			*/
			// _glfwInputMonitor(monitor, GLFW_CONNECTED, adapterType)
		}
		/*
			for i := 0; i < disconnectedCount; i++ {
				if disconnected[i] {
					_glfwInputMonitor(disconnected[i], GLFW_DISCONNECTED, 0)
				}
			}
		*/
	}
}

const (
	/* Printable keys */
	GLFW_KEY_SPACE         = 32
	GLFW_KEY_APOSTROPHE    = 39 /* ' */
	GLFW_KEY_COMMA         = 44 /* , */
	GLFW_KEY_MINUS         = 45 /* - */
	GLFW_KEY_PERIOD        = 46 /* . */
	GLFW_KEY_SLASH         = 47 /* / */
	GLFW_KEY_0             = 48
	GLFW_KEY_1             = 49
	GLFW_KEY_2             = 50
	GLFW_KEY_3             = 51
	GLFW_KEY_4             = 52
	GLFW_KEY_5             = 53
	GLFW_KEY_6             = 54
	GLFW_KEY_7             = 55
	GLFW_KEY_8             = 56
	GLFW_KEY_9             = 57
	GLFW_KEY_SEMICOLON     = 59 /* ; */
	GLFW_KEY_EQUAL         = 61 /* = */
	GLFW_KEY_A             = 65
	GLFW_KEY_B             = 66
	GLFW_KEY_C             = 67
	GLFW_KEY_D             = 68
	GLFW_KEY_E             = 69
	GLFW_KEY_F             = 70
	GLFW_KEY_G             = 71
	GLFW_KEY_H             = 72
	GLFW_KEY_I             = 73
	GLFW_KEY_J             = 74
	GLFW_KEY_K             = 75
	GLFW_KEY_L             = 76
	GLFW_KEY_M             = 77
	GLFW_KEY_N             = 78
	GLFW_KEY_O             = 79
	GLFW_KEY_P             = 80
	GLFW_KEY_Q             = 81
	GLFW_KEY_R             = 82
	GLFW_KEY_S             = 83
	GLFW_KEY_T             = 84
	GLFW_KEY_U             = 85
	GLFW_KEY_V             = 86
	GLFW_KEY_W             = 87
	GLFW_KEY_X             = 88
	GLFW_KEY_Y             = 89
	GLFW_KEY_Z             = 90
	GLFW_KEY_LEFT_BRACKET  = 91  /* [ */
	GLFW_KEY_BACKSLASH     = 92  /* \ */
	GLFW_KEY_RIGHT_BRACKET = 93  /* ] */
	GLFW_KEY_GRAVE_ACCENT  = 96  /* ` */
	GLFW_KEY_WORLD_1       = 161 /* non-US #1 */
	GLFW_KEY_WORLD_2       = 162 /* non-US #2 */

	/* Function keys */
	GLFW_KEY_ESCAPE        = 256
	GLFW_KEY_ENTER         = 257
	GLFW_KEY_TAB           = 258
	GLFW_KEY_BACKSPACE     = 259
	GLFW_KEY_INSERT        = 260
	GLFW_KEY_DELETE        = 261
	GLFW_KEY_RIGHT         = 262
	GLFW_KEY_LEFT          = 263
	GLFW_KEY_DOWN          = 264
	GLFW_KEY_UP            = 265
	GLFW_KEY_PAGE_UP       = 266
	GLFW_KEY_PAGE_DOWN     = 267
	GLFW_KEY_HOME          = 268
	GLFW_KEY_END           = 269
	GLFW_KEY_CAPS_LOCK     = 280
	GLFW_KEY_SCROLL_LOCK   = 281
	GLFW_KEY_NUM_LOCK      = 282
	GLFW_KEY_PRINT_SCREEN  = 283
	GLFW_KEY_PAUSE         = 284
	GLFW_KEY_F1            = 290
	GLFW_KEY_F2            = 291
	GLFW_KEY_F3            = 292
	GLFW_KEY_F4            = 293
	GLFW_KEY_F5            = 294
	GLFW_KEY_F6            = 295
	GLFW_KEY_F7            = 296
	GLFW_KEY_F8            = 297
	GLFW_KEY_F9            = 298
	GLFW_KEY_F10           = 299
	GLFW_KEY_F11           = 300
	GLFW_KEY_F12           = 301
	GLFW_KEY_KP_0          = 320
	GLFW_KEY_KP_1          = 321
	GLFW_KEY_KP_2          = 322
	GLFW_KEY_KP_3          = 323
	GLFW_KEY_KP_4          = 324
	GLFW_KEY_KP_5          = 325
	GLFW_KEY_KP_6          = 326
	GLFW_KEY_KP_7          = 327
	GLFW_KEY_KP_8          = 328
	GLFW_KEY_KP_9          = 329
	GLFW_KEY_KP_DECIMAL    = 330
	GLFW_KEY_KP_DIVIDE     = 331
	GLFW_KEY_KP_MULTIPLY   = 332
	GLFW_KEY_KP_SUBTRACT   = 333
	GLFW_KEY_KP_ADD        = 334
	GLFW_KEY_KP_ENTER      = 335
	GLFW_KEY_KP_EQUAL      = 336
	GLFW_KEY_LEFT_SHIFT    = 340
	GLFW_KEY_LEFT_CONTROL  = 341
	GLFW_KEY_LEFT_ALT      = 342
	GLFW_KEY_LEFT_SUPER    = 343
	GLFW_KEY_RIGHT_SHIFT   = 344
	GLFW_KEY_RIGHT_CONTROL = 345
	GLFW_KEY_RIGHT_ALT     = 346
	GLFW_KEY_RIGHT_SUPER   = 347
	GLFW_KEY_MENU          = 348
)

// func ToUnicode(vk uint32, scancode uint32, state *[512]byte , chars, len, 0) {
// r1,_,err := _ToUnicode.Call(uintptr(vk), uintptr(scancode), uintptr(state), uintptr(chars), size)
// }

// TODO :Updates key names according to the current keyboard layout
func glfwUpdateKeyNamesWin32() {
	for key := GLFW_KEY_SPACE; key <= GLFW_KEY_MENU; key++ {
		/* TODO: Make readable key names
		scancode := _glfw.win32.scancodes[key]
		var vk uint16
		if scancode == -1 {
			continue
		}
		if key >= GLFW_KEY_KP_0 && key <= GLFW_KEY_KP_ADD {
			vks := []uint16{VK_NUMPAD0, VK_NUMPAD1, VK_NUMPAD2, VK_NUMPAD3, VK_NUMPAD4, VK_NUMPAD5, VK_NUMPAD6, VK_NUMPAD7, VK_NUMPAD8, VK_NUMPAD9, VK_DECIMAL, VK_DIVIDE, VK_MULTIPLY, VK_SUBTRACT, VK_ADD}
			vk = vks[key-GLFW_KEY_KP_0]
		} else {
			r1, _, err := _MapVirtualKeyW.Call(uintptr(scancode), uintptr(MAPVK_VSC_TO_VK))
			if !errors.Is(err, syscall.Errno(0)) {
				panic("MapVirtualKeyW failed, " + err.Error())
			}
			vk = uint16(r1)
		}
		var state [256]uint8
		var vk uint16
		length := ToUnicode(vk, scancode, state, chars, sizeof(chars)/sizeof(WCHAR), 0);
		if length == -1 {
			// This is a dead key, so we need a second simulated key press
			// to make it output its own character (usually a diacritic)
			length = ToUnicode(vk, scancode, state, chars, sizeof(chars)/sizeof(WCHAR), 0);
		}

		if (length < 1) {
			continue;
		}
		WideCharToMultiByte(CP_UTF8, 0, chars, 1, _glfw.win32.keynames[key], sizeof(_glfw.win32.keynames[key]), NULL, NULL);
		*/
	}
}

// Notifies shared code of a monitor connection or disconnection
func _glfwInputMonitor(monitor *Monitor, action int, placement int) {
	if action == GLFW_CONNECTED {
		_glfw.monitorCount++
		if placement == _GLFW_INSERT_FIRST {
			_glfw.monitors = append([]*Monitor{monitor}, _glfw.monitors...)
		} else {
			_glfw.monitors = append(_glfw.monitors, monitor)
		}
	} else if action == GLFW_DISCONNECTED {
		for window := _glfw.windowListHead; window != nil; window = window.next {
			if window.monitor == monitor {
				// TODO var width, height, xoff, yoff int
				// glfwGetWindowSize(window, &width, &height);
				// _glfw.platform.setWindowMonitor(window, NULL, 0, 0, width, height, 0);
				// _glfw.platform.getWindowFrameSize(window, &xoff, &yoff, NULL, NULL);
				// _glfw.platform.setWindowPos(window, xoff, yoff);
			}
		}
		for i := 0; i < _glfw.monitorCount; i++ {
			if _glfw.monitors[i] == monitor {
				_glfw.monitors = append(_glfw.monitors[:i], _glfw.monitors[i+1:]...)
				_glfw.monitorCount--
				break
			}
		}
	}

	if _glfw.monitorCallback != nil {
		_glfw.monitorCallback(monitor, action)
	}

}
