package theme

import (
	"github.com/jkvatne/jkvgui/f32"
)

type UIRole uint8

// The odd roles are foreground colors, and the even are background.
// This means Primary.Bg() is the same as OnPrimary
const (
	// Transparent has Black/White as foreground, used for drawing on a transparent background
	Transparent UIRole = iota
	TransparentFg
	// Canvas is the extreme white/black colors
	Canvas
	OnCanvas
	// Surface is the default surface for windows.
	Surface
	OnSurface
	// Primary is for prominent buttons, active states etc
	Primary
	OnPrimary
	// PrimaryContainer is a light background ted with Primary color.
	PrimaryContainer
	OnPrimaryContainer
	// Secondary is for less prominent components and for variation
	Secondary
	OnSecondary
	SecondaryContainer // SecondaryContainer is a light background ted with Secondary color.
	OnSecondaryContainer
	// Tertiary is for less prominent components and for variation
	Tertiary
	OnTertiary
	TertiaryContainer // SecondaryContainer is a light background ted with Secondary color.
	OnTertiaryContainer
	// Outline is used for frames and buttons
	Outline
	OnOutline
	// SurfaceContainer is a darker variant of Surface
	SurfaceContainer
	OnSurfaceContainer
	// Error is usually red
	Error
	OnError
	// ErrorContainer is usually light red
	ErrorContainer
	OnErrorContainer
)

var (
	PrimaryColor   f32.Color
	SecondaryColor f32.Color
	TertiaryColor  f32.Color
	NeutralColor   f32.Color
	ErrorColor     f32.Color
	Colors         [32]f32.Color
)

func (u UIRole) Color() f32.Color {
	return Colors[u]
}

func (u UIRole) Bg() f32.Color {
	return Colors[u&0xFE]
}

func (u UIRole) Fg() f32.Color {
	return Colors[u|1]
}

// SetDefaultPalette will set primary, secondary, error and neutral colors
func SetDefaultPalette(light bool) {
	PrimaryColor = f32.FromRGB(0x5750C4)
	SecondaryColor = f32.FromRGB(0x925B51)
	TertiaryColor = f32.FromRGB(0x425B51)
	NeutralColor = f32.FromRGB(0x79747E)
	ErrorColor = f32.FromRGB(0xFF4242)
	SetupColors(light)
}

// SetPalette can be used to set all four base colors at once using the hex codes.
func SetPalette(light bool, p, s, t, n, e uint32) {
	PrimaryColor = f32.FromRGB(p)
	SecondaryColor = f32.FromRGB(s)
	TertiaryColor = f32.FromRGB(t)
	NeutralColor = f32.FromRGB(n)
	ErrorColor = f32.FromRGB(e)
	SetupColors(light)
}

// SetupColors will set the colors for the theme depending on light/dark mode.
func SetupColors(light bool) {
	if light {
		Colors[Canvas] = f32.White
		Colors[OnCanvas] = f32.Black
		Colors[TransparentFg] = NeutralColor.Tone(0)
		Colors[Transparent] = f32.Transparent
		Colors[TransparentFg] = NeutralColor.Tone(0)
		Colors[Primary] = PrimaryColor.Tone(40)
		Colors[OnPrimary] = PrimaryColor.Tone(100)
		Colors[PrimaryContainer] = PrimaryColor.Tone(90)
		Colors[OnPrimaryContainer] = PrimaryColor.Tone(10)
		Colors[Secondary] = SecondaryColor.Tone(40)
		Colors[OnSecondary] = SecondaryColor.Tone(100)
		Colors[SecondaryContainer] = SecondaryColor.Tone(85)
		Colors[OnSecondaryContainer] = SecondaryColor.Tone(10)
		Colors[Tertiary] = TertiaryColor.Tone(40)
		Colors[OnTertiary] = TertiaryColor.Tone(100)
		Colors[TertiaryContainer] = TertiaryColor.Tone(90)
		Colors[OnTertiaryContainer] = TertiaryColor.Tone(10)
		Colors[Outline] = NeutralColor.Tone(10)
		Colors[OnOutline] = NeutralColor.Tone(50)
		Colors[Surface] = NeutralColor.Tone(98)
		Colors[OnSurface] = NeutralColor.Tone(10)
		Colors[SurfaceContainer] = NeutralColor.Tone(90)
		Colors[OnSurfaceContainer] = NeutralColor.Tone(0)
		Colors[Error] = ErrorColor.Tone(40)
		Colors[OnError] = ErrorColor.Tone(100)
		Colors[ErrorContainer] = ErrorColor.Tone(80)
		Colors[OnErrorContainer] = ErrorColor.Tone(0)
	} else {
		Colors[Canvas] = f32.Black
		Colors[OnCanvas] = f32.White
		Colors[Transparent] = f32.Transparent
		Colors[TransparentFg] = NeutralColor.Tone(100)
		Colors[Primary] = PrimaryColor.Tone(80)
		Colors[OnPrimary] = PrimaryColor.Tone(20)
		Colors[PrimaryContainer] = PrimaryColor.Tone(30)
		Colors[OnPrimaryContainer] = PrimaryColor.Tone(90)
		Colors[Secondary] = SecondaryColor.Tone(80)
		Colors[OnSecondary] = SecondaryColor.Tone(20)
		Colors[SecondaryContainer] = SecondaryColor.Tone(30)
		Colors[OnSecondaryContainer] = SecondaryColor.Tone(90)
		Colors[Tertiary] = TertiaryColor.Tone(80)
		Colors[OnTertiary] = TertiaryColor.Tone(20)
		Colors[TertiaryContainer] = TertiaryColor.Tone(30)
		Colors[OnTertiaryContainer] = TertiaryColor.Tone(90)
		Colors[Outline] = NeutralColor.Tone(80)
		Colors[OnOutline] = NeutralColor.Tone(50)
		Colors[Surface] = NeutralColor.Tone(8)
		Colors[OnSurface] = NeutralColor.Tone(90)
		Colors[SurfaceContainer] = NeutralColor.Tone(20)
		Colors[OnSurfaceContainer] = NeutralColor.Tone(90)
		Colors[Error] = ErrorColor.Tone(80)
		Colors[OnError] = ErrorColor.Tone(20)
		Colors[ErrorContainer] = ErrorColor.Tone(20)
		Colors[OnErrorContainer] = ErrorColor.Tone(90)

	}
}
