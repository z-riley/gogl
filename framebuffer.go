package turdgl

import (
	"image/color"
	"math"
)

// Drawable is an interface for things that can be drawn onto a frame buffer.
type Drawable interface {
	Draw(*FrameBuffer)
}

const pxLen = 4

// Pixel represents a pixel, with bytes [red, green, blue, alpha]
type Pixel [pxLen]uint8

// NewPixel constructs a new coloured pixel.
func NewPixel(c color.Color) Pixel {
	r, g, b, a := c.RGBA()
	return [pxLen]uint8{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
		uint8(a >> 8),
	}
}

func (p Pixel) R() uint8 { return p[0] }
func (p Pixel) G() uint8 { return p[1] }
func (p Pixel) B() uint8 { return p[2] }
func (p Pixel) A() uint8 { return p[3] }

// FrameBuffer is a 2D slice of pixels which represents a screen.
type FrameBuffer [][]Pixel

// NewFrameBuffer constructs a new frame buffer with a particular width and height.
func NewFrameBuffer(width, height int) *FrameBuffer {
	f := FrameBuffer(make([][]Pixel, height))
	for i := 0; i < len(f); i++ {
		f[i] = make([]Pixel, width)
	}
	return &f
}

// SetPixel sets a pixel in the frame buffer. If the requested pixel is out of
// bounds, nothing happens.
func (f *FrameBuffer) SetPixel(y, x int, p Pixel) {
	if y > f.Height()-1 || y < 0 || x > f.Width()-1 || x < 0 {
		return
	}
	// Set pixel value
	(*f)[y][x] = alphaBlend(p, (*f)[y][x])
}

// addToPixel adds the values of each channel in a pixel to frame buffer.
func (f *FrameBuffer) addToPixel(y, x int, p Pixel) {
	if y > f.Height()-1 || y < 0 || x > f.Width()-1 || x < 0 {
		return
	}
	// Set pixel value
	(*f)[y][x] = additionBlend(p, (*f)[y][x])
}

// Clear sets every pixel in the frame buffer to zero.
func (f *FrameBuffer) Clear() {
	f.Fill(color.RGBA{0, 0, 0, 0})
}

// Fill sets every pixel in the frame buffer to the provided colour.
func (f *FrameBuffer) Fill(c color.Color) {
	for i := 0; i < len(*f); i++ {
		for j := 0; j < len((*f)[0]); j++ {
			(*f)[i][j] = NewPixel(c)
		}
	}
}

// Width returns the width of the frame buffer.
func (f *FrameBuffer) Width() int {
	return len((*f)[0])
}

// Height returns the height of the frame buffer.
func (f *FrameBuffer) Height() int {
	return len(*f)
}

// Bytes returns the frame buffer as a one-dimensional slice of bytes.
func (f *FrameBuffer) Bytes() []byte {
	buf := *f
	out := make([]byte, len(buf)*len(buf[0])*pxLen)
	offset := 0
	for i := 0; i < len(buf); i++ {
		slice := buf[i]
		for _, pixel := range slice {
			// Copy the bytes of each pixel in reverse order
			for k := pxLen - 1; k >= 0; k-- {
				out[offset] = pixel[k]
				offset++
			}
		}
	}
	return out
}

// WithinFrame returns true if the given point lies within the boundary of the frame
// buffer, taking the padding value into account.
func (f *FrameBuffer) WithinFrame(point Vec, padding float64) bool {
	w, h := float64(f.Width()), float64(f.Height())
	if padding > w || padding > h {
		panic("WithinFrame - padding cannot be greater than the frame width or height")
	}
	return point.X-1 < w-padding && point.Y-1 < h-padding && point.X > padding && point.Y > padding
}

// DrawLine draws a single pixel line using Bresenham's line drawing algorithm.
func DrawLine(v1, v2 Vec, buf *FrameBuffer) {
	dx := math.Abs(v2.X - v1.X)
	sx := 1
	if v1.X > v2.X {
		sx = -1
	}
	dy := -math.Abs(v2.Y - v1.Y)
	sy := 1
	if v1.Y > v2.Y {
		sy = -1
	}
	err := dx + dy
	x1, y1 := int(math.Round(v1.X)), int(math.Round(v1.Y))
	for {
		buf.SetPixel(y1, x1, NewPixel(color.White))
		x2, y2 := int(math.Round(v2.X)), int(math.Round(v2.Y))
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += sx
		}
		if e2 <= dx {
			err += dx
			y1 += sy
		}
	}
}

// alphaBlend blends a source pixel over a destination pixel using the Porter-Duff
// source-over operator. This results in the source pixel having a greater impact
// on the resulting colour.
func alphaBlend(src, dst Pixel) Pixel {
	srcR, dstR := uint32(src.R()), uint32(dst.R())
	srcG, dstG := uint32(src.G()), uint32(dst.G())
	srcB, dstB := uint32(src.B()), uint32(dst.B())
	srcA, dstA := uint32(src.A()), uint32(dst.A())

	invSrcA := math.MaxUint8 - srcA
	a := srcA + (dstA*(invSrcA))/math.MaxUint8

	// Handle fully transparent case
	if a == 0 {
		return Pixel{0, 0, 0, 0}
	}

	r := (srcR*srcA + dstR*dstA*(invSrcA)) / a
	g := (srcG*srcA + dstG*dstA*(invSrcA)) / a
	b := (srcB*srcA + dstB*dstA*(invSrcA)) / a

	return Pixel{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// additionBlend blends a source pixel with a destination pixel by adding the values
// of each channel.
func additionBlend(src, dst Pixel) Pixel {
	return Pixel{
		Clamp(dst[0]+src[0], 0, 255),
		Clamp(dst[1]+src[1], 0, 255),
		Clamp(dst[2]+src[2], 0, 255),
		Clamp(dst[3]+src[3], 0, 255),
	}
}
