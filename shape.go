package turdgl

import (
	"fmt"
	"image/color"
	"math"
	"reflect"
)

// Style contains style information for a shape.
type Style struct {
	Colour    color.Color
	Thickness float64 // leave 0 for solid
}

var DefaultStyle = Style{Colour: color.RGBA{0xff, 0xff, 0xff, 0xff}, Thickness: 0}

// Shape is an interface for shapes.
type Shape interface {
	GetPos() Vec
	Width() float64
	Height() float64
	Draw(*FrameBuffer)
}

// shape contains the generic attributes for a 2D shape.
type shape struct {
	Pos       Vec
	Direction Vec
	w, h      float64
	style     Style
}

// newShape constructs a new shape according to the supplied parameters.
func newShape(width, height float64, pos Vec, opts ...func(*shape)) *shape {
	s := defaultShape(width, height, pos)
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithStyle is used in the newShape constructor for setting a shape's style.
func WithStyle(style Style) func(*shape) {
	return func(s *shape) {
		s.style = style
	}
}

// defaultShape constructs a shape with default parameters.
func defaultShape(width, height float64, pos Vec) *shape {
	return &shape{
		Pos:       pos,
		Direction: Normalise(Vec{0, -1}), // upwards
		w:         width,
		h:         height,
		style:     DefaultStyle,
	}
}

// Width returns the width of the shape.
func (s *shape) Width() float64 {
	return s.w
}

// Height returns the height of the shape.
func (s *shape) Height() float64 {
	return s.h
}

// Move modifies the position of the shape by the given vector.
func (s *shape) Move(mov Vec) {
	s.Pos.X += mov.X
	s.Pos.Y += mov.Y
}

// Pos returns the Cartesian position of the shape.
func (s *shape) GetPos() Vec {
	return s.Pos
}

// SetPos sets the Cartesian position of the shape.
func (s *shape) SetPos(v Vec) {
	s.Pos = v
}

// GetStyle returns the style of the shape.
func (s *shape) GetStyle() Style {
	return s.style
}

// SetStyle sets style of the shape.
func (s *shape) SetStyle(style Style) {
	s.style = style
}

// Rect is a rectangle shape, aligned to the top-left corner.
type Rect struct{ *shape }

// NewRect constructs a new rectangle shape.
func NewRect(width, height float64, pos Vec, opts ...func(*shape)) *Rect {
	return &Rect{newShape(width, height, pos, opts...)}
}

// Draw draws the rectangle onto the provided frame buffer.
func (r *Rect) Draw(buf *FrameBuffer) {
	if r.style.Thickness == 0 {
		for i := 0; i <= int(math.Round(r.w)); i++ {
			for j := 0; j <= int(math.Round(r.h)); j++ {
				xInt, yInt := int(math.Round(r.Pos.X)), int(math.Round(r.Pos.Y))
				buf.SetPixel(yInt+j, xInt+i, NewPixel(r.style.Colour))
			}
		}
	} else {
		// Draw each edge as its own rectangle
		top := NewRect(r.w, r.style.Thickness, r.Pos,
			WithStyle(Style{r.style.Colour, 0}),
		)
		bottom := NewRect(
			r.w, r.style.Thickness, Vec{r.Pos.X, r.Pos.Y + float64(r.h) - float64(r.style.Thickness)},
			WithStyle(Style{r.style.Colour, 0}),
		)
		left := NewRect(r.style.Thickness, r.h, Vec{r.Pos.X, r.Pos.Y}, WithStyle(Style{r.style.Colour, 0}))
		right := NewRect(r.style.Thickness, r.h, Vec{r.Pos.X + float64(r.w) - float64(r.style.Thickness), r.Pos.Y},
			WithStyle(Style{r.style.Colour, 0}),
		)

		top.Draw(buf)
		bottom.Draw(buf)
		left.Draw(buf)
		right.Draw(buf)
	}
}

// Cicle is a circle shape, aligned to the centre of the circle.
type Circle struct{ *shape }

// NewCircle constructs a new circle.
func NewCircle(diameter float64, pos Vec, opts ...func(*shape)) *Circle {
	return &Circle{newShape(diameter, diameter, pos, opts...)}
}

// Draw draws the circle onto the provided frame buffer.
func (c *Circle) Draw(buf *FrameBuffer) {
	if c.w != c.h {
		fmt.Println("w:", c.w, "h:", c.h)
		panic("circle width and height must match")
	}

	// Construct bounding box
	radius := c.w / 2
	bbBoxPos := Vec{c.Pos.X - (radius), c.Pos.Y - (radius)}
	bbox := NewRect(c.w, c.h, bbBoxPos)

	// Iterate over every pixel in the bounding box
	for i := bbox.Pos.X; i <= bbox.Pos.X+bbox.w; i++ {
		for j := bbox.Pos.Y; j <= bbox.Pos.Y+bbox.h; j++ {
			// Draw pixel if it's close enough to centre
			dist := Dist(c.Pos, Vec{i, j})
			jInt, iInt := int(math.Round(j)), int(math.Round(i))
			if c.style.Thickness == 0 {
				// Solid fill
				if dist <= float64(radius) {
					buf.SetPixel(jInt, iInt, NewPixel(c.style.Colour))
				}
			} else {
				// Outline
				if dist >= float64(radius-c.style.Thickness) && dist <= float64(radius) {
					buf.SetPixel(jInt, iInt, NewPixel(c.style.Colour))

				}
			}
		}
	}
}

// EdgePoint generates a point the point on the circle's perimeter that is theta radians
// clockwise from the circle's direction.
func (c *Circle) EdgePoint(theta float64) Vec {
	return Add(c.Pos, (c.Direction.SetMag(c.Width() / 2).Rotate(theta)))
}

// IsColliding returns true if two shapes are colliding.
func IsColliding(s1, s2 Shape) bool {
	switch s1.(type) {
	case *Rect:
		switch s2.(type) {
		case *Rect:
			// Rect-Rect
			onLeft := s1.GetPos().X > s2.GetPos().X+s2.Width()
			onRight := s2.GetPos().X > s1.GetPos().X+s1.Width()
			above := s1.GetPos().Y > s2.GetPos().Y+s2.Height()
			below := s2.GetPos().Y > s1.GetPos().Y+s1.Height()
			if !onRight && !onLeft && !above && !below {
				return true
			} else {
				return false
			}
		case *Circle:
			// Rect-Circle:
			onLeft := s1.GetPos().X > s2.GetPos().X+s2.Width()/2
			onRight := s2.GetPos().X-s2.Width()/2 > s1.GetPos().X+s1.Width()
			above := s1.GetPos().Y > s2.GetPos().Y+s2.Height()/2
			below := s2.GetPos().Y-s2.Height()/2 > s1.GetPos().Y+s1.Height()
			if !onRight && !onLeft && !above && !below {
				return true
			} else {
				return false
			}
		default:
			t := reflect.TypeOf(s2).String()
			panic("collision detection is unsupported for type: %s" + t)
		}
	case *Circle:
		switch s2.(type) {
		case *Rect:
			// Circle-Rect
			onLeft := s2.GetPos().X > s1.GetPos().X+s1.Width()/2
			onRight := s1.GetPos().X-s1.Width()/2 > s2.GetPos().X+s2.Width()
			above := s2.GetPos().Y > s1.GetPos().Y+s1.Height()/2
			below := s1.GetPos().Y-s1.Height()/2 > s2.GetPos().Y+s2.Height()
			if !onRight && !onLeft && !above && !below {
				return true
			} else {
				return false
			}
		case *Circle:
			// Circle-Circle
			if Dist(s1.GetPos(), s2.GetPos()) <= s1.Width()/2+s2.Width()/2 {
				return true
			} else {
				return false
			}
		default:
			t := reflect.TypeOf(s2).String()
			panic("collision detection is unsupported for type: %s" + t)
		}
	default:
		t := reflect.TypeOf(s1).String()
		panic("collision detection is unsupported for type: %s" + t)
	}
}
