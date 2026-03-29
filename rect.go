package gogl

import (
	"image/color"
	"math"
)

// Rect is a rectangle shape, aligned to the top-left corner.
type Rect struct {
	Pos       Vec
	Direction Vec
	w, h      float64
	style     Style
}

var _ Shape = (*Rect)(nil)
var _ hoverable = (*Rect)(nil)

// NewRect constructs a new rectangle shape.
func NewRect(width, height float64, pos Vec) *Rect {
	return &Rect{
		Pos:       pos,
		Direction: Vec{},
		w:         width,
		h:         height,
		style:     DefaultStyle,
	}
}

// Draw draws the rectangle onto the provided frame buffer.
func (e *Rect) Draw(buf *FrameBuffer) {
	if e.style.Thickness == 0 {
		for x := 0; x <= int(math.Round(e.w)); x++ {
			for y := 0; y <= int(math.Round(e.h)); y++ {
				xInt, yInt := int(math.Round(e.Pos.X)), int(math.Round(e.Pos.Y))
				buf.SetPixel(xInt+x, yInt+y, NewPixel(e.style.Colour))
			}
		}
		if e.style.Bloom > 0 {
			e.drawBloom(buf)
		}
	} else {
		// Draw each edge as its own rectangle
		NewRect(e.w, e.style.Thickness, e.Pos).
			SetStyle(Style{e.style.Colour, 0, 0}).
			Draw(buf)
		NewRect(e.w, e.style.Thickness, Vec{e.Pos.X, e.Pos.Y + float64(e.h) - float64(e.style.Thickness)}).
			SetStyle(Style{e.style.Colour, 0, 0}).
			Draw(buf)
		NewRect(e.style.Thickness, e.h, e.Pos).
			SetStyle(Style{e.style.Colour, 0, 0}).
			Draw(buf)
		NewRect(e.style.Thickness, e.h, Vec{e.Pos.X + float64(e.w) - float64(e.style.Thickness), e.Pos.Y}).
			SetStyle(Style{e.style.Colour, 0, 0}).
			Draw(buf)

		if e.style.Bloom > 0 {
			e.drawBloom(buf)
		}
	}
}

// Width returns the pixel width of the rectangle.
func (e *Rect) Width() float64 {
	return e.w
}

// SetWidth sets the width of the rectangle.
func (r *Rect) SetWidth(px float64) *Rect {
	r.w = max(px, 0)
	return r
}

// Height returns the pixel height of the rectangle.
func (e *Rect) Height() float64 {
	return e.h
}

// SetHeight sets the height of the rectangle.
func (r *Rect) SetHeight(px float64) *Rect {
	r.h = max(px, 0)
	return r
}

// GetPos returns the position of the rectangle.
func (e *Rect) GetPos() Vec {
	return e.Pos
}

// SetPos sets the position of the rectangle.
func (e *Rect) SetPos(pos Vec) {
	e.Pos = pos
}

// GetStyle returns's the rectangle's style.
func (e *Rect) GetStyle() Style {
	return e.style
}

// SetStyle sets the style of the rectangle.
func (e *Rect) SetStyle(style Style) *Rect {
	e.style = style
	return e
}

// Move moves the rectangle by the given vector.
func (e *Rect) Move(px Vec) {
	e.Pos = Add(e.Pos, px)
}

// String returns the type of shape as a string.
func (e *Rect) String() string {
	return "rectangle"
}

// IsWithin returns whether a position lies within the rectangle's perimeter.
func (e *Rect) IsWithin(pos Vec) bool {
	return (pos.X >= e.Pos.X) && (pos.X <= e.Pos.X+e.Width()) &&
		(pos.Y >= e.Pos.Y) && (pos.Y <= e.Pos.Y+e.Height())
}

// drawBloom draws a bloom effect around the shape.
func (e *Rect) drawBloom(buf *FrameBuffer) {
	// Draw borders around the rectangle of increasing size and decreasing intensity.
	for dist := 1; dist <= e.style.Bloom; dist++ {
		x, y := int(math.Round(e.Pos.X)), int(math.Round(e.Pos.Y))
		topLeftX := x - dist
		topLeftY := y - dist
		topRightX := x + int(math.Round(e.w)) + dist
		bottomLeftY := y + int(math.Round(e.h)) + dist

		// Calculate colour from distance away from shape body
		brightness := 1 - (float64(dist) / float64(e.style.Bloom))
		r, g, b, a := RGBA8(e.style.Colour)
		bloomColour := color.RGBA{r, g, b, uint8(brightness * float64(a))}

		// Draw top and bottom bloom
		for i := topLeftX + 1; i < topRightX; i++ {
			buf.SetPixel(i, topLeftY, NewPixel(bloomColour))
			buf.SetPixel(i, bottomLeftY, NewPixel(bloomColour))
		}
		// Draw left and right bloom
		for i := topLeftY; i <= bottomLeftY; i++ {
			buf.SetPixel(topLeftX, i, NewPixel(bloomColour))
			buf.SetPixel(topRightX, i, NewPixel(bloomColour))
		}
	}
}

// CurvedRect is a rectangle with rounded corners, aligned to the top-left.
type CurvedRect struct {
	Pos       Vec
	Direction Vec
	w, h      float64
	style     Style
	radius    float64
}

var _ Shape = (*CurvedRect)(nil)
var _ hoverable = (*Rect)(nil)

// NewCurvedRect constructs a new curved rectangle.
func NewCurvedRect(width, height, radius float64, pos Vec) *CurvedRect {
	return &CurvedRect{
		Pos:       pos,
		Direction: Vec{},
		w:         width,
		h:         height,
		style:     DefaultStyle,
		radius:    radius,
	}
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
	NewRect(r.w-2*(r.radius), subRectHeight, Vec{r.Pos.X + r.radius, r.Pos.Y}).
		SetStyle(Style{r.style.Colour, 0, 0}).
		Draw(buf)
	NewRect(r.w-2*(r.radius), subRectHeight, Vec{r.Pos.X + r.radius, r.Pos.Y + r.h - subRectHeight}).
		SetStyle(Style{r.style.Colour, 0, 0}).
		Draw(buf)
	NewRect(subRectWidth, r.h-2*r.radius, Vec{r.Pos.X, r.Pos.Y + r.radius}).
		SetStyle(Style{r.style.Colour, 0, 0}).
		Draw(buf)
	NewRect(subRectWidth, r.h-2*r.radius, Vec{r.Pos.X + r.w - subRectWidth, r.Pos.Y + r.radius}).
		SetStyle(Style{r.style.Colour, 0, 0}).
		Draw(buf)

	// Draw rounded corners
	drawCorner := func(bbox *Rect, origin Vec) {
		// Iterate over every pixel in the bounding box
		for x := bbox.Pos.X; x <= bbox.Pos.X+bbox.w; x++ {
			for y := bbox.Pos.Y; y <= bbox.Pos.Y+bbox.h; y++ {
				dist := Dist(origin, Vec{x, y})

				withinCircle := func() bool {
					if r.style.Thickness == 0 {
						return dist <= r.radius
					}
					return dist <= r.radius && dist > r.radius-r.style.Thickness
				}()

				if withinCircle {
					buf.SetPixel(int(math.Round(x)), int(math.Round(y)), NewPixel(r.style.Colour))
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

// Width returns the pixel width of the curved rectangle.
func (r *CurvedRect) Width() float64 {
	return r.w
}

// SetWidth sets the width of the curved rectangle.
func (r *CurvedRect) SetWidth(px float64) *CurvedRect {
	r.w = max(px, 0)
	return r
}

// Height returns the pixel height of the curved rectangle.
func (r *CurvedRect) Height() float64 {
	return r.h
}

// SetHeight sets the height of the curved rectangle.
func (r *CurvedRect) SetHeight(px float64) *CurvedRect {
	r.h = max(px, 0)
	return r
}

// GetPos returns the position of the curved rectangle.
func (r *CurvedRect) GetPos() Vec {
	return r.Pos
}

// SetPos sets the position of the curved rectangle.
func (r *CurvedRect) SetPos(pos Vec) {
	r.Pos = pos
}

// GetStyle returns's the curved rectangle's style.
func (r *CurvedRect) GetStyle() Style {
	return r.style
}

// SetStyle sets the style of the curved rectangle.
func (r *CurvedRect) SetStyle(style Style) *CurvedRect {
	r.style = style
	return r
}

// Move moves the curved rectangle by the given vector.
func (r *CurvedRect) Move(px Vec) {
	r.Pos = Add(r.Pos, px)
}

// String returns the type of shape as a string.
func (r *CurvedRect) String() string {
	return "curved rectangle"
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
			buf.SetPixel(int(x), int(r.Pos.Y)-dist, NewPixel(bloomColour))
			buf.SetPixel(int(x), int(r.Pos.Y+r.h)+dist, NewPixel(bloomColour))
		}

		// Left and right
		for y := r.Pos.Y + r.radius + 1; y < r.Pos.Y+r.h-r.radius; y++ {
			buf.SetPixel(int(r.Pos.X)-dist, int(y), NewPixel(bloomColour))
			buf.SetPixel(int(r.Pos.X+r.w)+dist, int(y), NewPixel(bloomColour))
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
					buf.SetPixel(int(math.Round(x)), int(math.Round(y)), NewPixel(bloomColour))
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
