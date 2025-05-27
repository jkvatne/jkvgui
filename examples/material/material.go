package main

import (
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
)

var (
	lightMode = true
	MainRow   = wid.ContStyle.W(0.3)
	smallText wid.LabelStyle
	heading   wid.LabelStyle
	music     *wid.Img
	swift     *wid.Img
	entries   = []string{"Classic", "Jazz", "Rock", "Hiphop", "Opera", "Brass", "Soul"}
)

func do() {
	slog.Info("Save clicked")
}

var MenuStyle = (&wid.ContainerStyle{}).W(0.3)

// Menu demonstrates how to show a list that is generated while drawing it.
func Menu() wid.Wid {
	return wid.Col(MenuStyle,
		wid.Label("Genre", &smallText),
		func(ctx wid.Ctx) wid.Dim {
			widgets := make([]wid.Wid, len(entries))
			for i, s := range entries {
				widgets[i] = wid.Btn(s, gpu.Home, nil, wid.Text, "")
			}
			return wid.Col(wid.Secondary.W(0.3), widgets...)(ctx)
		},
	)
}

var ss = &wid.ScrollState{Width: 0.7}

func Items() wid.Wid {
	return wid.Scroller(ss,
		wid.Label("Articles", &smallText),
		wid.Col(&wid.Primary,
			wid.Label("Hiphop", nil),
			wid.Label("What Buttons are Artists Pushing When They Perform Live", &heading),
			wid.Label("12 hrs ago", &smallText),
			wid.Image(music, wid.DefImg.Bg(theme.PrimaryContainer), ""),
			wid.Row(nil,
				wid.Elastic(),
				wid.Btn("Save", gpu.ContentSave, do, nil, ""),
			),
		),
		wid.Col(&wid.Primary,
			wid.Label("More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("The new Beatles...", &heading),
		),
		wid.Col(&wid.Primary,
			wid.Label("1 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("2 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("3 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("4 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("5 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("6 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("7 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("8 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("9 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
	)
}

func Form() wid.Wid {
	return wid.Row(MainRow, Menu(), Items())
}

func main() {
	sys.InitWindow(500, 500, "Material demo", 2, 1.0)
	defer sys.Shutdown()

	music, _ = wid.NewImage("music.jpg")
	swift, _ = wid.NewImage("ts.jpg")
	smallText = wid.DefaultLabel
	smallText.FontNo = gpu.Normal10
	heading = *wid.H1L
	heading.Multiline = true
	heading.FontNo = gpu.Bold20
	theme.Colors[theme.OnPrimary] = f32.Yellow
	for sys.Running() {
		sys.StartFrame(theme.Surface.Bg())
		Form()(wid.NewCtx())
		dialog.ShowDialogue()
		sys.EndFrame()
	}
}
