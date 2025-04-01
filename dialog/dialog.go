package dialog

import (
	"github.com/jkvatne/jkvgui/btn"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"time"
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
	FontNo:          gpu.Normal,
	FontSize:        1.0,
	CornerRadius:    5,
	FontColor:       theme.OnSurface,
	BackgroundColor: theme.Surface,
	BorderColor:     theme.Outline,
	BorderWidth:     1,
	Padding:         f32.Padding{L: 25, T: 15, R: 25, B: 15},
	Delay:           time.Millisecond * 800,
}

var dialogStartTime = time.Now()

func Exit() {
	Current = nil
}

func YesNoDialog(heading string, text string, lbl1, lbl2 string, on1, on2 func()) wid.Wid {
	return wid.Col(
		nil,
		wid.Separator(0, 25, theme.Transparent),
		wid.Label(heading, wid.H1C),
		wid.Separator(0, 12, theme.Transparent),
		wid.Label(text, nil),
		wid.Separator(0, 25, theme.Transparent),
		wid.Row(nil,
			btn.Btn(lbl1, nil, on1, nil, ""),
			btn.Btn(lbl2, nil, on2, nil, ""),
		),
	)
}

var Current wid.Wid //  = YesNoDialog("Heading", "Some text", "Yes", "No", nil, nil)

func Show(style *DialogueStyle) {
	if Current == nil {
		return
	}
	if style == nil {
		style = &DefaultDialogueStyle
	}

	// f goes from 0 to 0.5 after ca 0.5 second
	f := min(1.0, float32(time.Since(dialogStartTime))/float32(time.Second))
	if f < 1.0 {
		gpu.Invalidate(0)
	}
	// Draw surface all over the underlying form with the transparent surface color
	rw := f32.Rect{W: gpu.WindowWidthDp, H: gpu.WindowHeightDp}
	gpu.Rect(rw, 0, f32.Black.Alpha(f*0.5), f32.Transparent)
	// Draw dialog
	w := float32(300)
	h := float32(180)
	x := (gpu.WindowWidthDp - w) / 2
	y := (gpu.WindowHeightDp - h) / 2
	ctx := wid.Ctx{Rect: f32.Rect{X: x, Y: y, W: w, H: h}, Baseline: 0}
	gpu.RoundedRect(ctx.Rect, 10, 2, theme.Colors[style.BackgroundColor], f32.Transparent)
	ctx.Rect = ctx.Rect.Inset(style.Padding, 0)
	_ = Current(ctx)

}
