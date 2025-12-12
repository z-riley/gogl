package gogl

import (
	"github.com/jupiterrider/purego-sdl3/sdl"
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
func (k *keyTracker) handleEvent(event sdl.KeyboardEvent) {
	if event.Down {
		// Execute on-press callback
		if event.Repeat {
			if fn, ok := k.keyBindingsPress[event.Key]; ok {
				fn()
			}
		}

		// Add key to tracker to handle instantaneous callbacks
		k.pressedKeys[event.Key] = struct{}{}
	} else {
		// Execute on-release callback
		if event.Repeat {
			if fn, ok := k.keyBindingsRelease[event.Key]; ok {
				fn()
			}
		}

		// Remove key from tracker to handle instantaneous callbacks
		delete(k.pressedKeys, event.Key)
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

const (
	// Key modifiers. See (https://wiki.libsdl.org/SDL_Keymod)

	KeyModeNone   = sdl.KeymodNone   // 0 (no modifier is applicable)
	KeyModeLShift = sdl.KeymodLShift // the left Shift key is down
	KeyModeRShift = sdl.KeymodRShift // the right Shift key is down
	KeyModeLCtrl  = sdl.KeymodLCtrl  // the left Ctrl (Control) key is down
	KeyModeRCtrl  = sdl.KeymodRCtrl  // the right Ctrl (Control) key is down
	KeyModeLAlt   = sdl.KeymodLAlt   // the left Alt key is down
	KeyModeRAlt   = sdl.KeymodRAlt   // the right Alt key is down
	KeyModeLGui   = sdl.KeymodLGui   // the left GUI key (often the Windows key) is down
	KeyModeRGui   = sdl.KeymodRGui   // the right GUI key (often the Windows key) is down
	KeyModeNum    = sdl.KeymodNum    // the Num Lock key (may be located on an extended keypad) is down
	KeyModeCaps   = sdl.KeymodCaps   // the Caps Lock key is down
	KeyModeMode   = sdl.KeymodMode   // the AltGr key is down
	KeyModeCtrl   = sdl.KeymodCtrl   // (KMOD_LCTRL|KMOD_RCTRL)
	KeyModeShift  = sdl.KeymodShift  // (KMOD_LSHIFT|KMOD_RSHIFT)
	KeyModeAlt    = sdl.KeymodAlt    // (KMOD_LALT|KMOD_RALT)
	KeyModeGui    = sdl.KeymodGui    // (KMOD_LGUI|KMOD_RGUI)

	// SDL virtual key representation. See https://wiki.libsdl.org/SDL_Keycode
	// and https://wiki.libsdl.org/SDLKeycodeLookup.

	KeyUnknown           = sdl.KeycodeUnknown       // "" (no name, empty string)
	KeyReturn            = sdl.KeycodeReturn        // "Return" (the Enter key (main keyboard))
	KeyEscape            = sdl.KeycodeEscape        // "Escape" (the Esc key)
	KeyBackspace         = sdl.KeycodeBackspace     // "Backspace"
	KeyTab               = sdl.KeycodeTab           // "Tab" (the Tab key)
	KeySpace             = sdl.KeycodeSpace         // "Space" (the Space Bar key(s))
	KeyExclaim           = sdl.KeycodeExclaim       // "!"
	KeycodeDblApostrophe = sdl.KeycodeDblApostrophe // """
	KeyHash              = sdl.KeycodeHash          // "#"
	KeyPercent           = sdl.KeycodePercent       // "%"
	KeyDollar            = sdl.KeycodeDollar        // "$"
	KeyAmpersand         = sdl.KeycodeAmpersand     // "&"
	KeycodeApostrophe    = sdl.KeycodeApostrophe    // "'"
	KeyLeftparen         = sdl.KeycodeLeftParen     // "("
	KeyRightparen        = sdl.KeycodeRightParen    // ")"
	KeyAsterisk          = sdl.KeycodeAsterisk      // "*"
	KeyPlus              = sdl.KeycodePlus          // "+"
	KeyComma             = sdl.KeycodeComma         // ","
	KeyMinus             = sdl.KeycodeMinus         // "-"
	KeyPeriod            = sdl.KeycodePeriod        // "."
	KeySlash             = sdl.KeycodeSlash         // "/"
	Key0                 = sdl.Keycode0             // "0"
	Key1                 = sdl.Keycode1             // "1"
	Key2                 = sdl.Keycode2             // "2"
	Key3                 = sdl.Keycode3             // "3"
	Key4                 = sdl.Keycode4             // "4"
	Key5                 = sdl.Keycode5             // "5"
	Key6                 = sdl.Keycode6             // "6"
	Key7                 = sdl.Keycode7             // "7"
	Key8                 = sdl.Keycode8             // "8"
	Key9                 = sdl.Keycode9             // "9"
	KeyColon             = sdl.KeycodeColon         // ":"
	KeySemicolon         = sdl.KeycodeSemicolon     // ";"
	KeyLess              = sdl.KeycodeLess          // "<"
	KeyEquals            = sdl.KeycodeEquals        // "="
	KeyGreater           = sdl.KeycodeGreater       // ">"
	KeyQuestion          = sdl.KeycodeQuestion      // "?"
	KeyAt                = sdl.KeycodeAt            // "@"

	KeyLeftBracket  = sdl.KeycodeLeftBracket  // "["
	KeyBackslash    = sdl.KeycodeBackslash    // "\"
	KeyRightBracket = sdl.KeycodeRightBracket // "]"
	KeyCaret        = sdl.KeycodeCaret        // "^"
	KeyUnderscore   = sdl.KeycodeUnderscore   // "_"

	KeyGrave = sdl.KeycodeGrave // "`"
	KeyA     = sdl.KeycodeA     // "A"
	KeyB     = sdl.KeycodeB     // "B"
	KeyC     = sdl.KeycodeC     // "C"
	KeyD     = sdl.KeycodeD     // "D"
	KeyE     = sdl.KeycodeE     // "E"
	KeyF     = sdl.KeycodeF     // "F"
	KeyG     = sdl.KeycodeG     // "G"
	KeyH     = sdl.KeycodeH     // "H"
	KeyI     = sdl.KeycodeI     // "I"
	KeyJ     = sdl.KeycodeJ     // "J"
	KeyK     = sdl.KeycodeK     // "K"
	KeyL     = sdl.KeycodeL     // "L"
	KeyM     = sdl.KeycodeM     // "M"
	KeyN     = sdl.KeycodeN     // "N"
	KeyO     = sdl.KeycodeO     // "O"
	KeyP     = sdl.KeycodeP     // "P"
	KeyQ     = sdl.KeycodeQ     // "Q"
	KeyR     = sdl.KeycodeR     // "R"
	KeyS     = sdl.KeycodeS     // "S"
	KeyT     = sdl.KeycodeT     // "T"
	KeyU     = sdl.KeycodeU     // "U"
	KeyV     = sdl.KeycodeV     // "V"
	KeyW     = sdl.KeycodeW     // "W"
	KeyX     = sdl.KeycodeX     // "X"
	KeyY     = sdl.KeycodeY     // "Y"
	KeyZ     = sdl.KeycodeZ     // "Z"

	KeyCapslock = sdl.KeycodeCapsLock // "CapsLock"

	KeyF1  = sdl.KeycodeF1  // "F1"
	KeyF2  = sdl.KeycodeF2  // "F2"
	KeyF3  = sdl.KeycodeF3  // "F3"
	KeyF4  = sdl.KeycodeF4  // "F4"
	KeyF5  = sdl.KeycodeF5  // "F5"
	KeyF6  = sdl.KeycodeF6  // "F6"
	KeyF7  = sdl.KeycodeF7  // "F7"
	KeyF8  = sdl.KeycodeF8  // "F8"
	KeyF9  = sdl.KeycodeF9  // "F9"
	KeyF10 = sdl.KeycodeF10 // "F10"
	KeyF11 = sdl.KeycodeF11 // "F11"
	KeyF12 = sdl.KeycodeF12 // "F12"

	KeyPrintscreen = sdl.KeycodePrintScreen // "PrintScreen"
	KeyScrolllock  = sdl.KeycodeScrollLock  // "ScrollLock"
	KeyPause       = sdl.KeycodePause       // "Pause" (the Pause / Break key)
	KeyInsert      = sdl.KeycodeInsert      // "Insert" (insert on PC, help on some Mac keyboards (but does send code 73, not 117))
	KeyHome        = sdl.KeycodeHome        // "Home"
	KeyPageup      = sdl.KeycodePageUp      // "PageUp"
	KeyDelete      = sdl.KeycodeDelete      // "Delete"
	KeyEnd         = sdl.KeycodeEnd         // "End"
	KeyPagedown    = sdl.KeycodePageDown    // "PageDown"
	KeyRight       = sdl.KeycodeRight       // "Right" (the Right arrow key (navigation keypad))
	KeyLeft        = sdl.KeycodeLeft        // "Left" (the Left arrow key (navigation keypad))
	KeyDown        = sdl.KeycodeDown        // "Down" (the Down arrow key (navigation keypad))
	KeyUp          = sdl.KeycodeUp          // "Up" (the Up arrow key (navigation keypad))

	KeyNumlockclear = sdl.KeycodeNumLockClear // "Numlock" (the Num Lock key (PC) / the Clear key (Mac))
	KeyKPDivide     = sdl.KeycodeKpDivide     // "Keypad /" (the / key (numeric keypad))
	KeyKPMultiply   = sdl.KeycodeKpMultiply   // "Keypad *" (the * key (numeric keypad))
	KeyKPMinus      = sdl.KeycodeKpMinus      // "Keypad -" (the - key (numeric keypad))
	KeyKPPlus       = sdl.KeycodeKpPlus       // "Keypad +" (the + key (numeric keypad))
	KeyKPEnter      = sdl.KeycodeKpEnter      // "Keypad Enter" (the Enter key (numeric keypad))
	KeyKP1          = sdl.KeycodeKp1          // "Keypad 1" (the 1 key (numeric keypad))
	KeyKP2          = sdl.KeycodeKp2          // "Keypad 2" (the 2 key (numeric keypad))
	KeyKP3          = sdl.KeycodeKp3          // "Keypad 3" (the 3 key (numeric keypad))
	KeyKP4          = sdl.KeycodeKp4          // "Keypad 4" (the 4 key (numeric keypad))
	KeyKP5          = sdl.KeycodeKp5          // "Keypad 5" (the 5 key (numeric keypad))
	KeyKP6          = sdl.KeycodeKp6          // "Keypad 6" (the 6 key (numeric keypad))
	KeyKP7          = sdl.KeycodeKp7          // "Keypad 7" (the 7 key (numeric keypad))
	KeyKP8          = sdl.KeycodeKp8          // "Keypad 8" (the 8 key (numeric keypad))
	KeyKP9          = sdl.KeycodeKp9          // "Keypad 9" (the 9 key (numeric keypad))
	KeyKP0          = sdl.KeycodeKp0          // "Keypad 0" (the 0 key (numeric keypad))
	KeyKPPeriod     = sdl.KeycodeKpPeriod     // "Keypad ." (the . key (numeric keypad))

	KeyApplication   = sdl.KeycodeApplication   // "Application" (the Application / Compose / Context Menu (Windows) key)
	KeyPower         = sdl.KeycodePower         // "Power" (The USB document says this is a status flag, not a physical key - but some Mac keyboards do have a power key)
	KeyKPEquals      = sdl.KeycodeKpEquals      // "Keypad =" (the = key (numeric keypad))
	KeyF13           = sdl.KeycodeF13           // "F13"
	KeyF14           = sdl.KeycodeF14           // "F14"
	KeyF15           = sdl.KeycodeF15           // "F15"
	KeyF16           = sdl.KeycodeF16           // "F16"
	KeyF17           = sdl.KeycodeF17           // "F17"
	KeyF18           = sdl.KeycodeF18           // "F18"
	KeyF19           = sdl.KeycodeF19           // "F19"
	KeyF20           = sdl.KeycodeF20           // "F20"
	KeyF21           = sdl.KeycodeF21           // "F21"
	KeyF22           = sdl.KeycodeF22           // "F22"
	KeyF23           = sdl.KeycodeF23           // "F23"
	KeyF24           = sdl.KeycodeF24           // "F24"
	KeyExecute       = sdl.KeycodeExecute       // "Execute"
	KeyHelp          = sdl.KeycodeHelp          // "Help"
	KeyMenu          = sdl.KeycodeMenu          // "Menu"
	KeySelect        = sdl.KeycodeSelect        // "Select"
	KeyStop          = sdl.KeycodeStop          // "Stop"
	KeyAgain         = sdl.KeycodeAgain         // "Again" (the Again key (Redo))
	KeyUndo          = sdl.KeycodeUndo          // "Undo"
	KeyCut           = sdl.KeycodeCut           // "Cut"
	KeyCopy          = sdl.KeycodeCopy          // "Copy"
	KeyPaste         = sdl.KeycodePaste         // "Paste"
	KeyFind          = sdl.KeycodeFind          // "Find"
	KeyMute          = sdl.KeycodeMute          // "Mute"
	KeyVolumeUp      = sdl.KeycodeVolumeUp      // "VolumeUp"
	KeyVolumeDown    = sdl.KeycodeVolumeDown    // "VolumeDown"
	KeyKPComma       = sdl.KeycodeKpComma       // "Keypad ," (the Comma key (numeric keypad))
	KeyKPEqualsAS400 = sdl.KeycodeKpEqualsAs400 // "Keypad = (AS400)" (the Equals AS400 key (numeric keypad))

	KeyAltErase   = sdl.KeycodeAltErase  // "AltErase" (Erase-Eaze)
	KeySysReq     = sdl.KeycodeSysReq    // "SysReq" (the SysReq key)
	KeyCancel     = sdl.KeycodeCancel    // "Cancel"
	KeyClear      = sdl.KeycodeClear     // "Clear"
	KeyPrior      = sdl.KeycodePrior     // "Prior"
	KeyReturn2    = sdl.KeycodeReturn    // "Return"
	KeySeparator  = sdl.KeycodeSeparator // "Separator"
	KeyOut        = sdl.KeycodeOut       // "Out"
	KeyOper       = sdl.KeycodeOper      // "Oper"
	KeyClearAgain = sdl.KeycodeClear     // "Clear / Again"
	KeyCrSel      = sdl.KeycodeCrSel     // "CrSel"
	KeyExSel      = sdl.KeycodeExSel     // "ExSel"

	KeyKP00               = sdl.KeycodeKp00               // "Keypad 00" (the 00 key (numeric keypad))
	KeyKP000              = sdl.KeycodeKp000              // "Keypad 000" (the 000 key (numeric keypad))
	KeyThousandsSeparator = sdl.KeycodeThousandsSeparator // "ThousandsSeparator" (the Thousands Separator key)
	KeyDecimalSeparator   = sdl.KeycodeDecimalSeparator   // "DecimalSeparator" (the Decimal Separator key)
	KeyCurrencyUnit       = sdl.KeycodeCurrencyUnit       // "CurrencyUnit" (the Currency Unit key)
	KeyCurrencySubunit    = sdl.KeycodeCurrencySubunit    // "CurrencySubUnit" (the Currency Subunit key)
	KeyKPLeftParen        = sdl.KeycodeKpLeftParen        // "Keypad (" (the Left Parenthesis key (numeric keypad))
	KeyKPRightParen       = sdl.KeycodeKpRightParen       // "Keypad )" (the Right Parenthesis key (numeric keypad))
	KeyKPLeftBrace        = sdl.KeycodeKpLeftBrace        // "Keypad {" (the Left Brace key (numeric keypad))
	KeyKPRightBrace       = sdl.KeycodeKpRightBrace       // "Keypad }" (the Right Brace key (numeric keypad))
	KeyKPTab              = sdl.KeycodeKpTab              // "Keypad Tab" (the Tab key (numeric keypad))
	KeyKPBackspace        = sdl.KeycodeKpBackspace        // "Keypad Backspace" (the Backspace key (numeric keypad))
	KeyKPA                = sdl.KeycodeKpA                // "Keypad A" (the A key (numeric keypad))
	KeyKPB                = sdl.KeycodeKpB                // "Keypad B" (the B key (numeric keypad))
	KeyKPC                = sdl.KeycodeKpC                // "Keypad C" (the C key (numeric keypad))
	KeyKPD                = sdl.KeycodeKpD                // "Keypad D" (the D key (numeric keypad))
	KeyKPE                = sdl.KeycodeKpE                // "Keypad E" (the E key (numeric keypad))
	KeyKPF                = sdl.KeycodeKpF                // "Keypad F" (the F key (numeric keypad))
	KeyKPXor              = sdl.KeycodeKpXor              // "Keypad XOR" (the XOR key (numeric keypad))
	KeyKPPower            = sdl.KeycodeKpPower            // "Keypad ^" (the Power key (numeric keypad))
	KeyKPPercent          = sdl.KeycodeKpPercent          // "Keypad %" (the Percent key (numeric keypad))
	KeyKPLess             = sdl.KeycodeKpLess             // "Keypad <" (the Less key (numeric keypad))
	KeyKPGreater          = sdl.KeycodeKpGreater          // "Keypad >" (the Greater key (numeric keypad))
	KeyKPAmpersand        = sdl.KeycodeKpAmpersand        // "Keypad &" (the & key (numeric keypad))
	KeyKPDblAmpersand     = sdl.KeycodeKpDblAmpersand     // "Keypad &&" (the && key (numeric keypad))
	KeyKPVerticalBar      = sdl.KeycodeKpVerticalBar      // "Keypad |" (the | key (numeric keypad))
	KeyKPDblVerticalBar   = sdl.KeycodeKpDblVerticalBar   // "Keypad ||" (the || key (numeric keypad))
	KeyKPColon            = sdl.KeycodeKpColon            // "Keypad :" (the : key (numeric keypad))
	KeyKPHash             = sdl.KeycodeKpHash             // "Keypad #" (the # key (numeric keypad))
	KeyKPSpace            = sdl.KeycodeKpSpace            // "Keypad Space" (the Space key (numeric keypad))
	KeyKPAt               = sdl.KeycodeKpAt               // "Keypad @" (the @ key (numeric keypad))
	KeyKPExclam           = sdl.KeycodeKpExclam           // "Keypad !" (the ! key (numeric keypad))
	KeyKPMemStore         = sdl.KeycodeKpMemStore         // "Keypad MemStore" (the Mem Store key (numeric keypad))
	KeyKPMemRecall        = sdl.KeycodeKpMemRecall        // "Keypad MemRecall" (the Mem Recall key (numeric keypad))
	KeyKPMemClear         = sdl.KeycodeKpMemClear         // "Keypad MemClear" (the Mem Clear key (numeric keypad))
	KeyKPMemAdd           = sdl.KeycodeKpMemAdd           // "Keypad MemAdd" (the Mem Add key (numeric keypad))
	KeyKPMemSubtract      = sdl.KeycodeKpMemSubtract      // "Keypad MemSubtract" (the Mem Subtract key (numeric keypad))
	KeyKPMemMultiply      = sdl.KeycodeKpMemMultiply      // "Keypad MemMultiply" (the Mem Multiply key (numeric keypad))
	KeyKPMemDivide        = sdl.KeycodeKpMemDivide        // "Keypad MemDivide" (the Mem Divide key (numeric keypad))
	KeyKPPlusMinus        = sdl.KeycodeKpPlusMinus        // "Keypad +/-" (the +/- key (numeric keypad))
	KeyKPClear            = sdl.KeycodeKpClear            // "Keypad Clear" (the Clear key (numeric keypad))
	KeyKPClearEntry       = sdl.KeycodeKpClearEntry       // "Keypad ClearEntry" (the Clear Entry key (numeric keypad))
	KeyKPBinary           = sdl.KeycodeKpBinary           // "Keypad Binary" (the Binary key (numeric keypad))
	KeyKPOctal            = sdl.KeycodeKpOctal            // "Keypad Octal" (the Octal key (numeric keypad))
	KeyKPDecimal          = sdl.KeycodeKpDecimal          // "Keypad Decimal" (the Decimal key (numeric keypad))
	KeyKPHexadecimal      = sdl.KeycodeKpHexadecimal      // "Keypad Hexadecimal" (the Hexadecimal key (numeric keypad))

	KeyLCtrl  = sdl.KeycodeLCtrl  // "Left Ctrl"
	KeyLShift = sdl.KeycodeLShift // "Left Shift"
	KeyLAlt   = sdl.KeycodeLAlt   // "Left Alt" (alt, option)
	KeyLGui   = sdl.KeycodeLGui   // "Left GUI" (windows, command (apple), meta)
	KeyRCtrl  = sdl.KeycodeRCtrl  // "Right Ctrl"
	KeyRShift = sdl.KeycodeRShift // "Right Shift"
	KeyRAlt   = sdl.KeycodeRAlt   // "Right Alt" (alt, option)
	KeyRGui   = sdl.KeycodeRGui   // "Right GUI" (windows, command (apple), meta)

	KeyMode = sdl.KeycodeMode // "ModeSwitch"

	KeyMediaNextTrack = sdl.KeycodeMediaNextTrack     // "AudioNext" (the Next Track media key)
	KeyMediaPrevTrack = sdl.KeycodeMediaPreviousTrack // "AudioPrev" (the Previous Track media key)
	KeyMediaStop      = sdl.KeycodeMediaStop          // "AudioStop" (the Stop media key)
	KeyMediaPlay      = sdl.KeycodeMediaPlay          // "AudioPlay" (the Play media key)
	KeyMediaSelect    = sdl.KeycodeMediaSelect        // "MediaSelect" (the Media Select key)
	KeyACSearch       = sdl.KeycodeAcSearch           // "AC Search" (the Search key (application control keypad))
	KeyACHome         = sdl.KeycodeAcHome             // "AC Home" (the Home key (application control keypad))
	KeyACBack         = sdl.KeycodeAcBack             // "AC Back" (the Back key (application control keypad))
	KeyACForward      = sdl.KeycodeAcForward          // "AC Forward" (the Forward key (application control keypad))
	KeyACStop         = sdl.KeycodeAcStop             // "AC Stop" (the Stop key (application control keypad))
	KeyACRefresh      = sdl.KeycodeAcRefresh          // "AC Refresh" (the Refresh key (application control keypad))
	KeyACBookmarks    = sdl.KeycodeAcBookmarks        // "AC Bookmarks" (the Bookmarks key (application control keypad))
	KeyMediaEject     = sdl.KeycodeMediaEject         // "Eject" (the Eject key)
	KeySleep          = sdl.KeycodeSleep              // "Sleep" (the Sleep key)
)
