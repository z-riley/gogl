package main

import (
	"image/color"
	"time"

	tgl "github.com/zac460/turdgl"
	"golang.org/x/exp/constraints"
)

const (
	maxNodeDistPx = 80
	numSegments   = 30
	segmentSize   = 20
)

type snake struct {
	head     *tgl.Circle
	body     []*tgl.Circle
	velocity *tgl.Vec // velocity in px/s
}

func NewSnake(headPos tgl.Vec) *snake {
	headStyle := tgl.Style{
		Colour:    color.RGBA{255, 255, 255, 255},
		Thickness: 0,
	}

	bodyStyle := tgl.Style{
		Colour:    color.RGBA{255, 255, 255, 255},
		Thickness: 4,
	}

	var b []*tgl.Circle
	for i := 0; i < numSegments-1; i++ {
		b = append(b, tgl.NewCircle(
			segmentSize, segmentSize,
			tgl.Vec{X: headPos.X, Y: headPos.Y + segmentSize*float64(i+1)}, bodyStyle))
	}

	return &snake{
		head: tgl.NewCircle(segmentSize, segmentSize, headPos, headStyle),
		body: b,
	}
}

func (s *snake) Draw(buf *tgl.FrameBuffer) {
	s.updateBodyPos()

	s.head.Draw(buf)
	for _, c := range s.body {
		c.Draw(buf)
	}
}

// Move moves the head of the snake. A reference to the frame buffer is needed
// so out of bounds pixels aren't referenced.
func (s *snake) Move(mov tgl.Vec) {
	s.head.Move(mov)
}

// Update recalculates the snake's position based on the current velocity and time interval.
// A reference to the frame buffer must be provided to check snake isn't out of bounds.
func (s *snake) Update(dt time.Duration, buf *tgl.FrameBuffer) {
	// Update the head
	newX := s.head.Pos.X + s.velocity.X*dt.Seconds()
	newY := s.head.Pos.Y + s.velocity.Y*dt.Seconds()
	const segmentRad float64 = segmentSize / 2
	newX = Constrain(newX, segmentRad, float64(buf.Width())-segmentRad-1)
	newY = Constrain(newY, segmentRad, float64(buf.Height())-segmentRad-1)
	s.head.Pos = tgl.Vec{X: newX, Y: newY}

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
		if tgl.Dist(node.Pos, nodeAhead.Pos) > segmentSize {
			// Move the node to be adjacent to the node ahead
			diff := tgl.Sub(nodeAhead.Pos, node.Pos)
			node.Move(tgl.Sub(diff, diff.SetMag(segmentSize)))
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
