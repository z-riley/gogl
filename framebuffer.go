package gogl

import (
	"image/color"
	"math"
	"unsafe"
)

// Drawable is an interface for things that can be drawn onto a frame buffer.
type Drawable interface {
	// Draw draws to a frame buffer.
	Draw(*FrameBuffer)
}

const pxLen = 4

// Pixel packs a pixel into a 32-bit word with bytes [alpha, blue, green, red].
// Matches SDL's RGBA8888 layout on little-endian systems when viewed as raw memory.
type Pixel uint32

func NewPixel(c color.Color) Pixel {
	r, g, b, a := RGBA8(c)
	return pack(a, b, g, r)
}

func (p Pixel) A() uint8 { return uint8(p) }
func (p Pixel) B() uint8 { return uint8(p >> 8) }
func (p Pixel) G() uint8 { return uint8(p >> 16) }
func (p Pixel) R() uint8 { return uint8(p >> 24) }

// pack packs individual channel bytes into Pixel.
func pack(a, b, g, r uint8) Pixel {
	return Pixel(uint32(a) | uint32(b)<<8 | uint32(g)<<16 | uint32(r)<<24)
}

type FrameBuffer struct {
	fb     []Pixel
	width  int
	height int
}

// NewFrameBuffer constructs a new frame buffer with a particular width and height.
func NewFrameBuffer(width, height int) *FrameBuffer {
	if width < 1 || height < 1 {
		panic("invalid frame buffer size")
	}

	return &FrameBuffer{
		fb:     make([]Pixel, width*height),
		width:  width,
		height: height,
	}
}

// GetPixel returns a copy of the pixel at the specified coordinates.
func (f *FrameBuffer) GetPixel(x, y int) Pixel {
	if y > f.Height()-1 || y < 0 || x > f.Width()-1 || x < 0 {
		panic("GetPixel out of bounds")
	}

	return f.getPixel(x, y)
}

func (f *FrameBuffer) getPixel(x, y int) Pixel {
	targetPix := x + f.width*y
	return f.fb[targetPix]
}

// BlendFunc is a function that blends a source pixel with a destination pixel.
type BlendFunc func(src, dst Pixel) Pixel

// SetPixelFunc sets a pixel in the frame buffer using the specified blend function.
func (f *FrameBuffer) SetPixelFunc(x, y int, p Pixel, blend BlendFunc) {
	if y > f.Height()-1 || y < 0 || x > f.Width()-1 || x < 0 {
		return
	}

	px := blend(p, f.getPixel(x, y))

	f.setPixel(x, y, px)
}

func (f *FrameBuffer) setPixel(x, y int, p Pixel) {
	targetPix := x + f.width*y
	f.fb[targetPix] = p
}

// SetPixel sets a pixel in the frame buffer. If the requested pixel is out of
// bounds, nothing happens. The default alpha blending technique is used. To use
// other blending methods, see SetPixelFunc.
func (f *FrameBuffer) SetPixel(x, y int, p Pixel) {
	f.SetPixelFunc(x, y, p, AlphaBlend)
}

// Clear sets every pixel in the frame buffer to zero.
func (f *FrameBuffer) Clear() {
	f.Fill(color.RGBA{0, 0, 0, 0})
}

// Fill sets every pixel in the frame buffer to the provided colour.
func (f *FrameBuffer) Fill(c color.Color) {
	p := NewPixel(c)

	for i := range f.fb {
		f.fb[i] = p
	}
}

// Width returns the width of the frame buffer.
func (f *FrameBuffer) Width() int {
	return f.width
}

// Height returns the height of the frame buffer.
func (f *FrameBuffer) Height() int {
	return f.height
}

// Bytes returns the frame buffer as a one-dimensional slice of bytes.
func (f *FrameBuffer) Bytes() []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&f.fb[0])), len(f.fb)*pxLen)
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
		buf.SetPixel(x1, y1, NewPixel(color.White))
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

// AlphaBlend blends a source pixel over a destination pixel using the Porter-Duff
// source-over operator. This results in the source pixel having a greater impact
// on the resulting colour.
func AlphaBlend(src, dst Pixel) Pixel {
	srcR, dstR := uint32(src.R()), uint32(dst.R())
	srcG, dstG := uint32(src.G()), uint32(dst.G())
	srcB, dstB := uint32(src.B()), uint32(dst.B())
	srcA, dstA := uint32(src.A()), uint32(dst.A())

	invSrcA := math.MaxUint8 - srcA

	// Resulting alpha calculation (Porter-Duff source-over)
	a := Clamp(srcA+(dstA*invSrcA)/math.MaxUint8, 0, math.MaxUint8)

	// Handle fully transparent case
	if a == 0 {
		return 0
	}

	// Resulting color channels calculation
	r := (srcR*srcA + dstR*dstA*invSrcA/math.MaxUint8) / a
	g := (srcG*srcA + dstG*dstA*invSrcA/math.MaxUint8) / a
	b := (srcB*srcA + dstB*dstA*invSrcA/math.MaxUint8) / a

	// Clamp the resulting colors to avoid overflow
	r = Clamp(r, 0, math.MaxUint8)
	g = Clamp(g, 0, math.MaxUint8)
	b = Clamp(b, 0, math.MaxUint8)

	return pack(uint8(a), uint8(b), uint8(g), uint8(r))
}

// AdditiveBlend blends a source pixel with a destination pixel by adding the values
// of each channel.
func AdditiveBlend(src, dst Pixel) Pixel {
	// Clamp to stop overflow
	b := Clamp(uint16(dst.B())+uint16(src.B()), 0, math.MaxUint8)
	g := Clamp(uint16(dst.G())+uint16(src.G()), 0, math.MaxUint8)
	r := Clamp(uint16(dst.R())+uint16(src.R()), 0, math.MaxUint8)

	return pack(src.A(), uint8(b), uint8(g), uint8(r))
}
