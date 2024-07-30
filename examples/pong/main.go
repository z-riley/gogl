package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/pixelgl"
	tgl "github.com/zac460/turdgl"
	"golang.org/x/exp/constraints"
)

var (
	frames = 0
	second = time.Tick(time.Second)
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Pong",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Screen state
	screenWidth := win.Canvas().Texture().Width()
	screenHeight := win.Canvas().Texture().Height()
	framebuf := tgl.NewFrameBuffer(screenWidth, screenHeight)
	prevSize := win.Bounds().Size()

	// Shapes
	paddleLeft := NewPaddle(tgl.Vec{X: 50, Y: 200})
	paddleRight := NewPaddle(tgl.Vec{X: float64(screenWidth) - 50, Y: 200})
	ball := NewBall(tgl.Vec{X: win.Bounds().Center().X, Y: win.Bounds().Center().Y})

	prevTime := time.Now()
	for {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// Handle user input
		if win.Closed() || win.JustPressed(pixelgl.KeyLeftControl) || win.JustPressed(pixelgl.KeyEscape) {
			return
		}
		if win.Pressed(pixelgl.KeyW) {
			paddleLeft.MovePos(dirUp, dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyS) {
			paddleLeft.MovePos(dirDown, dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyUp) {
			paddleRight.MovePos(dirUp, dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyDown) {
			paddleRight.MovePos(dirDown, dt, framebuf)
		}
		ball.Update(dt, framebuf)
		if tgl.IsColliding(ball.body, paddleLeft.body) ||
			tgl.IsColliding(ball.body, paddleRight.body) {
			ball.velocity.X *= -1
		}

		// Adjust frame buffer size if window size changes
		if !prevSize.Eq(win.Bounds().Size()) {
			framebuf = tgl.NewFrameBuffer(win.Canvas().Texture().Width(), win.Canvas().Texture().Height())
		}

		// Set background colour
		framebuf.SetBackground(color.RGBA{39, 45, 53, 255})

		// Modify frame buffer
		paddleLeft.Draw(framebuf)
		paddleRight.Draw(framebuf)
		ball.Draw(framebuf)

		// Render screen
		win.Canvas().SetPixels(framebuf.Bytes())
		win.Update()

		// Count FPS
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
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
