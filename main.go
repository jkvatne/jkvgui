package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/wid"
	"log"
)

var P = wid.Pad{2, 2, 2, 2}

func OkBtnClick() {
	log.Printf("Ok Btn Click\n")
}

func Form() wid.Wid {
	r := wid.RowSetup{}
	w := wid.Row(r,
		wid.Label("LucidaConsole", 13, P, 0),
		wid.Label("MpqyM", 13, P, 0),
		wid.Label("MafmrM", 13, P, 0),
		wid.Label("MqsdfyM", 13, P, 0),
		wid.Elastic(),
		wid.Button("Ok qyj", OkBtnClick, 0, 24, gpu.Lightgrey),
	)
	return w
}

func Draw() {
	// Calculate sizes
	form := Form()
	ctx := wid.Ctx{Rect: wid.Rect{X: 50, Y: 300, H: 260, W: 500, RR: 0}, Baseline: 0}
	_ = form(ctx)
}

func main() {
	window := gpu.InitWindow(91200, 99800, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	gpu.InitOpenGL(gpu.White)
	window.SetMouseButtonCallback(wid.MouseBtnCallback)

	for !window.ShouldClose() {
		wid.Clickables = wid.Clickables[0:0]
		gpu.StartFrame()
		gpu.RoundedRect(650, 50, 350, 50, 10, 2, gpu.Lightgrey, gpu.Blue)
		gpu.Fonts[0].Printf(50, 100, 16, "16 RobotoMedium")
		gpu.Fonts[1].Printf(50, 124, 22, "22 RobotoRegular")
		gpu.Fonts[2].Printf(50, 156, 32, "32 GoRegular")
		gpu.Fonts[3].Printf(50, 204, 45, "45 Gomedium")
		// Black frame around the whole window
		gpu.Rect(10, 10, float32(gpu.WindowWidth)-20, float32(gpu.WindowHeight)-20, 2, gpu.Transparent, gpu.Red)
		Draw()

		gpu.EndFrame(500, window)
	}
}
