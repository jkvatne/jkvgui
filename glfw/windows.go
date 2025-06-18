package glfw

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type Point struct {
	X, Y int32
}

type Msg struct {
	Hwnd     syscall.Handle
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}

const (
	TRUE                        = 1
	PM_REMOVE                   = 0x0001
	PM_NOREMOVE                 = 0x0000
	WM_CANCELMODE               = 0x001F
	WM_CHAR                     = 0x0102
	WM_SYSCHAR                  = 0x0106
	WM_CLOSE                    = 0x0010
	WM_CREATE                   = 0x0001
	WM_DPICHANGED               = 0x02E0
	WM_DESTROY                  = 0x0002
	WM_ERASEBKGND               = 0x0014
	WM_GETMINMAXINFO            = 0x0024
	WM_IME_COMPOSITION          = 0x010F
	WM_IME_ENDCOMPOSITION       = 0x010E
	WM_IME_STARTCOMPOSITION     = 0x010D
	WM_KEYDOWN                  = 0x0100
	WM_KEYUP                    = 0x0101
	WM_KILLFOCUS                = 0x0008
	WM_LBUTTONDOWN              = 0x0201
	WM_LBUTTONUP                = 0x0202
	WM_MBUTTONDOWN              = 0x0207
	WM_MBUTTONUP                = 0x0208
	WM_MOUSEMOVE                = 0x0200
	WM_MOUSEWHEEL               = 0x020A
	WM_MOUSEHWHEEL              = 0x020E
	WM_NCACTIVATE               = 0x0086
	WM_NCHITTEST                = 0x0084
	WM_NCCALCSIZE               = 0x0083
	WM_PAINT                    = 0x000F
	WM_QUIT                     = 0x0012
	WM_SETCURSOR                = 0x0020
	WM_SETFOCUS                 = 0x0007
	WM_SHOWWINDOW               = 0x0018
	WM_SIZE                     = 0x0005
	WM_STYLECHANGED             = 0x007D
	WM_SYSKEYDOWN               = 0x0104
	WM_SYSKEYUP                 = 0x0105
	WM_RBUTTONDOWN              = 0x0204
	WM_RBUTTONUP                = 0x0205
	WM_TIMER                    = 0x0113
	WM_UNICHAR                  = 0x0109
	WM_USER                     = 0x0400
	WM_WINDOWPOSCHANGED         = 0x0047
	UNICODE_NOCHAR              = 65535
	CW_USEDEFAULT               = -2147483648
	WS_CLIPCHILDREN             = 0x02000000
	WS_CLIPSIBLINGS             = 0x04000000
	WS_MAXIMIZE                 = 0x01000000
	WS_ICONIC                   = 0x20000000
	WS_VISIBLE                  = 0x10000000
	WS_OVERLAPPED               = 0x00000000
	WS_CAPTION                  = 0x00C00000
	WS_SYSMENU                  = 0x00080000
	WS_THICKFRAME               = 0x00040000
	WS_MINIMIZEBOX              = 0x00020000
	WS_MAXIMIZEBOX              = 0x00010000
	WS_OVERLAPPEDWINDOW         = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	WS_EX_APPWINDOW             = 0x40000
	CursorNormal            int = 0x00034001
	CursorHidden            int = 0x00034002
	CursorDisabled          int = 0x00034003
	IDC_APPSTARTING             = 32650 // Standard arrow and small hourglass
	IDC_ARROW                   = 32512 // Standard arrow
	IDC_CROSS                   = 32515 // Crosshair
	IDC_HAND                    = 32649 // Hand
	IDC_HELP                    = 32651 // Arrow and question mark
	IDC_IBEAM                   = 32513 // I-beam
	IDC_NO                      = 32648 // Slashed circle
	IDC_SIZEALL                 = 32646 // Four-pointed arrow pointing north, south, east, and west
	IDC_SIZENESW                = 32643 // Double-pointed arrow pointing northeast and southwest
	IDC_SIZENS                  = 32645 // Double-pointed arrow pointing north and south
	IDC_SIZENWSE                = 32642 // Double-pointed arrow pointing northwest and southeast
	IDC_SIZEWE                  = 32644 // Double-pointed arrow pointing west and east
	IDC_UPARROW                 = 32516 // Vertical arrow
	IDC_WAIT                    = 32514 // Hour

	VK_CONTROL  = 0x11
	VK_LWIN     = 0x5B
	VK_MENU     = 0x12
	VK_RWIN     = 0x5C
	VK_SHIFT    = 0x10
	VK_SNAPSHOT = 0x2C
	VK_CAPITAL  = 0x14
	VK_NUMLOCK  = 0x90
	VK_LSHIFT   = 0xA0
	VK_RSHIFT   = 0xA1
	VK_BACK     = 0x08
	VK_DELETE   = 0x2e
	VK_DOWN     = 0x28
	VK_END      = 0x23
	VK_ESCAPE   = 0x1b
	VK_HOME     = 0x24
	VK_LEFT     = 0x25
	VK_NEXT     = 0x22
	VK_PRIOR    = 0x21
	VK_RIGHT    = 0x27
	VK_RETURN   = 0x0d
	VK_SPACE    = 0x20
	VK_TAB      = 0x09
	VK_UP       = 0x26

	VK_F1  = 0x70
	VK_F2  = 0x71
	VK_F3  = 0x72
	VK_F4  = 0x73
	VK_F5  = 0x74
	VK_F6  = 0x75
	VK_F7  = 0x76
	VK_F8  = 0x77
	VK_F9  = 0x78
	VK_F10 = 0x79
	VK_F11 = 0x7A
	VK_F12 = 0x7B

	VK_OEM_1      = 0xba
	VK_OEM_PLUS   = 0xbb
	VK_OEM_COMMA  = 0xbc
	VK_OEM_MINUS  = 0xbd
	VK_OEM_PERIOD = 0xbe
	VK_OEM_2      = 0xbf
	VK_OEM_3      = 0xc0
	VK_OEM_4      = 0xdb
	VK_OEM_5      = 0xdc
	VK_OEM_6      = 0xdd
	VK_OEM_7      = 0xde
	VK_OEM_102    = 0xe2
)

var (
	kernel32                = windows.NewLazySystemDLL("kernel32.dll")
	_procGetModuleHandleExW = kernel32.NewProc("GetModuleHandleExW")
	_GetModuleHandleW       = kernel32.NewProc("GetModuleHandleW")
	_GlobalAlloc            = kernel32.NewProc("GlobalAlloc")
	_GlobalFree             = kernel32.NewProc("GlobalFree")
	_GlobalLock             = kernel32.NewProc("GlobalLock")
	_GlobalUnlock           = kernel32.NewProc("GlobalUnlock")

	user32                       = windows.NewLazySystemDLL("user32.dll")
	enumDisplayMonitors          = user32.NewProc("EnumDisplayMonitors")
	getMonitorInfo               = user32.NewProc("GetMonitorInfo")
	_AdjustWindowRectEx          = user32.NewProc("AdjustWindowRectEx")
	_CallMsgFilter               = user32.NewProc("CallMsgFilterW")
	_CloseClipboard              = user32.NewProc("CloseClipboard")
	_CreateWindowEx              = user32.NewProc("CreateWindowExW")
	_DefWindowProc               = user32.NewProc("DefWindowProcW")
	_DestroyWindow               = user32.NewProc("DestroyWindow")
	_DispatchMessage             = user32.NewProc("DispatchMessageW")
	_EmptyClipboard              = user32.NewProc("EmptyClipboard")
	_GetWindowRect               = user32.NewProc("GetWindowRect")
	_GetClientRect               = user32.NewProc("GetClientRect")
	_GetClipboardData            = user32.NewProc("GetClipboardData")
	_GetDC                       = user32.NewProc("GetDC")
	_GetDpiForWindow             = user32.NewProc("GetDpiForWindow")
	_GetKeyState                 = user32.NewProc("GetKeyState")
	_GetMessage                  = user32.NewProc("GetMessageW")
	_GetMessageTime              = user32.NewProc("GetMessageTime")
	_GetMonitorInfo              = user32.NewProc("GetMonitorInfoW")
	_GetSystemMetrics            = user32.NewProc("GetSystemMetrics")
	_GetWindowLong               = user32.NewProc("GetWindowLongPtrW")
	_GetWindowLong32             = user32.NewProc("GetWindowLongW")
	_GetWindowPlacement          = user32.NewProc("GetWindowPlacement")
	_KillTimer                   = user32.NewProc("KillTimer")
	_LoadCursor                  = user32.NewProc("LoadCursorW")
	_LoadImage                   = user32.NewProc("LoadImageW")
	_MonitorFromPoint            = user32.NewProc("MonitorFromPoint")
	_MonitorFromWindow           = user32.NewProc("MonitorFromWindow")
	_MoveWindow                  = user32.NewProc("MoveWindow")
	_MsgWaitForMultipleObjectsEx = user32.NewProc("MsgWaitForMultipleObjectsEx")
	_OpenClipboard               = user32.NewProc("OpenClipboard")
	_PeekMessage                 = user32.NewProc("PeekMessageW")
	_PostMessage                 = user32.NewProc("PostMessageW")
	_PostQuitMessage             = user32.NewProc("PostQuitMessage")
	_ReleaseCapture              = user32.NewProc("ReleaseCapture")
	_RegisterClassExW            = user32.NewProc("RegisterClassExW")
	_ReleaseDC                   = user32.NewProc("ReleaseDC")
	_ScreenToClient              = user32.NewProc("ScreenToClient")
	_ShowWindow                  = user32.NewProc("ShowWindow")
	_SetCapture                  = user32.NewProc("SetCapture")
	_SetCursor                   = user32.NewProc("SetCursor")
	_SetClipboardData            = user32.NewProc("SetClipboardData")
	_SetForegroundWindow         = user32.NewProc("SetForegroundWindow")
	_SetFocus                    = user32.NewProc("SetFocus")
	_SetProcessDPIAware          = user32.NewProc("SetProcessDPIAware")
	_SetTimer                    = user32.NewProc("SetTimer")
	_SetWindowLong               = user32.NewProc("SetWindowLongPtrW")
	_SetWindowLong32             = user32.NewProc("SetWindowLongW")
	_SetWindowPlacement          = user32.NewProc("SetWindowPlacement")
	_SetWindowPos                = user32.NewProc("SetWindowPos")
	_SetWindowText               = user32.NewProc("SetWindowTextW")
	_TranslateMessage            = user32.NewProc("TranslateMessage")
	_UnregisterClass             = user32.NewProc("UnregisterClassW")
	_UpdateWindow                = user32.NewProc("UpdateWindow")
	shcore                       = windows.NewLazySystemDLL("shcore")
	_GetDpiForMonitor            = shcore.NewProc("GetDpiForMonitor")
)

type WndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CnClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       syscall.Handle
}

func GetKeyState(nVirtKey int) uint16 {
	c, _, _ := _GetKeyState.Call(uintptr(nVirtKey))
	return uint16(c)
}
func GetModuleHandle() (syscall.Handle, error) {
	h, _, err := _GetModuleHandleW.Call(uintptr(0))
	if h == 0 {
		return 0, fmt.Errorf("GetModuleHandleW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func RegisterClassEx(cls *WndClassEx) (uint16, error) {
	a, _, err := _RegisterClassExW.Call(uintptr(unsafe.Pointer(cls)))
	if a == 0 {
		return 0, fmt.Errorf("RegisterClassExW failed: %v", err)
	}
	return uint16(a), nil
}

/*
func GetModuleHandle(modulename *uint16) (handle syscall.Handle, err error) {
	r0, _, e1 := syscall.SyscallN(_GetModuleHandleW.Addr(), 1, uintptr(unsafe.Pointer(modulename)), 0, 0)
	handle = syscall.Handle(r0)
	if handle == 0 {
		err = fmt.Errorf("GetModuleHandle error %v", e1)
	}
	return
}
*/
func LoadImage(hInst syscall.Handle, res uint32, typ uint32, cx, cy int, fuload uint32) (syscall.Handle, error) {
	h, _, err := _LoadImage.Call(uintptr(hInst), uintptr(res), uintptr(typ), uintptr(cx), uintptr(cy), uintptr(fuload))
	if h == 0 {
		return 0, fmt.Errorf("LoadImageW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func CreateWindowEx(dwExStyle uint32, lpClassName uint16, lpWindowName string, dwStyle uint32, x, y, w, h int32, hWndParent, hMenu, hInstance syscall.Handle, lpParam uintptr) (HANDLE, error) {
	wname, _ := syscall.UTF16PtrFromString(lpWindowName)
	hwnd, _, err := _CreateWindowEx.Call(
		uintptr(dwExStyle),
		uintptr(lpClassName),
		uintptr(unsafe.Pointer(wname)),
		uintptr(dwStyle),
		uintptr(x), uintptr(y),
		uintptr(w), uintptr(h),
		uintptr(hWndParent),
		uintptr(hMenu),
		uintptr(hInstance),
		uintptr(lpParam))
	if hwnd == 0 {
		return 0, fmt.Errorf("CreateWindowEx failed: %v", err)
	}
	return HANDLE(hwnd), nil
}

func PeekMessage(m *Msg, hwnd syscall.Handle, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) bool {
	r, _, _ := _PeekMessage.Call(uintptr(unsafe.Pointer(m)), uintptr(hwnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax), uintptr(wRemoveMsg))
	return r != 0
}

func TranslateMessage(m *Msg) {
	_TranslateMessage.Call(uintptr(unsafe.Pointer(m)))
}

func DispatchMessage(m *Msg) {
	_DispatchMessage.Call(uintptr(unsafe.Pointer(m)))
}

func glfwPollEvents() {
	var msg *Msg
	var window *_GLFWwindow
	for PeekMessage(msg, 0, 0, 0, PM_REMOVE) {
		if msg.Message == WM_QUIT {
			// NOTE: While GLFW does not itself post WM_QUIT, other processes
			//       may post it to this one, for example Task Manager
			// HACK: Treat WM_QUIT as a close on all windows
			window = _glfw.windowListHead
			for window != nil {
				// TODO _glfwInputWindowCloseRequest(window)
				window = window.next
			}
		} else {
			TranslateMessage(msg)
			DispatchMessage(msg)
		}
	}

	// HACK: Release modifier keys that the system did not emit KEYUP for
	// NOTE: Shift keys on Windows tend to "stick" when both are pressed as
	//       no key up message is generated by the first key release
	// NOTE: Windows key is not reported as released by the Win+V hotkey
	//       Other Win hotkeys are handled implicitly by _glfwInputWindowFocus
	//       because they change the input focus
	// NOTE: The other half of this is in the WM_*KEY* handler in windowProc
	/* TODO
	handle = GetActiveWindow()
	if handle != 0 {
		window = GetPropW(handle, "GLFW")
		if window != nil {
			keys := [4][2]int{{VK_LSHIFT, KeyLeftShift}, {VK_RSHIFT, KeyRightShift}, {VK_LWIN, KeyLeftSuper}, {VK_RWIN, KeyRightSuper}}
			for i := 0; i < 4; i++ {
				vk := keys[i][0]
				key := keys[i][1]
				scancode := _glfw.scancodes[key]
				if (GetKeyState(vk)&0x8000 != 0) || (window.keys[key] != GLFW_PRESS) {
					continue
				}
				_glfwInputKey(window, key, scancode, GLFW_RELEASE, getKeyMods())
			}
		}
	}
	window = _glfw.disabledCursorWindow
	if window != nil {
		var width, height int
		// TODO _glfwPlatformGetWindowSize(window, &width, &height);
		// NOTE: Re-center the cursor only if it has moved since the last call,
		//       to avoid breaking glfwWaitEvents with WM_MOUSEMOVE
		if window.Win32.lastCursorPosX != width/2 || window.Win32.lastCursorPosY != height/2 {
			// TODO _glfwPlatformSetCursorPos(window, width / 2, height / 2);
		}
	}*/
}
