package glfw

import (
	"fmt"
	"github.com/jkvatne/jkvgui/gl"
	"golang.org/x/sys/windows"
	"log/slog"
	"strings"
	"syscall"
	"unsafe"
)

func swapBuffersWGL(window *_GLFWwindow) {
	if window.monitor != nil {
		if IsWindowsVistaOrGreater() {
			/*
				// DWM Composition is always enabled on Win8+
				enabled := IsWindows8OrGreater()
				// HACK: Use DwmFlush when desktop composition is enabled
					if enabled || (SUCCEEDED(DwmIsCompositionEnabled(&enabled)) && enabled) {
						for count := abs(window.context.wgl.interval); count > 0; count-- {
							DwmFlush()
						}
					}
			*/
		}
	}
	swapBuffers(window.context.wgl.dc)
}

func swapBuffers(dc HDC) {
	r, _, err := _glfw.wgl.wglSwapBuffers.Call(uintptr(dc))
	if err != nil {
		panic(err)
	}
	if r == 0 {
		err = syscall.GetLastError()
		panic(err)
	}
}

func createContext(dc HDC) syscall.Handle {
	r1, _, err := _glfw.wgl.wglCreateContext.Call(uintptr(dc))
	if err != nil {
		panic(err)
	}
	return syscall.Handle(r1)
}

func deleteContext(handle HANDLE) {
	_, _, err := _glfw.wgl.wglDeleteContext.Call(uintptr(handle))
	if err != nil {
		panic(err)
	}
}

func getCurrentDC() HDC {
	r1, _, err := _glfw.wgl.wglGetCurrentDC.Call()
	if err != nil {
		panic("getCurrentDC failed, " + err.Error())
	}
	return HDC(r1)
}

func getCurrentContex() HANDLE {
	r1, _, err := _glfw.wgl.wglCreateContext.Call()
	if err != nil {
		panic("getCurrentDC failed, " + err.Error())
	}
	return HANDLE(r1)
}

func makeCurrent(dc HDC, handle HANDLE) bool {
	r1, _, err := _glfw.wgl.wglMakeCurrent.Call(uintptr(dc), uintptr(handle))
	if err != nil {
		panic("makeCurrent failed, " + err.Error())
	}
	return r1 != 0
}

func shareLists(dc HDC, handle HANDLE) bool {
	r1, _, err := _glfw.wgl.wglShareLists.Call(uintptr(dc), uintptr(handle))
	if err != nil {
		panic("wglShareLists failed, " + err.Error())
	}
	return r1 != 0
}

func glfwMakeContextCurrent(window *_GLFWwindow) error {
	// _GLFWwindow* window = (_GLFWwindow*) handle;
	// previous := _glfwPlatformGetTls(&_glfw.contextSlot);
	if window != nil && window.context.client == GLFW_NO_API {
		return fmt.Errorf("Cannot make current with a window that has no OpenGL or OpenGL ES context")
	}
	// if previous!=nil && w, r1indow!=nil || window.context.source != previous.context.source)
	//		previous.context.makeCurrent(NULL);
	// }

	if window != nil {
		// window.context.makeCurrent(&window.context)
		// window.context.wgl.dc, window.context.wgl.handle
		r1, _, err := _glfw.wgl.wglMakeCurrent.Call(uintptr(unsafe.Pointer(window)))
		slog.Info("Make current returned", "Err", err, "R1", r1)
	}
	return nil
}

// Initialize WGL
func _glfwInitWGL() error {
	var pfd PIXELFORMATDESCRIPTOR
	var pdc uintptr
	if _glfw.wgl.instance != nil {
		return nil
	}

	_glfw.wgl.wglSwapBuffers = gdi32.NewProc("SwapBuffers")

	_glfw.wgl.instance = windows.NewLazySystemDLL("opengl32.dll")
	// getProcAddress := opengl32.NewProc("wglGetProcAddress")
	_glfw.wgl.wglCreateContext = opengl32.NewProc("wglCreateContext")
	_glfw.wgl.wglDeleteContext = opengl32.NewProc("wglDeleteContext")
	_glfw.wgl.wglGetProcAddress = opengl32.NewProc("wglGetProcAddress")
	_glfw.wgl.wglGetCurrentDC = opengl32.NewProc("wglGetCurrentDC")
	_glfw.wgl.wglGetCurrentContext = opengl32.NewProc("wglGetCurrentContext")
	_glfw.wgl.wglMakeCurrent = opengl32.NewProc("wglMakeCurrent")
	_glfw.wgl.wglShareLists = opengl32.NewProc("wglShareLists")

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

func GetDC(w syscall.Handle) HDC {
	return 0
}

func _glfwCreateContextWGL(window *_GLFWwindow, ctxconfig *_GLFWctxconfig, fbconfig *_GLFWfbconfig) error {
	// var attribs [40]int
	// var pixelFormat int
	// var  pfd PIXELFORMATDESCRIPTOR
	// if ctxconfig.share != nil {
	//	share = ctxconfig.share.context.wgl.handle
	// }
	window.context.wgl.dc = GetDC(window.Win32.handle)
	if window.context.wgl.dc == 0 {
		return fmt.Errorf("WGL: Failed to retrieve DC for window")
	}
	/*
		pixelFormat := choosePixelFormat(window, ctxconfig, fbconfig);
		if pixelFormat==0 {
			return fmt.Errorf("WGL: Failed to retrieve PixelFormat for window")
		}

		if (!DescribePixelFormat(window.context.wgl.dc,	pixelFormat, sizeof(pfd), &pfd)) {
			return fmt.Errorf("WGL: Failed to retrieve PFD for selected pixel format");
		}

		if SetPixelFormat(window.context.wgl.dc, pixelFormat, &pfd)==0	{
			return fmt.Errorf("WGL: Failed to set selected pixel format");
		}
		if ctxconfig.client == GLFW_OPENGL_API {
			if ctxconfig.forward && (_glfw.wgl.ARB_create_context==nil) {
				return fmt.Errorf("WGL: A forward compatible OpenGL context requested but WGL_ARB_create_context is unavailable");
			}
			if (ctxconfig.profile==0) && _glfw.wgl.ARB_create_context_profile!=nil	{
				return fmt.Errorf("WGL: OpenGL profile requested but WGL_ARB_create_context_profile is unavailable");
			}
		} else {
			if (!_glfw.wgl.ARB_create_context || !_glfw.wgl.ARB_create_context_profile || !_glfw.wgl.EXT_create_context_es2_profile) {
				return fmt.Errorf("WGL: OpenGL ES requested but WGL_ARB_create_context_es2_profile is unavailable");
			}
		}

		if _glfw.wgl.ARB_create_context!=nil {
			mask := 0
			flags := 0;
			if (ctxconfig.client == GLFW_OPENGL_API) {
				if (ctxconfig.forward) {
					flags |= WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB;
				}
				if (ctxconfig.profile == GLFW_OPENGL_CORE_PROFILE) {
					mask |= WGL_CONTEXT_CORE_PROFILE_BIT_ARB;
				} else if (ctxconfig.profile == GLFW_OPENGL_COMPAT_PROFILE) {
					mask |= WGL_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB;
				}
			} else {
				mask |= WGL_CONTEXT_ES2_PROFILE_BIT_EXT;
			}
			if ctxconfig.debug {
				flags |= WGL_CONTEXT_DEBUG_BIT_ARB;
			}
			if ctxconfig.robustness!=0 {
				if _glfw.wgl.ARB_create_context_robustness {
					if (ctxconfig.robustness == GLFW_NO_RESET_NOTIFICATION) {
						setAttrib(WGL_CONTEXT_RESET_NOTIFICATION_STRATEGY_ARB, WGL_NO_RESET_NOTIFICATION_ARB);
					}
				} else if (ctxconfig.robustness == GLFW_LOSE_CONTEXT_ON_RESET) {
					setAttrib(WGL_CONTEXT_RESET_NOTIFICATION_STRATEGY_ARB, WGL_LOSE_CONTEXT_ON_RESET_ARB);
				}
				flags |= WGL_CONTEXT_ROBUST_ACCESS_BIT_ARB;
			}

			if ctxconfig.release!=0 {
				if (_glfw.wgl.ARB_context_flush_control) {
					if (ctxconfig.release == GLFW_RELEASE_BEHAVIOR_NONE) {
						setAttrib(WGL_CONTEXT_RELEASE_BEHAVIOR_ARB, WGL_CONTEXT_RELEASE_BEHAVIOR_NONE_ARB);
					} else if ctxconfig.release == GLFW_RELEASE_BEHAVIOR_FLUSH {
						setAttrib(WGL_CONTEXT_RELEASE_BEHAVIOR_ARB,
							WGL_CONTEXT_RELEASE_BEHAVIOR_FLUSH_ARB);
					}
				}
			}
			if (ctxconfig.noerror) {
				if (_glfw.wgl.ARB_create_context_no_error) {
					setAttrib(WGL_CONTEXT_OPENGL_NO_ERROR_ARB, GLFW_TRUE);
				}
			}
			// NOTE: Only request an explicitly versioned context when necessary, as
			//       explicitly requesting version 1.0 does not always return the
			//       highest version supported by the driver
			if ctxconfig.major != 1 || ctxconfig.minor != 0 {
				setAttrib(WGL_CONTEXT_MAJOR_VERSION_ARB, ctxconfig.major);
				setAttrib(WGL_CONTEXT_MINOR_VERSION_ARB, ctxconfig.minor);
			}
			if (flags!=0) {
				setAttrib(WGL_CONTEXT_FLAGS_ARB, flags)
			}
			if (mask!=0) {
				setAttrib(WGL_CONTEXT_PROFILE_MASK_ARB, mask)
			}
			setAttrib(0, 0);
			window.context.wgl.handle = wglCreateContextAttribsARB(window.context.wgl.dc, share, attribs);
			if (window.context.wgl.handle==0) {
				err := GetLastError();
				if err == (0xc0070000 | ERROR_INVALID_VERSION_ARB) {
					if (ctxconfig.client == GLFW_OPENGL_API) {
						return fmt.Errorf("WGL: Driver does not support OpenGL version %i.%i", ctxconfig.major, ctxconfig.minor);
					} else {
						return fmt.Errorf("WGL: Driver does not support OpenGL ES version %i.%i", ctxconfig.major, ctxconfig.minor);
					}
				}
			} else if (err == (0xc0070000 | ERROR_INVALID_PROFILE_ARB)) {
				return fmt.Errorf("WGL: Driver does not support the requested OpenGL profile");
			} else if (err == (0xc0070000 | ERROR_INCOMPATIBLE_DEVICE_CONTEXTS_ARB)){
				return fmt.Errorf("WGL: The share context is not compatible with the requested context");
			} else {
				if (ctxconfig.client == GLFW_OPENGL_API) {
					return fmt.Errorf("WGL: Failed to create OpenGL context");
				}
			} else {
				return fmt.Errorf("WGL: Failed to create OpenGL ES context");
			}
			return nil

		} else {
			window.context.wgl.handle = wglCreateContext(window.context.wgl.dc);
			if window.context.wgl.handle == 0 {
				return fmt.Errorf("WGL: Failed to create OpenGL context");
			}
		}
		if share != 0 {
			if (!wglShareLists(share, window.context.wgl.handle)) {
				return fmt.Errorf("WGL: Failed to enable sharing with specified OpenGL context");
			}
		}
	*/
	window.context.makeCurrent = makeContextCurrentWGL
	window.context.swapBuffers = swapBuffersWGL
	window.context.swapInterval = swapIntervalWGL
	window.context.extensionSupported = extensionSupportedWGL
	// window.context.getProcAddress = getProcAddressWGL
	window.context.destroy = destroyContextWGL
	return nil
}

func _glfwPlatformSetTls(g *_GLFWtls, w *_GLFWwindow) {
}

func wglMakeCurrent(g *_GLFWtls, w *_GLFWwindow) bool {
	return false
}

func makeContextCurrentWGL(window *_GLFWwindow) error {
	if window != nil {
		if makeCurrent(window.context.wgl.dc, window.context.wgl.handle) {
			_glfwPlatformSetTls(&_glfw.contextSlot, window)
		} else {
			_glfwPlatformSetTls(&_glfw.contextSlot, nil)
			return fmt.Errorf("WGL: Failed to make context current")
		}
	} else {
		if !wglMakeCurrent(nil, nil) {
			return fmt.Errorf("WGL: Failed to clear current context")
		}
		_glfwPlatformSetTls(&_glfw.contextSlot, nil)
	}
	return nil
}

func IsWindowsVistaOrGreater() bool {
	return true
}

func _glfwPlatformGetTls(s *_GLFWtls) *_GLFWwindow {
	return nil
}

func IsWindows8OrGreater() bool {
	return true
}

func DwmIsCompositionEnabled(enabled *bool) bool {
	*enabled = true
	return true
}

func swapIntervalWGL(interval int) {
	window := _glfwPlatformGetTls(&_glfw.contextSlot)
	window.context.wgl.interval = interval
	if window.monitor != nil {
		if IsWindowsVistaOrGreater() {
			// DWM Composition is always enabled on Win8+
			enabled := IsWindows8OrGreater()
			// HACK: Disable WGL swap interval when desktop composition is enabled to
			//       avoid interfering with DWM vsync
			if enabled || (DwmIsCompositionEnabled(&enabled) && enabled) {
				interval = 0
			}
		}
	}
	/*
		if _glfw.wgl.EXT_swap_control {
			wglSwapIntervalEXT(interval)
		}*/
}

func extensionSupportedWGL(extension byte) bool {
	/*if (_glfw.wgl.GetExtensionsStringARB) {
		extensions = wglGetExtensionsStringARB(wglGetCurrentDC());
	} else if (_glfw.wgl.GetExtensionsStringEXT) {
		extensions = wglGetExtensionsStringEXT();
	}
	if (!extensions) {
		return false
	}
	return _glfwStringInExtensionString(extension, extensions);
	*/
	return false
}

/*
func getProcAddressWGL(procname string) GLFWglproc {
	proc := wglGetProcAddress(procname)
	if proc != nil {
		return proc
	}
	return syscall.GetProcAddress(_glfw.wgl.instance, procname)
}
*/
func destroyContextWGL(window *_GLFWwindow) {
	if window.context.wgl.handle != 0 {
		deleteContext(window.context.wgl.handle)
		window.context.wgl.handle = 0
	}
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

// GoStr takes a null-terminated string returned by OpenGL and constructs a
// corresponding Go string.
func GoStr(cstr *uint8) string {
	str := ""
	for {
		if *cstr == 0 {
			break
		}
		str += string(*cstr)
		cstr = (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(cstr)) + 1))
	}
	return str
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
	version := GoStr(gl.GetString(gl.VERSION))
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
