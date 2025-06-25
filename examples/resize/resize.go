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
	stop := false
	sys.CreateWindow(100, 100, 400, 400, "Resizing1", 1, 2)
	sys.CreateWindow(400, 400, 400, 400, "Resizing2", 2, 2)
	defer sys.Shutdown()
	image[0], _ = wid.NewImage("music.jpg")
	image[1], _ = wid.NewImage("ts.jpg")
	for !stop {
		stop = true
		for wno, _ := range sys.WindowList {
			if !sys.Running(wno) {
				continue
			}
			stop = false
			sys.StartFrame(wno, theme.Surface.Bg())
			ctx := wid.NewCtx(wno)
			wid.HorResizer(
				&posH[wno], nil,
				wid.Image(image[wno], nil, ""),
				wid.VertResizer(&posV[wno], nil,
					wid.Btn("Left", nil, func() {}, nil, ""),
					wid.Btn("Right", nil, func() {}, nil, ""),
				),
			)(ctx)
			// EndFrame will swap buffers and limit the maximum framerate.
			sys.EndFrame(wno)
		}
		sys.PollEvents()
	}
}
