package main

import (
	"github.com/jkvatne/jkvgui/btn"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	lightMode = true
	MainRow   = wid.RowStyle{
		W: []float32{0.3, 0.7},
	}
	smallText wid.LabelStyle
	heading   wid.LabelStyle
	music     *wid.Img
	entries   = []string{"Classic", "Jazz", "Rock", "Hiphop", "Opera", "Brass", "Soul"}
)

// Menu demonstrates how to show a list that is generated while drawing it.
func Menu() wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		widgets := make([]wid.Wid, len(entries))
		// Note that "no" has to be atomic because it is concurrently updated in "ticker"
		// entries[4] = strconv.Itoa(int(no.Load()))
		for i, s := range entries {
			widgets[i] = btn.Btn(s, icon.Home, nil, &btn.Text, "")
		}
		return wid.Col(&wid.Secondary, widgets...)(ctx)
	}
}

func Items() wid.Wid {
	return wid.Col(nil,
		wid.Col(&wid.Primary,
			wid.Label("Music", nil),
			wid.Label("What Buttons are Artists Pushing When They Perform Live", &heading),
			wid.W(music, wid.SCALE, ""),
			wid.Row(nil,
				wid.Label("12 hrs ago", &smallText),
				wid.Elastic(),
				btn.Btn("Save", icon.ContentSave, nil, nil, ""),
			),
		),
		wid.Col(&wid.Primary,
			wid.Label("Click Save btn to test the confirmation dialog", nil),
		),
	)
}

func Form() wid.Wid {
	return wid.Row(&MainRow, Menu(), Items())
}

func main() {
	// Setting this true will draw a light blue frame around widgets.
	gpu.DebugWidgets = false
	theme.SetDefaultPallete(lightMode)
	window := gpu.InitWindow(500, 800, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	sys.Initialize(window, 14)
	music, _ = wid.New("music.jpg")
	smallText = wid.DefaultLabel
	smallText.FontSize = 0.6
	heading = *wid.H1L
	heading.Multiline = true
	heading.FontNo = gpu.Normal
	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		Form()(wid.NewCtx())
		wid.ShowHint(nil)
		dialog.Show(nil)
		sys.EndFrame(50)
	}
}
