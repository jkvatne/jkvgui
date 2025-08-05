package main

import (
	"github.com/jkvatne/jkvgui/gl"
	"github.com/jkvatne/jkvgui/glfw"
	"os"
	"runtime"
)

// UTF-8 window title test

func main() {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	window, err := glfw.CreateWindow(800, 400, "English 日本語 русский язык 官話", nil, nil)
	window.SetPos(100, 100)
	if err != nil {
		glfw.Terminate()
		os.Exit(1)
	}
	window.MakeContextCurrent()
	gl.Init()
	glfw.SwapInterval(1)

	for !window.ShouldClose() {
		gl.ClearColor(100, 100, 0, 256)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		window.SwapBuffers()
		glfw.WaitEvents()
	}
	glfw.Terminate()
}
