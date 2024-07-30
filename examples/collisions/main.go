package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/pixelgl"
	tgl "github.com/zac460/turdgl"
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
	rect1 := tgl.NewRect(100, 60, tgl.Vec{X: 500, Y: 200})
	rect2 := tgl.NewRect(130, 50, tgl.Vec{X: 500, Y: 300},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{0, 0, 255, 255}}),
	)
	circle1 := tgl.NewCircle(88, tgl.Vec{X: 500, Y: 600})
	circle2 := tgl.NewCircle(130, tgl.Vec{X: 600, Y: 500},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{0, 0, 255, 255}}),
	)

	for {
		// Handle user input
		if win.Closed() || win.JustPressed(pixelgl.KeyLeftControl) || win.JustPressed(pixelgl.KeyEscape) {
			return
		}
		if win.Pressed(pixelgl.KeyW) {
			rect2.Pos.Y--
			circle2.Pos.Y--
		}
		if win.Pressed(pixelgl.KeyS) {
			rect2.Pos.Y++
			circle2.Pos.Y++
		}
		if win.Pressed(pixelgl.KeyA) {
			rect2.Pos.X--
			circle2.Pos.X--
		}
		if win.Pressed(pixelgl.KeyD) {
			rect2.Pos.X++
			circle2.Pos.X++
		}

		// Adjust frame buffer size if window size changes
		if !prevSize.Eq(win.Bounds().Size()) {
			framebuf = tgl.NewFrameBuffer(win.Canvas().Texture().Width(), win.Canvas().Texture().Height())
		}

		// Set background colour
		framebuf.SetBackground(color.Black)

		// Adjust shape colours to react to collisions
		rect1.SetStyle(tgl.DefaultStyle)
		if tgl.IsColliding(rect1, rect2) {
			rect1.SetStyle(tgl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}
		if tgl.IsColliding(circle2, rect1) {
			rect1.SetStyle(tgl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}
		circle1.SetStyle(tgl.DefaultStyle)
		if tgl.IsColliding(circle1, circle2) {
			circle1.SetStyle(tgl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}
		if tgl.IsColliding(rect2, circle1) {
			circle1.SetStyle(tgl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}

		// Modify frame buffer
		rect1.Draw(framebuf)
		rect2.Draw(framebuf)
		circle1.Draw(framebuf)
		circle2.Draw(framebuf)

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
