// Window re-opener (open/close stress test)
package main

import (
	"fmt"
	"github.com/jkvatne/jkvgui/gl"
	"github.com/jkvatne/jkvgui/glfw"
	"os"
	"runtime"
	"unsafe"
)

// It closes and re-opens the GLFW window every five seconds, alternating
// between windowed and full screen mode
// It also times and logs opening and closing actions and attempts to separate
// user initiated window closing from its own

var vertex_shader_text = `#version 110
uniform mat4 MVP;
attribute vec2 vPos;
void main()
{
    gl_Position = MVP * vec4(vPos, 0.0, 1.0);	
}` + "\x00"

var fragment_shader_text = `#version 110
void main()
{
    gl_FragColor = vec4(1.0);
}` + "\x00"

var vertices = [4][2]float32{{-0.5, -0.5}, {0.5, -0.5}, {0.5, 0.5}, {-0.5, 0.5}}

func window_close_callback(window *glfw.Window) {
	fmt.Printf("Close callback triggered\n")
}

func key_callback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.GLFW_PRESS {
		return
	}
	switch key {
	case glfw.GLFW_KEY_Q:
	case glfw.GLFW_KEY_ESCAPE:
		window.SetWindowShouldClose(true)
	}
}

func close_window(window *glfw.Window) {
	base := glfw.GetTime()
	window.Destroy()
	fmt.Printf("Closing window took %0.3f seconds\n", glfw.GetTime()-base)
}

type Mat4 [16]float32

func mat4x4_ortho(left, right, bottom, top, near, far float32) Mat4 {
	var m Mat4
	// Implementation details would involve calculating the matrix elements
	// based on the orthographic projection formula.
	// This is a simplified representation.
	m[0] = 2 / (right - left)
	m[5] = 2 / (top - bottom)
	m[10] = -2 / (far - near)
	m[12] = -(right + left) / (right - left)
	m[13] = -(top + bottom) / (top - bottom)
	m[14] = -(far + near) / (far - near)
	m[15] = 1
	return m
}

func main() {
	runtime.LockOSThread()
	count := 0
	var monitor *glfw.Monitor
	err := glfw.Init()
	if err != nil {
		panic("glfwInit err: " + err.Error())
	}
	glfw.WindowHint(glfw.GLFW_CONTEXT_VERSION_MAJOR, 2)
	glfw.WindowHint(glfw.GLFW_CONTEXT_VERSION_MINOR, 0)
	for {
		monitor = nil
		if count&1 != 0 {
			// monitors := glfw.GetMonitors()
			// monitor = monitors[rand.Int()%len(monitors)]
		}
		// mode := glfw.GetVideoMode(monitor)
		width := 400  // mode.width
		height := 400 // mode.height
		base := glfw.GetTime()
		window, err := glfw.CreateWindow(width, height, "Window Re-opener", monitor, nil)
		if err != nil {
			glfw.Terminate()
			fmt.Printf("Could not create window: %v\n", err)
		}
		if monitor != nil {
			fmt.Printf("Opening full screen window on monitor %s took %0.3f seconds\n",
				monitor.GetMonitorName(),
				glfw.GetTime()-base)
		} else {
			fmt.Printf("Opening regular window took %0.3f seconds\n", glfw.GetTime()-base)
		}
		window.MakeContextCurrent()
		gl.Init()
		glfw.SwapInterval(1)

		gl.ClearColor(100, 100, 0, 256)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		// window.SwapBuffers()
		window.SetWindowCloseCallback(window_close_callback)
		window.SetKeyCallback(key_callback)
		window.MakeContextCurrent()
		gl.Init()
		glfw.SwapInterval(1)
		vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
		csources, free := gl.Strs(vertex_shader_text)
		gl.ShaderSource(vertex_shader, 1, csources, nil)
		free()
		gl.CompileShader(vertex_shader)

		fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
		csources, free = gl.Strs(vertex_shader_text)
		gl.ShaderSource(fragment_shader, 1, csources, nil)
		free()
		gl.CompileShader(fragment_shader)

		program := gl.CreateProgram()
		gl.AttachShader(program, vertex_shader)
		gl.AttachShader(program, fragment_shader)
		gl.LinkProgram(program)

		mvp_location := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
		vpos_location := gl.GetAttribLocation(program, gl.Str("vPos\x00"))
		var vertex_buffer uint32
		gl.GenBuffers(1, &vertex_buffer)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertex_buffer)
		gl.BufferData(gl.ARRAY_BUFFER, 2*4*4, unsafe.Pointer(&vertices[0][0]), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(uint32(vpos_location))
		gl.VertexAttribPointer(uint32(vpos_location), 2, gl.FLOAT, false, 2*4*4, nil)

		glfw.SetTime(0.0)
		for {
			t := glfw.GetTime()
			if t > 1.0 {
				break
			}
			var w, h int32
			glfw.GetFramebufferSize(window, &w, &h)
			// ratio := float32(w) / float32(h)
			gl.Viewport(0, 0, w, h)
			gl.Clear(gl.COLOR_BUFFER_BIT)
			// p := mat4x4_ortho(-ratio, ratio, -1, 1, 0, 1)
			m := Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
			// mat4x4_rotate_Z(m, m, glfw.GetTime())
			// mat4x4_mul(mvp, p, m);
			mvp := m
			gl.UseProgram(program)
			gl.UniformMatrix4fv(mvp_location, 1, false, &mvp[0])
			gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
			window.SwapBuffers()
			glfw.PollEvents()
			if window.ShouldClose() {
				close_window(window)
				fmt.Printf("User closed window\n")
				glfw.Terminate()
				os.Exit(0)
			}
		}
		fmt.Printf("Closing window\n")
		close_window(window)
		count++
	}
	glfw.Terminate()
}
