package theme

import (
	"github.com/jkvatne/jkvgui/f32"
)

type UIRole uint8

// The odd roles are foreground colors, and the even are background.
// This means Primary.Bg() is the same as OnPrimary
// The zero value of UiRole is transparent both foreground and background.
const (
	// Transparent has Black/WHite as foreground, used for drawing on a transparent background
	Transparent   UIRole = iota
	TransparentFg UIRole = iota
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
	// Outline is used for frames and buttons
	Outline
	OnOutline
	// SurfaceContainer is a darker variant of Surface
	SurfaceContainer
	OnSurfaceContainer
	// Error          // Error is usualy red
	// ErrorContainer // ErrorContainer is usualy light red
	// OnErrorContainer
	// OutlineVariant
	// SurfaceContainerHigh // SurfaceContainerHighest is the grayest surface
	// SurfaceContainerLow // SurfaceContainerLowest is almost white/black
)

var Colors [32]f32.Color

var (
	PrimaryColor   f32.Color
	SecondaryColor f32.Color
	NeutralColor   f32.Color
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

// SetDefaultPallete will set primary,secondary and neutral colors
// and initialize th colors
func SetDefaultPallete(light bool) {
	PrimaryColor = f32.FromRGB(0x5750C4)
	SecondaryColor = f32.FromRGB(0x925B51)
	NeutralColor = f32.FromRGB(0x79747E)
	SetupColors(light)
}

func SetPallete(light bool, p, s, n uint32) {
	PrimaryColor = f32.FromRGB(p)
	SecondaryColor = f32.FromRGB(s)
	NeutralColor = f32.FromRGB(n)
	SetupColors(light)
}

func SetupColors(light bool) {
	if light {
		Colors[Transparent] = f32.Transparent
		Colors[TransparentFg] = NeutralColor.Tone(0)
		Colors[Primary] = PrimaryColor.Tone(40)
		Colors[OnPrimary] = PrimaryColor.Tone(100)
		Colors[PrimaryContainer] = PrimaryColor.Tone(90)
		Colors[OnPrimaryContainer] = PrimaryColor.Tone(10)
		Colors[Secondary] = SecondaryColor.Tone(40)
		Colors[OnSecondary] = SecondaryColor.Tone(100)
		Colors[SecondaryContainer] = SecondaryColor.Tone(90)
		Colors[OnSecondaryContainer] = SecondaryColor.Tone(10)
		Colors[Outline] = NeutralColor.Tone(30)
		Colors[OnOutline] = NeutralColor.Tone(30)
		Colors[Surface] = NeutralColor.Tone(98)
		Colors[OnSurface] = NeutralColor.Tone(10)
		Colors[SurfaceContainer] = NeutralColor.Tone(90)
		Colors[OnSurfaceContainer] = NeutralColor.Tone(0)
	} else {
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
		Colors[Outline] = NeutralColor.Tone(80)
		Colors[OnOutline] = NeutralColor.Tone(80)
		Colors[Surface] = NeutralColor.Tone(8)
		Colors[OnSurface] = NeutralColor.Tone(90)
		Colors[SurfaceContainer] = NeutralColor.Tone(20)
		Colors[OnSurfaceContainer] = NeutralColor.Tone(90)
	}
}
