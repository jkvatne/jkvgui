package theme

import (
	"github.com/jkvatne/jkvgui/f32"
)

type UIRole uint8

const (
	Transparent UIRole = iota
	Surface            // Surface is the default surface for windows.
	OnSurface
	Primary // Primary is for prominent buttons, active states etc
	OnPrimary
	PrimaryContainer // PrimaryContainer is a light background ted with Primary color.
	OnPrimaryContainer
	Secondary // Secondary is for less prominent components
	OnSecondary
	SecondaryContainer // SecondaryContainer is a light background ted with Secondary color.
	OnSecondaryContainer
	Outline // Outline is used for frames and buttons
	SurfaceContainer
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

// SetDefaultPallete will set primary,secondary and neutral collors
// and initialize th colors
func SetDefaultPallete(light bool) {
	PrimaryColor = f32.FromRGB(0x6750A4)
	SecondaryColor = f32.FromRGB(0x625B71)
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
		Colors[Primary] = f32.Tone(PrimaryColor, 40)
		Colors[OnPrimary] = f32.Tone(PrimaryColor, 100)
		Colors[PrimaryContainer] = f32.Tone(PrimaryColor, 90)
		Colors[OnPrimaryContainer] = f32.Tone(PrimaryColor, 10)
		Colors[Secondary] = f32.Tone(SecondaryColor, 40)
		Colors[OnSecondary] = f32.Tone(SecondaryColor, 100)
		Colors[SecondaryContainer] = f32.Tone(SecondaryColor, 90)
		Colors[OnSecondaryContainer] = f32.Tone(SecondaryColor, 10)
		Colors[Outline] = f32.Tone(NeutralColor, 50)
		Colors[Surface] = f32.Tone(NeutralColor, 98)
		Colors[OnSurface] = f32.Tone(NeutralColor, 10)
		Colors[SurfaceContainer] = f32.Tone(NeutralColor, 90)
	} else {
		Colors[Primary] = f32.Tone(PrimaryColor, 80)
		Colors[OnPrimary] = f32.Tone(PrimaryColor, 20)
		Colors[PrimaryContainer] = f32.Tone(PrimaryColor, 30)
		Colors[OnPrimaryContainer] = f32.Tone(PrimaryColor, 90)
		Colors[Secondary] = f32.Tone(SecondaryColor, 80)
		Colors[OnSecondary] = f32.Tone(SecondaryColor, 20)
		Colors[SecondaryContainer] = f32.Tone(SecondaryColor, 30)
		Colors[OnSecondaryContainer] = f32.Tone(SecondaryColor, 90)
		Colors[Outline] = f32.Tone(NeutralColor, 60)
		Colors[Surface] = f32.Tone(NeutralColor, 4)
		Colors[OnSurface] = f32.Tone(NeutralColor, 90)
		Colors[SurfaceContainer] = f32.Tone(NeutralColor, 22)
	}
}
