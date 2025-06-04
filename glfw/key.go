package glfw

type ModifierKey int
type Key int

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
const KeySpace = 32
const KeyApostrophe = 39 /* ' */
const KeyComma = 44      /* , */
const KeyMINUS = 45      /* - */
const KeyPeriode = 46    /* . */
const KeySlash = 47      /* / */
const Key0 = 48
const Key1 = 49
const Key2 = 50
const Key3 = 51
const Key4 = 52
const Key5 = 53
const Key6 = 54
const Key7 = 55
const Key8 = 56
const Key9 = 57
const KeySemicolon = 59 /* ; */
const KeyEqual = 61     /* = */
const KeyA = 65
const KeyB = 66
const KeyC = 67
const KeyD = 68
const KeyE = 69
const KeyF = 70
const KeyG = 71
const KeyH = 72
const KeyI = 73
const KeyJ = 74
const KeyK = 75
const KeyL = 76
const KeyM = 77
const KeyN = 78
const KeyO = 79
const KeyP = 80
const KeyQ = 81
const KeyR = 82
const KeyS = 83
const KeyT = 84
const KeyU = 85
const KeyV = 86
const KeyW = 87
const KeyX = 88
const KeyY = 89
const KeyZ = 90
const KeyLeftBracket = 91  /* [ */
const KeyBackslash = 92    /* \ */
const KeyRightBracket = 93 /* ] */
const KeyGraveAccent = 96  /* ` */
const KeyWorld1 = 161      /* non-US #1 */
const KeyWorld2 = 162      /* non-US #2 */

/* Function keys */
const KeyEscape = 256
const KeyEnter = 257
const KeyTab = 258
const KeyBackspace = 259
const KeyInsert = 260
const KeyDelete = 261
const KeyRight = 262
const KeyLeft = 263
const KeyDown = 264
const KeyUp = 265
const KeyPageUp = 266
const KeyPageDown = 267
const KeyHome = 268
const KeyEnd = 269
const KeyCapsLock = 280
const KeyScropllLock = 281
const KeyNumLock = 282
const KeyPrintScreen = 283
const KeyPause = 284
const KeyF1 = 290
const KeyF2 = 291
const KeyF3 = 292
const KeyF4 = 293
const KeyF5 = 294
const KeyF6 = 295
const KeyF7 = 296
const KeyF8 = 297
const KeyF9 = 298
const KeyF10 = 299
const KeyF11 = 300
const KeyF12 = 301
const KeyF13 = 302
const KeyF14 = 303
const KeyF15 = 304
const KeyF16 = 305
const KeyF17 = 306
const KeyF18 = 307
const KeyF19 = 308
const KeyF20 = 309
const KeyF21 = 310
const KeyF22 = 311
const KeyF23 = 312
const KeyF24 = 313
const KeyF25 = 314
const KeyKP_0 = 320
const KeyKP_1 = 321
const KeyKP_2 = 322
const KeyKP_3 = 323
const KeyKP_4 = 324
const KeyKP_5 = 325
const KeyKP_6 = 326
const KeyKP_7 = 327
const KeyKP_8 = 328
const KeyKP_9 = 329
const KeyKPDecimal = 330
const KeyKPDivide = 331
const KeyKPMultiply = 332
const KeyKPSubtract = 333
const KeyKPAdd = 334
const KeyKPEnter = 335
const KeyKPEqual = 336
const KeyLeftShift = 340
const KeyLeftControl = 341
const KeyLeftAlt = 342
const KeyLeftSuper = 343
const KeyRightControl = 345
const KeyRightShift = 344
const KeyRightAlt = 346
const KeyRightSuper = 347
const KeyMenu = 348
const KeyLast = KeyMenu
