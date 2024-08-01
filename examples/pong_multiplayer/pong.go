package main

import (
	"time"

	tgl "github.com/zac460/turdgl"
	"golang.org/x/exp/constraints"
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

const (
	ballDiameter = 10
	ballSpeed    = 300
)

type ball struct {
	body     *tgl.Circle
	velocity tgl.Vec // velocity in px/s
}

// NewBall constructs a new ball.
func NewBall(pos tgl.Vec) *ball {
	return &ball{
		body:     tgl.NewCircle(10, pos),
		velocity: tgl.Normalise(tgl.Vec{X: 1, Y: 1}).SetMag(ballSpeed),
	}
}

// Draw draws the ball on the provided frame buffer.
func (b *ball) Draw(buf *tgl.FrameBuffer) {
	b.body.Draw(buf)
}

// Update recalculates the balls's position based on the current velocity and time interval.
// A reference to the frame buffer must be provided to check the ball isn't out of bounds.
func (b *ball) Update(dt time.Duration, buf *tgl.FrameBuffer) {
	// Update the position
	pos := b.body.GetPos()
	newX := b.body.GetPos().X + b.velocity.X*dt.Seconds()
	newY := b.body.GetPos().Y + b.velocity.Y*dt.Seconds()

	// Deflect ball if it hits the edge of the screen
	if pos.X < ballDiameter || pos.X > float64(buf.Width())-ballDiameter {
		newX = Constrain(newX, ballDiameter, float64(buf.Width())-ballDiameter)
		b.velocity.X *= -1
	}
	if b.body.GetPos().Y < ballDiameter || pos.Y > float64(buf.Height())-ballDiameter {
		newY = Constrain(newY, ballDiameter, float64(buf.Height())-ballDiameter)
		b.velocity.Y *= -1
	}

	b.body.SetPos(tgl.Vec{X: newX, Y: newY})
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
