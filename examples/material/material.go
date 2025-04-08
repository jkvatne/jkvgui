package main

import (
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	lightMode = true
	MainRow   = wid.ContStyle.W(0.3)
	smallText *wid.LabelStyle
	heading   *wid.LabelStyle
	music     *wid.Img
	entries   = []string{"Classic", "Jazz", "Rock", "Hiphop", "Opera", "Brass", "Soul"}
)

// Menu demonstrates how to show a list that is generated while drawing it.
func Menu() wid.Wid {
	return wid.Col((&wid.ContainerStyle{}).W(0.3),
		wid.Label("Genre", smallText),
		func(ctx wid.Ctx) wid.Dim {
			widgets := make([]wid.Wid, len(entries))
			// Note that "no" has to be atomic because it is concurrently updated in "ticker"
			// entries[4] = strconv.Itoa(int(no.Load()))
			for i, s := range entries {
				widgets[i] = wid.Btn(s, gpu.Home, nil, wid.Text, "")
			}
			return wid.Col(wid.Secondary.W(0.3), widgets...)(ctx)
		},
	)
}

var ES wid.ContainerStyle

func Items() wid.Wid {
	return wid.Col((&wid.ContainerStyle{}).W(0.7),
		wid.Label("Articles", smallText),
		wid.Col(&wid.Primary,
			wid.Label("Hiphop", nil),
			wid.Label("What Buttons are Artists Pushing When They Perform Live", heading),
			wid.Label("12 hrs ago", smallText),
			wid.Image(music, nil, ""),
			wid.Row(nil,
				wid.Elastic(),
				wid.Btn("Save", gpu.ContentSave, nil, nil, ""),
			),
		),
		wid.Col(&wid.Primary,
			wid.Label("More about Taylor Swift...", heading),
		),
		wid.Col(&wid.Primary,
			wid.Label("The new Beatles...", heading),
		),
	)
}

func Form() wid.Wid {
	return wid.Row(MainRow, Menu(), Items())
}

func main() {
	// Setting this true will draw a light blue frame around widgets.
	gpu.DebugWidgets = false
	theme.SetDefaultPallete(lightMode)
	window := gpu.InitWindow(500, 700, "Rounded rectangle demo", 2)
	defer gpu.Shutdown()
	sys.Initialize(window, 14)
	music, _ = wid.NewImage("music.jpg")
	smallText = &wid.DefaultLabel
	smallText.FontSize = 0.8
	heading = wid.H1L
	heading.Multiline = true
	heading.FontSize = 1.5
	heading.FontNo = gpu.Normal
	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		Form()(wid.NewCtx())
		wid.ShowHint(nil)
		dialog.ShowDialogue(nil)
		sys.EndFrame(50)
	}
}
