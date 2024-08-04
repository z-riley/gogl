package turdgl

import (
	"math"
	"slices"
)

type Triangle struct {
	v1, v2, v3 Vec
	style      Style
}

// NewTriangle constructs a new triangle from the provided vertices.
func NewTriangle(v1, v2, v3 Vec) *Triangle {
	// Sort by Y value (required for rasterisation)
	vs := []Vec{v1, v2, v3}
	slices.SortFunc(vs, func(a, b Vec) int {
		switch {
		case a.Y < b.Y:
			return -1
		case a.Y > b.Y:
			return 1
		default:
			return 0
		}
	})
	return &Triangle{v1: vs[2], v2: vs[1], v3: vs[0]}
}

// Style returns a copy of the triangle's style.
func (t *Triangle) Style() Style {
	return t.style
}

// SetStyle sets the style of a triangle.
func (t *Triangle) SetStyle(s Style) *Triangle {
	t.style = s
	return t
}

// Draw rasterises and draws the triangle onto the provided frame buffer.
func (t *Triangle) Draw(buf *FrameBuffer) {
	// Construct bounding box
	maxX := math.Max(math.Max(t.v1.X, t.v2.X), t.v3.X)
	minX := math.Min(math.Min(t.v1.X, t.v2.X), t.v3.X)
	maxY := math.Max(math.Max(t.v1.Y, t.v2.Y), t.v3.Y)
	minY := math.Min(math.Min(t.v1.Y, t.v2.Y), t.v3.Y)
	bboxPos := Vec{minX, minY}
	bbox := NewRect(maxX-minX, maxY-minY, bboxPos)

	// Iterate over pixels in bounding box
	for i := bbox.Pos.X; i <= bbox.Pos.X+bbox.w; i++ {
		for j := bbox.Pos.Y; j <= bbox.Pos.Y+bbox.h; j++ {
			p := Vec{float64(i), float64(j)}
			ABP := edgeFunction(t.v1, t.v2, p)
			BCP := edgeFunction(t.v2, t.v3, p)
			CAP := edgeFunction(t.v3, t.v1, p)
			// If all edge functions are positive, the point is inside the triangle
			if ABP >= 0 && BCP >= 0 && CAP >= 0 {
				jInt, iInt := int(math.Round(j)), int(math.Round(i))
				buf.SetPixel(jInt, iInt, NewPixel(t.style.Colour))
			}
		}
	}
}

// edgeFunction returns double the signed area of a triangle which has
// vertices a, b and c.
func edgeFunction(a, b, c Vec) float64 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}
