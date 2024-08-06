package main

import (
	"fmt"
	"image/color"
	"time"

	tgl "github.com/zac460/turdgl"
)

func main() {
	win, err := tgl.NewWindow(tgl.WindowCfg{
		Title:  "Basic Shapes Example",
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
	polygon := tgl.NewPolygon([]tgl.Vec{
		{X: 560, Y: 120},
		{X: 450, Y: 340},
		{X: 250, Y: 220},
		{X: 100, Y: 420},
		{X: 400, Y: 580},
		{X: 560, Y: 470},
		{X: 800, Y: 600},
		{X: 830, Y: 240},
		{X: 680, Y: 250},
	})
	txt := tgl.NewText("Hello there", tgl.Vec{X: 100, Y: 600}).
		SetColour(color.RGBA{255, 255, 255, 255})

	// Keybinds
	win.RegisterKeybind(tgl.KeyEscape, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyLCtrl, func() { win.Quit() })

	for win.IsRunning() {
		win.Framebuffer.SetBackground(color.RGBA{39, 45, 53, 255})

		// Draw shapes
		win.Draw(polygon)
		win.Draw(txt)

		win.SetTitle(fmt.Sprint(win.MouseLocation()))

		win.Update()

		// Count FPS
		frames++
		select {
		case <-second:
			txt.SetText(fmt.Sprintf("FPS: %d", frames))
			frames = 0
		default:
		}
	}
}
