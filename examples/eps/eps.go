package main

import (
	"github.com/jkvatne/jkvgui/button"
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/scroller"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	Cpu1Nl = "300"
	Cpu1Ni = "400"
	Cpu1Nh string
	Cpu0Ns string

	Cpu2Nl string
	Cpu2Ni string
	Cpu3Ni string
	Cpu2Nh string

	Cpu3Nh     string
	Cpu1Nldot  string
	Cpu2Nldot  string
	Cpu1status string

	Cpu2status string
	Cpu3status string
	Status4    string
	Crc        string

	Status1txt string
	Status2txt string
	Status3txt string
	Status4txt string
	schedule   int

	hb           int
	MainStatus   string
	BackupStatus string
	freq         = [5]float32{100.1, 100.1, 100.1, 100.1, 100.1}
	ad           = [6]float32{100.1, 100.1, 100.1, 100.1, 100.1, 100.1}
	t            = [6]int{99, 99, 99, 99, 99, 99}
)

func set0() {
	// n1.WriteObject(0x4000, 0, 1, 0, "Set schedule 0")
}

func set1() {
	// n1.WriteObject(0x4000, 0, 1, 1, "Set schedule 1")
}

func set2() {
	// n1.WriteObject(0x4000, 0, 1, 2, "Set schedule 2")
}

func set3() {
	// n1.WriteObject(0x4000, 0, 1, 3, "Set schedule 3")
}

func set4() {
	// n1.WriteObject(0x4000, 0, 1, 4, "Set schedule 4")
}

func set5() {
}

var MainForm = &scroller.State{}

// Foirm  setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func epsForm() wid.Wid {
	stsStyle1 := wid.DefaultEdit
	stsStyle1.LabelFraction = 0.1
	stsStyle2 := wid.DefaultEdit
	stsStyle2.LabelFraction = 0.5
	return scroller.W(MainForm,
		wid.Label("EPS Test", wid.H1C),
		wid.Separator(0, 1.0, theme.OnSurface),
		wid.Separator(0, 5.0, theme.Transparent),
		wid.Row(wid.Distribute,
			wid.Col(nil,
				wid.Edit("CPU1 NL", &Cpu1Nl, nil, &stsStyle2),
				wid.Edit("CPU1 NI", &Cpu1Ni, nil, &stsStyle2),
				wid.Edit("CPU1 NH", &Cpu1Nh, nil, &stsStyle2),
				wid.Edit("CPU0 NS", &Cpu0Ns, nil, &stsStyle2),
			),
			wid.Col(nil,
				wid.Edit("CPU2 NL", &Cpu2Nl, nil, &stsStyle2),
				wid.Edit("CPU2 NI", &Cpu2Ni, nil, &stsStyle2),
				wid.Edit("CPU3 NI", &Cpu3Ni, nil, &stsStyle2),
				wid.Edit("CPU2 NH", &Cpu2Nh, nil, &stsStyle2),
			),
			wid.Col(nil,
				wid.Edit("CPU3 NH", &Cpu3Nh, nil, &stsStyle2),
				wid.Edit("CPU1 NLDOT", &Cpu1Nldot, nil, &stsStyle2),
				wid.Edit("CPU2 NLDOT", &Cpu2Nldot, nil, &stsStyle2),
				wid.Edit("CPU1 STATUS", &Cpu1status, nil, &stsStyle2),
			),
			wid.Col(nil,
				wid.Edit("CPU2 STATUS", &Cpu2status, nil, &stsStyle2),
				wid.Edit("CPU3 STATUS", &Cpu3status, nil, &stsStyle2),
				wid.Edit("STATUS4", &Status4, nil, &stsStyle2),
				wid.Edit("Program CRC", &Crc, nil, &stsStyle2),
			),
		),
		wid.Edit("CPU1 STATUS", &Status1txt, nil, &stsStyle1),
		wid.Edit("CPU2 STATUS", &Status2txt, nil, &stsStyle1),
		wid.Edit("CPU3 STATUS", &Status3txt, nil, &stsStyle1),
		wid.Edit("STATUS4", &Status4txt, nil, &stsStyle1),
		wid.Separator(0, 16.0, theme.Surface),
		wid.Row(wid.Distribute,
			wid.Col(nil,
				wid.Label("Measured speed [Hz]", wid.H2R),
				wid.Edit("NL (Hz)", &freq[0], nil, &stsStyle2),
				wid.Edit("NI (Hz)", &freq[1], nil, &stsStyle2),
				wid.Edit("NH (Hz)", &freq[2], nil, &stsStyle2),
				wid.Edit("NS (Hz)", &freq[3], nil, &stsStyle2),
				wid.Separator(0, 9, theme.Transparent),
				wid.Edit("Heartbeats", &hb, nil, &stsStyle2),
			),
			wid.Col(nil,
				wid.Label("Internal measurements", wid.H2R),
				wid.Edit("ESOV current (A)", &ad[0], nil, &stsStyle2),
				wid.Edit("ESOV Lo (V)", &ad[1], nil, &stsStyle2),
				wid.Edit("Supply (V)", &ad[2], nil, &stsStyle2),
				wid.Edit("ESOV Hi (V)", &ad[3], nil, &stsStyle2),
				wid.Edit("RF gnd (V)", &ad[4], nil, &stsStyle2),
				wid.Edit("Internal (V)", &ad[5], nil, &stsStyle2),
			),
			wid.Col(nil,
				wid.Label("Internal timers ", wid.H2R),
				wid.Edit("calculate+prepare_com", &t[0], nil, &stsStyle2),
				wid.Edit("End of master TX", &t[1], nil, &stsStyle2),
				wid.Edit("End of slave RX", &t[2], nil, &stsStyle2),
				wid.Edit("After process_can_pdo", &t[3], nil, &stsStyle2),
				wid.Edit("Last time in SEND", &t[4], nil, &stsStyle2),
				wid.Edit("Time spent in SEND", &t[5], nil, &stsStyle2),
			),
			wid.Col(nil,
				wid.Label("Schedule", wid.H1C),
				wid.Edit("Selected schedule", &schedule, nil, &stsStyle2),
				wid.Row(wid.Left,
					wid.Elastic(),
					button.Filled("0", nil, set0, nil, ""),
					wid.Elastic(),
					button.Filled("1", nil, set1, nil, ""),
					wid.Elastic(),
					button.Filled("2", nil, set2, nil, ""),
					wid.Elastic(),
					button.Filled("3", nil, set3, nil, ""),
					wid.Elastic(),
					button.Filled("4", nil, set4, nil, ""),
					wid.Elastic(),
				),
				wid.Label("Click a button to change schedule", wid.Center),
				wid.Separator(0, 2, theme.OnSurface),
				wid.Edit("", &MainStatus, nil, nil),
				wid.Edit("", &BackupStatus, nil, nil),
			),
		),
		wid.Row(wid.Distribute,
			button.Filled("Primary", icon.Home, set0, button.Btn.Role(theme.Primary), ""),
			button.Filled("Secondary", icon.ContentOpen, set1, button.Btn.Role(theme.Secondary), ""),
			button.Filled("Surface", icon.ContentSave, set2, button.Btn.Role(theme.Surface), ""),
			button.Filled("Container", icon.RadioChecked, set3, button.Btn.Role(theme.SurfaceContainer), ""),
			button.Filled("Round", nil, set5, &button.RoundBtn, ""),
		),
		wid.Row(wid.Left,
			wid.Elastic(),
			button.Filled("Primary", icon.Home, set0, button.Btn.Role(theme.Primary), ""),
			wid.Elastic(),
			button.Filled("Secondary", icon.ContentOpen, set1, button.Btn.Role(theme.Secondary), ""),
			wid.Elastic(),
			button.Filled("Surface", icon.ContentSave, set2, button.Btn.Role(theme.Surface), ""),
			wid.Elastic(),
			button.Filled("Container", icon.RadioChecked, set3, button.Btn.Role(theme.SurfaceContainer), ""),
			wid.Elastic(),
			button.Filled("Round", nil, set5, &button.RoundBtn, ""),
			wid.Elastic(),
		),
		wid.Label("EPS Test1", wid.H1C),
		wid.Label("EPS Test2", wid.H1C),
		wid.Label("EPS Test3", wid.H1C),
		wid.Label("EPS Test4", wid.H1C),
		wid.Label("EPS Test5", wid.H1C),
		wid.Label("EPS Test6", wid.H1C),
		wid.Label("EPS Test7", wid.H1C),
		wid.Label("EPS Test8", wid.H1C),

	)
}

func main() {
	theme.SetDefaultPallete(true)
	window := gpu.InitWindow(0, 0, "EPS test", 1)
	defer gpu.Shutdown()
	Status1txt = "Status1 text"
	Status2txt = "Status2 text"
	Status3txt = "Status3 text"
	Status4txt = "Status4 text"
	callback.Initialize(window)
	for !window.ShouldClose() {
		gpu.BackgroundRole(theme.Surface)
		ctx := wid.Ctx{Rect: f32.Rect{X: 0, Y: 0, W: gpu.WindowWidthDp, H: gpu.WindowHeightDp}, Baseline: 0}
		gpu.StartFrame(theme.Surface.Bg())
		form := epsForm()
		_ = form(ctx)
		wid.ShowHint(nil)
		dialog.Show(nil)
		gpu.EndFrame(30)
	}
}
