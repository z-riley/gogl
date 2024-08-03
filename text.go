package turdgl

import (
	"image"
	"image/color"
	"math"
	"os"

	"github.com/golang/freetype"
)

// Text is a customisable block of text.
type Text struct {
	body               string
	pos                Vec
	colour             color.Color
	fontPath           string      // .ttf file for generating mask
	dpi, size, spacing float64     // settings for generating mask
	width, height      int         // dimensions of the generated mask
	mask               *image.RGBA // pixel image to be drawn
}

// NewText constructs a new text object with default parameters.
func NewText(body string, pos Vec) *Text {
	t := Text{
		body:     body,
		pos:      pos,
		colour:   color.RGBA{0xff, 0, 0, 0xff},
		fontPath: "../../fonts/luxisr.ttf",
		dpi:      72,
		size:     20,
		spacing:  1.5,
		width:    1024,
		height:   768,
	}
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return &t
}

// Draw draws the text onto the provided frame buffer.
func (t *Text) Draw(buf *FrameBuffer) {
	for i := 0; i < t.mask.Bounds().Max.Y; i++ {
		for j := 0; j < t.mask.Bounds().Max.X; j++ {
			rgba := t.mask.RGBAAt(j, i)
			if rgba.A > 0 {
				buf.SetPixel(i, j, NewPixel(rgba))
			}
		}
	}
}

// Text returns the current text content.
func (t *Text) Text() string {
	return t.body
}

// SetText sets the text content.
func (t *Text) SetText(txt string) *Text {
	t.body = txt
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// Pos returns the current position.
func (t *Text) Pos() Vec { return t.pos }

// SetPos sets the text's position.
func (t *Text) SetPos(pos Vec) *Text {
	t.pos = pos
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// Colour returns the current colour.
func (t *Text) Colour() color.Color { return t.colour }

// SetColour sets the text's colour.
func (t *Text) SetColour(c color.Color) *Text {
	t.colour = c
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// FontPath returns the path of the .tff file that was used to generate the text.
func (t *Text) FontPath() string { return t.fontPath }

// SetFontPath sets the path fo the .tff file that is used to generate the text.
func (t *Text) SetFontPath(path string) *Text {
	t.fontPath = path
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// DPI returns the current DPI of the font.
func (t *Text) DPI() float64 { return t.dpi }

// SetDPI sets the DPI of the font.
func (t *Text) SetDPI(dpi float64) *Text {
	t.dpi = dpi
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// Size returns the current font size.
func (t *Text) Size() float64 { return t.size }

// SetSize sets the size of the font.
func (t *Text) SetSize(size float64) *Text {
	t.size = size
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// Spacing returns the line spacing.
func (t *Text) Spacing() float64 { return t.spacing }

// SetSpacing sets the line spacing.
func (t *Text) SetSpacing(spacing float64) *Text {
	t.spacing = spacing
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// MaskSize returns the dimensions of the generated mask.
func (t *Text) MaskSize() (int, int) { return t.width, t.height }

// SetMaskSize sets the size of the mask used to generate the text on.
func (t *Text) SetMaskSize(w, h int) *Text {
	t.width, t.height = w, h
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// generateMask regenerates the mask used to generate the font pixel grid.
// It should be called any time the text settings change.
func (t *Text) generateMask() error {
	// Load font into memory.
	// Note: reading the file each time could be avoided if the *truetype.Font
	// is stored in cached instead.
	fontBytes, err := os.ReadFile(t.fontPath)
	if err != nil {
		return err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	// Configure settings
	img := image.NewRGBA(image.Rect(0, 0, t.width, t.height))
	ctx := freetype.NewContext()
	ctx.SetDPI(t.dpi)
	ctx.SetFont(font)
	ctx.SetFontSize(t.size)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(image.NewUniform(t.colour))

	// Draw the text
	x, y := int(math.Round(t.pos.X)), int(math.Round(t.pos.Y))
	pt := freetype.Pt(x, y+int(ctx.PointToFixed(t.size)>>6))
	_, err = ctx.DrawString(t.body, pt)
	if err != nil {
		panic(err)
	}
	pt.Y += ctx.PointToFixed(t.size * t.spacing)

	t.mask = img
	return nil
}
