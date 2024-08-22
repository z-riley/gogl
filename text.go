package turdgl

import (
	"image"
	"image/color"
	"math"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
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
	AlignCustom
)

// Text is a customisable block of text.
type Text struct {
	body               string
	pos                Vec
	alignment          Alignment
	labelOffset        Vec // label offset used when in AlignCustom mode
	colour             color.Color
	font               *sfnt.Font
	dpi, size, spacing float64     // settings for generating mask
	width, height      int         // dimensions of the generated mask
	mask               *image.RGBA // pixel image to be drawn
}

// NewText constructs a new text object with default parameters.
func NewText(body string, pos Vec, fontPath string) *Text {
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
	t.font, err = loadFont(fontPath)
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
	bbox := t.textBoundry()
	textOffset := func() Vec {
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
		case AlignCustom:
			return Vec{(-bbox.w / 2) + t.labelOffset.X, (-bbox.h / 2) + t.labelOffset.Y}
		default:
			panic("Unsupported text alignment")
		}
	}()

	// Draw pixels to frame buffer
	startX, startY := int(bbox.Pos.X), int(bbox.Pos.Y)
	endX, endY := startX+int(bbox.w), startY+int(bbox.h)
	for i := startY; i < endY; i++ {
		for j := startX; j < endX; j++ {
			rgba := t.mask.RGBAAt(j, i)
			if rgba.A > 0 {
				x := j + int(math.Round(textOffset.X))
				y := i + int(math.Round(textOffset.Y))
				buf.SetPixel(y, x, NewPixel(rgba))
			}
		}
	}
}

// Move moves the text by a given vector.
func (t *Text) Move(mov Vec) {
	t.pos = Add(t.pos, mov)
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
// Warning: use with caution becuase the alignment is consistent, but often innaccurate.
func (t *Text) SetAlignment(align Alignment) *Text {
	t.alignment = align
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// Offset returns the text's offset.
func (t *Text) Offset() Vec {
	return t.labelOffset
}

// SetOffset sets the text's offset. Note: this overwrites previous alignment settings.
func (t *Text) SetOffset(offset Vec) *Text {
	t.SetAlignment(AlignCustom)
	t.labelOffset = offset
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

// SetPath sets the path fo the .ttf file that is used to generate the text.
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
	face, err := opentype.NewFace(t.font, &opentype.FaceOptions{
		Size:    t.size,
		DPI:     t.dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}
	defer face.Close()

	img := image.NewRGBA(image.Rect(0, 0, t.width, t.height))
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(t.colour),
		Face: face,
		Dot:  fixed.P(int(t.pos.X), int(t.pos.Y)), // Set the position (x, y)
	}

	// Draw lines seperately
	y := int(t.pos.Y)
	lineHeight := face.Metrics().Height.Ceil()
	for _, line := range strings.Split(t.body, "\n") {
		d.Dot = fixed.P(int(t.pos.X), y)
		d.DrawString(line)
		y += lineHeight // move y position to next line
	}

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

// loadFont loads an OpenType font from a .ttf file.
func loadFont(path string) (*sfnt.Font, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	font, err := opentype.Parse(fontBytes)
	// font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return font, nil
}
