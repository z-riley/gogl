package main

import (
	"time"

	tgl "github.com/zac460/turdgl"
)

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
