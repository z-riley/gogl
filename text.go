package turdgl

import (
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
	customOffset       Vec // only applies in AlignCustom mode
	colour             color.Color
	font               *sfnt.Font
	dpi, size, spacing float64     // settings for generating mask
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
	// Calculate the offset caused by alignment option
	xAlignmentOffset, yAlignmentOffset := func() (int, int) {
		switch t.alignment {
		case AlignTopLeft:
			return 0, 0
		case AlignTopCentre:
			return -t.mask.Rect.Dx() / 2, 0
		case AlignTopRight:
			return -t.mask.Rect.Dx(), 0
		case AlignCentreLeft:
			return 0, -t.mask.Rect.Dy() / 2
		case AlignCentre:
			return -t.mask.Rect.Dx() / 2, -t.mask.Rect.Dy() / 2
		case AlignCentreRight:
			return -t.mask.Rect.Dx(), -t.mask.Rect.Dy() / 2
		case AlignBottomLeft:
			return 0, -t.mask.Rect.Dy()
		case AlignBottomCentre:
			return -t.mask.Rect.Dx() / 2, -t.mask.Rect.Dy()
		case AlignBottomRight:
			return -t.mask.Rect.Dx(), -t.mask.Rect.Dy()
		case AlignCustom:
			return -t.mask.Rect.Dx()/2 + int(math.Round(t.customOffset.X)),
				-t.mask.Rect.Dy()/2 + int(math.Round(t.customOffset.Y))
		default:
			panic(fmt.Errorf("unsupported text alignment: %v", t.alignment))
		}
	}()

	// Write pixels to frame buffer
	r, g, b, a := RGBA8(t.colour)
	for y := t.mask.Rect.Min.Y; y < t.mask.Rect.Max.Y; y++ {
		for x := t.mask.Rect.Min.X; x < t.mask.Rect.Max.X; x++ {
			maskRGBA := t.mask.RGBAAt(x, y)
			if maskRGBA.A > 0 {
				posX := int(t.pos.X) + x + xAlignmentOffset
				posY := int(t.pos.Y) + y + yAlignmentOffset

				// OpenType reduces the RGB values of border pixels for anti-aliasing.
				// However, we are using alpha blending, so reset the RBG values to their
				// original, and only keep the alpha channel.
				isBorderPixel := maskRGBA.A < a
				if isBorderPixel {
					buf.SetPixel(posY, posX, NewPixel(color.RGBA{r, g, b, maskRGBA.A}))
				} else {
					buf.SetPixel(posY, posX, NewPixel(maskRGBA))
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
func (t *Text) SetAlignment(align Alignment) *Text {
	t.alignment = align
	if err := t.generateMask(); err != nil {
		panic(err)
	}
	return t
}

// Offset returns the text's offset.
func (t *Text) Offset() Vec {
	return t.customOffset
}

// SetOffset sets the text's offset. Note: this overwrites existing alignment settings.
func (t *Text) SetOffset(offset Vec) *Text {
	t.SetAlignment(AlignCustom)
	t.customOffset = offset
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

	splitLines := strings.Split(t.body, "\n")

	// Make the mask width the size of the longest line of text
	maskWidth := func() (w int) {
		for _, line := range splitLines {
			_, adv := font.BoundString(face, line)
			if adv.Ceil() > w {
				w = adv.Ceil()
			}
		}
		return w
	}()

	// Make the mask height slightly larger to allow for characters that are drawn
	// below the line (like underscore)
	faceHeight := face.Metrics().Height.Ceil()
	maskHeight := float64(faceHeight*len(splitLines)) + (0.4 * float64(faceHeight))

	// Draw the font into the mask
	mask := image.NewRGBA(image.Rect(0, 0, maskWidth, int(maskHeight)))
	drawer := &font.Drawer{
		Dst:  mask,
		Src:  image.NewUniform(t.colour),
		Face: face,
		Dot:  fixed.Point26_6{}, // set later
	}

	// Draw each line of text separately
	for i, line := range splitLines {
		// Move the dot to correct position
		switch t.alignment {
		case AlignTopLeft, AlignCentreLeft, AlignBottomLeft:
			drawer.Dot = fixed.P(0, faceHeight*(1+i))
		case AlignCentre, AlignTopCentre, AlignBottomCentre, AlignCustom:
			_, adv := font.BoundString(face, line)
			drawer.Dot = fixed.P((maskWidth-adv.Ceil())/2, faceHeight*(1+i))
		case AlignTopRight, AlignCentreRight, AlignBottomRight:
			_, adv := font.BoundString(face, line)
			drawer.Dot = fixed.P(maskWidth-adv.Ceil(), faceHeight*(1+i))
		}

		// Draw the line
		drawer.DrawString(line)
	}

	t.mask = mask

	return nil
}

// loadFont loads an OpenType font from a .ttf file.
func loadFont(path string) (*sfnt.Font, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return opentype.Parse(fontBytes)
}

// saveImageAsPNG saves an image as a PNG file. Used for testing and development.
func saveImageAsPNG(img image.Image, filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}
