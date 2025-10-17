package dialog

import (
	"log/slog"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

type DialogueStyle struct {
	FontNo          int
	FontSize        float32
	FontColor       theme.UIRole
	CornerRadius    float32
	BorderColor     theme.UIRole
	BackgroundColor theme.UIRole
	BorderWidth     float32
	Padding         f32.Padding
	Delay           time.Duration
}

var DefaultDialogueStyle = DialogueStyle{
	FontNo:          gpu.Normal14,
	FontSize:        1.0,
	CornerRadius:    5,
	FontColor:       theme.OnSurface,
	BackgroundColor: theme.SurfaceContainer,
	BorderColor:     theme.Outline,
	BorderWidth:     1,
	Padding:         f32.Padding{L: 25, T: 15, R: 25, B: 15},
	Delay:           time.Millisecond * 800,
}

type DialogMap map[*sys.Window]*wid.Wid

var Dialogs DialogMap

var dialogStartTime = time.Now()

func YesNoDialog(heading string, text string, lbl1, lbl2 string, on1, on2 func()) wid.Wid {
	return wid.Col(
		nil,
		wid.Separator(0, 25),
		wid.Label(heading, wid.H1C),
		wid.Separator(0, 12),
		wid.Label(text, wid.C),
		wid.Separator(0, 25),
		wid.Row(nil,
			wid.Btn(lbl1, nil, on1, nil, ""),
			wid.Btn(lbl2, nil, on2, nil, ""),
		),
	)
}

func Hide() {
	win := sys.GetCurrentWindow()
	delete(Dialogs, win)
	sys.GetCurrentWindow().SuppressEvents = false
}

func Show(w *wid.Wid) {
	win := sys.GetCurrentWindow()
	Dialogs[win] = w
}

func Display() {
	win := sys.GetCurrentWindow()
	if win == nil {
		slog.Error("Dialog Display(), Window Not Found")
		return
	}
	CurrentDialog := Dialogs[win]
	win.DialogVisible = CurrentDialog != nil
	if !win.DialogVisible {
		return
	}
	style := &DefaultDialogueStyle
	// f goes from 0 to 0.5 after ca 0.5 second
	f := min(1.0, float32(time.Since(dialogStartTime))/float32(time.Second))
	// Draw surface all over the underlying form with the transparent surface color
	win.Gd.SolidRect(win.ClientRectDp(), f32.Black.MultAlpha(f*0.5))
	// Draw dialog
	w := float32(300)
	h := float32(180)
	x := (win.WidthDp - w) / 2
	y := (win.HeightDp - h) / 2
	ctx := wid.Ctx{Rect: f32.Rect{X: x, Y: y, W: w, H: h}, Baseline: 0}
	ctx.Win = win
	ctx.Win.SuppressEvents = false
	if f < 1.0 {
		ctx.Win.Invalidate()
	}
	win.Gd.RoundedRect(ctx.Rect, 10, 2, theme.Colors[style.BackgroundColor], f32.Transparent)
	ctx.Rect = ctx.Rect.Inset(style.Padding, 0)
	_ = (*CurrentDialog)(ctx)
	ctx.Win.SuppressEvents = true
}

func init() {
	Dialogs = make(DialogMap)
}
