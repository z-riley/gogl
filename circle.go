package turdgl

import (
	"fmt"
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
		Direction: Vec{},
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
	for i := bbox.Pos.X; i <= bbox.Pos.X+bbox.w; i++ {
		for j := bbox.Pos.Y; j <= bbox.Pos.Y+bbox.h; j++ {
			// Draw pixel if it's close enough to centre
			dist := Dist(c.Pos, Vec{i, j})
			if dist >= float64(radius-thickness) && dist <= float64(radius) {
				buf.SetPixel(int(j), int(i), NewPixel(c.style.Colour))
			}
		}
	}

	// Draw bloom if it exists
	if c.style.Bloom > 0 {
		c.drawBloom(buf)
	}
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
	for i := bbox.Pos.X; i <= bbox.Pos.X+bbox.w; i++ {
		for j := bbox.Pos.Y; j <= bbox.Pos.Y+bbox.h; j++ {
			dist := Dist(c.Pos, Vec{i, j})
			if dist >= radius && dist <= radius+bloom {
				brightness := 1 - ((dist - radius) / bloom)
				if brightness < 0 {
					fmt.Println(brightness)
				}
				bloomColour := color.RGBA{r, g, b, uint8(brightness * float64(a))}
				buf.SetPixel(int(j), int(i), NewPixel(bloomColour))
			}
		}
	}
}

func (c *Circle) DrawCircleSegment(limitDir Vec, buf *FrameBuffer) {
	// Construct bounding box
	radius := c.d / 2
	bbBoxPos := Vec{c.Pos.X - (radius), c.Pos.Y - (radius)}
	bbox := NewRect(c.d, c.d, bbBoxPos)

	// Iterate over every pixel in the bounding box
	for i := bbox.Pos.X; i <= bbox.Pos.X+bbox.w; i++ {
		for j := bbox.Pos.Y; j <= bbox.Pos.Y+bbox.h; j++ {
			// Draw pixel if it's close enough to centre
			dist := Dist(c.Pos, Vec{i, j})
			jInt, iInt := int(math.Round(j)), int(math.Round(i))
			if c.style.Thickness == 0 {
				// Solid fill
				if dist <= float64(radius) && Theta(c.Direction, Sub(Vec{i, j}, c.Pos)) >= 0 {
					buf.SetPixel(jInt, iInt, NewPixel(c.style.Colour))
				}
			} else {
				// Outline
				if dist >= float64(radius-c.style.Thickness) && dist <= float64(radius) &&
					Theta(c.Direction, Sub(Vec{i, j}, c.Pos)) >= Theta(Upwards, limitDir) {
					buf.SetPixel(jInt, iInt, NewPixel(color.White))
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

func (c *Circle) Width() float64 {
	return c.d
}

func (c *Circle) SetDiameter(px float64) *Circle {
	c.d = px
	return c
}

func (c *Circle) Height() float64 {
	return c.d
}

func (c *Circle) GetPos() Vec {
	return c.Pos
}

func (c *Circle) SetPos(pos Vec) *Circle {
	c.Pos = pos
	return c
}

func (c *Circle) GetStyle() Style {
	return c.style
}

func (c *Circle) SetStyle(style Style) *Circle {
	c.style = style
	return c
}

func (c *Circle) Move(px Vec) {
	c.Pos = Add(c.Pos, px)
}

func (c *Circle) String() string {
	return "circle"
}

var (
	Upwards    = Vec{0, -1}
	Downwards  = Vec{0, 1}
	Leftwards  = Vec{-1, 0}
	Rightwards = Vec{1, 0}
)
