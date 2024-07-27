package turdgl

import (
	"image/color"
	"math"
)

// Style contains style information for a shape.
type Style struct {
	Colour    color.Color
	Thickness float64 // leave 0 for solid
}

// shape contains the generic attributes for a 2D shape.
type shape struct {
	Pos   Vec
	w, h  float64 // dimensions
	style Style
}

// newShape constructs a new shape according to the supplied parameters.
func newShape(width, height float64, pos Vec, style Style) shape {
	return shape{
		Pos:   pos,
		w:     width,
		h:     height,
		style: style,
	}
}

// Move modifies the position of the shape by the given vector.
func (s *shape) Move(mov Vec) {
	s.Pos.X += mov.X
	s.Pos.Y += mov.Y
}

// SetPos sets the Cartesian position of the shape.
func (s *shape) SetPos(v Vec) {
	s.Pos = v
}

// Rect is a rectangle shape, aligned to the top-left corner.
type Rect struct{ shape }

// NewRect constructs a new rectangle shape.
func NewRect(width, height float64, pos Vec, style Style) *Rect {
	return &Rect{newShape(width, height, pos, style)}
}

// Draw draws the rectangle onto the provided frame buffer.
func (r *Rect) Draw(buf *FrameBuffer) {
	b := *buf
	colourPixel := NewPixel(r.style.Colour)

	if r.style.Thickness == 0 {
		for i := 0; i <= int(math.Round(r.w)); i++ {
			for j := 0; j <= int(math.Round(r.h)); j++ {
				xInt, yInt := int(math.Round(r.Pos.X)), int(math.Round(r.Pos.Y))
				b[yInt+j][xInt+i] = colourPixel
			}
		}
	} else {
		// Draw each edge as its own rectangle
		top := NewRect(r.w, r.style.Thickness, r.Pos, Style{r.style.Colour, 0})
		bottom := NewRect(r.w, r.style.Thickness, Vec{r.Pos.X, r.Pos.Y + float64(r.h) - float64(r.style.Thickness)}, Style{r.style.Colour, 0})
		left := NewRect(r.style.Thickness, r.h, Vec{r.Pos.X, r.Pos.Y}, Style{r.style.Colour, 0})
		right := NewRect(r.style.Thickness, r.h, Vec{r.Pos.X + float64(r.w) - float64(r.style.Thickness), r.Pos.Y}, Style{r.style.Colour, 0})

		top.Draw(buf)
		bottom.Draw(buf)
		left.Draw(buf)
		right.Draw(buf)
	}
}

// Cicle is a circle shape, aligned to the centre of the circle.
type Circle struct{ shape }

// NewCircle constructs a new circle.
func NewCircle(width, height float64, pos Vec, style Style) *Circle {
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
	bbBoxPos := Vec{c.Pos.X - (radius), c.Pos.Y - (radius)}
	bbox := NewRect(c.w, c.h, bbBoxPos, Style{})

	// Iterate over every pixel in the bounding box
	for i := bbox.Pos.X; i <= bbox.Pos.X+bbox.w; i++ {
		for j := bbox.Pos.Y; j <= bbox.Pos.Y+bbox.h; j++ {
			// Draw pixel if it's close enough to centre
			dist := Dist(c.Pos, Vec{i, j})
			jInt, iInt := int(math.Round(j)), int(math.Round(i))
			if c.style.Thickness == 0 {
				// Solid fill
				if dist <= float64(radius) {
					fb[jInt][iInt] = NewPixel(c.style.Colour)
				}
			} else {
				// Outline
				if dist >= float64(radius-c.style.Thickness) && dist <= float64(radius) {
					fb[jInt][iInt] = NewPixel(c.style.Colour)
				}
			}
		}
	}
}
