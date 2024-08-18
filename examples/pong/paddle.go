package main

import (
	"time"

	tgl "github.com/z-riley/turdgl"
)

const (
	paddleWidth  = 10
	paddleHeight = 60
	paddleSpeed  = 700
)

type direction int

const (
	dirUp direction = iota
	dirDown
)

type paddle struct {
	body     *tgl.Rect
	velocity *tgl.Vec // velocity in px/s
}

// NewPaddle constructs a new paddle.
func NewPaddle(pos tgl.Vec) *paddle {
	return &paddle{
		body:     tgl.NewRect(paddleWidth, paddleHeight, pos),
		velocity: &tgl.Vec{},
	}
}

// Draw draws the paddle on the provided frame buffer.
func (p *paddle) Draw(buf *tgl.FrameBuffer) {
	p.body.Draw(buf)
}

// MovePos recalculates the paddles's position based on the current velocity and time interval.
// A reference to the frame buffer must be provided to check the paddle isn't out of bounds.
func (p *paddle) MovePos(dir direction, dt time.Duration, buf *tgl.FrameBuffer) {
	switch dir {
	case dirUp:
		p.velocity = &tgl.Vec{Y: -paddleSpeed}
	case dirDown:
		p.velocity = &tgl.Vec{Y: paddleSpeed}
	}

	// Update the position
	newX := p.body.GetPos().X + p.velocity.X*dt.Seconds()
	newY := p.body.GetPos().Y + p.velocity.Y*dt.Seconds()
	// Make sure the paddle on the screen
	newY = Constrain(newY, 0, float64(buf.Height())-(paddleHeight)-1)

	p.body.SetPos(tgl.Vec{X: newX, Y: newY})
}
