package wid

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	utf8 "golang.org/x/exp/utf8string"
	"time"
)

type ComboStyle struct {
	FontSize           float32
	FontNo             int
	FontColor          f32.Color
	InsideColor        f32.Color
	BorderColor        f32.Color
	BorderWidth        float32
	BorderCornerRadius float32
	InsidePadding      f32.Padding
	OutsidePadding     f32.Padding
	CursorWidth        float32
}

type ComboState struct {
	SelStart int
	SelEnd   int
	Buffer   utf8.String
	index    int
	expanded bool
}

var DefaultCombo = ComboStyle{
	FontSize:           1.5,
	FontNo:             gpu.DefaultFont,
	InsideColor:        f32.Color{1.0, 1.0, 1.0, 1.0},
	BorderColor:        f32.Color{0, 0, 0, 1},
	FontColor:          f32.Color{0, 0, 0, 1},
	OutsidePadding:     f32.Padding{4, 4, 4, 4},
	InsidePadding:      f32.Padding{5, 2, 2, 2},
	BorderWidth:        1,
	BorderCornerRadius: 5,
	CursorWidth:        1.5,
}

type Theme struct {
	FontSize        float32
	FontNo          int
	BackgroundColor f32.Color
	FontColor       f32.Color
	PrimaryColor    f32.Color
	SecondaryColor  f32.Color
	SurfaceColor    f32.Color
}

var DefaultTheme = Theme{
	FontSize: 1.0,
}

func NewComboStyle(th *Theme) *ComboStyle {
	return &ComboStyle{
		FontSize: th.FontSize,
		FontNo:   th.FontNo,
	}
}

func setValue(i int, s *ComboState, list []string) {
	s.index = i
	s.Buffer.Init(list[i])
	s.expanded = false
	gpu.Invalidate(0)
}

var ComboStateMap = make(map[*string]*ComboState)

func Combo(text *string, list []string, style *ComboStyle) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultCombo
		}
		s := ComboStateMap[text]
		if s == nil {
			ComboStateMap[text] = &ComboState{}
			s = ComboStateMap[text]
			s.Buffer.Init(*text)
		}
		focused := focus.At(text)

		r := ctx.Rect.Inset(style.OutsidePadding)
		dho := style.OutsidePadding.T + style.OutsidePadding.B
		dhi := style.InsidePadding.T + style.InsidePadding.B + 2*style.BorderWidth
		dwi := style.InsidePadding.L + style.InsidePadding.R + 2*style.BorderWidth
		dwo := style.OutsidePadding.R + style.OutsidePadding.L
		fh := font.Fonts[style.FontNo].Height(style.FontSize)
		height := fh + dho + dhi
		width := r.W + dwo
		baseline := font.Fonts[style.FontNo].Baseline(style.FontSize) + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth
		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		}
		f := font.Fonts[style.FontNo]
		f.SetColor(style.FontColor)
		rpd := f32.Rect{
			ctx.Rect.X + ctx.Rect.W - style.OutsidePadding.R - style.BorderWidth - style.InsidePadding.R - fh,
			ctx.Rect.Y + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth,
			fh,
			fh,
		}
		col := style.InsideColor
		outline := ctx.Rect
		outline.W = width
		outline.H = height

		if focus.LeftMouseBtnReleased(rpd) {
			s.expanded = !s.expanded
			gpu.Invalidate(0)
			focus.Set(text)
			focused = true

		}

		if !focused {
			s.expanded = false
		}
		if s.expanded {
			if gpu.LastKey == glfw.KeyDown {
				s.index = min(s.index+1, len(list)-1)
			} else if gpu.LastKey == glfw.KeyUp {
				s.index = max(s.index-1, 0)
			} else if gpu.LastKey == glfw.KeyEnter {
				setValue(s.index, s, list)
				gpu.LastKey = 0
			}

			lh := fh + style.InsidePadding.T + style.InsidePadding.B
			boxHeight := float32(len(list)) * lh
			r := f32.Rect{
				ctx.Rect.X + style.OutsidePadding.L,
				ctx.Rect.Y + height - style.OutsidePadding.B,
				width - style.OutsidePadding.L - style.OutsidePadding.R,
				boxHeight}
			for i := range len(list) {
				itemRect := f32.Rect{r.X, r.Y + float32(i)*lh, r.W, lh}
				if focus.LeftMouseBtnReleased(itemRect) {
					setValue(i, s, list)
				}
			}

			dropDownBox := func() {

				// Box surrounding the droptown items
				gpu.RoundedRect(r, 0, 1, f32.White, f32.Black)
				baseline := font.Fonts[style.FontNo].Baseline(style.FontSize) + style.InsidePadding.T

				for i := range len(list) {
					bgColor := f32.White
					if i == s.index {
						bgColor = f32.LightGrey
					}
					y := height + ctx.Rect.Y + float32(i)*lh - style.OutsidePadding.B
					box := f32.Rect{r.X, r.Y + float32(i)*lh, r.W, lh}
					gpu.RoundedRect(box, 0, 0.5, bgColor, f32.Black)
					x := r.X + style.InsidePadding.L + style.BorderWidth
					f.Printf(x, y+baseline, style.FontSize, r.W-dwi-style.BorderWidth*2-fh, list[i])
				}
			}
			gpu.Defer(dropDownBox)
		}
		if focused {
			col.A *= 0.3
			gpu.Invalidate(111 * time.Millisecond)
			if gpu.LastRune != 0 {
				s1 := s.Buffer.Slice(0, s.SelStart)
				s2 := s.Buffer.Slice(s.SelEnd, s.Buffer.RuneCount())
				s.Buffer.Init(s1 + string(gpu.LastRune) + s2)
				gpu.LastRune = 0
				s.SelStart++
				s.SelEnd++
			}
			if gpu.LastKey == glfw.KeyBackspace {
				str := *text
				*text = str[0 : len(str)-1]
			} else if gpu.LastKey == glfw.KeyLeft {
				s.SelStart--
				s.SelStart = max(0, s.SelStart)
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyRight {
				s.SelStart++
				s.SelStart = min(s.SelStart, s.Buffer.RuneCount())
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyEnd {
				s.SelStart = s.Buffer.RuneCount()
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyHome {
				s.SelStart = 0
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyEnter {
				if s.expanded {
					setValue(s.index, s, list)
				} else {
					s.expanded = true
				}
				gpu.Invalidate(0)
			} else if gpu.LastKey != 0 {
				gpu.Invalidate(0)
			}
		} else if focus.Hovered(outline) {
			col.A *= 0.1
		}

		focus.Move(text)
		if focus.LeftMouseBtnPressed(outline) {
			gpu.Invalidate(0)
			col.A = 1
		}

		if focus.LeftMouseBtnReleased(outline) {
			halfUnit = time.Now().UnixMilli() % 333
			focus.Set(text)
			s.SelStart = f.RuneNo(focus.MousePos.X-(r.X), style.FontSize, s.Buffer.String())
			s.SelEnd = s.SelStart
			focus.MouseBtnReleased = false
			gpu.Invalidate(0)
		}

		gpu.RoundedRect(r, style.BorderCornerRadius, style.BorderWidth, col, style.BorderColor)
		f.SetColor(style.FontColor)
		x := ctx.Rect.X + style.OutsidePadding.L + style.InsidePadding.L + style.BorderWidth
		f.Printf(
			x,
			ctx.Rect.Y+baseline,
			style.FontSize,
			r.W-dwi-style.BorderWidth*2-fh,
			s.Buffer.String())
		f.SetColor(f32.Black)
		s.SelStart = max(0, s.Buffer.RuneCount())
		s.SelEnd = s.SelStart
		if focused && (time.Now().UnixMilli()-halfUnit)/333&1 == 1 {
			dx := f.Width(style.FontSize, s.Buffer.Slice(0, s.SelStart))
			gpu.VertLine(x+dx, r.Y+style.InsidePadding.T, r.Y+baseline, 1, f32.Black)
		}

		DrawIcon(
			ctx.Rect.X+ctx.Rect.W-style.OutsidePadding.R-style.BorderWidth-style.InsidePadding.R-fh,
			ctx.Rect.Y+style.OutsidePadding.T+style.InsidePadding.T+style.BorderWidth,
			fh,
			ArrowDropDown, f32.Black)

		return Dim{w: width, h: height, baseline: baseline}
	}
}
