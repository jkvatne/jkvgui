package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/lib"
	"github.com/jkvatne/jkvgui/wid"
	"log"
	"math"
)

var P = wid.Padding{2, 2, 2, 2}

func YesBtnClick() {
	log.Printf("Yes Btn Click\n")
}
func CancelBtnClick() {
	log.Printf("Cancel Btn Click\n")
}

func NoBtnClick() {
	log.Printf("No Btn Click\n")
}

func Form() wid.Wid {
	r := wid.RowSetup{}
	w := wid.Row(r,
		wid.Label("Mpqy", 13, P, 0),
		wid.Label("MpqyM", 24, P, 0),
		wid.Elastic(),
		wid.Button("Cancel", CancelBtnClick, wid.OkBtn),
		wid.Button("No", NoBtnClick, wid.OkBtn),
		wid.Button("Yes", YesBtnClick, wid.OkBtn),
	)
	return w
}

func Draw() {
	// Calculate sizes
	form := Form()
	ctx := wid.Ctx{Rect: lib.Rect{X: 50, Y: 300, H: 260, W: 800, RR: 0}, Baseline: 0}
	_ = form(ctx)
}

func main() {
	window := gpu.InitWindow(math.MaxInt, math.MaxInt, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	gpu.InitOpenGL(gpu.White)

	for !window.ShouldClose() {

		gpu.StartFrame()
		gpu.RoundedRect(850, 50, 350, 50, 10, 2, gpu.Lightgrey, gpu.Blue)
		gpu.Fonts[0].Printf(50, 100, 24, "24 RobotoMedium")
		gpu.Fonts[1].Printf(50, 130, 24, "24 RobotoRegular")
		gpu.Fonts[2].Printf(50, 160, 24, "24 RobotoMono")
		gpu.Fonts[3].Printf(50, 190, 24, "24 GoRegular")
		gpu.Fonts[4].Printf(50, 220, 24, "24 GoMedium")
		gpu.Fonts[5].Printf(50, 250, 24, "24 GoBold")
		gpu.Fonts[6].Printf(50, 280, 24, "24 GoMono")
		// Black frame around the whole window
		gpu.Rect(10, 10, float32(gpu.WindowWidth)-20, float32(gpu.WindowHeight)-20, 2, gpu.Transparent, gpu.Red)
		Draw()

		gpu.EndFrame(500, window)
	}
}
