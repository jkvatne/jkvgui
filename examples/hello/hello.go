package main

import (
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

func main() {
	sys.Init()
	defer sys.Shutdown()
	// Create a window with a title and size.
	w := sys.CreateWindow(100, 100, 200, 100, "Hello world", 0, 2)

	// Loop until the window is closed.
	for sys.Running() {
		w.StartFrame(theme.Surface.Bg())
		// Show just a single widget and call it with a new Ctx.
		wid.Show(wid.Label("Hello world!", nil))
		// EndFrame do housekeeping and swap buffers
		w.EndFrame()
		sys.PollEvents()
	}
}
