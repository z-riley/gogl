package turdgl

import (
	"image"
	"image/color"
	"math"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// Alignment is used to specify alignment relative to the position coordinate.
type Alignment int

const (
	AlignTopLeft Alignment = iota
	AlignTopCentre
	AlignTopRight
	AlignCentreLeft
	AlignCentre
	AlignCentreRight
	AlignBottomLeft
	AlignBottomCentre
	AlignBottomRight
)

// Text is a customisable block of text.
type Text struct {
	body               string
	pos                Vec
	alignment          Alignment
	colour             color.Color
	font               *truetype.Font
	dpi, size, spacing float64     // settings for generating mask
	width, height      int         // dimensions of the generated mask
	mask               *image.RGBA // pixel image to be drawn
}

// NewText constructs a new text object with default parameters.
func NewText(body string, pos Vec) *Text {
	t := Text{
		body:      body,
		pos:       pos,
		alignment: AlignTopLeft,
		colour:    color.RGBA{0xff, 0, 0, 0xff},
		dpi:       80,
		size:      20,
		spacing:   1.5,
		width:     1200, // FIXME: this information should come from the window
		height:    768,  // FIXME: this information should come from the window
	}
	var err error
	t.font, err = loadFont("../../fonts/arial.ttf")
	if err != nil {
		panic(err)
	}
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return &t
}

// Draw draws the text onto the provided frame buffer.
func (t *Text) Draw(buf *FrameBuffer) {
	// Calculate position offset depending on alignment config
	textOffset := func() Vec {
		bbox := t.textBoundry()
		switch t.alignment {
		case AlignTopLeft:
			return Vec{0, 0}
		case AlignTopCentre:
			return Vec{-bbox.w / 2, 0}
		case AlignTopRight:
			return Vec{-bbox.w, 0}
		case AlignCentreLeft:
			return Vec{0, -bbox.h / 2}
		case AlignCentre:
			return Vec{-bbox.w / 2, -bbox.h / 2}
		case AlignCentreRight:
			return Vec{-bbox.w, 0 - bbox.h/2}
		case AlignBottomLeft:
			return Vec{0, -bbox.h}
		case AlignBottomCentre:
			return Vec{-bbox.w / 2, -bbox.h}
		case AlignBottomRight:
			return Vec{-bbox.w, -bbox.h}
		default:
			panic("Unsupported text alignment")
		}
	}()

	// Draw pixels to frame buffer
	// FIXME: Coloured artifacts and edges appear when drawing text onto a coloured background.
	// green background -> pink artifacts
	// blue background -> yellow artifacts
	// red background -> blue artifacts
	for i := 0; i < t.mask.Rect.Dy(); i++ {
		for j := 0; j < t.mask.Rect.Dx(); j++ {
			rgba := t.mask.RGBAAt(j, i)
			if rgba.A > 0 {
				x := j + int(math.Round(textOffset.X))
				y := i + int(math.Round(textOffset.Y))
				buf.SetPixel(y, x, NewPixel(rgba))
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

// Text returns the current alignment.
func (t *Text) Alignment() Alignment {
	return t.alignment
}

// SetText sets the text alignment.
func (t *Text) SetAlignment(align Alignment) *Text {
	t.alignment = align
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

// SetPath sets the path fo the .tff file that is used to generate the text.
func (t *Text) SetFont(path string) error {
	// Load font into memory
	var err error
	t.font, err = loadFont(path)
	if err != nil {
		return err
	}
	// Regenerate text
	if err := t.generateMask(); err != nil {
		return err
	}
	return nil
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
	// Configure settings
	img := image.NewRGBA(image.Rect(0, 0, t.width, t.height))
	ctx := freetype.NewContext()
	ctx.SetDPI(t.dpi)
	ctx.SetFont(t.font)
	ctx.SetFontSize(t.size)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(image.NewUniform(t.colour))

	// Draw the text
	x, y := int(math.Round(t.pos.X)), int(math.Round(t.pos.Y))
	pt := freetype.Pt(x, y+int(ctx.PointToFixed(t.size)>>6))
	_, err := ctx.DrawString(t.body, pt)
	if err != nil {
		return err
	}
	pt.Y += ctx.PointToFixed(t.size * t.spacing)

	t.mask = img
	return nil
}

// textBoundry returns a bounding box precisely surrounding the rendered text.
func (t *Text) textBoundry() *Rect {
	minX := float64(t.mask.Rect.Dx())
	maxY := 0.0
	minY := float64(t.mask.Rect.Dy())
	maxX := 0.0

	// Iterate over text mask bytes; save min and max X and Y values
	for y := 0; y < t.mask.Rect.Dy(); y++ {
		rowLen := t.mask.Rect.Dx()
		rowStart := 4 * (y * rowLen)
		rowEnd := 4 * (((y + 1) * rowLen) - 1)
		row := t.mask.Pix[rowStart:rowEnd]
		newText := false
		for x := 3; x < len(row); x += 4 {
			textFound := row[x] != 0x00
			if textFound {
				if !newText {
					newText = true
					minY = math.Min(minY, float64(y))
					minX = math.Min(minX, float64(x/4))
				}
				maxY = float64(y)
				maxX = math.Max(maxX, float64(x/4))
			}
		}
	}

	return NewRect(
		maxX-minX,
		maxY-minY,
		Vec{minX, minY},
		WithStyle(Style{Colour: color.White, Thickness: 1}),
	)
}

// loadFont loads a Truetype font from a .tff file.
func loadFont(path string) (*truetype.Font, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return font, nil
}
