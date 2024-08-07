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
	}).SetStyle(tgl.Style{Colour: color.RGBA{20, 70, 20, 255}})

	txt := tgl.NewText("Press E to move", tgl.Vec{X: 100, Y: 600}).
		SetColour(color.RGBA{255, 255, 255, 255})

	// Keybinds
	win.RegisterKeybind(tgl.KeyEscape, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyLCtrl, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyE, func() { polygon.Move(tgl.Vec{X: 1, Y: 1}) })

	for win.IsRunning() {
		win.Framebuffer.SetBackground(color.Black)

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
