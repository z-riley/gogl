package turdgl

import "image/color"

// RGBA8 returns a colour in uint8 format.
func RGBA8(c color.Color) (r, g, b, a uint8) {
	r32, g32, b32, a32 := c.RGBA()
	r = uint8(r32 >> 8)
	g = uint8(g32 >> 8)
	b = uint8(b32 >> 8)
	a = uint8(a32 >> 8)
	return
}

// Common colour set.
var (
	Black   = color.RGBA{0, 0, 0, 255}
	White   = color.RGBA{255, 255, 255, 255}
	Red     = color.RGBA{255, 0, 0, 255}
	Lime    = color.RGBA{0, 255, 0, 255}
	Blue    = color.RGBA{0, 0, 255, 255}
	Yellow  = color.RGBA{255, 255, 0, 255}
	Cyan    = color.RGBA{0, 255, 255, 255}
	Magenta = color.RGBA{255, 0, 255, 255}
	Silver  = color.RGBA{192, 192, 192, 255}
	Grey    = color.RGBA{128, 128, 128, 255}
	Maroon  = color.RGBA{128, 0, 0, 255}
	Olive   = color.RGBA{128, 128, 0, 255}
	Green   = color.RGBA{0, 128, 0, 255}
	Purple  = color.RGBA{128, 0, 128, 255}
	Teal    = color.RGBA{0, 128, 128, 255}
	Navy    = color.RGBA{0, 0, 128, 255}
)

// Extra colours for if you're feeling artistic.
var (
	DarkRed              = color.RGBA{139, 0, 0, 255}
	Brown                = color.RGBA{165, 42, 42, 255}
	Firebrick            = color.RGBA{178, 34, 34, 255}
	Crimson              = color.RGBA{220, 20, 60, 255}
	Tomato               = color.RGBA{255, 99, 71, 255}
	Coral                = color.RGBA{255, 127, 80, 255}
	IndianRed            = color.RGBA{205, 92, 92, 255}
	LightCoral           = color.RGBA{240, 128, 128, 255}
	DarkSalmon           = color.RGBA{233, 150, 122, 255}
	Salmon               = color.RGBA{250, 128, 114, 255}
	LightSalmon          = color.RGBA{255, 160, 122, 255}
	OrangeRed            = color.RGBA{255, 69, 0, 255}
	DarkOrange           = color.RGBA{255, 140, 0, 255}
	Orange               = color.RGBA{255, 165, 0, 255}
	Gold                 = color.RGBA{255, 215, 0, 255}
	DarkGoldenRod        = color.RGBA{184, 134, 11, 255}
	GoldenRod            = color.RGBA{218, 165, 32, 255}
	PaleGoldenRod        = color.RGBA{238, 232, 170, 255}
	DarkKhaki            = color.RGBA{189, 183, 107, 255}
	Khaki                = color.RGBA{240, 230, 140, 255}
	YellowGreen          = color.RGBA{154, 205, 50, 255}
	DarkOliveGreen       = color.RGBA{85, 107, 47, 255}
	OliveDrab            = color.RGBA{107, 142, 35, 255}
	LawnGreen            = color.RGBA{124, 252, 0, 255}
	Chartreuse           = color.RGBA{127, 255, 0, 255}
	GreenYellow          = color.RGBA{173, 255, 47, 255}
	DarkGreen            = color.RGBA{0, 100, 0, 255}
	ForestGreen          = color.RGBA{34, 139, 34, 255}
	LimeGreen            = color.RGBA{50, 205, 50, 255}
	LightGreen           = color.RGBA{144, 238, 144, 255}
	PaleGreen            = color.RGBA{152, 251, 152, 255}
	DarkSeaGreen         = color.RGBA{143, 188, 143, 255}
	MediumSpringGreen    = color.RGBA{0, 250, 154, 255}
	SpringGreen          = color.RGBA{0, 255, 127, 255}
	SeaGreen             = color.RGBA{46, 139, 87, 255}
	MediumAquaMarine     = color.RGBA{102, 205, 170, 255}
	MediumSeaGreen       = color.RGBA{60, 179, 113, 255}
	LightSeaGreen        = color.RGBA{32, 178, 170, 255}
	DarkSlateGrey        = color.RGBA{47, 79, 79, 255}
	DarkCyan             = color.RGBA{0, 139, 139, 255}
	Aqua                 = color.RGBA{0, 255, 255, 255}
	LightCyan            = color.RGBA{224, 255, 255, 255}
	DarkTurquoise        = color.RGBA{0, 206, 209, 255}
	Turquoise            = color.RGBA{64, 224, 208, 255}
	MediumTurquoise      = color.RGBA{72, 209, 204, 255}
	PaleTurquoise        = color.RGBA{175, 238, 238, 255}
	AquaMarine           = color.RGBA{127, 255, 212, 255}
	PowderBlue           = color.RGBA{176, 224, 230, 255}
	CadetBlue            = color.RGBA{95, 158, 160, 255}
	SteelBlue            = color.RGBA{70, 130, 180, 255}
	CornBlowerBlue       = color.RGBA{100, 149, 237, 255}
	DeepSkyBlue          = color.RGBA{0, 191, 255, 255}
	DodgerBlue           = color.RGBA{30, 144, 255, 255}
	LightBlue            = color.RGBA{173, 216, 230, 255}
	SkyBlue              = color.RGBA{135, 206, 235, 255}
	LighSkyBlue          = color.RGBA{135, 206, 250, 255}
	MidnightBlue         = color.RGBA{25, 25, 112, 255}
	DarkBlue             = color.RGBA{0, 0, 139, 255}
	MediumBlue           = color.RGBA{0, 0, 205, 255}
	RoyalBlue            = color.RGBA{65, 105, 225, 255}
	BlueViolet           = color.RGBA{138, 43, 226, 255}
	Indigo               = color.RGBA{75, 0, 130, 255}
	DarkSlateBlue        = color.RGBA{72, 61, 139, 255}
	SlateBlue            = color.RGBA{106, 90, 205, 255}
	MediumSlateBlue      = color.RGBA{123, 104, 238, 255}
	MediumPurple         = color.RGBA{147, 112, 219, 255}
	DarkMagenta          = color.RGBA{139, 0, 139, 255}
	DarkViolet           = color.RGBA{148, 0, 211, 255}
	DarkOrchid           = color.RGBA{153, 50, 204, 255}
	MediumOrchid         = color.RGBA{186, 85, 211, 255}
	Thistle              = color.RGBA{216, 191, 216, 255}
	Plum                 = color.RGBA{221, 160, 221, 255}
	Violet               = color.RGBA{238, 130, 238, 255}
	Orchid               = color.RGBA{218, 112, 214, 255}
	MediumVioletRed      = color.RGBA{199, 21, 133, 255}
	PaleVioletRed        = color.RGBA{219, 112, 147, 255}
	DeepPink             = color.RGBA{255, 20, 147, 255}
	HotPink              = color.RGBA{255, 105, 180, 255}
	LightPink            = color.RGBA{255, 182, 193, 255}
	Pink                 = color.RGBA{255, 192, 203, 255}
	AntiqueWhite         = color.RGBA{250, 235, 215, 255}
	Beige                = color.RGBA{245, 245, 220, 255}
	Bisque               = color.RGBA{255, 228, 196, 255}
	BlanchedAlmond       = color.RGBA{255, 235, 205, 255}
	Wheat                = color.RGBA{245, 222, 179, 255}
	CornSilk             = color.RGBA{255, 248, 220, 255}
	LemonChiffon         = color.RGBA{255, 250, 205, 255}
	LightGoldenRodYellow = color.RGBA{250, 250, 210, 255}
	LightYellow          = color.RGBA{255, 255, 224, 255}
	SaddleBrown          = color.RGBA{139, 69, 19, 255}
	Sienna               = color.RGBA{160, 82, 45, 255}
	Chocolate            = color.RGBA{210, 105, 30, 255}
	Peru                 = color.RGBA{205, 133, 63, 255}
	SandyBrown           = color.RGBA{244, 164, 96, 255}
	BurlyWood            = color.RGBA{222, 184, 135, 255}
	Tan                  = color.RGBA{210, 180, 140, 255}
	RosyBrown            = color.RGBA{188, 143, 143, 255}
	Moccasin             = color.RGBA{255, 228, 181, 255}
	NavajoWhite          = color.RGBA{255, 222, 173, 255}
	PeachPuff            = color.RGBA{255, 218, 185, 255}
	MistyRose            = color.RGBA{255, 228, 225, 255}
	LavenderBlush        = color.RGBA{255, 240, 245, 255}
	Linen                = color.RGBA{250, 240, 230, 255}
	OldLace              = color.RGBA{253, 245, 230, 255}
	PapayaWhip           = color.RGBA{255, 239, 213, 255}
	SeaShell             = color.RGBA{255, 245, 238, 255}
	MintCream            = color.RGBA{245, 255, 250, 255}
	SlateGrey            = color.RGBA{112, 128, 144, 255}
	LightSlateGrey       = color.RGBA{119, 136, 153, 255}
	LightSteelBlue       = color.RGBA{176, 196, 222, 255}
	Lavender             = color.RGBA{230, 230, 250, 255}
	FloralWhite          = color.RGBA{255, 250, 240, 255}
	AliceBlue            = color.RGBA{240, 248, 255, 255}
	GhostWhite           = color.RGBA{248, 248, 255, 255}
	Honeydew             = color.RGBA{240, 255, 240, 255}
	Ivory                = color.RGBA{255, 255, 240, 255}
	Azure                = color.RGBA{240, 255, 255, 255}
	Snow                 = color.RGBA{255, 250, 250, 255}
	DimGrey              = color.RGBA{105, 105, 105, 255}
	DarkGrey             = color.RGBA{169, 169, 169, 255}
	LightGrey            = color.RGBA{211, 211, 211, 255}
	Gainsboro            = color.RGBA{220, 220, 220, 255}
	WhiteSmoke           = color.RGBA{245, 245, 245, 255}
)
