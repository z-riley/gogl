package turdgl

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
func (r *Rect) Draw(buf *FrameBuffer) {
	if r.style.Thickness == 0 {
		width := int(math.Round(r.w))
		height := int(math.Round(r.h))
		for i := 0; i <= width; i++ {
			for j := 0; j <= height; j++ {
				xInt, yInt := int(math.Round(r.Pos.X)), int(math.Round(r.Pos.Y))
				buf.SetPixel(xInt+i, yInt+j, NewPixel(r.style.Colour))
			}
		}
		if r.style.Bloom > 0 {
			r.drawBloom(buf)
		}
	} else {
		// Draw each edge as its own rectangle
		NewRect(r.w, r.style.Thickness, r.Pos).
			SetStyle(Style{r.style.Colour, 0, 0}).
			Draw(buf)
		NewRect(r.w, r.style.Thickness, Vec{r.Pos.X, r.Pos.Y + float64(r.h) - float64(r.style.Thickness)}).
			SetStyle(Style{r.style.Colour, 0, 0}).
			Draw(buf)
		NewRect(r.style.Thickness, r.h, r.Pos).
			SetStyle(Style{r.style.Colour, 0, 0}).
			Draw(buf)
		NewRect(r.style.Thickness, r.h, Vec{r.Pos.X + float64(r.w) - float64(r.style.Thickness), r.Pos.Y}).
			SetStyle(Style{r.style.Colour, 0, 0}).
			Draw(buf)

		if r.style.Bloom > 0 {
			r.drawBloom(buf)
		}
	}
}

func (r *Rect) Width() float64 {
	return r.w
}

func (r *Rect) SetWidth(px float64) *Rect {
	r.w = px
	return r
}

func (r *Rect) Height() float64 {
	return r.h
}

func (r *Rect) SetHeight(px float64) *Rect {
	r.h = px
	return r
}

func (r *Rect) GetPos() Vec {
	return r.Pos
}

func (r *Rect) SetPos(pos Vec) {
	r.Pos = pos
}

func (r *Rect) GetStyle() Style {
	return r.style
}

func (r *Rect) SetStyle(style Style) *Rect {
	r.style = style
	return r
}

func (r *Rect) Move(px Vec) {
	r.Pos = Add(r.Pos, px)
}

func (r *Rect) String() string {
	return "rectangle"
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

func (r *CurvedRect) Width() float64 {
	return r.w
}

func (r *CurvedRect) SetWidth(px float64) *CurvedRect {
	r.w = px
	return r
}

func (r *CurvedRect) Height() float64 {
	return r.h
}

func (r *CurvedRect) SetHeight(px float64) *CurvedRect {
	r.h = px
	return r
}

func (r *CurvedRect) GetPos() Vec {
	return r.Pos
}

func (r *CurvedRect) SetPos(pos Vec) {
	r.Pos = pos
}

func (r *CurvedRect) GetStyle() Style {
	return r.style
}

func (r *CurvedRect) SetStyle(style Style) *CurvedRect {
	r.style = style
	return r
}

func (r *CurvedRect) Move(px Vec) {
	r.Pos = Add(r.Pos, px)
}

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
