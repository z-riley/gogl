package turdgl

import (
	"image/color"
)

// Style contains style information for a shape.
type Style struct {
	Colour    color.Color
	Thickness int // leave 0 for solid
}

// shape contains the generic attributes for a 2D shape.
type shape struct {
	pos   Vec
	w, h  int // dimensions
	style Style
}

// newShape constructs a new shape according to the supplied parameters.
func newShape(width, height int, pos Vec, style Style) shape {
	return shape{
		pos:   pos,
		w:     width,
		h:     height,
		style: style,
	}
}

// Move modifies the position of the shape by the given vector.
func (r *shape) Move(mov Vec) {
	r.pos.X += mov.X
	r.pos.Y += mov.Y
}

// Pos returns the Cartesian position of the shape.
func (r *shape) Pos() Vec {
	return r.pos
}

// Rect is a rectangle shape, aligned to the top-left corner.
type Rect struct{ shape }

// NewRect constructs a new rectangle shape.
func NewRect(width, height int, pos Vec, style Style) *Rect {
	return &Rect{newShape(width, height, pos, style)}
}

// Draw draws the rectangle onto the provided frame buffer.
func (r *Rect) Draw(buf *FrameBuffer) {
	b := *buf
	colourPixel := NewPixel(r.style.Colour)

	if r.style.Thickness == 0 {
		for i := 0; i <= r.w; i++ {
			for j := 0; j <= r.h; j++ {
				b[r.pos.Y+j][r.pos.X+i] = colourPixel
			}
		}
	} else {
		// Draw each edge as its own rectangle
		top := NewRect(r.w, r.style.Thickness, r.pos, Style{r.style.Colour, 0})
		bottom := NewRect(r.w, r.style.Thickness, Vec{r.pos.X, r.pos.Y + r.h - r.style.Thickness}, Style{r.style.Colour, 0})
		left := NewRect(r.style.Thickness, r.h, Vec{r.pos.X, r.pos.Y}, Style{r.style.Colour, 0})
		right := NewRect(r.style.Thickness, r.h, Vec{r.pos.X + r.w - r.style.Thickness, r.pos.Y}, Style{r.style.Colour, 0})

		top.Draw(buf)
		bottom.Draw(buf)
		left.Draw(buf)
		right.Draw(buf)
	}
}

// Cicle is a circle shape, aligned to the centre of the circle.
type Circle struct{ shape }

// NewCircle constructs a new circle.
func NewCircle(width, height int, pos Vec, style Style) *Circle {
	return &Circle{newShape(width, height, pos, style)}
}

// Draw draws the circle onto the provided frame buffer.
func (c *Circle) Draw(buf *FrameBuffer) {
	if c.w != c.h {
		panic("circle width and height must match")
	}

	fb := *buf
	radius := c.w / 2

	// Construct bounding box
	bbBoxPos := Vec{c.pos.X - (radius), c.pos.Y - (radius)}
	bbox := NewRect(c.w, c.h, bbBoxPos, Style{})

	// Iterate over every pixel in the bounding box
	for i := bbox.pos.X; i <= bbox.pos.X+bbox.w; i++ {
		for j := bbox.pos.Y; j <= bbox.pos.Y+bbox.h; j++ {
			// Draw pixel if it's close enough to centre
			dist := Dist(c.pos, Vec{i, j})
			if c.style.Thickness == 0 {
				// Solid fill
				if dist <= float64(radius) {
					fb[j][i] = NewPixel(c.style.Colour)
				}
			} else {
				// Outline
				if dist >= float64(radius-c.style.Thickness) && dist <= float64(radius) {
					fb[j][i] = NewPixel(c.style.Colour)
				}
			}
		}
	}
}
