package gogl

import (
	"image/color"
	"math"
)

// Cicle is a circle shape, aligned to the centre of the circle.
type Circle struct {
	Pos       Vec
	Direction Vec
	d         float64
	style     Style
}

var _ Shape = (*Circle)(nil)
var _ hoverable = (*Circle)(nil)

// NewCircle constructs a new circle.
func NewCircle(diameter float64, pos Vec) *Circle {
	return &Circle{
		Pos:       pos,
		Direction: Vec{0, 0},
		d:         diameter,
		style:     DefaultStyle,
	}
}

// IsWithin returns whether a position lies within the circle's perimeter.
func (c *Circle) IsWithin(pos Vec) bool {
	return Dist(c.Pos, pos) <= c.Width()/2
}

// Draw draws the circle onto the provided frame buffer.
func (c *Circle) Draw(buf *FrameBuffer) {
	thickness := c.style.Thickness
	if c.style.Thickness == 0 { // for filled shape
		thickness = c.d / 2
	}

	// Construct bounding box
	radius := c.d / 2
	bbBoxPos := Vec{c.Pos.X - radius, c.Pos.Y - radius}
	bbox := NewRect(c.d, c.d, bbBoxPos)

	// Iterate over every pixel in the bounding box
	for x := bbox.Pos.X; x <= bbox.Pos.X+bbox.w; x++ {
		for y := bbox.Pos.Y; y <= bbox.Pos.Y+bbox.h; y++ {
			// Draw pixel if it's close enough to centre
			dist := Dist(c.Pos, Vec{x, y})
			if dist >= float64(radius-thickness) && dist <= float64(radius) {
				buf.SetPixel(int(x), int(y), NewPixel(c.style.Colour))
			}
		}
	}

	// Draw bloom if it exists
	if c.style.Bloom > 0 {
		c.drawBloom(buf)
	}
}

// DrawCircleSegment draws only a segment of the circle to the frame buffer, limited by the
// provided vector.
func (c *Circle) DrawCircleSegment(limitDir Vec, buf *FrameBuffer) {
	// Construct bounding box
	radius := c.d / 2
	bbBoxPos := Vec{c.Pos.X - (radius), c.Pos.Y - (radius)}
	bbox := NewRect(c.d, c.d, bbBoxPos)

	// Iterate over every pixel in the bounding box
	for x := bbox.Pos.X; x <= bbox.Pos.X+bbox.w; x++ {
		for y := bbox.Pos.Y; y <= bbox.Pos.Y+bbox.h; y++ {
			// Draw pixel if it's close enough to centre
			dist := Dist(c.Pos, Vec{x, y})
			jInt, iInt := int(math.Round(y)), int(math.Round(x))
			if c.style.Thickness == 0 {
				// Solid fill
				if dist <= float64(radius) && Theta(c.Direction, Sub(Vec{x, y}, c.Pos)) >= 0 {
					buf.SetPixel(iInt, jInt, NewPixel(c.style.Colour))
				}
			} else {
				// Outline
				if dist >= float64(radius-c.style.Thickness) && dist <= float64(radius) &&
					Theta(c.Direction, Sub(Vec{x, y}, c.Pos)) >= Theta(Upwards, limitDir) {
					buf.SetPixel(iInt, jInt, NewPixel(color.White))
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

// Width returns the pixel width of the circle.
func (c *Circle) Width() float64 {
	return c.d
}

// SetDiameter sets the diameter of the circle, in pixels.
func (c *Circle) SetDiameter(px float64) *Circle {
	c.d = px
	return c
}

// Height returns the pixel height of the circle.
func (c *Circle) Height() float64 {
	return c.d
}

// GetPos returns the position of the circle.
func (c *Circle) GetPos() Vec {
	return c.Pos
}

// SetPos sets the position of the circle.
func (c *Circle) SetPos(pos Vec) {
	c.Pos = pos
}

// GetStyle returns the style of the circle.
func (c *Circle) GetStyle() Style {
	return c.style
}

// SetStyle sets the style of the circle.
func (c *Circle) SetStyle(style Style) *Circle {
	c.style = style
	return c
}

// Move moves the circle's position by the given pixel vector.
func (c *Circle) Move(px Vec) {
	c.Pos = Add(c.Pos, px)
}

// String returns the type of shape as a string.
func (c *Circle) String() string {
	return "circle"
}

// drawBloom draws a bloom effect around a circle.
func (c *Circle) drawBloom(buf *FrameBuffer) {
	bloom := float64(c.style.Bloom)

	// Construct bounding box
	radius := c.d / 2
	bbBoxPos := Vec{c.Pos.X - radius - bloom, c.Pos.Y - radius - bloom}
	bbox := NewRect(c.d+bloom+bloom, c.d+bloom+bloom, bbBoxPos)

	// Iterate over every pixel in the bounding box
	r, g, b, a := RGBA8(c.style.Colour)
	for x := bbox.Pos.X; x <= bbox.Pos.X+bbox.w; x++ {
		for y := bbox.Pos.Y; y <= bbox.Pos.Y+bbox.h; y++ {
			dist := Dist(c.Pos, Vec{x, y})
			if dist >= radius && dist <= radius+bloom {
				brightness := 1 - ((dist - radius) / bloom)
				bloomColour := color.RGBA{r, g, b, uint8(brightness * float64(a))}
				buf.SetPixel(int(x), int(y), NewPixel(bloomColour))
			}
		}
	}
}

var (
	Upwards    = Vec{0, -1}
	Downwards  = Vec{0, 1}
	Leftwards  = Vec{-1, 0}
	Rightwards = Vec{1, 0}
)
