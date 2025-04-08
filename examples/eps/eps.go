package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
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

var ss = &wid.ScrollState{}

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

// Form  setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func epsForm() wid.Wid {
	stsStyle := wid.DefaultEdit
	stsStyle.LabelSize = 0.1
	stsStyle.EditSize = 0.9
	ValueStyle := wid.DefaultEdit
	ValueStyle.EditSize = 0.6
	ValueStyle.LabelSize = 0.4
	ValueStyle.LabelRightAdjust = true
	return wid.ScrollPane(ss,
		wid.Label("EPS Test", wid.H1C),
		wid.Separator(0, 1.0, theme.OnSurface),
		wid.Separator(0, 5.0, theme.Transparent),
		wid.Row(nil,
			wid.Col(nil,
				wid.Edit(&Cpu1Nl, "CPU1 NL", nil, &ValueStyle),
				wid.Edit(&Cpu1Ni, "CPU1 NI", nil, &ValueStyle),
				wid.Edit(&Cpu1Nh, "CPU1 NH", nil, &ValueStyle),
				wid.Edit(&Cpu0Ns, "CPU0 NS", nil, &ValueStyle),
			),
			wid.Col(nil,
				wid.Edit(&Cpu2Nl, "CPU2 NL", nil, &ValueStyle),
				wid.Edit(&Cpu2Ni, "CPU2 NI", nil, &ValueStyle),
				wid.Edit(&Cpu3Ni, "CPU3 NI", nil, &ValueStyle),
				wid.Edit(&Cpu2Nh, "CPU2 NH", nil, &ValueStyle),
			),
			wid.Col(nil,
				wid.Edit(&Cpu3Nh, "CPU3 NH", nil, &ValueStyle),
				wid.Edit(&Cpu1Nldot, "CPU1 NLDOT", nil, &ValueStyle),
				wid.Edit(&Cpu2Nldot, "CPU2 NLDOT", nil, &ValueStyle),
				wid.Edit(&Cpu1status, "CPU1 STATUS", nil, &ValueStyle),
			),
			wid.Col(nil,
				wid.Edit(&Cpu2status, "CPU2 STATUS", nil, &ValueStyle),
				wid.Edit(&Cpu3status, "CPU3 STATUS", nil, &ValueStyle),
				wid.Edit(&Status4, "STATUS4", nil, &ValueStyle),
				wid.Edit(&Crc, "Program CRC", nil, &ValueStyle),
			),
		),
		wid.Edit(&Status1txt, "CPU1 STATUS", nil, &stsStyle),
		wid.Edit(&Status2txt, "CPU2 STATUS", nil, &stsStyle),
		wid.Edit(&Status3txt, "CPU3 STATUS", nil, &stsStyle),
		wid.Edit(&Status4txt, "STATUS4", nil, &stsStyle),
		wid.Separator(0, 16.0, theme.Surface),
		wid.Row(nil,
			wid.Col(nil,
				wid.Label("Measured speed [Hz]", wid.H2R),
				wid.Edit(&freq[0], "NL (Hz)", nil, &ValueStyle),
				wid.Edit(&freq[1], "NI (Hz)", nil, &ValueStyle),
				wid.Edit(&freq[2], "NH (Hz)", nil, &ValueStyle),
				wid.Edit(&freq[3], "NS (Hz)", nil, &ValueStyle),
				wid.Separator(0, 9, theme.Transparent),
				wid.Edit(&hb, "Heartbeats", nil, &ValueStyle),
			),
			wid.Col(nil,
				wid.Label("Internal measurements", wid.H2R),
				wid.Edit(&ad[0], "ESOV current (A)", nil, &ValueStyle),
				wid.Edit(&ad[1], "ESOV Lo (V)", nil, &ValueStyle),
				wid.Edit(&ad[2], "Supply (V)", nil, &ValueStyle),
				wid.Edit(&ad[3], "ESOV Hi (V)", nil, &ValueStyle),
				wid.Edit(&ad[4], "RF gnd (V)", nil, &ValueStyle),
				wid.Edit(&ad[5], "Internal (V)", nil, &ValueStyle),
			),
			wid.Col(nil,
				wid.Label("Internal timers ", wid.H2R),
				wid.Edit(&t[0], "calculate+prepare_com", nil, &ValueStyle),
				wid.Edit(&t[1], "End of master TX", nil, &ValueStyle),
				wid.Edit(&t[2], "End of slave RX", nil, &ValueStyle),
				wid.Edit(&t[3], "After process_can_pdo", nil, &ValueStyle),
				wid.Edit(&t[4], "Last time in SEND", nil, &ValueStyle),
				wid.Edit(&t[5], "Time spent in SEND", nil, &ValueStyle),
			),
			wid.Col(nil,
				wid.Label("Schedule", wid.H1C),
				wid.Edit(&schedule, "Selected schedule", nil, &ValueStyle),
				wid.Row(nil,
					wid.Elastic(),
					wid.Btn("0", nil, set0, nil, ""),
					wid.Elastic(),
					wid.Btn("1", nil, set1, nil, ""),
					wid.Elastic(),
					wid.Btn("2", nil, set2, nil, ""),
					wid.Elastic(),
					wid.Btn("3", nil, set3, nil, ""),
					wid.Elastic(),
					wid.Btn("4", nil, set4, nil, ""),
					wid.Elastic(),
				),
				wid.Label("Click a btn to change schedule", wid.Center),
				wid.Separator(0, 2, theme.OnSurface),
				wid.Edit(&MainStatus, "", nil, nil),
				wid.Edit(&BackupStatus, "", nil, nil),
			),
		),
	)
}

func main() {
	theme.SetDefaultPallete(true)
	gpu.UserScale = 1.5
	// gpu.DebugWidgets = true
	window := gpu.InitWindow(0, 0, "EPS test", 1)
	defer gpu.Shutdown()
	Status1txt = "Status1 text"
	Status2txt = "Status2 text"
	Status3txt = "Status3 text"
	Status4txt = "Status4 text"
	sys.Initialize(window, 16)
	for !window.ShouldClose() {
		ctx := wid.Ctx{Rect: f32.Rect{X: 0, Y: 0, W: gpu.WindowWidthDp, H: gpu.WindowHeightDp}, Baseline: 0}
		sys.StartFrame(theme.Surface.Bg())
		_ = epsForm()(ctx)
		wid.ShowHint(nil)
		sys.EndFrame(30)
	}
}
