package glfw

import (
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
	// previous := glfwPlatformGetTls(&_glfw.contextSlot)
	_ = glfwMakeContextCurrent(window)
	if _glfwPlatformGetTls(&_glfw.contextSlot) != window {
		return nil
	}

	window.context.GetIntegerv = window.context.getProcAddress("glGetIntegerv")
	window.context.GetString = window.context.getProcAddress("glGetString")
	if window.context.GetIntegerv == 0 || window.context.GetString == 0 {
		return fmt.Errorf("_glfwRefreshContextAttribs: Entry point retrieval is broken")
	}
	window.context.major = 3
	window.context.minor = 3
	window.context.revision = 3
	if window.context.major == 0 {
		return fmt.Errorf("No version found in OpenGL version string")
	}
	if window.context.major < ctxconfig.major || window.context.major == ctxconfig.major && window.context.minor < ctxconfig.minor {
		// The desired OpenGL version is greater than the actual version
		// This only happens if the machine lacks {GLX|WGL}_ARB_create_context
		// /and/ the user has requested an OpenGL version greater than 1.0
		return fmt.Errorf("Requested OpenGL version %d.%d, got version %d.%d", ctxconfig.major, ctxconfig.minor, window.context.major, window.context.minor)
	}
	return nil
}

func glfwMakeContextCurrent(window *_GLFWwindow) error {
	// _GLFWwindow* window = (_GLFWwindow*) hMonitor;
	// previous := _glfwPlatformGetTls(&_glfw.contextSlot);
	// if previous!=nil && w, r1indow!=nil || window.context.source != previous.context.source)
	//		previous.context.makeCurrent(NULL);
	// }
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

func glfwSwapInterval(interval int) {
	window := glfwGetCurrentContext()
	if window == nil {
		panic("glfwSwapInterval: window == nil")
	}
	window.context.swapInterval(interval)
}

func glfwExtensionSupported(extension string) bool {
	window := _glfwPlatformGetTls(&_glfw.contextSlot)
	if window == nil {
		return false
	}
	if extension == "" {
		return false
	}
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
	window := _glfwPlatformGetTls(&_glfw.contextSlot)
	if window == nil {
		panic("glfwGetProcAddress: window == nil")
	}
	p := window.context.getProcAddress(procname)
	return (*Window)(unsafe.Pointer(p))
}
