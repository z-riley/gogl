package main

import (
	"time"

	"github.com/z-riley/gogl"
)

const (
	ballDiameter = 10
	ballSpeed    = 500
)

type ball struct {
	body     *gogl.Circle
	velocity gogl.Vec // velocity in px/s
}

// NewBall constructs a new ball.
func NewBall(pos gogl.Vec) *ball {
	return &ball{
		body:     gogl.NewCircle(10, pos),
		velocity: gogl.Normalise(gogl.Vec{X: 1, Y: 1}).SetMag(ballSpeed),
	}
}

// Draw draws the ball on the provided frame buffer.
func (b *ball) Draw(buf *gogl.FrameBuffer) {
	b.body.Draw(buf)
}

type pongEvent int

const (
	noWin pongEvent = iota
	leftWin
	rightWin
)

// Update recalculates the balls's position based on the current velocity and time interval.
// A reference to the frame buffer must be provided to check the ball isn't out of bounds.
func (b *ball) Update(dt time.Duration, buf *gogl.FrameBuffer) pongEvent {
	var event pongEvent

	// Update the position
	pos := b.body.GetPos()
	newX := b.body.GetPos().X + b.velocity.X*dt.Seconds()
	newY := b.body.GetPos().Y + b.velocity.Y*dt.Seconds()

	// Deflect ball if it hits the edge of the screen
	if pos.X < ballDiameter || pos.X > float64(buf.Width())-ballDiameter {
		newX = Constrain(newX, ballDiameter, float64(buf.Width())-ballDiameter)
		b.velocity.X *= -1

		// Generate win/lose events
		if b.body.GetPos().X < float64(buf.Width())/2 {
			event = rightWin
		} else {
			event = leftWin
		}
	}
	if b.body.GetPos().Y < ballDiameter || pos.Y > float64(buf.Height())-ballDiameter {
		newY = Constrain(newY, ballDiameter, float64(buf.Height())-ballDiameter)
		b.velocity.Y *= -1
	}

	b.body.SetPos(gogl.Vec{X: newX, Y: newY})

	return event
}
