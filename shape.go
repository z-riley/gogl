package turdgl

import (
	"fmt"
	"image/color"
	"math"
	"reflect"

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
	Width() float64
	SetWidth(float64)
	Height() float64
	SetHeight(float64)
	Draw(*FrameBuffer)
	GetPos() Vec
	SetPos(Vec)
	GetStyle() Style
	SetStyle(style Style)
	Move(Vec)
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

var (
	Upwards    = Vec{0, -1}
	Downwards  = Vec{0, 1}
	Leftwards  = Vec{-1, 0}
	Rightwards = Vec{1, 0}
)

// WithDirection is used in the newShape constructor for setting a shape's starting direction.
func WithDirection(direction Vec) func(*shape) {
	return func(s *shape) {
		s.Direction = direction
	}
}

// defaultShape constructs a shape with default parameters.
func defaultShape(width, height float64, pos Vec) *shape {
	return &shape{
		Pos:       pos,
		Direction: Upwards,
		w:         width,
		h:         height,
		style:     DefaultStyle,
	}
}

// Width returns the width of the shape.
func (s *shape) Width() float64 {
	return s.w
}

// SetWidth sets the width of the shape, in pixels.
func (s *shape) SetWidth(w float64) {
	s.w = w
}

// Height returns the height of the shape.
func (s *shape) Height() float64 {
	return s.h
}

// SetHeight sets the height of the shape, in pixels.
func (s *shape) SetHeight(h float64) {
	s.h = h
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
		width := int(math.Round(r.w))
		height := int(math.Round(r.h))
		for i := 0; i <= width; i++ {
			for j := 0; j <= height; j++ {
				xInt, yInt := int(math.Round(r.Pos.X)), int(math.Round(r.Pos.Y))
				buf.SetPixel(yInt+j, xInt+i, NewPixel(r.style.Colour))
			}
		}
		if r.style.Bloom > 0 {
			r.drawBloom(buf)
		}
	} else {
		// Draw each edge as its own rectangle
		NewRect(r.w, r.style.Thickness, r.Pos,
			WithStyle(Style{r.style.Colour, 0, 0}),
		).Draw(buf)
		NewRect(
			r.w, r.style.Thickness, Vec{r.Pos.X, r.Pos.Y + float64(r.h) - float64(r.style.Thickness)},
			WithStyle(Style{r.style.Colour, 0, 0}),
		).Draw(buf)
		NewRect(r.style.Thickness, r.h, Vec{r.Pos.X, r.Pos.Y},
			WithStyle(Style{r.style.Colour, 0, 0}),
		).Draw(buf)
		NewRect(r.style.Thickness, r.h, Vec{r.Pos.X + float64(r.w) - float64(r.style.Thickness), r.Pos.Y},
			WithStyle(Style{r.style.Colour, 0, 0}),
		).Draw(buf)

		if r.style.Bloom > 0 {
			r.drawBloom(buf)
		}
	}
}

// IsWithin returns whether a position lies within the rectangle's perimeter.
func (r *Rect) IsWithin(pos Vec) bool {
	return (pos.X >= r.Pos.X) && (pos.X <= r.Pos.X+r.Width()) &&
		(pos.Y >= r.Pos.Y) && (pos.Y <= r.Pos.Y+r.Height())
}

// drawBloom draws a bloom effect around the shape.
func (r *Rect) drawBloom(buf *FrameBuffer) {
	// Draw borders around the rectangle of increasing size and decreasing intensity.
	for dist := 1; dist <= r.style.Bloom; dist++ {
		x, y := int(math.Round(r.Pos.X)), int(math.Round(r.Pos.Y))
		topLeftX := x - dist
		topLeftY := y - dist
		topRightX := x + int(math.Round(r.w)) + dist
		bottomLeftY := y + int(math.Round(r.h)) + dist

		// Calculate colour from distance away from shape body
		brightness := 1 - (float64(dist) / float64(r.style.Bloom))
		r, g, b, a := RGBA8(r.style.Colour)
		bloomColour := color.RGBA{r, g, b, uint8(brightness * float64(a))}

		// Draw top and bottom bloom
		for i := topLeftX + 1; i < topRightX; i++ {
			buf.SetPixel(topLeftY, i, NewPixel(bloomColour))
			buf.SetPixel(bottomLeftY, i, NewPixel(bloomColour))
		}
		// Draw left and right bloom
		for i := topLeftY; i <= bottomLeftY; i++ {
			buf.SetPixel(i, topLeftX, NewPixel(bloomColour))
			buf.SetPixel(i, topRightX, NewPixel(bloomColour))
		}
	}
}

// Cicle is a circle shape, aligned to the centre of the circle.
type Circle struct{ *shape }

// NewCircle constructs a new circle.
func NewCircle(diameter float64, pos Vec, opts ...func(*shape)) *Circle {
	return &Circle{newShape(diameter, diameter, pos, opts...)}
}

// IsWithin returns whether a position lies within the circle's perimeter.
func (c *Circle) IsWithin(pos Vec) bool {
	return Dist(c.Pos, pos) <= c.Width()/2
}

// Draw draws the circle onto the provided frame buffer.
func (c *Circle) Draw(buf *FrameBuffer) {
	thickness := c.style.Thickness
	if c.style.Thickness == 0 { // for filled shape
		thickness = c.w / 2
	}

	// Construct bounding box
	radius := c.w / 2
	bbBoxPos := Vec{c.Pos.X - radius, c.Pos.Y - radius}
	bbox := NewRect(c.w, c.h, bbBoxPos)

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
	radius := c.w / 2
	bbBoxPos := Vec{c.Pos.X - radius - bloom, c.Pos.Y - radius - bloom}
	bbox := NewRect(c.w+bloom+bloom, c.h+bloom+bloom, bbBoxPos)

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

// CurvedRect is a rectangle with rounded corners, aligned to the top-left.
type CurvedRect struct {
	*shape
	radius float64
}

// NewCurvedRect constructs a new curved rectangle.
func NewCurvedRect(width, height, radius float64, pos Vec, opts ...func(*shape)) *CurvedRect {
	return &CurvedRect{newShape(width, height, pos, opts...), radius}
}

// IsWithin returns whether a position lies within the curved rectangle's perimeter.
func (r *CurvedRect) IsWithin(pos Vec) bool {
	// Note: this doesn't account for the rounded corners
	return (pos.X >= r.Pos.X) && (pos.X <= r.Pos.X+r.Width()) &&
		(pos.Y >= r.Pos.Y) && (pos.Y <= r.Pos.Y+r.Height())
}

// Draw draws the curved rectangle onto the provided frame buffer.
func (r *CurvedRect) Draw(buf *FrameBuffer) {
	subRectHeight := r.style.Thickness
	subRectWidth := r.style.Thickness
	if r.style.Thickness == 0 {
		// For filled shape
		subRectWidth = r.w / 2
		subRectHeight = r.h / 2
	}

	// Draw each edge as its own rectangle
	NewRect(
		r.w-2*(r.radius), subRectHeight, Vec{r.Pos.X + r.radius, r.Pos.Y},
		WithStyle(Style{r.style.Colour, 0, 0}),
	).Draw(buf)
	NewRect(
		r.w-2*(r.radius), subRectHeight, Vec{r.Pos.X + r.radius, r.Pos.Y + r.h - subRectHeight},
		WithStyle(Style{r.style.Colour, 0, 0}),
	).Draw(buf)
	NewRect(
		subRectWidth, r.h-2*r.radius, Vec{r.Pos.X, r.Pos.Y + r.radius},
		WithStyle(Style{r.style.Colour, 0, 0}),
	).Draw(buf)
	NewRect(
		subRectWidth, r.h-2*r.radius, Vec{r.Pos.X + r.w - subRectWidth, r.Pos.Y + r.radius},
		WithStyle(Style{r.style.Colour, 0, 0}),
	).Draw(buf)

	// Draw rounded corners
	drawCorner := func(bbox *Rect, origin Vec) {
		// Iterate over every pixel in the bounding box
		for x := bbox.Pos.X; x <= bbox.Pos.X+bbox.w; x++ {
			for y := bbox.Pos.Y; y <= bbox.Pos.Y+bbox.h; y++ {
				dist := Dist(origin, Vec{x, y})

				withinCircle := func() bool {
					if r.style.Thickness == 0 {
						return dist <= r.radius
					} else {
						return dist <= r.radius && dist > r.radius-r.style.Thickness
					}
				}()

				if withinCircle {
					buf.SetPixel(int(math.Round(y)), int(math.Round(x)), NewPixel(r.style.Colour))
				}
			}
		}
	}

	// Top left
	bboxSize := r.radius
	bbox := NewRect(bboxSize, bboxSize, Vec{r.Pos.X, r.Pos.Y})
	drawCorner(bbox, Vec{bbox.Pos.X + bbox.w, bbox.Pos.Y + bbox.h})

	// Top right
	bbox = NewRect(bboxSize, bboxSize, Vec{r.Pos.X + r.w - r.radius, r.Pos.Y})
	drawCorner(bbox, Vec{bbox.Pos.X, bbox.Pos.Y + bbox.h})

	// Bottom left
	bbox = NewRect(bboxSize, bboxSize, Vec{r.Pos.X, r.Pos.Y + r.h - r.radius})
	drawCorner(bbox, Vec{bbox.Pos.X + bbox.w, bbox.Pos.Y})

	// Bottom right
	bbox = NewRect(bboxSize, bboxSize, Vec{r.Pos.X + r.w - r.radius, r.Pos.Y + r.h - r.radius})
	drawCorner(bbox, Vec{bbox.Pos.X, bbox.Pos.Y})

	if r.style.Bloom > 0 {
		r.drawBloom(buf)
	}
}

// drawBloom draws a bloom effect around the shape.
func (r *CurvedRect) drawBloom(buf *FrameBuffer) {
	bloom := float64(r.style.Bloom)

	// Draw straight edge bloom
	for dist := 1; dist <= r.style.Bloom; dist++ {
		// Calculate colour from distance away from shape body
		brightness := 1 - (float64(dist) / bloom)
		R, G, B, A := RGBA8(r.style.Colour)
		bloomColour := color.RGBA{R, G, B, uint8(brightness * float64(A))}

		// Top and bottom bloom
		for x := r.Pos.X + r.radius + 1; x < r.Pos.X+r.w-r.radius; x++ {
			buf.SetPixel(int(r.Pos.Y)-dist, int(x), NewPixel(bloomColour))
			buf.SetPixel(int(r.Pos.Y+r.h)+dist, int(x), NewPixel(bloomColour))
		}

		// Left and right
		for y := r.Pos.Y + r.radius + 1; y < r.Pos.Y+r.h-r.radius; y++ {
			buf.SetPixel(int(y), int(r.Pos.X)-dist, NewPixel(bloomColour))
			buf.SetPixel(int(y), int(r.Pos.X+r.w)+dist, NewPixel(bloomColour))
		}
	}

	// Draw rounded corner bloom
	drawCorner := func(bbox *Rect, origin Vec) {
		// Iterate over every pixel in the bounding box
		for x := bbox.Pos.X; x <= bbox.Pos.X+bbox.w; x++ {
			for y := bbox.Pos.Y; y <= bbox.Pos.Y+bbox.h; y++ {
				dist := Dist(origin, Vec{x, y})
				withinCircle := dist > r.radius && dist <= r.radius+bloom
				if withinCircle {
					// Calculate colour from distance away from shape body
					brightness := 1 - (dist-r.radius)/(bloom)
					r, g, b, a := RGBA8(r.style.Colour)

					bloomColour := color.RGBA{r, g, b, uint8(brightness * float64(a))}
					buf.SetPixel(int(math.Round(y)), int(math.Round(x)), NewPixel(bloomColour))
				}
			}
		}
	}

	// Top left
	bboxSize := bloom + r.radius
	bbox := NewRect(bboxSize, bboxSize, Vec{r.Pos.X - bloom, r.Pos.Y - bloom})
	drawCorner(bbox, Vec{bbox.Pos.X + bbox.w, bbox.Pos.Y + bbox.h})

	// Top right
	bbox = NewRect(bboxSize, bboxSize, Vec{r.Pos.X + r.w - r.radius, r.Pos.Y - bloom})
	drawCorner(bbox, Vec{bbox.Pos.X, bbox.Pos.Y + bbox.h})

	// Bottom left
	bbox = NewRect(bboxSize, bboxSize, Vec{r.Pos.X - bloom, r.Pos.Y + r.h - r.radius})
	drawCorner(bbox, Vec{bbox.Pos.X + bbox.w, bbox.Pos.Y})

	// Bottom right
	bbox = NewRect(bboxSize, bboxSize, Vec{r.Pos.X + r.w - r.radius, r.Pos.Y + r.h - r.radius})
	drawCorner(bbox, Vec{bbox.Pos.X, bbox.Pos.Y})
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
	case *CurvedRect:
		panic("todo")
	default:
		t := reflect.TypeOf(s1).String()
		panic("collision detection is unsupported for type: %s" + t)
	}
}
