package dialog

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/wid"
	"time"
)

type DialogueStyle struct {
	FontNo          int
	FontSize        float32
	FontColor       f32.Color
	CornerRadius    float32
	BorderColor     f32.Color
	BackgroundColor f32.Color
	BorderWidth     float32
	Padding         f32.Padding
	Delay           time.Duration
}

var DefaultDialogueStyle = DialogueStyle{
	FontNo:          gpu.DefaultFont,
	FontSize:        1.0,
	FontColor:       f32.Color{0.0, 0.0, 0.0, 1.0},
	CornerRadius:    5,
	BorderColor:     f32.Color{R: 0.4, G: 0.4, B: 0.5, A: 1.0},
	BackgroundColor: f32.Shade,
	BorderWidth:     1,
	Padding:         f32.Padding{15, 15, 15, 15},
	Delay:           time.Millisecond * 800,
}

var dialogStartTime time.Time = time.Now()

func YesNoDialog(heading string, text string, lbl1, lbl2 string, on1, on2 func()) wid.Wid {
	return wid.Col(
		nil,
		wid.Row(nil,
			wid.Elastic(),
			wid.Label("Heading", 3, nil, 0),
			wid.Elastic(),
		),
		wid.Label("Some text", 1, nil, 0),
		wid.Label("More text", 1, nil, 0),
		wid.Row(nil,
			wid.Elastic(),
			wid.Button(lbl1, on1, wid.PrimaryBtn, ""),
			wid.Button(lbl2, on2, wid.PrimaryBtn, ""),
		),
	)
}

var CurrentDialogue wid.Wid = YesNoDialog("Heading", "Some text", "Yes", "No", nil, nil)

func Show(style *DialogueStyle) {
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

	// form := YesNoDialog("Heading", "Some text", "Yes", "No", nil, nil)
	ctx := wid.Ctx{Rect: f32.Rect{X: 200, Y: 100, W: 600, H: 400}, Baseline: 0}
	gpu.RoundedRect(ctx.Rect, 10, 2, f32.MultAlpha(f32.White, f), f32.Transparent, 10, 0.3)
	ctx.Rect = ctx.Rect.Inset(style.Padding)
	_ = CurrentDialogue(ctx)

}
