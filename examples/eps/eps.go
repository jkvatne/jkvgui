package main

import (
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/scroller"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	dummy  string
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

var MainForm *scroller.State = &scroller.State{}

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
				wid.EditF32("NL (Hz)", &freq[0], nil, &stsStyle2),
				wid.EditF32("NI (Hz)", &freq[1], nil, &stsStyle2),
				wid.EditF32("NH (Hz)", &freq[2], nil, &stsStyle2),
				wid.EditF32("NS (Hz)", &freq[3], nil, &stsStyle2),
				wid.Separator(0, 9, theme.Transparent),
				wid.EditInt("Heartbeats", &hb, nil, &stsStyle2),
			),
			wid.Col(nil,
				wid.Label("Internal measurements", wid.H2R),
				wid.EditF32("ESOV current (A)", &ad[0], nil, &stsStyle2),
				wid.EditF32("ESOV Lo (V)", &ad[1], nil, &stsStyle2),
				wid.EditF32("Supply (V)", &ad[2], nil, &stsStyle2),
				wid.EditF32("ESOV Hi (V)", &ad[3], nil, &stsStyle2),
				wid.EditF32("RF gnd (V)", &ad[4], nil, &stsStyle2),
				wid.EditF32("Internal (V)", &ad[5], nil, &stsStyle2),
			),
			wid.Col(nil,
				wid.Label("Internal timers ", wid.H2R),
				wid.EditInt("calculate+prepare_com", &t[0], nil, &stsStyle2),
				wid.EditInt("End of master TX", &t[1], nil, &stsStyle2),
				wid.EditInt("End of slave RX", &t[2], nil, &stsStyle2),
				wid.EditInt("After process_can_pdo", &t[3], nil, &stsStyle2),
				wid.EditInt("Last time in SEND", &t[4], nil, &stsStyle2),
				wid.EditInt("Time spent in SEND", &t[5], nil, &stsStyle2),
			),
			wid.Col(nil,
				wid.Label("Schedule", wid.H1C),
				wid.EditInt("Selected schedule", &schedule, nil, &stsStyle2),
				wid.Row(wid.Left,
					wid.Elastic(),
					wid.Button("0", set0, nil, ""),
					wid.Elastic(),
					wid.Button("1", set1, nil, ""),
					wid.Elastic(),
					wid.Button("2", set2, nil, ""),
					wid.Elastic(),
					wid.Button("3", set3, nil, ""),
					wid.Elastic(),
					wid.Button("4", set4, nil, ""),
					wid.Elastic(),
				),
				wid.Label("Click a button to change schedule", wid.Center),
				wid.Separator(0, 2, theme.OnSurface),
				wid.Edit("", &MainStatus, nil, nil),
				wid.Edit("", &BackupStatus, nil, nil),
			),
		),
		wid.Row(wid.Left,
			wid.Elastic(),
			wid.Button("Primary", set0, wid.Btn.Role(theme.Primary), ""),
			wid.Elastic(),
			wid.Button("Secondary", set1, wid.Btn.Role(theme.Secondary), ""),
			wid.Elastic(),
			wid.Button("Surface", set2, wid.Btn.Role(theme.Surface), ""),
			wid.Elastic(),
			wid.Button("Container", set3, wid.Btn.Role(theme.SurfaceContainer), ""),
			wid.Elastic(),
			wid.Button("Round", set5, &wid.RoundBtn, ""),
			wid.Elastic(),
		),
	)
}

func main() {
	theme.SetDefaultPallete(true)
	window := gpu.InitWindow(0, 0, "EPS test", 1)
	defer gpu.Shutdown()
	callback.Initialize(window)
	for !window.ShouldClose() {
		gpu.BackgroundColor(theme.Surface)
		ctx := wid.Ctx{Rect: f32.Rect{X: 0, Y: 0, W: gpu.WindowWidthDp, H: gpu.WindowHeightDp}, Baseline: 0}
		gpu.StartFrame()
		form := epsForm()
		_ = form(ctx)
		wid.ShowHint(nil)
		dialog.Show(nil)
		gpu.EndFrame(30)
	}
}
