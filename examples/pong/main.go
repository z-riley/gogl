package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/z-riley/turdgl"
	"golang.org/x/exp/constraints"
)

func main() {
	win, err := turdgl.NewWindow(turdgl.WindowCfg{
		Title:  "Pong Example",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// For measuring FPS
	frames := 0
	second := time.Tick(time.Second)

	// Shapes
	paddleLeft := NewPaddle(turdgl.Vec{X: 50, Y: 200})
	paddleRight := NewPaddle(turdgl.Vec{X: float64(win.GetConfig().Width) - 50, Y: 200})
	ball := NewBall(turdgl.Vec{
		X: float64(win.GetConfig().Width / 2),
		Y: float64(win.GetConfig().Height / 2),
	})
	scores := turdgl.NewText("0 | 0", turdgl.Vec{X: 470, Y: 20}, "../../fonts/arial.ttf").
		SetSize(34).SetColour(color.White)
	setScore := func(left, right int) {
		scores.SetText(fmt.Sprintf("%d | %d", left, right))
	}

	// Keybinds
	win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(turdgl.KeyLCtrl, turdgl.KeyPress, func() { win.Quit() })

	// Game state
	leftScore := 0
	rightScore := 0

	prevTime := time.Now()
	for win.IsRunning() {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// React to pressed keys
		if win.KeyIsPressed(turdgl.KeyW) {
			paddleLeft.MovePos(dirUp, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(turdgl.KeyS) {
			paddleLeft.MovePos(dirDown, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(turdgl.KeyUp) {
			paddleRight.MovePos(dirUp, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(turdgl.KeyDown) {
			paddleRight.MovePos(dirDown, dt, win.Framebuffer)
		}

		// Ball movement
		event := ball.Update(dt, win.Framebuffer)
		switch event {
		case noWin:
		case leftWin:
			leftScore++
			setScore(leftScore, rightScore)
		case rightWin:
			rightScore++
			setScore(leftScore, rightScore)
		}
		if turdgl.IsColliding(ball.body, paddleLeft.body) ||
			turdgl.IsColliding(ball.body, paddleRight.body) {
			ball.velocity.X *= -1
		}

		// Set background colour
		win.SetBackground(color.RGBA{39, 45, 53, 255})

		// Modify frame buffer
		win.Draw(scores)
		win.Draw(paddleLeft)
		win.Draw(paddleRight)
		win.Draw(ball)

		// Render screen
		win.Update()

		// Count FPS
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", win.GetConfig().Title, frames))
			frames = 0
		default:
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
