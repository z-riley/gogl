package turdgl

import (
	"fmt"
	"image/color"

	"golang.org/x/exp/rand"
)

// Style contains style information for a shape.
type Style struct {
	Colour    color.Color
	Thickness float64 // leave 0 for solid
	Bloom     int     // bloom reach, in pixels
}

var DefaultStyle = Style{
	Colour:    color.RGBA{0xff, 0xff, 0xff, 0xff},
	Thickness: 0,
	Bloom:     0,
}

// RandomStyle generates a style of random colour and thickness 0.
func RandomStyle() Style {
	return Style{
		Colour: color.RGBA{
			R: byte(rand.Intn(256)),
			G: byte(rand.Intn(256)),
			B: byte(rand.Intn(256)),
			A: byte(rand.Intn(256)),
		},
		Thickness: 0,
	}
}

// Shape is an interface for shapes.
type Shape interface {
	Drawable

	// Width returns the width of the shape in pixels.
	Width() float64
	// SetWidth sets the width of the shape in pixels.
	// SetWidth(float64)
	// Height returns the height of she shape in pixels.
	Height() float64
	// SetHeight sets the height of the shape in pixels.
	// SetHeight(float64)
	// GetPos returns the position of the shape.
	GetPos() Vec
	// SetPos sets the position of the shape.
	// SetPos(Vec)
	// GetStyle returns the shape's style.
	// GetStyle() Style
	// SetStyle sets the shape's style.
	// SetStyle(Style)
	// Move moves the shape by a pixel vector.
	Move(Vec)
	// String returns the name of the shape.
	String() string
}

var (
	Upwards    = Vec{0, -1}
	Downwards  = Vec{0, 1}
	Leftwards  = Vec{-1, 0}
	Rightwards = Vec{1, 0}
)

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
			panic(fmt.Sprintf("collision detection is unsupported for type: %s" + s2.String()))
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
			panic(fmt.Sprintf("collision detection is unsupported for type: %s" + s2.String()))
		}
	case *CurvedRect:
		panic("todo")
	default:
		panic(fmt.Sprintf("collision detection is unsupported for type: %s" + s1.String()))
	}
}
