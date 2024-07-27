package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/pixelgl"
	tgl "github.com/zac460/turdgl"
)

const (
	dirUp = iota
	dirDown
	dirLeft
	dirRight
)

var (
	frames = 0
	second = time.Tick(time.Second)
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Title",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Screen state
	framebuf := tgl.NewFrameBuffer(win.Canvas().Texture().Width(), win.Canvas().Texture().Height())
	prevSize := win.Bounds().Size()

	// Shapes
	snake := NewSnake(tgl.Vec{X: 400, Y: 100})

	prevTime := time.Now()
	for {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// Handle user input
		if win.Closed() || win.JustPressed(pixelgl.KeyLeftControl) || win.JustPressed(pixelgl.KeyEscape) {
			return
		}
		const speed = 1000

		if win.Pressed(pixelgl.KeyW) {
			snake.velocity = &tgl.Vec{Y: -speed}
			snake.Update(dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyA) {
			snake.velocity = &tgl.Vec{X: -speed}
			snake.Update(dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyS) {
			snake.velocity = &tgl.Vec{Y: speed}
			snake.Update(dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyD) {
			snake.velocity = &tgl.Vec{X: speed}
			snake.Update(dt, framebuf)
		}

		// Adjust frame buffer size if window size changes
		if !prevSize.Eq(win.Bounds().Size()) {
			framebuf = tgl.NewFrameBuffer(win.Canvas().Texture().Width(), win.Canvas().Texture().Height())
		}

		// Set background colour
		framebuf.SetColour(color.RGBA{39, 45, 53, 255})

		// Modify frame buffer
		snake.Draw(framebuf)

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
