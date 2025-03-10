package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/wid"
	"log"
)

var P = f32.Padding{2, 2, 2, 2}

func YesBtnClick() {
	log.Printf("Yes Btn Click\n")
}
func CancelBtnClick() {
	log.Printf("Cancel Btn Click\n")
}

func NoBtnClick() {
	log.Printf("No Btn Click\n")
}

var name = "jkvgui"

func Form() wid.Wid {
	r := wid.RowSetup{}
	w := wid.Row(r,
		wid.Label("Mpqy", 13, P, 4),
		wid.Label("MpqyM", 24, P, 4),
		wid.Elastic(),
		wid.Button("Cancel", CancelBtnClick, wid.OkBtn),
		wid.Button("No", NoBtnClick, wid.OkBtn),
		wid.Button("Yes", YesBtnClick, wid.OkBtn),
		wid.Edit(&name, 10, nil, wid.DefaultEdit),
	)
	return w
}

func Draw() {
	// Calculate sizes
	form := Form()
	ctx := wid.Ctx{Rect: f32.Rect{X: 50, Y: 400, H: 260, W: 800}, Baseline: 0}
	_ = form(ctx)
}

func main() {
	window := gpu.InitWindow(2508, 1270, "Rounded rectangle demo", 1, f32.White)
	defer gpu.Shutdown()

	for !window.ShouldClose() {
		gpu.StartFrame()
		gpu.RoundedRect(300, 50, 100, 100, 10, 2, f32.Lightgrey, f32.Blue)
		gpu.Fonts[2].Printf(50, 100, 24, 0, "24 Roboto100")
		gpu.Fonts[3].Printf(50, 130, 24, 0, "24 Roboto200")
		gpu.Fonts[4].Printf(50, 160, 24, 0, "24 Roboto300")
		gpu.Fonts[5].Printf(50, 190, 24, 0, "24 Roboto500")
		gpu.Fonts[6].Printf(50, 250, 24, 0, "24 Roboto600")
		gpu.Fonts[1].Printf(50, 280, 24, 0, "24 Roboto700")
		gpu.Fonts[7].Printf(50, 310, 24, 0, "24 Roboto800")
		// Red frame around the whole window
		gpu.Rect(10, 10, float32(gpu.WindowWidthDp)-20, float32(gpu.WindowHeightDp)-20, 2, f32.Transparent, f32.Red)
		Draw()
		wid.ShowHint(nil)
		gpu.EndFrame(500, window)
	}
}
