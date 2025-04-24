package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/wid"
)

func main() {
	gpu.DebugWidgets = true
	window := gpu.InitWindow(200, 50, "Hello world", 0)
	sys.Initialize(window, 14)
	for !window.ShouldClose() {
		sys.StartFrame(f32.White)
		wid.Label("Hello world!", nil)(wid.NewCtx())
		sys.EndFrame(50)
	}
}
