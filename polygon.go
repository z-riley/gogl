package turdgl

import (
	"container/ring"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/netgusto/poly2tri-go"
)

// Polygon is a 2D shape with 3 or more sides.
type Polygon struct {
	vertices []Vec
	style    Style
	segments []*Triangle
}

// NewPolygon constructs a polygon from the specified vertices.
// The order of the vertices dictates the edges of the polygon.
// The final vertex is always linked to the first.
func NewPolygon(vecs []Vec) *Polygon {
	return &Polygon{
		vertices: vecs,
		style:    DefaultStyle,
		segments: triangulatePoly2Tri(vecs),
	}
}

// Move modifies the position of the polygon by the given vector.
func (p *Polygon) Move(mov Vec) {
	// Translate every vertex by the move vector
	for i := range p.vertices {
		p.vertices[i] = Add(p.vertices[i], mov)
	}
	// Refresh the triangle segments
	p.segments = triangulatePoly2Tri(p.vertices)
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
		if segment == nil {
			fmt.Println("Segment nil error", time.Now())
		} else {
			segment.SetStyle(p.style)
			segment.Draw(buf)
		}
	}
}

// vertex contains information about a vertex of a polygon.
type vertex struct {
	pos   Vec
	isEar bool
}

// getVertex returns the value of the vertex in the given ring buffer node.
func getVertex(node *ring.Ring) vertex {
	val, ok := node.Value.(vertex)
	if !ok {
		panic("failed to take ring.Value as type vertex")
	}
	return val
}

// triangulateEarClipping triangulates a polygon defined by a slice of vectors
// into a slice of drawable triangles using the ear clipping method.
func triangulateEarClipping(vecs []Vec) []*Triangle {
	// Construct ring list of vertices for easier manipulation
	r := ring.New(len(vecs))
	for i := range r.Len() {
		r.Value = vertex{pos: vecs[i]}
		r = r.Next()
	}

	// Mark which vertices are ears
	for range r.Len() {
		calculateIsEar(r)
		r = r.Next()
	}

	// Remove remove ears one by one, saving the triangle segment each time.
	// Stop when there is only one triangle remaining. There are always 2 less
	// triangles than polygon vertices
	var segments []*Triangle

	safetyCounter := 0
	for r.Len() >= 3 {
		safetyCounter++
		if safetyCounter > 100 {
			fmt.Println("Ear clipping emergency exit triggered due to not enough ears found")
			break
		}

		// FIXME: this branch is sometimes not triggered enough times, meaning the loop would
		// run forever without the emergency exit. This is because the algorithm fails to
		// triangulate complex geometry properly. Run examples/snake/main.go to trigger the bug.
		if getVertex(r).isEar {
			// Save triangle
			segments = append(segments, NewTriangle(
				getVertex(r).pos,
				getVertex(r.Prev()).pos,
				getVertex(r.Next()).pos,
			).SetStyle(RandomStyle()))

			// Remove the ear vertex
			r = r.Unlink(r.Len() - 1)

			// Update states of adjacent vertices now the ear vertex has been removed.
			calculateIsEar(r.Prev())
			calculateIsEar(r.Next())
		}
		r = r.Next()
	}

	return segments
}

// isError keeps track of whether there's a polygon error.
var isError bool

// triangulatePoly2Tri triangulates a polygon defined by a slice of vectors
// into a slice of drawable triangles using Delauney triangulation.
// https://github.com/ByteArena/poly2tri-go
func triangulatePoly2Tri(vecs []Vec) []*Triangle {
	// Convert vertices to poly2tri format
	contour := make([]*poly2tri.Point, len(vecs))
	for i, v := range vecs {
		contour[i] = poly2tri.NewPoint(v.X, v.Y)
	}

	// Note: polygon holes can be added if needed using swctx.AddHole()
	swctx := poly2tri.NewSweepContext(contour, false)

	// Keep going if Triangulate() fails due to intersecting edges
	defer func() {
		msg := recover()
		// This logic exists to not spam the logs. The error will only be
		// printed once each time the polygon enters an error state
		if msg == nil {
			isError = false
		} else if !isError {
			isError = true
			fmt.Println("Warning: Polygon error:", msg)
		}
	}()
	swctx.Triangulate()
	trianglesRaw := swctx.GetTriangles()

	// Convert library format to turdgl triangles
	var triangles []*Triangle
	for _, t := range trianglesRaw {
		a := Vec{t.Points[0].X, t.Points[0].Y}
		b := Vec{t.Points[1].X, t.Points[1].Y}
		c := Vec{t.Points[2].X, t.Points[2].Y}
		triangles = append(triangles, NewTriangle(a, b, c))
	}

	return triangles
}

// calculateIsEar populates the isEar field of a single ring list element.
func calculateIsEar(vert *ring.Ring) {
	current := getVertex(vert)
	prev := getVertex(vert.Prev())
	next := getVertex(vert.Next())

	if !isConvex(current.pos, next.pos, prev.pos) {
		return // ear vertices must be convex
	}

	// Check if any other vertices lie within the triangle formed by the
	// vertex and its two neighbours.
	triangle := NewTriangle(
		getVertex(vert).pos,
		getVertex(vert.Next()).pos,
		getVertex(vert.Prev()).pos,
	)
	v := vert.Next()
	for range v.Len() - 1 {
		current := getVertex(v)
		switch {
		case slices.Contains([]Vec{triangle.v1, triangle.v2, triangle.v3}, current.pos):
			// Ignore the triangle's own vertices
		case triangle.pointInTriangle(current.pos):
			// The vertex is not an ear if any other vertices exist in the triangle
			v := getVertex(vert)
			v.isEar = true
			vert.Value = v
		default:
			// The vertex is an ear; set the isEar flag
			v := getVertex(vert)
			v.isEar = true
			vert.Value = v
			return
		}
		v = v.Next()
	}
}

// isConvex returns true if the angle made by vectors ab and bc is convex,
// assuming that vectors ab and bc are oriented such that when you rotate the
// ab towards bc, the rotation is counterclockwise.
func isConvex(a, b, c Vec) bool {
	return Cross(Sub(a, b), Sub(c, a)) > 0
}

// Triangle is triangle shape, defined by the position of its vertices.
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
				buf.SetPixel(iInt, jInt, NewPixel(t.style.Colour))
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
