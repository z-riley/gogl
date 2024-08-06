package turdgl

import (
	"container/ring"
	"fmt"
	"math"
	"slices"
)

type vertex struct {
	pos   Vec
	isEar bool
}

// calculateIsEar populates the isEar field of a ring element.
func calculateIsEar(ring *ring.Ring) {
	v := ring.Value.(vertex)
	prev := ring.Prev().Value.(vertex)
	next := ring.Next().Value.(vertex)

	// Return early if vertex isn't convex because it can't be an ear
	if !isConvex(v.pos, next.pos, prev.pos) {
		return
	}

	triangle := NewTriangle(
		ring.Value.(vertex).pos,
		ring.Next().Value.(vertex).pos,
		ring.Prev().Value.(vertex).pos,
	)

	// Check if any other vertices lie within the triangle formed by the
	// vertex and its two neighbours.
	r := ring.Next()
	for i := 0; i < r.Len()-1; i++ {

		current := r.Value.(vertex)
		switch {
		case slices.Contains([]Vec{triangle.v1, triangle.v2, triangle.v3}, current.pos):
			// Ignore the triangle's own vertices
		case triangle.pointInTriangle(current.pos):
			// v is not an ear if other vertices exist in the triangle
		default:
			// vertex is an ear. Set the isEar flag for the vertex
			v := ring.Value.(vertex)
			v.isEar = true
			ring.Value = v
			return
		}

		r = r.Next()
	}
}

// Polygon is a 2D shape with 3 or more sides.
type Polygon struct {
	vertices []Vec
	style    Style
	segments []*Triangle
}

// NewPolygon constructs a polygon according to the supplied vertices.
// The order of the vertices dictates the edges of the polygon. The final
// vertex is always linked to the first.
func NewPolygon(vecs []Vec) *Polygon {

	// Construct ring list of vertices for easier manipulation
	r := ring.New(len(vecs))
	for i := 0; i < r.Len(); i++ {
		// TODO use Do(func()) for this instead
		r.Value = vertex{pos: vecs[i]}
		r = r.Next()
	}

	// Mark which vertices are ears
	// TODO use Do(func()) for this instead
	for i := 0; i < r.Len(); i++ {
		calculateIsEar(r)
		r = r.Next()
	}

	var segments = []*Triangle{} // TODO THIS COULD BE make(.., len(vecs)-2)?
	for r.Len() > 1 {
		if r.Value.(vertex).isEar {
			// Add triangle to polygon segments
			segments = append(segments, NewTriangle(
				r.Value.(vertex).pos,
				r.Prev().Value.(vertex).pos,
				r.Next().Value.(vertex).pos,
			).SetStyle(RandomStyle()))

			// Remove current vertex
			r = r.Unlink(r.Len() - 1)

			// Update states of adjacent vertices
			calculateIsEar(r.Prev())
			calculateIsEar(r.Next())

			// TODO add optimisations:
			// If Vi is an ear that is removed, then the edge configuration at the adjacent vertices
			// Vi−1 and Vi+1 can change. If an adjacent vertex is convex, a quick sketch will convince
			// you that it remains convex.
			// If an adjacent vertex is an ear, it does not necessarily remain an ear after Vi is removed.
			// If the adjacent vertex is reflex, it is possible that it becomes convex and possibly an ear.

		}

		r = r.Next()
	}
	return &Polygon{
		vertices: vecs,
		style:    DefaultStyle,
		segments: segments,
	}
}

func PrintRing(ring *ring.Ring) {
	s := "ring "
	for i := 0; i < ring.Len(); i++ {
		s += fmt.Sprint(ring.Value.(vertex).pos) + " "
		ring = ring.Next()
	}
	fmt.Println(s)
}

// isConvex returns true if the angle of vertices `before`->`p`->`after` is convex.
func isConvex(p, before, after Vec) bool {
	dBefore := Sub(p, before)
	dAfter := Sub(after, p)
	return Cross(dBefore, dAfter) > 0
}

// Style returns a copy of the polygon's style.
func (p *Polygon) Style() Style {
	return p.style
}

// SetStyle sets the style of a polygon.
func (p *Polygon) SetStyle(s Style) *Polygon {
	p.style = s
	return p
}

// Draw draws the polygon onto the provided frame buffer.
func (p *Polygon) Draw(buf *FrameBuffer) {
	for _, segment := range p.segments {
		segment.Draw(buf)
	}

	// Overlay debug geometry
	for i, r := range p.vertices {
		NewCircle(10, r).Draw(buf)
		NewText(fmt.Sprintf("%d %v", i, r), r).Draw(buf)
	}
}

type Triangle struct {
	v1, v2, v3 Vec
	style      Style
}

// NewTriangle constructs a new triangle from the provided vertices.
func NewTriangle(v1, v2, v3 Vec) *Triangle {
	return &Triangle{
		v1: v1, v2: v2, v3: v3,
		style: DefaultStyle,
	}
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

	isClockwise := edgeFunction(t.v1, t.v2, t.v3) > 0

	// Iterate over pixels in bounding box
	for i := bbox.Pos.X; i <= bbox.Pos.X+bbox.w; i++ {
		for j := bbox.Pos.Y; j <= bbox.Pos.Y+bbox.h; j++ {
			p := Vec{float64(i), float64(j)}
			ABP := edgeFunction(t.v1, t.v2, p)
			BCP := edgeFunction(t.v2, t.v3, p)
			CAP := edgeFunction(t.v3, t.v1, p)

			// If the triangle is clockwise, check for all edge functions being >= 0.
			// If it's anticlockwise, check for all edge functions being <= 0.
			if (isClockwise && ABP >= 0 && BCP >= 0 && CAP >= 0) ||
				(!isClockwise && ABP <= 0 && BCP <= 0 && CAP <= 0) {
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

// pointInTriangle returns true if point p exists within the area of the triangle.
// This function uses the barycentric coordinate method.
func (t *Triangle) pointInTriangle(p Vec) bool {
	denominator := ((t.v2.Y-t.v3.Y)*(t.v1.X-t.v3.X) + (t.v3.X-t.v2.X)*(t.v1.Y-t.v3.Y))
	a := ((t.v2.Y-t.v3.Y)*(p.X-t.v3.X) + (t.v3.X-t.v2.X)*(p.Y-t.v3.Y)) / denominator
	b := ((t.v3.Y-t.v1.Y)*(p.X-t.v3.X) + (t.v1.X-t.v3.X)*(p.Y-t.v3.Y)) / denominator
	c := 1 - a - b
	return 0 <= a && a <= 1 && 0 <= b && b <= 1 && 0 <= c && c <= 1
}
