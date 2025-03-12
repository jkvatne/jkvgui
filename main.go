package main

import (
	_ "embed"
	"fmt"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/wid"
	"github.com/jkvatne/jkvgui/svg"
	"log"
	"strings"
)

//go:embed gpu/fonts/Roboto-Thin.ttf
var Roboto100 []byte // 100

//go:embed gpu/fonts/Roboto-ExtraLight.ttf
var Roboto200 []byte // 200

//go:embed gpu/fonts/Roboto-Light.ttf
var Roboto300 []byte // 300

//go:embed gpu/fonts/Roboto-Regular.ttf
var Roboto400 []byte // 400

//go:embed gpu/fonts/Roboto-Medium.ttf
var Roboto500 []byte // 500

//go:embed gpu/fonts/Roboto-SemiBold.ttf
var Roboto600 []byte // 600

//go:embed gpu/fonts/Roboto-Bold.ttf
var Roboto700 []byte // 700

//go:embed gpu/fonts/Roboto-Bold.ttf
var Roboto800 []byte // 800

//go:embed gpu/fonts/Roboto-Bold.ttf
var Roboto900 []byte // 900

//go:embed gpu/fonts/RobotoMono-Regular.ttf
var RobotoMono400 []byte

//go:embed gpu/fonts/RobotoMono-Bold.ttf
var RobotoMono700 []byte

//go:embed gpu/fonts/RobotoMono-Light.ttf
var RobotoMono300 []byte

var P = f32.Padding{2, 2, 2, 2}

func YesBtnClick() {
	log.Printf("Yes Btn Click\n")
}
func CancelBtnClick() {
	log.Printf("Cancel Btn Click\n")
}

func NoBtnClick() {
	log.Printf("No Btn Click\n")
}

var name = "jkvgui"
var address = "Mo i Rana"
var hint = "This is a hint word5 word6 word7 word8 qYyM9 qYyM10"

func Form() wid.Wid {
	r := wid.RowSetup{}
	w := wid.Row(r,
		wid.Label("Mpqy", 13, P, 4),
		wid.Label("MpqyM", 24, P, 4),
		wid.Elastic(),
		wid.Button("Cancel", CancelBtnClick, wid.OkBtn, hint),
		wid.Button("Noqy", NoBtnClick, wid.OkBtn, hint),
		wid.Button("Yes", YesBtnClick, wid.OkBtn, hint),
		wid.Edit(&name, nil, &wid.DefaultEdit),
		wid.Edit(&address, nil, &wid.DefaultEdit),
	)
	return w
}

func Draw() {
	// Calculate sizes
	form := Form()
	ctx := wid.Ctx{Rect: f32.Rect{X: 50, Y: 450, H: 260, W: 800}, Baseline: 0}
	_ = form(ctx)
}

func main() {

	svg := `
		<svg width="100" height="100">
			<circle cx="50" cy="50" r="40" fill="red" />
		</svg>
	`
	reader := strings.NewReader(svg)

	element, _ := svg.Parse(reader, false)

	fmt.Printf("SVG width: %s", element.Attributes["width"])
	fmt.Printf("Circle fill: %s", element.Children[0].Attributes["fill"])

	window := gpu.InitWindow(2508, 1270, "Rounded rectangle demo", 1, f32.Lightgrey)
	defer gpu.Shutdown()

	R0 := gpu.LoadFont(Roboto100, gpu.InitialSize, "Roboto", 100)
	R1 := gpu.LoadFont(Roboto200, gpu.InitialSize, "Roboto", 200)
	R2 := gpu.LoadFont(Roboto300, gpu.InitialSize, "Roboto", 300)
	R3 := gpu.LoadFont(Roboto400, gpu.InitialSize, "Roboto", 400)
	R4 := gpu.LoadFont(Roboto500, gpu.InitialSize, "Roboto", 500)
	R5 := gpu.LoadFont(Roboto600, gpu.InitialSize, "Roboto", 600)
	R6 := gpu.LoadFont(Roboto700, gpu.InitialSize, "Roboto", 700)
	R7 := gpu.LoadFont(Roboto800, gpu.InitialSize, "Roboto", 800)
	R8 := gpu.LoadFont(RobotoMono300, gpu.InitialSize, "RobotoMono", 300)
	R9 := gpu.LoadFont(RobotoMono400, gpu.InitialSize, "RobotoMono", 400)
	R10 := gpu.LoadFont(RobotoMono700, gpu.InitialSize, "RobotoMono", 700)
	wid.LoadIcons()

	w, h := window.GetSize()
	gpu.SizeCallback(window, w, h)
	for !window.ShouldClose() {
		gpu.StartFrame()
		gpu.Fonts[R0].Printf(50, 100, 24, 0, "24 Roboto100")
		gpu.Fonts[R1].Printf(50, 130, 24, 0, "24 Roboto200")
		gpu.Fonts[R2].Printf(50, 160, 24, 0, "24 Roboto300")
		gpu.Fonts[R3].Printf(50, 190, 24, 0, "24 Roboto400") // Regular
		gpu.Fonts[R4].Printf(50, 220, 24, 0, "24 Roboto500")
		gpu.Fonts[R5].Printf(50, 250, 24, 0, "24 Roboto600")
		gpu.Fonts[R6].Printf(50, 280, 24, 0, "24 Roboto700")
		gpu.Fonts[R7].Printf(50, 310, 24, 0, "24 Roboto800")
		gpu.Fonts[R8].Printf(50, 340, 24, 0, "24 RobotoMono300")
		gpu.Fonts[R9].Printf(50, 370, 24, 0, "24 RobotoMono400")  // Regular
		gpu.Fonts[R10].Printf(50, 400, 24, 0, "24 RobotoMono700") // Bold

		gpu.Fonts[0].Printf(350, 100, 12, 0, "12 Roboto100")
		gpu.Fonts[1].Printf(350, 130, 12, 0, "12 Roboto200")
		gpu.Fonts[2].Printf(350, 160, 12, 0, "12 Roboto300")
		gpu.Fonts[3].Printf(350, 190, 12, 0, "12 Roboto400  PCAN error 512") // Regular
		gpu.Fonts[4].Printf(350, 220, 12, 0, "12 Roboto500")
		gpu.Fonts[5].Printf(350, 250, 12, 0, "12 Roboto600")
		gpu.Fonts[6].Printf(350, 280, 12, 0, "12 Roboto700")
		gpu.Fonts[7].Printf(350, 310, 12, 0, "12 Roboto800")
		gpu.Fonts[8].Printf(350, 340, 12, 0, "12 RobotoMono300")
		gpu.Fonts[9].Printf(350, 370, 12, 0, "12 RobotoMono400")  // Regular
		gpu.Fonts[10].Printf(350, 400, 12, 0, "12 RobotoMono700") // Bold

		// Paint a red frame around the whole window
		gpu.Rect(10, 10, float32(gpu.WindowWidthDp)-20, float32(gpu.WindowHeightDp)-20, 2, f32.Transparent, f32.Red)
		// Draw the screen widgets
		Draw()
		// Show hints if any is active
		wid.ShowHint(nil)

		wid.DrawIcon(50, 420, 24, wid.Home, f32.Blue)
		wid.DrawIcon(75, 420, 24, wid.BoxChecked, f32.Black)
		wid.DrawIcon(100, 420, 24, wid.BoxUnchecked, f32.Black)
		wid.DrawIcon(125, 420, 24, wid.RadioChecked, f32.Black)
		wid.DrawIcon(150, 420, 24, wid.RadioUnchecked, f32.Black)
		wid.DrawIcon(175, 420, 24, wid.ContentSave, f32.Black)
		wid.DrawIcon(200, 420, 24, wid.NavigationArrowDownward, f32.Black)
		wid.DrawIcon(225, 420, 24, wid.NavigationArrowUpward, f32.Black)
		wid.DrawIcon(250, 420, 24, wid.NavigationUnfoldMore, f32.Black)
		wid.DrawIcon(275, 420, 24, wid.NavigationArrowDropDown, f32.Black)
		wid.DrawIcon(300, 420, 24, wid.NavigationArrowDropUp, f32.Black)
		wid.DrawIcon(325, 420, 24, wid.ContentCreate, f32.Black)

		gpu.EndFrame(30, window)
	}
}
