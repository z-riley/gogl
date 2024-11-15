package main

import (
	"image/color"

	"github.com/z-riley/turdgl"
)

func main() {
	win, err := turdgl.NewWindow(turdgl.WindowCfg{
		Title:  "Turdgl Polygon Example",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

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

	win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyPress, func() { win.Quit() })

	for win.IsRunning() {
		win.SetBackground(color.Black)
		win.Draw(polygon)
		win.Update()
	}
}
