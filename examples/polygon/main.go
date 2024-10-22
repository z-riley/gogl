package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/z-riley/turdgl"
)

func main() {
	win, err := turdgl.NewWindow(turdgl.WindowCfg{
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

	polygon := turdgl.NewPolygon([]turdgl.Vec{
		{X: 560, Y: 120},
		{X: 450, Y: 340},
		{X: 250, Y: 220},
		{X: 100, Y: 420},
		{X: 400, Y: 580},
		{X: 560, Y: 470},
		{X: 800, Y: 600},
		{X: 830, Y: 240},
		{X: 680, Y: 250},
	}).SetStyle(turdgl.Style{Colour: color.RGBA{20, 70, 20, 255}})

	txt := turdgl.NewText("Press E to move", turdgl.Vec{X: 100, Y: 600}, "../../fonts/arial.ttf").
		SetColour(color.RGBA{255, 255, 255, 255})

	// Keybinds
	win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(turdgl.KeyLCtrl, turdgl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(turdgl.KeyE, turdgl.Instantaneous, func() { polygon.Move(turdgl.Vec{X: 1, Y: 1}) })

	for win.IsRunning() {
		win.SetBackground(color.Black)

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
