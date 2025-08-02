package glfw

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

func glfwIsValidContextConfig(ctxconfig *_GLFWctxconfig) error {
	if (ctxconfig.major < 1 || ctxconfig.minor < 0) ||
		(ctxconfig.major == 1 && ctxconfig.minor > 5) ||
		(ctxconfig.major == 2 && ctxconfig.minor > 1) ||
		(ctxconfig.major == 3 && ctxconfig.minor > 3) {
		return fmt.Errorf("Invalid OpenGL version %d.%d", ctxconfig.major, ctxconfig.minor)
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
	return nil
}

func glfwChooseFBConfig(desired *_GLFWfbconfig, alternatives []_GLFWfbconfig, count int) *_GLFWfbconfig {
	var i int
	var missing, leastMissing = INT_MAX, INT_MAX
	var colorDiff, leastColorDiff = INT_MAX, INT_MAX
	var extraDiff, leastExtraDiff = INT_MAX, INT_MAX
	var closest *_GLFWfbconfig

	for i = 0; i < count; i++ {
		current := &alternatives[i]
		// Count number of missing buffers
		missing = 0
		if desired.alphaBits > 0 && current.alphaBits == 0 {
			missing++
		}
		if desired.depthBits > 0 && current.depthBits == 0 {
			missing++
		}

		if desired.stencilBits > 0 && current.stencilBits == 0 {
			missing++
		}
		if desired.auxBuffers > 0 && current.auxBuffers < desired.auxBuffers {
			missing += desired.auxBuffers - current.auxBuffers
		}
		if desired.samples > 0 && current.samples == 0 {
			missing++
		}
		if desired.transparent != current.transparent {
			missing++
		}
		colorDiff = 0
		if desired.redBits != GLFW_DONT_CARE {
			colorDiff += (desired.redBits - current.redBits) * (desired.redBits - current.redBits)
		}
		if desired.greenBits != GLFW_DONT_CARE {
			colorDiff += (desired.greenBits - current.greenBits) * (desired.greenBits - current.greenBits)
		}
		if desired.blueBits != GLFW_DONT_CARE {
			colorDiff += (desired.blueBits - current.blueBits) * (desired.blueBits - current.blueBits)
		}

		// Calculate non-color channel size difference value
		extraDiff = 0
		if desired.alphaBits != GLFW_DONT_CARE {
			extraDiff += (desired.alphaBits - current.alphaBits) * (desired.alphaBits - current.alphaBits)
		}
		if desired.depthBits != GLFW_DONT_CARE {
			extraDiff += (desired.depthBits - current.depthBits) * (desired.depthBits - current.depthBits)
		}
		if desired.stencilBits != GLFW_DONT_CARE {
			extraDiff += (desired.stencilBits - current.stencilBits) * (desired.stencilBits - current.stencilBits)
		}
		if desired.accumRedBits != GLFW_DONT_CARE {
			extraDiff += (desired.accumRedBits - current.accumRedBits) * (desired.accumRedBits - current.accumRedBits)
		}
		if desired.accumGreenBits != GLFW_DONT_CARE {
			extraDiff += (desired.accumGreenBits - current.accumGreenBits) * (desired.accumGreenBits - current.accumGreenBits)
		}
		if desired.accumBlueBits != GLFW_DONT_CARE {
			extraDiff += (desired.accumBlueBits - current.accumBlueBits) * (desired.accumBlueBits - current.accumBlueBits)
		}
		if desired.accumAlphaBits != GLFW_DONT_CARE {
			extraDiff += (desired.accumAlphaBits - current.accumAlphaBits) * (desired.accumAlphaBits - current.accumAlphaBits)
		}
		if desired.samples != GLFW_DONT_CARE {
			extraDiff += (desired.samples - current.samples) * (desired.samples - current.samples)
		}
		if desired.sRGB && !current.sRGB {
			extraDiff++
		}

		// Figure out if the current one is better than the best one found so far
		// Least number of missing buffers is the most important heuristic,
		// then color buffer size match and lastly size match for other buffers
		if missing < leastMissing {
			closest = current
		} else if missing == leastMissing {
			if (colorDiff < leastColorDiff) || (colorDiff == leastColorDiff && extraDiff < leastExtraDiff) {
				closest = current
			}
		}

		if current == closest {
			leastMissing = missing
			leastColorDiff = colorDiff
			leastExtraDiff = extraDiff
		}
	}
	return closest
}

func _glfwRefreshContextAttribs(window *_GLFWwindow, ctxconfig *_GLFWctxconfig) error {
	window.context.source = ctxconfig.source
	window.context.client = GLFW_OPENGL_API
	previous := (*Window)(unsafe.Pointer(glfwPlatformGetTls(&_glfw.contextSlot)))
	_ = glfwMakeContextCurrent(window)
	if glfwPlatformGetTls(&_glfw.contextSlot) != uintptr(unsafe.Pointer(window)) {
		return fmt.Errorf("glfwRefrechContext got Tls slot error")
	}

	window.context.GetIntegerv = window.context.getProcAddress("glGetIntegerv")
	window.context.GetString = window.context.getProcAddress("glGetString")
	if window.context.GetIntegerv == 0 || window.context.GetString == 0 {
		return fmt.Errorf("_glfwRefreshContextAttribs: Entry point retrieval is broken")
	}
	r, _, err := syscall.SyscallN(window.context.GetString, uintptr(GL_VERSION))
	if !errors.Is(err, syscall.Errno(0)) {
		return fmt.Errorf("String retrieval is broken, " + err.Error())
	}
	version := GoStr((*uint8)(unsafe.Pointer(uintptr(r))))
	prefixes := []string{"OpenGL ES-CM ", "OpenGL ES-CL ", "OpenGL ES ", ""}
	for s := range prefixes {
		if strings.HasPrefix(version, prefixes[s]) {
			version = strings.TrimPrefix(version, prefixes[s])
		}
	}
	i := 0
	var v [3]int
	for _, ch := range version {
		if ch >= '0' && ch <= '9' {
			v[i] = v[i]*10 + int(ch) - int('0')
		} else {
			i++
			if i >= 3 {
				break
			}
		}
	}
	window.context.major = v[0]
	window.context.minor = v[1]
	window.context.revision = v[2]
	if window.context.major == 0 {
		return fmt.Errorf("No version found in OpenGL version string")
	}
	if window.context.major < ctxconfig.major || window.context.major == ctxconfig.major && window.context.minor < ctxconfig.minor {
		// The desired OpenGL version is greater than the actual version
		// This only happens if the machine lacks {GLX|WGL}_ARB_create_context
		// /and/ the user has requested an OpenGL version greater than 1.0
		if window.context.client == GLFW_OPENGL_API {
			return fmt.Errorf("Requested OpenGL version %d.%d, got version %d.%d", ctxconfig.major, ctxconfig.minor, window.context.major, window.context.minor)
		} else {
			return fmt.Errorf("Requested OpenGL ES version %d.%d, got version %d.%d", ctxconfig.major, ctxconfig.minor, window.context.major, window.context.minor)
		}
		// makeContextCurrentWGL(previous)
	}
	if window.context.major >= 3 {
		// OpenGL 3.0+ uses a different function for extension string retrieval
		// We cache it here instead of in glfwExtensionSupported mostly to alert
		// users as early as possible that their build may be broken
		window.context.GetStringi = window.context.getProcAddress("glGetStringi")
		if window.context.GetStringi == 0 {
			return fmt.Errorf("Entry point retrieval is broken for glGetStringi")
		}
	}
	if window.context.client == GLFW_OPENGL_API {
		// Read back context flags (OpenGL 3.0 and above)
		if window.context.major >= 3 {
			var flags int
			GetIntegerv(window, GL_CONTEXT_FLAGS, &flags)
			if flags&GL_CONTEXT_FLAG_FORWARD_COMPATIBLE_BIT != 0 {
				window.context.forward = true
			}
			if (flags & GL_CONTEXT_FLAG_DEBUG_BIT) != 0 {
				window.context.debug = true
			} else if ExtensionSupported("GL_ARB_debug_output") && ctxconfig.debug {
				// HACK: This is a workaround for older drivers (pre KHR_debug)
				//       not setting the debug bit in the context flags for
				//       debug contexts
				window.context.debug = true
			}
			if flags&GL_CONTEXT_FLAG_NO_ERROR_BIT_KHR != 0 {
				window.context.noerror = true
			}
		}
		// Read back OpenGL context profile (OpenGL 3.2 and above)
		if window.context.major >= 4 || window.context.major == 3 && window.context.minor >= 2 {
			var mask int
			GetIntegerv(window, GL_CONTEXT_PROFILE_MASK, &mask)
			if (mask & GL_CONTEXT_COMPATIBILITY_PROFILE_BIT) != 0 {
				window.context.profile = GLFW_OPENGL_COMPAT_PROFILE
			} else if (mask & GL_CONTEXT_CORE_PROFILE_BIT) != 0 {
				window.context.profile = GLFW_OPENGL_CORE_PROFILE
			} else if ExtensionSupported("GL_ARB_compatibility") {
				// HACK: This is a workaround for the compatibility profile bit
				//       not being set in the context flags if an OpenGL 3.2+
				//       context was created without having requested a specific
				//       version
				window.context.profile = GLFW_OPENGL_COMPAT_PROFILE
			}
		}
		// Read back robustness strategy
		if ExtensionSupported("GL_ARB_robustness") {
			// NOTE: We avoid using the context flags for detection, as they are
			//       only present from 3.0 while the extension applies from 1.1
			var strategy int
			GetIntegerv(window, GL_RESET_NOTIFICATION_STRATEGY_ARB, &strategy)
			if strategy == GL_LOSE_CONTEXT_ON_RESET_ARB {
				window.context.robustness = GLFW_LOSE_CONTEXT_ON_RESET
			} else if strategy == GL_NO_RESET_NOTIFICATION_ARB {
				window.context.robustness = GLFW_NO_RESET_NOTIFICATION
			}
		}
	} else {
		// Read back robustness strategy
		if ExtensionSupported("GL_EXT_robustness") {
			// NOTE: The values of these constants match those of the OpenGL ARB one, so we can reuse them here
			var strategy int
			GetIntegerv(window, GL_RESET_NOTIFICATION_STRATEGY_ARB, &strategy)
			if strategy == GL_LOSE_CONTEXT_ON_RESET_ARB {
				window.context.robustness = GLFW_LOSE_CONTEXT_ON_RESET
			} else if strategy == GL_NO_RESET_NOTIFICATION_ARB {
				window.context.robustness = GLFW_NO_RESET_NOTIFICATION
			}
		}
	}

	if ExtensionSupported("GL_KHR_context_flush_control") {
		var behavior int
		GetIntegerv(window, GL_CONTEXT_RELEASE_BEHAVIOR, &behavior)
		if behavior == 0 {
			window.context.release = GLFW_RELEASE_BEHAVIOR_NONE
		} else if behavior == GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH {
			window.context.release = GLFW_RELEASE_BEHAVIOR_FLUSH
		}
	}

	// Clearing the front buffer to black to avoid garbage pixels left over from
	// previous uses of our bit of VRAM
	glClear := window.context.getProcAddress("glClear")
	syscall.SyscallN(glClear, GL_COLOR_BUFFER_BIT)
	if window.doublebuffer {
		window.context.swapBuffers(window)
	}
	glfwMakeContextCurrent(previous)
	return nil
}

func GetIntegerv(window *Window, name int, value *int) {
	_, _, err := syscall.SyscallN(window.context.GetIntegerv, uintptr(name), uintptr(unsafe.Pointer(value)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetIntegerv failed, " + err.Error())
	}
}

func glfwMakeContextCurrent(window *_GLFWwindow) error {
	previous := (*Window)(unsafe.Pointer(glfwPlatformGetTls(&_glfw.contextSlot)))
	if window != nil && window.context.client == GLFW_NO_API {
		return fmt.Errorf("glfwMakeContextCurrent failed: Cannot make current with a window that has no OpenGL or OpenGL ES context")
	}
	if previous != nil && (window == nil || window.context.source != previous.context.source) {
		previous.context.makeCurrent(nil)
	}
	if window != nil {
		window.context.makeCurrent(window)
	}
	return nil
}

func glfwGetCurrentContext() *Window {
	p := glfwPlatformGetTls(&_glfw.contextSlot)
	return (*Window)(unsafe.Pointer(p))
}

func glfwSwapBuffers(window *_GLFWwindow) {
	if window == nil {
		panic("glfwSwapBuffers: window == nil")
	}
	window.context.swapBuffers(window)
}

func SwapInterval(interval int) {
	window := glfwGetCurrentContext()
	if window == nil {
		panic("glfwSwapInterval: window == nil")
	}
	window.context.swapInterval(interval)
}

func ExtensionSupported(extension string) bool {
	p := glfwPlatformGetTls(&_glfw.contextSlot)
	if p == 0 {
		return false
	}
	if extension == "" {
		return false
	}
	window := (*_GLFWwindow)(unsafe.Pointer(p))
	if window.context.major >= 3 {
		// Check if extension is in the modern OpenGL extensions string list
		// count := window.context.GetIntegerv(GL_NUM_EXTENSIONS)
		r, _, _ := syscall.SyscallN(window.context.GetIntegerv, uintptr(GL_NUM_EXTENSIONS))
		count := int(r)
		for i := 0; i < count; i++ {
			// en := window.context.GetStringi(GL_EXTENSIONS, i)
			r, _, _ := syscall.SyscallN(window.context.GetStringi, uintptr(GL_EXTENSIONS), uintptr(i))
			en := GoStr((*uint8)(unsafe.Pointer(r)))
			if en == extension {
				return true
			}
		}
	} else {
		// Check if extension is in the old style OpenGL extensions string
		// extensions := window.context.GetString(GL_EXTENSIONS)
		r, _, _ := syscall.SyscallN(window.context.GetStringi, uintptr(GL_EXTENSIONS))
		extensions := GoStr((*uint8)(unsafe.Pointer(r)))
		if strings.Contains(extensions, extension) {
			return true
		}
	}
	// Check if extension is in the platform-specific string
	return window.context.extensionSupported(extension)
}

func glfwGetProcAddress(procname string) *Window {
	p := glfwPlatformGetTls(&_glfw.contextSlot)
	window := (*Window)(unsafe.Pointer(p))
	if window == nil {
		panic("glfwGetProcAddress: window == nil")
	}
	p = window.context.getProcAddress(procname)
	return (*Window)(unsafe.Pointer(p))
}
