package glfw

type Key int

// ModifierKey corresponds to a modifier key.
type ModifierKey int

// Modifier keys.
const (
	ModShift    ModifierKey = 1
	ModControl  ModifierKey = 2
	ModAlt      ModifierKey = 4
	ModSuper    ModifierKey = 8
	ModCapsLock ModifierKey = 16
	ModNumLock  ModifierKey = 32
)

// Action types.
const (
	Release Action = 0 // The key or button was released.
	Press   Action = 1 // The key or button was pressed.
	Repeat  Action = 2 // The key was held down until it repeated.
)

/* Printable keys */
const (
	KeySpace        = 32
	KeyApostrophe   = 39 /* ' */
	KeyComma        = 44 /* , */
	KeyMinus        = 45 /* - */
	KeyPeriode      = 46 /* . */
	KeySlash        = 47 /* / */
	Key0            = 48
	Key1            = 49
	Key2            = 50
	Key3            = 51
	Key4            = 52
	Key5            = 53
	Key6            = 54
	Key7            = 55
	Key8            = 56
	Key9            = 57
	KeySemicolon    = 59 /* ; */
	KeyEqual        = 61 /* = */
	KeyA            = 65
	KeyB            = 66
	KeyC            = 67
	KeyD            = 68
	KeyE            = 69
	KeyF            = 70
	KeyG            = 71
	KeyH            = 72
	KeyI            = 73
	KeyJ            = 74
	KeyK            = 75
	KeyL            = 76
	KeyM            = 77
	KeyN            = 78
	KeyO            = 79
	KeyP            = 80
	KeyQ            = 81
	KeyR            = 82
	KeyS            = 83
	KeyT            = 84
	KeyU            = 85
	KeyV            = 86
	KeyW            = 87
	KeyX            = 88
	KeyY            = 89
	KeyZ            = 90
	KeyLeftBracket  = 91  /* [ */
	KeyBackslash    = 92  /* \ */
	KeyRightBracket = 93  /* ] */
	KeyGraveAccent  = 96  /* ` */
	KeyWorld1       = 161 /* non-US #1 */
	KeyWorld2       = 162 /* non-US #2 */

	/* Function keys */
	KeyEscape       = 256
	KeyEnter        = 257
	KeyTab          = 258
	KeyBackspace    = 259
	KeyInsert       = 260
	KeyDelete       = 261
	KeyRight        = 262
	KeyLeft         = 263
	KeyDown         = 264
	KeyUp           = 265
	KeyPageUp       = 266
	KeyPageDown     = 267
	KeyHome         = 268
	KeyEnd          = 269
	KeyCapsLock     = 280
	KeyScrollLock   = 281
	KeyNumLock      = 282
	KeyPrintScreen  = 283
	KeyPause        = 284
	KeyF1           = 290
	KeyF2           = 291
	KeyF3           = 292
	KeyF4           = 293
	KeyF5           = 294
	KeyF6           = 295
	KeyF7           = 296
	KeyF8           = 297
	KeyF9           = 298
	KeyF10          = 299
	KeyF11          = 300
	KeyF12          = 301
	KeyKP_0         = 320
	KeyKP_1         = 321
	KeyKP_2         = 322
	KeyKP_3         = 323
	KeyKP_4         = 324
	KeyKP_5         = 325
	KeyKP_6         = 326
	KeyKP_7         = 327
	KeyKP_8         = 328
	KeyKP_9         = 329
	KeyKPDecimal    = 330
	KeyKPDivide     = 331
	KeyKPMultiply   = 332
	KeyKPSubtract   = 333
	KeyKPAdd        = 334
	KeyKPEnter      = 335
	KeyKPEqual      = 336
	KeyLeftShift    = 340
	KeyLeftControl  = 341
	KeyLeftAlt      = 342
	KeyLeftSuper    = 343
	KeyRightControl = 345
	KeyRightShift   = 344
	KeyRightAlt     = 346
	KeyRightSuper   = 347
	KeyMenu         = 348
	KeyLast         = KeyMenu
)

const (
	VK_CONTROL  = 0x11
	VK_LWIN     = 0x5B
	VK_MENU     = 0x12
	VK_RWIN     = 0x5C
	VK_SHIFT    = 0x10
	VK_SNAPSHOT = 0x2C
	VK_CAPITAL  = 0x14
	VK_NUMLOCK  = 0x90
)

// createKeyTables will generate the tables keycodes and scancodes (in _glfw.win32)
// They are used to translate between keycodes and scancodes.

func createKeyTables() {
	_glfw.win32.keycodes[0x00B] = Key0
	_glfw.win32.keycodes[0x002] = Key1
	_glfw.win32.keycodes[0x003] = Key2
	_glfw.win32.keycodes[0x004] = Key3
	_glfw.win32.keycodes[0x005] = Key4
	_glfw.win32.keycodes[0x006] = Key5
	_glfw.win32.keycodes[0x007] = Key6
	_glfw.win32.keycodes[0x008] = Key7
	_glfw.win32.keycodes[0x009] = Key8
	_glfw.win32.keycodes[0x00A] = Key9
	_glfw.win32.keycodes[0x01E] = KeyA
	_glfw.win32.keycodes[0x030] = KeyB
	_glfw.win32.keycodes[0x02E] = KeyC
	_glfw.win32.keycodes[0x020] = KeyD
	_glfw.win32.keycodes[0x012] = KeyE
	_glfw.win32.keycodes[0x021] = KeyF
	_glfw.win32.keycodes[0x022] = KeyG
	_glfw.win32.keycodes[0x023] = KeyH
	_glfw.win32.keycodes[0x017] = KeyI
	_glfw.win32.keycodes[0x024] = KeyJ
	_glfw.win32.keycodes[0x025] = KeyK
	_glfw.win32.keycodes[0x026] = KeyL
	_glfw.win32.keycodes[0x032] = KeyM
	_glfw.win32.keycodes[0x031] = KeyN
	_glfw.win32.keycodes[0x018] = KeyO
	_glfw.win32.keycodes[0x019] = KeyP
	_glfw.win32.keycodes[0x010] = KeyQ
	_glfw.win32.keycodes[0x013] = KeyR
	_glfw.win32.keycodes[0x01F] = KeyS
	_glfw.win32.keycodes[0x014] = KeyT
	_glfw.win32.keycodes[0x016] = KeyU
	_glfw.win32.keycodes[0x02F] = KeyV
	_glfw.win32.keycodes[0x011] = KeyW
	_glfw.win32.keycodes[0x02D] = KeyX
	_glfw.win32.keycodes[0x015] = KeyY
	_glfw.win32.keycodes[0x02C] = KeyZ

	_glfw.win32.keycodes[0x028] = KeyApostrophe
	_glfw.win32.keycodes[0x02B] = KeyBackslash
	_glfw.win32.keycodes[0x033] = KeyComma
	_glfw.win32.keycodes[0x00D] = KeyEqual
	_glfw.win32.keycodes[0x029] = KeyGraveAccent
	_glfw.win32.keycodes[0x01A] = KeyLeftBracket
	_glfw.win32.keycodes[0x00C] = KeyMinus
	_glfw.win32.keycodes[0x034] = KeyPeriode
	_glfw.win32.keycodes[0x01B] = KeyRightBracket
	_glfw.win32.keycodes[0x027] = KeySemicolon
	_glfw.win32.keycodes[0x035] = KeySlash
	_glfw.win32.keycodes[0x056] = KeyWorld2

	_glfw.win32.keycodes[0x00E] = KeyBackspace
	_glfw.win32.keycodes[0x153] = KeyDelete
	_glfw.win32.keycodes[0x14F] = KeyEnd
	_glfw.win32.keycodes[0x01C] = KeyEnter
	_glfw.win32.keycodes[0x001] = KeyEscape
	_glfw.win32.keycodes[0x147] = KeyHome
	_glfw.win32.keycodes[0x152] = KeyInsert
	_glfw.win32.keycodes[0x15D] = KeyMenu
	_glfw.win32.keycodes[0x151] = KeyPageDown
	_glfw.win32.keycodes[0x149] = KeyPageUp
	_glfw.win32.keycodes[0x045] = KeyPause
	_glfw.win32.keycodes[0x039] = KeySpace
	_glfw.win32.keycodes[0x00F] = KeyTab
	_glfw.win32.keycodes[0x03A] = KeyCapsLock
	_glfw.win32.keycodes[0x145] = KeyNumLock
	_glfw.win32.keycodes[0x046] = KeyScrollLock
	_glfw.win32.keycodes[0x03B] = KeyF1
	_glfw.win32.keycodes[0x03C] = KeyF2
	_glfw.win32.keycodes[0x03D] = KeyF3
	_glfw.win32.keycodes[0x03E] = KeyF4
	_glfw.win32.keycodes[0x03F] = KeyF5
	_glfw.win32.keycodes[0x040] = KeyF6
	_glfw.win32.keycodes[0x041] = KeyF7
	_glfw.win32.keycodes[0x042] = KeyF8
	_glfw.win32.keycodes[0x043] = KeyF9
	_glfw.win32.keycodes[0x044] = KeyF10
	_glfw.win32.keycodes[0x057] = KeyF11
	_glfw.win32.keycodes[0x058] = KeyF12
	_glfw.win32.keycodes[0x038] = KeyLeftAlt
	_glfw.win32.keycodes[0x01D] = KeyLeftControl
	_glfw.win32.keycodes[0x02A] = KeyLeftShift
	_glfw.win32.keycodes[0x15B] = KeyLeftSuper
	_glfw.win32.keycodes[0x137] = KeyPrintScreen
	_glfw.win32.keycodes[0x138] = KeyRightAlt
	_glfw.win32.keycodes[0x11D] = KeyRightControl
	_glfw.win32.keycodes[0x036] = KeyRightShift
	_glfw.win32.keycodes[0x15C] = KeyRightSuper
	_glfw.win32.keycodes[0x150] = KeyDown
	_glfw.win32.keycodes[0x14B] = KeyLeft
	_glfw.win32.keycodes[0x14D] = KeyRight
	_glfw.win32.keycodes[0x148] = KeyUp
	_glfw.win32.keycodes[0x052] = KeyKP_0
	_glfw.win32.keycodes[0x04F] = KeyKP_1
	_glfw.win32.keycodes[0x050] = KeyKP_2
	_glfw.win32.keycodes[0x051] = KeyKP_3
	_glfw.win32.keycodes[0x04B] = KeyKP_4
	_glfw.win32.keycodes[0x04C] = KeyKP_5
	_glfw.win32.keycodes[0x04D] = KeyKP_6
	_glfw.win32.keycodes[0x047] = KeyKP_7
	_glfw.win32.keycodes[0x048] = KeyKP_8
	_glfw.win32.keycodes[0x049] = KeyKP_9
	_glfw.win32.keycodes[0x04E] = KeyKPAdd
	_glfw.win32.keycodes[0x053] = KeyKPDecimal
	_glfw.win32.keycodes[0x135] = KeyKPDivide
	_glfw.win32.keycodes[0x11C] = KeyKPEnter
	_glfw.win32.keycodes[0x059] = KeyKPEqual
	_glfw.win32.keycodes[0x037] = KeyKPMultiply
	_glfw.win32.keycodes[0x04A] = KeyKPSubtract
	for scancode := int16(0); scancode < 512; scancode++ {
		if _glfw.win32.keycodes[scancode] > 0 {
			_glfw.win32.scancodes[_glfw.win32.keycodes[scancode]] = scancode
		}
	}
}
