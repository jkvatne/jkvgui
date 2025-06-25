package main

import (
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var posV [16]wid.ResizerState
var posH [16]wid.ResizerState
var image [16]*wid.Img

func main() {
	sys.CreateWindow(100, 100, 400, 400, "Resizing1", 1, 2)
	sys.CreateWindow(400, 400, 400, 400, "Resizing2", 2, 2)
	defer sys.Shutdown()
	image[0], _ = wid.NewImage("music.jpg")
	image[1], _ = wid.NewImage("ts.jpg")
	for sys.Running() {
		for sys.CurrentWno, _ = range sys.WindowList {
			sys.StartFrame(theme.Surface.Bg())
			ctx := wid.NewCtx()
			wid.HorResizer(
				&posH[sys.CurrentWno], nil,
				wid.Image(image[sys.CurrentWno], nil, ""),
				wid.VertResizer(&posV[sys.CurrentWno], nil,
					wid.Btn("Left", nil, func() {}, nil, ""),
					wid.Btn("Right", nil, func() {}, nil, ""),
				),
			)(ctx)
			// EndFrame will swap buffers and limit the maximum framerate.
			sys.EndFrame()
		}
		sys.PollEvents()
	}
}
