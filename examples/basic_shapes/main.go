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
		Title:     "Title",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	// Screen state
	framebuf := tgl.NewFrameBuffer(win.Canvas().Texture().Width(), win.Canvas().Texture().Height())
	prevSize := win.Bounds().Size()

	// Shapes
	rect := tgl.NewRect(
		120, 90,
		tgl.Vec{X: 200, Y: 200},
		tgl.Style{Colour: color.RGBA{0, 255, 0, 1}, Thickness: 0},
	)
	rect2 := tgl.NewRect(
		120, 90,
		tgl.Vec{X: 200, Y: 200},
		tgl.Style{Colour: color.RGBA{255, 0, 0, 1}, Thickness: 4},
	)
	circle := tgl.NewCircle(
		100, 100,
		tgl.Vec{X: 500, Y: 200},
		tgl.Style{Colour: color.RGBA{255, 0, 0, 1}, Thickness: 10},
	)

	for !win.Closed() && !win.JustPressed(pixelgl.KeyLeftControl) && !win.JustPressed(pixelgl.KeyEscape) {
		framebuf.Clear()

		// Adjust frame buffer size if window size changes
		if !prevSize.Eq(win.Bounds().Size()) {
			framebuf = tgl.NewFrameBuffer(win.Canvas().Texture().Width(), win.Canvas().Texture().Height())
		}

		// Modify frame buffer
		if win.JustPressed(pixelgl.KeyE) {
			rect.Move(tgl.Vec{X: 2, Y: 2})
		}
		rect.Draw(framebuf)
		rect2.Draw(framebuf)
		circle.Draw(framebuf)

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
