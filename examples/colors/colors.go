package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"strconv"
)

var ShowRoles bool

var Box = wid.LabelStyle{
	Height: 30,
	Width:  0.1,
}

func showTones(c f32.Color) wid.Wid {
	return wid.Row(nil,
		wid.BoxText("00", f32.White, c.Tone(00), &Box),
		wid.BoxText("10", f32.White, c.Tone(10), &Box),
		wid.BoxText("20", f32.White, c.Tone(20), &Box),
		wid.BoxText("30", f32.White, c.Tone(30), &Box),
		wid.BoxText("40", f32.White, c.Tone(40), &Box),
		wid.BoxText("50", f32.White, c.Tone(50), &Box),
		wid.BoxText("60", f32.Black, c.Tone(60), &Box),
		wid.BoxText("70", f32.Black, c.Tone(70), &Box),
		wid.BoxText("80", f32.Black, c.Tone(80), &Box),
		wid.BoxText("88", f32.Black, c.Tone(88), &Box),
		wid.BoxText("93", f32.Black, c.Tone(93), &Box),
		wid.BoxText("97", f32.Black, c.Tone(97), &Box),
		wid.BoxText("100", f32.Black, c.Tone(100), &Box),
	)
}

func setDefault() {
	theme.SetPallete(true, 0x5750C4, 0x925B51, 0x27624E, 0x79747E, 0xAF1515)
}

func setPalette1() {
	theme.SetPallete(true, 0x67622E, 0x27622E, 0x27624E, 0x1D5D7D, 0xAF1515)
}

func setPalette2() {
	theme.SetPallete(true, 0x17624E, 0x27624E, 0x27624E, 0x1D4D7D, 0xAF1515)
}

func setPalette3() {
	theme.SetPallete(true, 0x27624E, 0x15625E, 0x27624E, 0x1D1D7D, 0xBF0000)
}

func setColorsRoles() {
	ShowRoles = !ShowRoles
}

var lightMode bool

func setDarkLight() {
	lightMode = !lightMode
	theme.SetupColors(lightMode)
}

func form2() wid.Wid {
	ld := "Set dark"
	if lightMode {
		ld = "Set light"
	}
	cr := "Show Colors"
	if ShowRoles {
		cr = "Show Roles"
	}
	return wid.Col(nil,
		wid.Label("Show all UI roles", wid.H1C),
		wid.Label("FPS="+strconv.Itoa(sys.RedrawsPrSec), nil),
		wid.Row(nil,
			wid.Btn("Set default", nil, setDefault, nil, "Set the default palette on all widgets"),
			wid.Btn("Set palette 1", nil, setPalette1, nil, "Select palette 1"),
			wid.Btn("Set palette 2", nil, setPalette2, nil, "Select palette nr 2"),
			wid.Btn("Set palette 3", nil, setPalette3, nil, "Select palette nr. 3"),
			wid.Btn(cr, nil, setColorsRoles, nil, "Change between showing color tones and role palette"),
			wid.Btn(ld, nil, setDarkLight, nil, "Select light or dark mode"),
		),
		wid.Separator(1.0, 3.0),
		wid.Row(nil,
			wid.Col(nil,
				wid.Col(wid.ContStyle.R(theme.Primary),
					wid.Label("Primary", wid.C.R(theme.Primary))),
				wid.Col(wid.ContStyle.R(theme.Secondary),
					wid.Label("Secondary", wid.C.R(theme.Secondary))),
				wid.Col(wid.ContStyle.R(theme.Error),
					wid.Label("Error", wid.C.R(theme.Error))),
				wid.Col(wid.ContStyle.R(theme.Outline),
					wid.Label("Outline", wid.C.R(theme.Outline))),
			),
			wid.Col(nil,
				wid.Col(wid.ContStyle.R(theme.PrimaryContainer),
					wid.Label("PrimaryContainer.", wid.C.R(theme.PrimaryContainer))),
				wid.Col(wid.ContStyle.R(theme.SecondaryContainer),
					wid.Label("SecondaryContainer.", wid.C.R(theme.SecondaryContainer))),
				wid.Col(wid.ContStyle.R(theme.SurfaceContainer),
					wid.Label("SurfaceContainer.", wid.C.R(theme.SurfaceContainer))),
				wid.Col(wid.ContStyle.R(theme.ErrorContainer),
					wid.Label("ErrorContainer.", wid.C.R(theme.ErrorContainer))),
				wid.Col(wid.ContStyle.R(theme.Surface),
					wid.Label("Surface", wid.C.R(theme.Surface))),
			),
		),
	)
}

func form1() wid.Wid {
	var ld string
	var cr string
	if lightMode {
		ld = "Set dark"
	} else {
		ld = "Set light"
	}
	if ShowRoles {
		cr = "Show Colors"
	} else {
		cr = "Show Roles"
	}
	return wid.Col(nil,
		wid.Label("Show all tones for some palettes", wid.H1C),
		wid.Label("FPS="+strconv.Itoa(sys.RedrawsPrSec), nil),
		wid.Row(nil,
			wid.Btn("Set default", nil, setDefault, nil, "Set the default palette on all widgets"),
			wid.Btn("Set palette 1", nil, setPalette1, nil, "Use a palette 1"),
			wid.Btn("Set palette 2", nil, setPalette2, nil, "Select palette nr 2"),
			wid.Btn("Set palette 3", nil, setPalette3, nil, "Select palette nr. 3"),
			wid.Btn(cr, nil, setColorsRoles, nil, "Change between showing color tones and role palette"),
			wid.Btn(ld, nil, setDarkLight, nil, "Select light or dark mode"),
		),
		wid.Separator(0.0, 1.0),
		wid.Label("PrimaryColor", nil),
		showTones(theme.PrimaryColor),
		wid.Label("SecondaryColor", nil),
		showTones(theme.SecondaryColor),
		wid.Label("TertiaryColor", nil),
		showTones(theme.TertiaryColor),
		wid.Label("ErrorColor", nil),
		showTones(theme.ErrorColor),
		wid.Label("NeutralColor", nil),
		showTones(theme.NeutralColor),
	)
}

func main() {
	sys.Initialize()
	theme.SetDefaultPallete(lightMode)
	gpu.InitWindow(0, 0, "Rounded rectangle demo", 2, 2.0)
	defer sys.Shutdown()
	sys.InitializeWindow()
	for !gpu.ShouldClose() {
		// We want a fully white or black background, so we use the Canvas role
		sys.StartFrame(theme.Canvas.Bg())
		// Draw form
		if ShowRoles == true {
			form2()(wid.NewCtx())
		} else {
			form1()(wid.NewCtx())
		}
		sys.EndFrame(0)
	}
}
