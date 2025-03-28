package main

import (
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/wid"
)

func main() {
	window := gpu.InitWindow(150, 50, "Hello world", 0)
	callback.Initialize(window)
	for !window.ShouldClose() {
		gpu.StartFrame(f32.White)
		wid.Label("Hello world!", wid.H1C)(wid.NewCtx())
		gpu.EndFrame(50)
	}
}
