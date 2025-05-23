package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

func main() {
	// Initialize the GUI system and parse arguments
	sys.Initialize()
	// Create a window with a title and size.
	gpu.InitWindow(200, 100, "Hello world", 0, 2)
	// Initialize the window and the GUI system, including callbacks.
	sys.InitializeWindow()
	// Loop until the window is closed.
	for !gpu.ShouldClose() {
		sys.StartFrame(theme.Surface)
		// Show just a single widget and call it with a new Ctx.
		wid.Label("Hello world!", nil)(wid.NewCtx())
		// EndFrame will swap buffers and limit the maximum framerate.
		sys.EndFrame(50)
	}
}
