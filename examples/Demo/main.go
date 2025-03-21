package main

import (
	"github.com/jkvatne/jkvgui/button"
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"strconv"
)

var (
	lightMode        = true
	gender    string = "Male"
	genders          = []string{"Male", "Female", "Both", "qyjpy"}
	name             = "Ole Petter Olsen"
	address          = "Mo i Rana"
	hint1            = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"
	hint2            = "This is a hint"
	hint3            = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10 Word11 word12 jyword13"
)

func YesBtnClick() {
	lightMode = true
	theme.SetDefaultPallete(lightMode)
	slog.Info("Yes Btn Clicked")
}

func do() {
	dialog.Exit()
}

func DlgBtnClick() {
	dialog.Current = dialog.YesNoDialog("Heading", "Some text", "Yes", "No", do, do)
	slog.Info("Cancel Btn clicked")
}

func NoBtnClick() {
	lightMode = false
	theme.SetDefaultPallete(lightMode)
	slog.Info("No Btn Click\n")
}

func Form() wid.Wid {
	return wid.Col(nil,
		wid.Label("Edit user information", wid.H1C),
		wid.Label("Use TAB to move focus, and Enter to save data", wid.I),
		wid.Edit("", &name, nil, nil),
		wid.Edit("", &address, nil, nil),
		wid.Combo(&gender, genders, nil),
		wid.Label("MpqyM2", nil),
		wid.Label(strconv.Itoa(gpu.RedrawsPrSec), nil),
		wid.Checkbox("Darkmode", &lightMode, nil, ""),
		wid.Row(1,
			wid.Label("Buttons", nil),
			wid.Elastic(),
			button.Filled("Show dialogue", nil, DlgBtnClick, nil, hint1),
			button.Filled("No", nil, NoBtnClick, &button.Btn, hint2),
			button.Filled("Yes", nil, YesBtnClick, &button.Btn, hint3),
		),
	)
}

func Draw() {
	// Calculate sizes

}

func main() {
	theme.SetDefaultPallete(lightMode)
	window := gpu.InitWindow(0, 0, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()

	callback.Initialize(window)

	for !window.ShouldClose() {
		gpu.BackgroundColor(theme.Surface)
		gpu.StartFrame()
		form := Form()
		ctx := wid.Ctx{Rect: f32.Rect{X: 0, Y: 0, W: gpu.WindowWidthDp, H: gpu.WindowHeightDp}, Baseline: 0}
		_ = form(ctx)
		wid.ShowHint(nil)
		dialog.Show(nil)

		gpu.EndFrame(30)
	}
}
