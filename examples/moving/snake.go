package main

import (
	"image/color"
	"math"
	"time"

	tgl "github.com/zac460/turdgl"
	"golang.org/x/exp/constraints"
)

const (
	maxNodeDistPx   = 80
	numSegments     = 10
	headSize        = 30
	bodyScaleFactor = 0.97
)

type snake struct {
	head     *tgl.Circle
	body     []*tgl.Circle
	velocity *tgl.Vec // velocity in px/s
}

// Snake constructs a new snake based on the given head position.
func NewSnake(headPos tgl.Vec) *snake {
	headStyle := tgl.Style{
		Colour:    color.RGBA{255, 255, 255, 255},
		Thickness: 0,
	}

	bodyStyle := tgl.Style{
		Colour:    color.RGBA{255, 255, 255, 255},
		Thickness: 4,
	}

	// Construct body in with segments stretched out...
	var b []*tgl.Circle
	for i := 1; i < numSegments; i++ {
		segmentDiameter := headSize * math.Pow(bodyScaleFactor, float64(i))
		segment := tgl.NewCircle(
			segmentDiameter, segmentDiameter,
			tgl.Vec{X: headPos.X, Y: headPos.Y + headSize*float64(i) + 1},
			bodyStyle,
		)
		b = append(b, segment)
	}
	h := tgl.NewCircle(headSize, headSize, headPos, headStyle)
	h.Direction = tgl.Vec{X: 0, Y: -1} // faces upwards

	s := snake{
		head: h,
		body: b,
	}
	// ...then align the body segments
	s.updateBodyPos()

	return &s
}

// Draw draws the snake on the provided frame buffer.
func (s *snake) Draw(buf *tgl.FrameBuffer) {
	const markerSize = 4
	markerStyle := tgl.Style{Colour: color.RGBA{255, 0, 0, 0}, Thickness: 0}

	// Draw head segment
	s.head.Draw(buf)
	// Draw head marker
	lMarkerHead := s.head.EdgePoint(math.Pi / 2)
	lMarker := tgl.NewCircle(markerSize, markerSize, lMarkerHead, markerStyle)
	lMarker.Draw(buf)
	rMarkerHead := s.head.EdgePoint(math.Pi / 2 * 3)
	rMarker := tgl.NewCircle(markerSize, markerSize, rMarkerHead, markerStyle)
	rMarker.Draw(buf)

	// Draw body
	markers := make([]tgl.Vec, 2+2*len(s.body))
	for i, c := range s.body {
		// Draw segment
		c.Draw(buf)

		// Draw markers
		lPos := c.EdgePoint(math.Pi / 2)
		lMarker := tgl.NewCircle(markerSize, markerSize, lPos, markerStyle)
		lMarker.Draw(buf)

		rPos := c.EdgePoint(math.Pi / 2 * 3)
		rMarker := tgl.NewCircle(markerSize, markerSize, rPos, markerStyle)
		rMarker.Draw(buf)

		markers[i+1] = lPos
		markers[2*(len(s.body))-i] = rPos
	}

	markers[0] = lMarkerHead
	markers[len(markers)-1] = rMarkerHead

	splinePoints := tgl.GenerateCatmullRomSpline(markers, 5)
	for _, point := range splinePoints {
		pointStyle := tgl.Style{Colour: color.RGBA{0, 255, 0, 255}, Thickness: 0}
		const pointSize = 3
		point := tgl.NewCircle(pointSize, pointSize, point, pointStyle)
		point.Draw(buf)
	}
}

// Update recalculates the snake's position based on the current velocity and time interval.
// A reference to the frame buffer must be provided to check snake isn't out of bounds.
func (s *snake) Update(dt time.Duration, buf *tgl.FrameBuffer) {
	// Update the head
	newX := s.head.Pos.X + s.velocity.X*dt.Seconds()
	newY := s.head.Pos.Y + s.velocity.Y*dt.Seconds()
	const segmentRad float64 = headSize / 2
	newX = Constrain(newX, segmentRad, float64(buf.Width())-segmentRad-1)
	newY = Constrain(newY, segmentRad, float64(buf.Height())-segmentRad-1)
	s.head.Pos = tgl.Vec{X: newX, Y: newY}
	s.head.Direction = tgl.Normalise(*s.velocity)

	// Update the body
	s.updateBodyPos()
}

func (s *snake) updateBodyPos() {
	for i, node := range s.body {
		var nodeAhead *tgl.Circle
		if i == 0 {
			nodeAhead = s.head
		} else {
			nodeAhead = s.body[i-1]
		}
		// If node is too far away from the node ahead of it...
		if tgl.Dist(node.Pos, nodeAhead.Pos) > nodeAhead.Width() {
			// Move the node to be adjacent to the node ahead
			diff := tgl.Sub(nodeAhead.Pos, node.Pos)
			node.Move(tgl.Sub(diff, diff.SetMag(nodeAhead.Width())))
			node.Direction = tgl.Normalise(tgl.Sub(nodeAhead.Pos, node.Pos))
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
