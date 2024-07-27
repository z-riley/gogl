package turdgl

import "image/color"

const pxLen = 4

type Pixel [pxLen]byte // red, green, blue, alpha

// NewPixel constructs a new coloured pixel.
func NewPixel(c color.Color) Pixel {
	r, g, b, a := c.RGBA()
	return [pxLen]byte{
		byte(r >> 8),
		byte(g >> 8),
		byte(b >> 8),
		byte(a >> 8),
	}
}

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

// Clear sets every pixel in the frame buffer to zero.
func (f *FrameBuffer) Clear() {
	for i := 0; i < len(*f); i++ {
		for j := 0; j < len((*f)[0]); j++ {
			(*f)[i][j] = Pixel{}
		}
	}
}

// SetColour sets every pixel in the frame buffer to the provided colour.
func (f *FrameBuffer) SetColour(c color.Color) {
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

// Length returns the length of the frame buffer.
func (f *FrameBuffer) Height() int {
	return len(*f)
}

// Bytes returns the frame buffer as a one-dimensional slice of bytes.
func (f *FrameBuffer) Bytes() []byte {
	buf := *f
	out := make([]byte, len(buf)*len(buf[0])*pxLen)
	offset := 0
	for i := len(buf) - 1; i >= 0; i-- {
		slice := buf[i]
		for _, pixel := range slice {
			copy(out[offset:], pixel[:]) // copy the bytes of each Pixel
			offset += len(pixel)
		}
	}
	return out
}
