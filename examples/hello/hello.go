package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/wid"
)

func main() {
	window := gpu.InitWindow(200, 100, "Hello world", 0, 2)
	sys.Initialize(window)
	for !window.ShouldClose() {
		sys.StartFrame(f32.White)
		wid.Label("Hello world!", nil)(wid.NewCtx())
		sys.EndFrame(50)
	}
}
