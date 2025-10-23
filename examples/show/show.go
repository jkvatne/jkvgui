package main

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/sys"
)

func ShowIcons(Gd gpu.GlData, x float32, y float32) {
	Gd.DrawIcon(x+25, y, 24, gpu.ArrowDropDown, f32.Black)
	Gd.DrawIcon(x+50, y, 24, gpu.Home, f32.Black)
	Gd.DrawIcon(x+75, y, 24, gpu.BoxChecked, f32.Black)
	Gd.DrawIcon(x+100, y, 24, gpu.BoxUnchecked, f32.Black)
	Gd.DrawIcon(x+125, y, 24, gpu.RadioChecked, f32.Black)
	Gd.DrawIcon(x+150, y, 24, gpu.RadioUnchecked, f32.Black)
	Gd.DrawIcon(x+175, y, 24, gpu.ContentOpen, f32.Black)
	Gd.DrawIcon(x+200, y, 24, gpu.ContentSave, f32.Black)
	Gd.DrawIcon(x+225, y, 24, gpu.NavigationArrowDownward, f32.Black)
	Gd.DrawIcon(x+250, y, 24, gpu.NavigationArrowUpward, f32.Black)
	Gd.DrawIcon(x+275, y, 24, gpu.NavigationUnfoldMore, f32.Black)
	Gd.DrawIcon(x+300, y, 24, gpu.NavigationArrowDropDown, f32.Black)
	Gd.DrawIcon(x+325, y, 24, gpu.NavigationArrowDropUp, f32.Black)
}

func ShowFonts(Gd gpu.GlData, x float32, y float32) {
	for _, f := range font.Fonts {
		if f != nil {
			f.DrawText(Gd, x, y, f32.Black, 0, gpu.LTR, strconv.Itoa(f.No)+" "+f.Name+" "+strconv.Itoa(f.Size))
			y += 25
		}
	}
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	slog.Info("Show fonts and icons")
	sys.Init()
	defer sys.Shutdown()
	w := sys.CreateWindow(0, 0, 0, 0, "Fonts and images", 1, 2.0)
	for sys.Running() {
		w.StartFrame()
		ShowIcons(w.Gd, 0, 10)
		ShowFonts(w.Gd, 10, 100)

		font.Fonts[gpu.Normal14].DrawText(w.Gd, 400, 250, f32.Black, 250, gpu.BTT, "BottomToTopBottomToTop")
		font.Fonts[gpu.Normal14].DrawText(w.Gd, 400, 100, f32.Black, 250, gpu.TTB, "TopToBottomTopToBottom")

		for i := range 14 {
			ww := float32(i)*5.0 + 120
			x := float32(450)
			y := 100 + float32(i)*15
			font.Fonts[gpu.Normal14].DrawText(w.Gd, x, y, f32.Black, ww, gpu.LTR, "TruncatedTruncatedTruncatedTruncated")
			w.Gd.VertLine(x+ww, y-15, y, 1, f32.Blue)
		}
		font.Fonts[gpu.Normal14].DrawText(w.Gd, 400, 25, f32.Black, 0, gpu.LTR, fmt.Sprintf("FPS=%0.1f", w.Fps()))

		font.Fonts[gpu.Bold20].DrawText(w.Gd, 350, 350, f32.Black, 0, 0, "R")
		font.Fonts[gpu.Bold20].DrawText(w.Gd, 370, 350, f32.Black, 0, 1, "R")
		font.Fonts[gpu.Bold20].DrawText(w.Gd, 390, 350, f32.Black, 0, 2, "R")
		font.Fonts[gpu.Bold20].DrawText(w.Gd, 410, 350, f32.Black, 0, 3, "R")
		font.Fonts[gpu.Bold20].DrawText(w.Gd, 430, 350, f32.Black, 0, 4, "R")
		font.Fonts[gpu.Bold20].DrawText(w.Gd, 450, 350, f32.Black, 0, 5, "R")
		font.Fonts[gpu.Bold20].DrawText(w.Gd, 470, 350, f32.Black, 0, 6, "R")
		font.Fonts[gpu.Bold20].DrawText(w.Gd, 490, 350, f32.Black, 0, 7, "R")

		w.EndFrame()
		sys.PollEvents()
	}
}
