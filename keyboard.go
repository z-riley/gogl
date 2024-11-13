package turdgl

import (
	"github.com/veandco/go-sdl2/sdl"
)

// keyTracker is used to track which keys are pressed by the user and
// react accordingly.
type keyTracker struct {
	keyBindingsInstant map[sdl.Keycode]func()
	keyBindingsPress   map[sdl.Keycode]func()
	keyBindingsRelease map[sdl.Keycode]func()

	pressedKeys map[sdl.Keycode]struct{}
}

// newKeyTracker constructs a key tracker object.
func newKeyTracker() *keyTracker {
	return &keyTracker{
		keyBindingsInstant: make(map[sdl.Keycode]func()),
		keyBindingsPress:   make(map[sdl.Keycode]func()),
		keyBindingsRelease: make(map[sdl.Keycode]func()),
		pressedKeys:        make(map[sdl.Keycode]struct{}),
	}
}

// KeybindMode describes how a callback is executed after a bound key is pressed.
type KeybindMode int

const (
	Instantaneous KeybindMode = iota // callback executed continuously while is pressed
	KeyPress                         // callback executed on press
	KeyRelease                       // callback executed on release
)

// registerKeybind sets a callback function which is executed when a key is pressed.
// The callback is executed according to the KeybindMode.
func (k *keyTracker) registerKeybind(key sdl.Keycode, mode KeybindMode, callback func()) {
	switch mode {
	case Instantaneous:
		k.keyBindingsInstant[key] = callback
	case KeyPress:
		k.keyBindingsPress[key] = callback
	case KeyRelease:
		k.keyBindingsRelease[key] = callback
	}
}

// registerKeybind sets a callback function which is executed when a key is pressed.
// The callback can be executed always while the key is pressed, on press, or on release.
func (k *keyTracker) unregisterKeybind(key sdl.Keycode, mode KeybindMode) {
	switch mode {
	case Instantaneous:
		delete(k.keyBindingsInstant, key)
	case KeyPress:
		delete(k.keyBindingsPress, key)
	case KeyRelease:
		delete(k.keyBindingsRelease, key)
	}
}

// DropKeybinds unregisters all keybinds.
func (k *keyTracker) dropKeybinds() {
	k.keyBindingsInstant = make(map[sdl.Keycode]func())
	k.keyBindingsPress = make(map[sdl.Keycode]func())
	k.keyBindingsRelease = make(map[sdl.Keycode]func())
}

// Is pressed returns true if a key is currently pressed.
func (k *keyTracker) isPressed(key sdl.Keycode) bool {
	_, ok := k.pressedKeys[key]
	return ok
}

// handleEvent processes key press events and keeps track of pressed keys.
func (k *keyTracker) handleEvent(event *sdl.KeyboardEvent) {
	switch event.State {
	case sdl.PRESSED:
		// Execute on-press callback
		if event.Repeat == 0 {
			if fn, ok := k.keyBindingsPress[event.Keysym.Sym]; ok {
				fn()
			}
		}

		// Add key to tracker to handle instantaneous callbacks
		k.pressedKeys[event.Keysym.Sym] = struct{}{}

	case sdl.RELEASED:
		// Execute on-release callback
		if event.Repeat == 0 {
			if fn, ok := k.keyBindingsRelease[event.Keysym.Sym]; ok {
				fn()
			}
		}

		// Remove key from tracker to handle instantaneous callbacks
		delete(k.pressedKeys, event.Keysym.Sym)
	}
}

// update executes callbacks for keys registered in instantaneous mode. If the key
// is pressed down at the time update is called, the callback is executed.
func (k keyTracker) update() {
	for key, fn := range k.keyBindingsInstant {
		if k.isPressed(key) {
			fn()
		}
	}
}

// textMutator allows the user to modify text.
type textMutator struct {
	buffer string
}

// newTextTracker constructs a text tracker object.
func newTextTracker() *textMutator {
	return &textMutator{}
}

// Load loads a string into the text mutator so it can be modified.
func (t *textMutator) Load(s string) {
	t.buffer = s
}

// Append appends a string to the text mutator.
func (t *textMutator) Append(s string) {
	t.buffer += s
}

// String returns the string stored in the text mutator.
func (t *textMutator) String() string {
	return t.buffer
}

// Flush returns the string stored in the text mutator, and empties it.
func (t *textMutator) Flush() string {
	s := t.buffer
	t.buffer = ""
	return s
}

// handleEvent processes key press events.
func (t *textMutator) handleEvent(event *sdl.KeyboardEvent) {
	if event.Keysym.Sym == KeyBackspace && event.State == sdl.PRESSED {
		t.backspace()
	}
}

// backspace removes the last character from the text.
func (t *textMutator) backspace() {
	if len(t.buffer) == 0 {
		return
	}
	t.buffer = t.buffer[:len(t.buffer)-1]
}

const (
	// Key modifiers. See (https://wiki.libsdl.org/SDL_Keymod)

	KeyModeNone     = sdl.KMOD_NONE     // 0 (no modifier is applicable)
	KeyModeLShift   = sdl.KMOD_LSHIFT   // the left Shift key is down
	KeyModeRShift   = sdl.KMOD_RSHIFT   // the right Shift key is down
	KeyModeLCtrl    = sdl.KMOD_LCTRL    // the left Ctrl (Control) key is down
	KeyModeRCtrl    = sdl.KMOD_RCTRL    // the right Ctrl (Control) key is down
	KeyModeLAlt     = sdl.KMOD_LALT     // the left Alt key is down
	KeyModeRAlt     = sdl.KMOD_RALT     // the right Alt key is down
	KeyModeLGui     = sdl.KMOD_LGUI     // the left GUI key (often the Windows key) is down
	KeyModeRGui     = sdl.KMOD_RGUI     // the right GUI key (often the Windows key) is down
	KeyModeNum      = sdl.KMOD_NUM      // the Num Lock key (may be located on an extended keypad) is down
	KeyModeCaps     = sdl.KMOD_CAPS     // the Caps Lock key is down
	KeyModeMode     = sdl.KMOD_MODE     // the AltGr key is down
	KeyModeCtrl     = sdl.KMOD_CTRL     // (KMOD_LCTRL|KMOD_RCTRL)
	KeyModeShift    = sdl.KMOD_SHIFT    // (KMOD_LSHIFT|KMOD_RSHIFT)
	KeyModeAlt      = sdl.KMOD_ALT      // (KMOD_LALT|KMOD_RALT)
	KeyModeGui      = sdl.KMOD_GUI      // (KMOD_LGUI|KMOD_RGUI)
	KeyModeReserved = sdl.KMOD_RESERVED // reserved for future use

	// SDL virtual key representation. See https://wiki.libsdl.org/SDL_Keycode
	// and https://wiki.libsdl.org/SDLKeycodeLookup.

	KeyUnknown    = sdl.K_UNKNOWN    // "" (no name, empty string)
	KeyReturn     = sdl.K_RETURN     // "Return" (the Enter key (main keyboard))
	KeyEscape     = sdl.K_ESCAPE     // "Escape" (the Esc key)
	KeyBackspace  = sdl.K_BACKSPACE  // "Backspace"
	KeyTab        = sdl.K_TAB        // "Tab" (the Tab key)
	KeySpace      = sdl.K_SPACE      // "Space" (the Space Bar key(s))
	KeyExclaim    = sdl.K_EXCLAIM    // "!"
	KeyQuotedBL   = sdl.K_QUOTEDBL   // """
	KeyHash       = sdl.K_HASH       // "#"
	KeyPercent    = sdl.K_PERCENT    // "%"
	KeyDollar     = sdl.K_DOLLAR     // "$"
	KeyAmpersand  = sdl.K_AMPERSAND  // "&"
	KeyQuote      = sdl.K_QUOTE      // "'"
	KeyLeftparen  = sdl.K_LEFTPAREN  // "("
	KeyRightparen = sdl.K_RIGHTPAREN // ")"
	KeyAsterisk   = sdl.K_ASTERISK   // "*"
	KeyPlus       = sdl.K_PLUS       // "+"
	KeyComma      = sdl.K_COMMA      // ","
	KeyMinus      = sdl.K_MINUS      // "-"
	KeyPeriod     = sdl.K_PERIOD     // "."
	KeySlash      = sdl.K_SLASH      // "/"
	Key0          = sdl.K_0          // "0"
	Key1          = sdl.K_1          // "1"
	Key2          = sdl.K_2          // "2"
	Key3          = sdl.K_3          // "3"
	Key4          = sdl.K_4          // "4"
	Key5          = sdl.K_5          // "5"
	Key6          = sdl.K_6          // "6"
	Key7          = sdl.K_7          // "7"
	Key8          = sdl.K_8          // "8"
	Key9          = sdl.K_9          // "9"
	KeyColon      = sdl.K_COLON      // ":"
	KeySemicolon  = sdl.K_SEMICOLON  // ";"
	KeyLess       = sdl.K_LESS       // "<"
	KeyEquals     = sdl.K_EQUALS     // "="
	KeyGreater    = sdl.K_GREATER    // ">"
	KeyQuestion   = sdl.K_QUESTION   // "?"
	KeyAt         = sdl.K_AT         // "@"

	KeyLeftBracket  = sdl.K_LEFTBRACKET  // "["
	KeyBackslash    = sdl.K_BACKSLASH    // "\"
	KeyRightBracket = sdl.K_RIGHTBRACKET // "]"
	KeyCaret        = sdl.K_CARET        // "^"
	KeyUnderscore   = sdl.K_UNDERSCORE   // "_"
	KeyBackquote    = sdl.K_BACKQUOTE    // "`"
	KeyA            = sdl.K_a            // "A"
	KeyB            = sdl.K_b            // "B"
	KeyC            = sdl.K_c            // "C"
	KeyD            = sdl.K_d            // "D"
	KeyE            = sdl.K_e            // "E"
	KeyF            = sdl.K_f            // "F"
	KeyG            = sdl.K_g            // "G"
	KeyH            = sdl.K_h            // "H"
	KeyI            = sdl.K_i            // "I"
	KeyJ            = sdl.K_j            // "J"
	KeyK            = sdl.K_k            // "K"
	KeyL            = sdl.K_l            // "L"
	KeyM            = sdl.K_m            // "M"
	KeyN            = sdl.K_n            // "N"
	KeyO            = sdl.K_o            // "O"
	KeyP            = sdl.K_p            // "P"
	KeyQ            = sdl.K_q            // "Q"
	KeyR            = sdl.K_r            // "R"
	KeyS            = sdl.K_s            // "S"
	KeyT            = sdl.K_t            // "T"
	KeyU            = sdl.K_u            // "U"
	KeyV            = sdl.K_v            // "V"
	KeyW            = sdl.K_w            // "W"
	KeyX            = sdl.K_x            // "X"
	KeyY            = sdl.K_y            // "Y"
	KeyZ            = sdl.K_z            // "Z"

	KeyCapslock = sdl.K_CAPSLOCK // "CapsLock"

	KeyF1  = sdl.K_F1  // "F1"
	KeyF2  = sdl.K_F2  // "F2"
	KeyF3  = sdl.K_F3  // "F3"
	KeyF4  = sdl.K_F4  // "F4"
	KeyF5  = sdl.K_F5  // "F5"
	KeyF6  = sdl.K_F6  // "F6"
	KeyF7  = sdl.K_F7  // "F7"
	KeyF8  = sdl.K_F8  // "F8"
	KeyF9  = sdl.K_F9  // "F9"
	KeyF10 = sdl.K_F10 // "F10"
	KeyF11 = sdl.K_F11 // "F11"
	KeyF12 = sdl.K_F12 // "F12"

	KeyPrintscreen = sdl.K_PRINTSCREEN // "PrintScreen"
	KeyScrolllock  = sdl.K_SCROLLLOCK  // "ScrollLock"
	KeyPause       = sdl.K_PAUSE       // "Pause" (the Pause / Break key)
	KeyInsert      = sdl.K_INSERT      // "Insert" (insert on PC, help on some Mac keyboards (but does send code 73, not 117))
	KeyHome        = sdl.K_HOME        // "Home"
	KeyPageup      = sdl.K_PAGEUP      // "PageUp"
	KeyDelete      = sdl.K_DELETE      // "Delete"
	KeyEnd         = sdl.K_END         // "End"
	KeyPagedown    = sdl.K_PAGEDOWN    // "PageDown"
	KeyRight       = sdl.K_RIGHT       // "Right" (the Right arrow key (navigation keypad))
	KeyLeft        = sdl.K_LEFT        // "Left" (the Left arrow key (navigation keypad))
	KeyDown        = sdl.K_DOWN        // "Down" (the Down arrow key (navigation keypad))
	KeyUp          = sdl.K_UP          // "Up" (the Up arrow key (navigation keypad))

	KeyNumlockclear = sdl.K_NUMLOCKCLEAR // "Numlock" (the Num Lock key (PC) / the Clear key (Mac))
	KeyKPDivide     = sdl.K_KP_DIVIDE    // "Keypad /" (the / key (numeric keypad))
	KeyKPMultiply   = sdl.K_KP_MULTIPLY  // "Keypad *" (the * key (numeric keypad))
	KeyKPMinus      = sdl.K_KP_MINUS     // "Keypad -" (the - key (numeric keypad))
	KeyKPPlus       = sdl.K_KP_PLUS      // "Keypad +" (the + key (numeric keypad))
	KeyKPEnter      = sdl.K_KP_ENTER     // "Keypad Enter" (the Enter key (numeric keypad))
	KeyKP1          = sdl.K_KP_1         // "Keypad 1" (the 1 key (numeric keypad))
	KeyKP2          = sdl.K_KP_2         // "Keypad 2" (the 2 key (numeric keypad))
	KeyKP3          = sdl.K_KP_3         // "Keypad 3" (the 3 key (numeric keypad))
	KeyKP4          = sdl.K_KP_4         // "Keypad 4" (the 4 key (numeric keypad))
	KeyKP5          = sdl.K_KP_5         // "Keypad 5" (the 5 key (numeric keypad))
	KeyKP6          = sdl.K_KP_6         // "Keypad 6" (the 6 key (numeric keypad))
	KeyKP7          = sdl.K_KP_7         // "Keypad 7" (the 7 key (numeric keypad))
	KeyKP8          = sdl.K_KP_8         // "Keypad 8" (the 8 key (numeric keypad))
	KeyKP9          = sdl.K_KP_9         // "Keypad 9" (the 9 key (numeric keypad))
	KeyKP0          = sdl.K_KP_0         // "Keypad 0" (the 0 key (numeric keypad))
	KeyKPPeriod     = sdl.K_KP_PERIOD    // "Keypad ." (the . key (numeric keypad))

	KeyApplication   = sdl.K_APPLICATION    // "Application" (the Application / Compose / Context Menu (Windows) key)
	KeyPower         = sdl.K_POWER          // "Power" (The USB document says this is a status flag, not a physical key - but some Mac keyboards do have a power key)
	KeyKPEquals      = sdl.K_KP_EQUALS      // "Keypad =" (the = key (numeric keypad))
	KeyF13           = sdl.K_F13            // "F13"
	KeyF14           = sdl.K_F14            // "F14"
	KeyF15           = sdl.K_F15            // "F15"
	KeyF16           = sdl.K_F16            // "F16"
	KeyF17           = sdl.K_F17            // "F17"
	KeyF18           = sdl.K_F18            // "F18"
	KeyF19           = sdl.K_F19            // "F19"
	KeyF20           = sdl.K_F20            // "F20"
	KeyF21           = sdl.K_F21            // "F21"
	KeyF22           = sdl.K_F22            // "F22"
	KeyF23           = sdl.K_F23            // "F23"
	KeyF24           = sdl.K_F24            // "F24"
	KeyExecute       = sdl.K_EXECUTE        // "Execute"
	KeyHelp          = sdl.K_HELP           // "Help"
	KeyMenu          = sdl.K_MENU           // "Menu"
	KeySelect        = sdl.K_SELECT         // "Select"
	KeyStop          = sdl.K_STOP           // "Stop"
	KeyAgain         = sdl.K_AGAIN          // "Again" (the Again key (Redo))
	KeyUndo          = sdl.K_UNDO           // "Undo"
	KeyCut           = sdl.K_CUT            // "Cut"
	KeyCopy          = sdl.K_COPY           // "Copy"
	KeyPaste         = sdl.K_PASTE          // "Paste"
	KeyFind          = sdl.K_FIND           // "Find"
	KeyMute          = sdl.K_MUTE           // "Mute"
	KeyVolumeUp      = sdl.K_VOLUMEUP       // "VolumeUp"
	KeyVolumeDown    = sdl.K_VOLUMEDOWN     // "VolumeDown"
	KeyKPComma       = sdl.K_KP_COMMA       // "Keypad ," (the Comma key (numeric keypad))
	KeyKPEqualsAS400 = sdl.K_KP_EQUALSAS400 // "Keypad = (AS400)" (the Equals AS400 key (numeric keypad))

	KeyAltErase   = sdl.K_ALTERASE   // "AltErase" (Erase-Eaze)
	KeySysReq     = sdl.K_SYSREQ     // "SysReq" (the SysReq key)
	KeyCancel     = sdl.K_CANCEL     // "Cancel"
	KeyClear      = sdl.K_CLEAR      // "Clear"
	KeyPrior      = sdl.K_PRIOR      // "Prior"
	KeyReturn2    = sdl.K_RETURN2    // "Return"
	KeySeparator  = sdl.K_SEPARATOR  // "Separator"
	KeyOut        = sdl.K_OUT        // "Out"
	KeyOper       = sdl.K_OPER       // "Oper"
	KeyClearAgain = sdl.K_CLEARAGAIN // "Clear / Again"
	KeyCrSel      = sdl.K_CRSEL      // "CrSel"
	KeyExSel      = sdl.K_EXSEL      // "ExSel"

	KeyKP00               = sdl.K_KP_00              // "Keypad 00" (the 00 key (numeric keypad))
	KeyKP000              = sdl.K_KP_000             // "Keypad 000" (the 000 key (numeric keypad))
	KeyThousandsSeparator = sdl.K_THOUSANDSSEPARATOR // "ThousandsSeparator" (the Thousands Separator key)
	KeyDecimalSeparator   = sdl.K_DECIMALSEPARATOR   // "DecimalSeparator" (the Decimal Separator key)
	KeyCurrencyUnit       = sdl.K_CURRENCYUNIT       // "CurrencyUnit" (the Currency Unit key)
	KeyCurrencySubUnit    = sdl.K_CURRENCYSUBUNIT    // "CurrencySubUnit" (the Currency Subunit key)
	KeyKPLeftParen        = sdl.K_KP_LEFTPAREN       // "Keypad (" (the Left Parenthesis key (numeric keypad))
	KeyKPRightParen       = sdl.K_KP_RIGHTPAREN      // "Keypad )" (the Right Parenthesis key (numeric keypad))
	KeyKPLeftBrace        = sdl.K_KP_LEFTBRACE       // "Keypad {" (the Left Brace key (numeric keypad))
	KeyKPRightBrace       = sdl.K_KP_RIGHTBRACE      // "Keypad }" (the Right Brace key (numeric keypad))
	KeyKPTab              = sdl.K_KP_TAB             // "Keypad Tab" (the Tab key (numeric keypad))
	KeyKPBackspace        = sdl.K_KP_BACKSPACE       // "Keypad Backspace" (the Backspace key (numeric keypad))
	KeyKPA                = sdl.K_KP_A               // "Keypad A" (the A key (numeric keypad))
	KeyKPB                = sdl.K_KP_B               // "Keypad B" (the B key (numeric keypad))
	KeyKPC                = sdl.K_KP_C               // "Keypad C" (the C key (numeric keypad))
	KeyKPD                = sdl.K_KP_D               // "Keypad D" (the D key (numeric keypad))
	KeyKPE                = sdl.K_KP_E               // "Keypad E" (the E key (numeric keypad))
	KeyKPF                = sdl.K_KP_F               // "Keypad F" (the F key (numeric keypad))
	KeyKPXor              = sdl.K_KP_XOR             // "Keypad XOR" (the XOR key (numeric keypad))
	KeyKPPower            = sdl.K_KP_POWER           // "Keypad ^" (the Power key (numeric keypad))
	KeyKPPercent          = sdl.K_KP_PERCENT         // "Keypad %" (the Percent key (numeric keypad))
	KeyKPLess             = sdl.K_KP_LESS            // "Keypad <" (the Less key (numeric keypad))
	KeyKPGreater          = sdl.K_KP_GREATER         // "Keypad >" (the Greater key (numeric keypad))
	KeyKPAmpersand        = sdl.K_KP_AMPERSAND       // "Keypad &" (the & key (numeric keypad))
	KeyKPDblAmpersand     = sdl.K_KP_DBLAMPERSAND    // "Keypad &&" (the && key (numeric keypad))
	KeyKPVerticalBar      = sdl.K_KP_VERTICALBAR     // "Keypad |" (the | key (numeric keypad))
	KeyKPDblVerticalBar   = sdl.K_KP_DBLVERTICALBAR  // "Keypad ||" (the || key (numeric keypad))
	KeyKPColon            = sdl.K_KP_COLON           // "Keypad :" (the : key (numeric keypad))
	KeyKPHash             = sdl.K_KP_HASH            // "Keypad #" (the # key (numeric keypad))
	KeyKPSpace            = sdl.K_KP_SPACE           // "Keypad Space" (the Space key (numeric keypad))
	KeyKPAt               = sdl.K_KP_AT              // "Keypad @" (the @ key (numeric keypad))
	KeyKPExclam           = sdl.K_KP_EXCLAM          // "Keypad !" (the ! key (numeric keypad))
	KeyKPMemStore         = sdl.K_KP_MEMSTORE        // "Keypad MemStore" (the Mem Store key (numeric keypad))
	KeyKPMemRecall        = sdl.K_KP_MEMRECALL       // "Keypad MemRecall" (the Mem Recall key (numeric keypad))
	KeyKPMemClear         = sdl.K_KP_MEMCLEAR        // "Keypad MemClear" (the Mem Clear key (numeric keypad))
	KeyKPMemAdd           = sdl.K_KP_MEMADD          // "Keypad MemAdd" (the Mem Add key (numeric keypad))
	KeyKPMemSubtract      = sdl.K_KP_MEMSUBTRACT     // "Keypad MemSubtract" (the Mem Subtract key (numeric keypad))
	KeyKPMemMultiply      = sdl.K_KP_MEMMULTIPLY     // "Keypad MemMultiply" (the Mem Multiply key (numeric keypad))
	KeyKPMemDivide        = sdl.K_KP_MEMDIVIDE       // "Keypad MemDivide" (the Mem Divide key (numeric keypad))
	KeyKPPlusMinus        = sdl.K_KP_PLUSMINUS       // "Keypad +/-" (the +/- key (numeric keypad))
	KeyKPClear            = sdl.K_KP_CLEAR           // "Keypad Clear" (the Clear key (numeric keypad))
	KeyKPClearEntry       = sdl.K_KP_CLEARENTRY      // "Keypad ClearEntry" (the Clear Entry key (numeric keypad))
	KeyKPBinary           = sdl.K_KP_BINARY          // "Keypad Binary" (the Binary key (numeric keypad))
	KeyKPOctal            = sdl.K_KP_OCTAL           // "Keypad Octal" (the Octal key (numeric keypad))
	KeyKPDecimal          = sdl.K_KP_DECIMAL         // "Keypad Decimal" (the Decimal key (numeric keypad))
	KeyKPHexadecimal      = sdl.K_KP_HEXADECIMAL     // "Keypad Hexadecimal" (the Hexadecimal key (numeric keypad))

	KeyLCtrl  = sdl.K_LCTRL  // "Left Ctrl"
	KeyLShift = sdl.K_LSHIFT // "Left Shift"
	KeyLAlt   = sdl.K_LALT   // "Left Alt" (alt, option)
	KeyLGui   = sdl.K_LGUI   // "Left GUI" (windows, command (apple), meta)
	KeyRCtrl  = sdl.K_RCTRL  // "Right Ctrl"
	KeyRShift = sdl.K_RSHIFT // "Right Shift"
	KeyRAlt   = sdl.K_RALT   // "Right Alt" (alt, option)
	KeyRGui   = sdl.K_RGUI   // "Right GUI" (windows, command (apple), meta)

	KeyMode = sdl.K_MODE // "ModeSwitch"

	KeyAudioNext   = sdl.K_AUDIONEXT    // "AudioNext" (the Next Track media key)
	KeyAudioPrev   = sdl.K_AUDIOPREV    // "AudioPrev" (the Previous Track media key)
	KeyAudioStop   = sdl.K_AUDIOSTOP    // "AudioStop" (the Stop media key)
	KeyAudioPlay   = sdl.K_AUDIOPLAY    // "AudioPlay" (the Play media key)
	KeyAudioMute   = sdl.K_AUDIOMUTE    // "AudioMute" (the Mute volume key)
	KeyMediaSelect = sdl.K_MEDIASELECT  // "MediaSelect" (the Media Select key)
	KeyWWW         = sdl.K_WWW          // "WWW" (the WWW/World Wide Web key)
	KeyMail        = sdl.K_MAIL         // "Mail" (the Mail/eMail key)
	KeyCalculator  = sdl.K_CALCULATOR   // "Calculator" (the Calculator key)
	KeyComputer    = sdl.K_COMPUTER     // "Computer" (the My Computer key)
	KeyACSearch    = sdl.K_AC_SEARCH    // "AC Search" (the Search key (application control keypad))
	KeyACHome      = sdl.K_AC_HOME      // "AC Home" (the Home key (application control keypad))
	KeyACBack      = sdl.K_AC_BACK      // "AC Back" (the Back key (application control keypad))
	KeyACForward   = sdl.K_AC_FORWARD   // "AC Forward" (the Forward key (application control keypad))
	KeyACStop      = sdl.K_AC_STOP      // "AC Stop" (the Stop key (application control keypad))
	KeyACRefresh   = sdl.K_AC_REFRESH   // "AC Refresh" (the Refresh key (application control keypad))
	KeyACBookmarks = sdl.K_AC_BOOKMARKS // "AC Bookmarks" (the Bookmarks key (application control keypad))

	KeyBrightnessDown = sdl.K_BRIGHTNESSDOWN // "BrightnessDown" (the Brightness Down key)
	KeyBrightnessUp   = sdl.K_BRIGHTNESSUP   // "BrightnessUp" (the Brightness Up key)
	KeyDisplaySwitch  = sdl.K_DISPLAYSWITCH  // "DisplaySwitch" (display mirroring/dual display switch, video mode switch)
	KeyKbdIllumtoggle = sdl.K_KBDILLUMTOGGLE // "KBDIllumToggle" (the Keyboard Illumination Toggle key)
	KeyKbdIllumdown   = sdl.K_KBDILLUMDOWN   // "KBDIllumDown" (the Keyboard Illumination Down key)
	KeyKbdIllumup     = sdl.K_KBDILLUMUP     // "KBDIllumUp" (the Keyboard Illumination Up key)
	KeyEject          = sdl.K_EJECT          // "Eject" (the Eject key)
	KeySleep          = sdl.K_SLEEP          // "Sleep" (the Sleep key)
)
