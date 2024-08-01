package main

import (
	"fmt"
	"image/color"
	"time"

	tgl "github.com/zac460/turdgl"
	"golang.org/x/exp/constraints"
)

func main() {
	win, err := tgl.NewWindow(tgl.WindowCfg{
		Title:  "Pong Example",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}

	// For measuring FPS
	frames := 0
	second := time.Tick(time.Second)

	// Shapes
	paddleLeft := NewPaddle(tgl.Vec{X: 50, Y: 200})
	paddleRight := NewPaddle(tgl.Vec{X: float64(win.GetConfig().Width) - 50, Y: 200})
	ball := NewBall(tgl.Vec{
		X: float64(win.GetConfig().Width / 2),
		Y: float64(win.GetConfig().Height / 2),
	})

	// Keybinds
	win.RegisterKeybind(tgl.KeyEscape, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyLCtrl, func() { win.Quit() })

	prevTime := time.Now()
	for win.IsRunning() {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// React to pressed keys
		if win.KeyIsPressed(tgl.KeyW) {
			paddleLeft.MovePos(dirUp, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(tgl.KeyS) {
			paddleLeft.MovePos(dirDown, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(tgl.KeyUp) {
			paddleRight.MovePos(dirUp, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(tgl.KeyDown) {
			paddleRight.MovePos(dirDown, dt, win.Framebuffer)
		}

		// Ball movement
		ball.Update(dt, win.Framebuffer)
		if tgl.IsColliding(ball.body, paddleLeft.body) ||
			tgl.IsColliding(ball.body, paddleRight.body) {
			ball.velocity.X *= -1
		}

		// Set background colour
		win.Framebuffer.SetBackground(color.RGBA{39, 45, 53, 255})

		// Modify frame buffer
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
