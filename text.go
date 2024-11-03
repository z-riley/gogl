package turdgl

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
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

// NewText constructs a new text object with default parameters. The default font
// size is 20.
func NewText(body string, pos Vec, fontPath string) *Text {
	t := Text{
		body:      body,
		pos:       pos,
		alignment: AlignTopLeft,
		colour:    Red,
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
	// Write pixels to frame buffer
	bbox := t.textBoundry()
	r, g, b, a := RGBA8(t.colour)
	startX, startY := int(bbox.Pos.X), int(bbox.Pos.Y)
	endX, endY := startX+int(bbox.w)+100, startY+int(bbox.h)+100
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			maskRGBA := t.mask.RGBAAt(x, y)
			if maskRGBA.A > 0 {
				// OpenType reduces the RGB values of border pixels for anti-aliasing.
				// However, we are using alpha blending, so reset the RBG values to their
				// original, and only keep the alpha channel.
				aaPixel := maskRGBA.A < a
				if aaPixel {
					buf.SetPixel(y, x, NewPixel(color.RGBA{r, g, b, maskRGBA.A}))
				} else {
					buf.SetPixel(y, x, NewPixel(maskRGBA))
				}
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
	// Load the font face
	face, err := opentype.NewFace(t.font, &opentype.FaceOptions{
		Size:    t.size,
		DPI:     t.dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("failed to create font face: %w", err)
	}
	defer face.Close()

	// TODO: This could be optimised to only be the size of the text instead of the entire window
	mask := image.NewRGBA(image.Rect(
		0, 0,
		t.width, t.height,
	))

	d := &font.Drawer{
		Dst:  mask,
		Src:  image.NewUniform(t.colour),
		Face: face,
		Dot:  fixed.Point26_6{}, // set later
	}

	faceHeight := face.Metrics().Height.Ceil()
	bounds, _ := font.BoundString(face, t.body)
	lineHeight := bounds.Max.Y.Ceil() - bounds.Min.Y.Floor()

	// Draw each line of text seperately
	splitLines := strings.Split(t.body, "\n")
	numLines := len(splitLines)
	for i, line := range splitLines {
		// Calculate offset caused by alignment option
		xOffset, yOffset, err := func() (int, int, error) {
			bounds, _ := font.BoundString(face, line)
			w := bounds.Max.X.Ceil() - bounds.Min.X.Floor()

			bounds, _ = font.BoundString(face, t.body)
			h := lineHeight

			switch t.alignment {
			case AlignTopLeft:
				return 0, h, nil
			case AlignTopCentre:
				return -w / 2, h, nil
			case AlignTopRight:
				return -w, h, nil
			case AlignCentreLeft:
				return 0, h - (faceHeight * numLines / 2), nil
			case AlignCentre:
				return -w / 2, h - (faceHeight * numLines / 2), nil
			case AlignCentreRight:
				return -w, h - (faceHeight * numLines / 2), nil
			case AlignBottomLeft:
				return 0, 0 - faceHeight*(numLines-1), nil
			case AlignBottomCentre:
				return -w / 2, 0 - faceHeight*(numLines-1), nil
			case AlignBottomRight:
				return -w, 0 - faceHeight*(numLines-1), nil
			case AlignCustom:
				return -w/2 + int(math.Round(t.labelOffset.X)),
					h - (faceHeight * numLines / 2) + int(math.Round(t.labelOffset.Y)),
					nil
			default:
				return 0, 0, errors.New("unsupported text alignment")
			}
		}()
		if err != nil {
			return fmt.Errorf("failed to generate text mask: %w", err)
		}

		// Move the dot to correct position and draw the line
		d.Dot = fixed.P(
			int(t.pos.X)+xOffset,
			int(t.pos.Y)+(i*faceHeight)+yOffset,
		)
		d.DrawString(line)
	}

	t.mask = mask

	return nil
}

func SaveImageAsPNG(img image.Image, filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
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
		maxX-minX+1,
		maxY-minY+1,
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
	if err != nil {
		return nil, err
	}
	return font, nil
}
