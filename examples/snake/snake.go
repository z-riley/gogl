package main

import (
	"image/color"
	"math"
	"time"

	"github.com/z-riley/turdgl"
	"golang.org/x/exp/constraints"
)

const (
	maxNodeDistPx   = 80
	numSegments     = 40
	headSize        = 30
	bodyScaleFactor = 0.97
)

type snake struct {
	head     *turdgl.Circle
	body     []*turdgl.Circle
	velocity *turdgl.Vec // velocity in px/s
}

// Snake constructs a new snake based on the given head position.
func NewSnake(headPos turdgl.Vec) *snake {
	// Construct head
	h := turdgl.NewCircle(
		headSize,
		headPos,
	).SetStyle(turdgl.Style{Colour: color.RGBA{255, 255, 255, 255}, Thickness: 0})
	h.Direction = turdgl.Vec{X: 0, Y: -1} // faces upwards

	// Construct body with segments stretched out...
	var b []*turdgl.Circle
	for i := 1; i < numSegments; i++ {
		segmentDiameter := headSize * math.Pow(bodyScaleFactor, float64(i))
		segment := turdgl.NewCircle(
			segmentDiameter,
			turdgl.Vec{X: headPos.X, Y: headPos.Y + headSize*float64(i) + 1},
		).SetStyle(turdgl.Style{Colour: color.RGBA{255, 255, 255, 255}, Thickness: 4})
		b = append(b, segment)
	}
	s := snake{
		head: h,
		body: b,
	}
	// ...then align the body segments
	s.updateBodyPos()

	return &s
}

// GetPos returns the position of the snake from the centre of the head.
func (s *snake) GetPos() turdgl.Vec {
	return s.head.Pos
}

// Draw draws the snake on the provided frame buffer.
func (s *snake) Draw(buf *turdgl.FrameBuffer) {
	const markerSize = 4
	markerStyle := turdgl.Style{Colour: color.RGBA{255, 0, 0, 0}, Thickness: 0}

	// Draw head segment
	s.head.Draw(buf)
	// Draw head markers
	lMarkerHead := s.head.EdgePoint(math.Pi / 2)
	lMarker := turdgl.NewCircle(markerSize, lMarkerHead).SetStyle(markerStyle)
	lMarker.Draw(buf)
	rMarkerHead := s.head.EdgePoint(math.Pi / 2 * 3)
	rMarker := turdgl.NewCircle(markerSize, rMarkerHead).SetStyle(markerStyle)
	rMarker.Draw(buf)
	fMarkerHead := s.head.EdgePoint(0)
	fMarker := turdgl.NewCircle(markerSize, fMarkerHead).SetStyle(markerStyle)
	fMarker.Draw(buf)

	// Draw body
	markers := make([]turdgl.Vec, 3+2*len(s.body))
	for i, c := range s.body {
		// Draw segment
		c.Draw(buf)

		// Draw body markers
		lPos := c.EdgePoint(math.Pi / 2)
		turdgl.NewCircle(markerSize, lPos).SetStyle(markerStyle).Draw(buf)

		rPos := c.EdgePoint(math.Pi / 2 * 3)
		turdgl.NewCircle(markerSize, rPos).SetStyle(markerStyle).Draw(buf)

		markers[i+1] = lPos
		markers[2*(len(s.body))-i] = rPos
	}

	markers[0] = lMarkerHead
	markers[len(markers)-2] = rMarkerHead
	markers[len(markers)-1] = fMarkerHead

	// Draw body fill
	turdgl.NewPolygon(markers).
		SetStyle(turdgl.Style{Colour: color.RGBA{20, 70, 20, 255}}).
		Draw(buf)

	// Draw body outline
	splinePoints := turdgl.GenerateCatmullRomSpline(markers, 20)
	for _, point := range splinePoints {
		pointStyle := turdgl.Style{Colour: color.RGBA{0, 255, 0, 255}, Thickness: 0}
		const pointSize = 3
		turdgl.NewCircle(pointSize, point).SetStyle(pointStyle).Draw(buf)
	}

}

// Update recalculates the snake's position based on the current velocity and time interval.
// A reference to the frame buffer must be provided to check snake isn't out of bounds.
func (s *snake) Update(dt time.Duration, buf *turdgl.FrameBuffer) {
	// Update the head
	newX := s.head.GetPos().X + s.velocity.X*dt.Seconds()
	newY := s.head.GetPos().Y + s.velocity.Y*dt.Seconds()
	const segmentRad float64 = headSize / 2
	newX = Constrain(newX, segmentRad, float64(buf.Width())-segmentRad-1)
	newY = Constrain(newY, segmentRad, float64(buf.Height())-segmentRad-1)
	s.head.SetPos(turdgl.Vec{X: newX, Y: newY})
	s.head.Direction = turdgl.Normalise(*s.velocity)

	// Update the body
	s.updateBodyPos()
}

func (s *snake) updateBodyPos() {
	for i, node := range s.body {
		var nodeAhead *turdgl.Circle
		if i == 0 {
			nodeAhead = s.head
		} else {
			nodeAhead = s.body[i-1]
		}
		// If node is too far away from the node ahead of it...
		if turdgl.Dist(node.GetPos(), nodeAhead.GetPos()) > nodeAhead.Width() {
			// Move the node to be adjacent to the node ahead
			diff := turdgl.Sub(nodeAhead.GetPos(), node.GetPos())
			node.Move(turdgl.Sub(diff, diff.SetMag(nodeAhead.Width())))
			node.Direction = turdgl.Normalise(turdgl.Sub(nodeAhead.GetPos(), node.GetPos()))
		}
	}
}

// Constrain keeps a number between lower and upper bounds.
func Constrain[T constraints.Ordered](x, lower, upper T) T {
	switch {
	case x < lower:
		return lower
	case x > upper:
		return upper
	default:
		return x
	}
}
