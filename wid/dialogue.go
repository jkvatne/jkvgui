package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
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
	FontSize:        gpu.InitialSize * 0.5,
	FontColor:       f32.Color{0.0, 0.0, 0.0, 1.0},
	CornerRadius:    5,
	BorderColor:     f32.Color{R: 0.4, G: 0.4, B: 0.5, A: 1.0},
	BackgroundColor: f32.Color{R: 1.0, G: 1.0, B: 0.9, A: 1.0},
	BorderWidth:     1,
	Padding:         f32.Padding{3, 0, 1, 0},
	Delay:           time.Millisecond * 800,
}

var CurrentDialogue any
var dialogStartTime time.Time

func YesNoDialog(heading string, text string, lbl1, lbl2 string, on1, on2 func()) Wid {
	return Col(
		nil,
		Label("Mpqy", 13, nil, 0),
		Label("MpqyM", 24, nil, 4),
		Label("Mpqy", 13, nil, 4),
		/*
			Row(th, nil, SpaceRightAdjust,
				Button(lbl1, Do(on1)),
				Button(lbl2, Do(on2)),
			),
		*/
	)
}

func ShowDialogue(style *DialogueStyle) {
	if CurrentDialogue == nil {
		return
	}
	if style == nil {
		style = &DefaultDialogueStyle
	}
	// f goes from 0 to 1 after ca 0.5 second
	f := float32(min(1.0, float64(time.Since(dialogStartTime))/float64(time.Second/2)))
	// Draw surface all over the underlying form with the transparent surface color
	gpu.Rect(0, 0, float32(gpu.WindowWidthDp), float32(gpu.WindowHeightDp), 0, f32.WithAlpha(f32.Shade, f), f32.Transparent)

	scale := style.FontSize / gpu.InitialSize
	textHeight := (gpu.Fonts[style.FontNo].Ascent + gpu.Fonts[style.FontNo].Descent) * scale * 1.2

	w := textHeight * 8
	x := min(CurrentHint.Pos.X+w+style.Padding.L+style.Padding.R, gpu.WindowWidthDp)
	x = max(float32(0), x-w)

	lines := split(CurrentHint.Text, w-style.Padding.L-style.Padding.R, gpu.Fonts[style.FontNo], scale)
	gpu.Fonts[style.FontNo].SetColor(style.FontColor)

	h := textHeight*float32(len(lines)) + style.Padding.T + style.Padding.B
	y := min(CurrentHint.Pos.Y+h, gpu.WindowHeightDp)
	y = max(0, y-h)
	yb := y + style.Padding.T + textHeight
	gpu.RoundedRect(x, y, w, h, style.CornerRadius, style.BorderWidth, style.BackgroundColor, style.BorderColor, 5, 0)
	for _, line := range lines {
		gpu.Fonts[style.FontNo].Printf(
			x+style.Padding.L+style.Padding.L+style.BorderWidth,
			yb, style.FontSize,
			style.FontSize, line)
		yb = yb + style.FontSize
	}
	CurrentHint.Active = false
}

/*
func Dialog(th *Theme, role UIRole, widgets ...Wid) Wid {
	return func(gtx C) D {
		pt := Px(gtx, th.DialogPadding.Top)
		pb := Px(gtx, th.DialogPadding.Bottom)
		pl := Px(gtx, th.DialogPadding.Left)
		pr := Px(gtx, th.DialogPadding.Right)
		f := Min(1.0, float64(time.Since(dialogStartTime))/float64(time.Second/4))
		// Shade underlying form
		// Draw surface all over the underlying form with the transparent surface color
		outline := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		defer clip.Rect(outline).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, WithAlpha(Black, uint8(f*200)))
		// Calculate dialog constraints
		ctx := gtx
		ctx.Constraints.Min.Y = 0
		// Margins left and right for a constant maximum dialog size
		margin := Max(12, (ctx.Constraints.Max.X - pl - pr - Px(gtx, th.DialogTextWidth)))
		ctx.Constraints.Max.X = gtx.Constraints.Max.X - margin - pl - pr
		calls := make([]op.CallOp, len(widgets))
		dims := make([]D, len(widgets))
		size := 0
		for i, child := range widgets {
			macro := op.Record(gtx.Ops)
			dims[i] = child(ctx)
			calls[i] = macro.Stop()
			size += dims[i].Size.Y
		}
		mt := (gtx.Constraints.Max.Y - size) / 2
		// Calculate posision and size of the dialog box
		x := int(f*float64(margin/2) + (1-f)*float64(startX))
		y := int(f*float64(mt) + (1-f)*float64(startY))
		outline = image.Rect(0, 0, int(f*float64(gtx.Constraints.Max.X-margin)), int(f*float64(size+pt+pb)))
		// Draw the dialog surface with caclculated margins
		defer op.Offset(image.Pt(x, y)).Push(gtx.Ops).Pop()
		defer clip.UniformRRect(outline, Px(gtx, th.DialogCorners)).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, th.Bg[role])
		if f < 1.0 {
			// While animating, no widgets are drawn, but we invalidate to force a new redraw
			Invalidate()
		} else {
			// Now do the actual drawing of the widgets, with offsets
			y := pt
			for i := range widgets {
				trans := op.Offset(image.Pt(pl, int(math.Round(float64(y))))).Push(gtx.Ops)
				calls[i].Add(gtx.Ops)
				trans.Pop()
				y += dims[i].Size.Y
			}
		}
		sz := gtx.Constraints.Constrain(image.Pt(gtx.Constraints.Max.X, size+pb+pt))
		return D{Size: sz, Baseline: sz.Y}
	}
}
*/
