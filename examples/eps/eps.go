package main

import (
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
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

// Foirm  setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func epsForm() wid.Wid {
	stsStyle := wid.DefaultEdit
	stsStyle.LabelFraction = 0.1
	/*return wid.Col(nil,
	wid.Label("EPS Test", wid.H1),
	wid.Separator(0, 1.0, theme.OnSurface),
	wid.Separator(0, 5.0, theme.Transparent),
	wid.Row(wid.Distribute,
		wid.Col(nil,
			wid.Edit("CPU1 NL", &Cpu1Nl, nil, nil),
			wid.Edit("CPU1 NI", &Cpu1Ni, nil, nil),
			wid.Edit("CPU1 NH", &Cpu1Nh, nil, nil),
			wid.Edit("CPU0 NS", &Cpu0Ns, nil, nil),
		),
		wid.Col(nil,
			wid.Edit("CPU2 NL", &Cpu2Nl, nil, nil),
			wid.Edit("CPU2 NI", &Cpu2Ni, nil, nil),
			wid.Edit("CPU3 NI", &Cpu3Ni, nil, nil),
			wid.Edit("CPU2 NH", &Cpu2Nh, nil, nil),
		),
		wid.Col(nil,
			wid.Edit("CPU3 NH", &Cpu3Nh, nil, nil),
			wid.Edit("CPU1 NLDOT", &Cpu1Nldot, nil, nil),
			wid.Edit("CPU2 NLDOT", &Cpu2Nldot, nil, nil),
			wid.Edit("CPU1 STATUS", &Cpu1status, nil, nil),
		),
		wid.Col(nil,
			wid.Edit("CPU2 STATUS", &Cpu2status, nil, nil),
			wid.Edit("CPU3 STATUS", &Cpu3status, nil, nil),
			wid.Edit("STATUS4", &Status4, nil, nil),
			wid.Edit("Program CRC", &Crc, nil, nil),
		),
	),
	wid.Edit("CPU1 STATUS", &Status1txt, nil, &stsStyle),
	wid.Edit("CPU2 STATUS", &Status2txt, nil, &stsStyle),
	wid.Edit("CPU3 STATUS", &Status3txt, nil, &stsStyle),
	wid.Edit("STATUS4", &Status4txt, nil, &stsStyle),
	wid.Separator(0, 1.0, theme.OnSurface), */
	return wid.Row(wid.Distribute,
		wid.Col(nil,
			wid.Label("Measured speed [Hz]", wid.H1R),
			wid.EditF32("NL (Hz)", &freq[0], nil, nil),
			wid.EditF32("NI (Hz)", &freq[1], nil, nil),
			wid.EditF32("NH (Hz)", &freq[2], nil, nil),
			wid.EditF32("NS (Hz)", &freq[3], nil, nil),
			wid.Separator(0, 9, theme.Transparent),
			wid.EditInt("Heartbeats", &hb, nil, nil),
		),
		wid.Col(nil,
			wid.Label("Internal measurements", wid.H1R),
			wid.EditF32("ESOV current (A)", &ad[0], nil, nil),
			wid.EditF32("ESOV Lo (V)", &ad[1], nil, nil),
			wid.EditF32("Supply (V)", &ad[2], nil, nil),
			wid.EditF32("ESOV Hi (V)", &ad[3], nil, nil),
			wid.EditF32("RF gnd (V)", &ad[4], nil, nil),
			wid.EditF32("Internal (V)", &ad[5], nil, nil),
		),
		wid.Col(nil,
			wid.Label("Internal timers ", wid.H1R),
			wid.EditInt("calculate+prepare_com", &t[0], nil, nil),
			wid.EditInt("End of master TX", &t[1], nil, nil),
			wid.EditInt("End of slave RX", &t[2], nil, nil),
			wid.EditInt("After process_can_pdo", &t[3], nil, nil),
			wid.EditInt("Last time in SEND", &t[4], nil, nil),
			wid.EditInt("Time spent in SEND", &t[5], nil, nil),
		),
		/*wid.Col(nil,
			wid.Label("Schedule", nil),
			wid.EditInt("Selected schedule", &schedule, nil, nil),
			wid.Row(nil,
				wid.Separator(2, 0, theme.OnSurface),
				wid.Button("0 ", set0, nil, ""),
				wid.Separator(2, 0, theme.OnSurface),
				wid.Button("1 ", set1, nil, ""),
				wid.Separator(2, 0, theme.OnSurface),
				wid.Button("2 ", set2, nil, ""),
				wid.Separator(2, 0, theme.OnSurface),
				wid.Button("3 ", set3, nil, ""),
				wid.Separator(2, 0, theme.OnSurface),
				wid.Button("4 ", set4, nil, ""),
				wid.Separator(2, 0, theme.OnSurface),
			),
		),
		wid.Label("Click a button to change schedule", nil),
		wid.Separator(0, 2, theme.OnSurface),
		wid.Edit("Main Status", &MainStatus, nil, nil),
		wid.Edit("Backup sts", &BackupStatus, nil, nil),*/
	)
	// )
}

func main() {
	theme.SetDefaultPallete(true)
	window := gpu.InitWindow(500, 500, "Rounded rectangle demo", 2)
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
