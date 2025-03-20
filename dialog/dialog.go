package dialog

import (
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
	BackgroundColor: theme.SurfaceContainer,
	BorderColor:     theme.Outline,
	BorderWidth:     1,
	Padding:         f32.Padding{25, 15, 25, 15},
	Delay:           time.Millisecond * 800,
}

var dialogStartTime time.Time = time.Now()

func Exit() {
	Current = nil
}

func YesNoDialog(heading string, text string, lbl1, lbl2 string, on1, on2 func()) wid.Wid {
	return wid.Col(
		nil,
		wid.Elastic(),
		wid.Row(wid.Distribute,
			wid.Elastic(),
			wid.Label(heading, wid.H1C),
			wid.Elastic(),
		),
		wid.Label(text, nil),
		wid.Elastic(),
		wid.Row(wid.Right,
			wid.Button(lbl1, on1, &wid.Btn, ""),
			wid.Button(lbl2, on2, &wid.Btn, ""),
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
	f := float32(min(1.0, float32(time.Since(dialogStartTime))/float32(time.Second)))
	if f < 1.0 {
		gpu.Invalidate(0)
	}
	// Draw surface all over the underlying form with the transparent surface color
	rw := f32.Rect{0, 0, float32(gpu.WindowWidthDp), float32(gpu.WindowHeightDp)}
	gpu.Rect(rw, 0, f32.MultAlpha(f32.Shade, f), f32.Transparent)
	// Draw dialog
	w := float32(300)
	h := float32(250)
	x := (gpu.WindowWidthDp - w) / 2
	y := (gpu.WindowHeightDp - h) / 2
	ctx := wid.Ctx{Rect: f32.Rect{X: x, Y: y, W: w, H: h}, Baseline: 0}
	gpu.RoundedRect(ctx.Rect, 10, 2, theme.Colors[style.BackgroundColor], f32.Transparent)
	ctx.Rect = ctx.Rect.Inset(style.Padding)
	_ = Current(ctx)

}
