package main

import (
	"image/color"

	"github.com/z-riley/turdgl"
)

func main() {
	win, err := turdgl.NewWindow(turdgl.WindowCfg{
		Title:  "Turdgl Collision Example",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// Initialise shapes
	rect1 := turdgl.NewRect(100, 60, turdgl.Vec{X: 500, Y: 200})
	rect2 := turdgl.NewRect(130, 50, turdgl.Vec{X: 500, Y: 300}).
		SetStyle(turdgl.Style{Colour: color.RGBA{0, 0, 255, 255}})
	circle1 := turdgl.NewCircle(88, turdgl.Vec{X: 500, Y: 600})
	circle2 := turdgl.NewCircle(130, turdgl.Vec{X: 600, Y: 500}).
		SetStyle(turdgl.Style{Colour: color.RGBA{0, 0, 255, 255}})
	instruction := turdgl.NewText("Use WASD to move the shapes", turdgl.Vec{X: 10}, "../../fonts/arial.ttf")

	// Set up keybinds
	win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(turdgl.KeyW, turdgl.Instantaneous,
		func() {
			rect2.Move(turdgl.Vec{Y: -1})
			circle2.Move(turdgl.Vec{Y: -1})
		})
	win.RegisterKeybind(turdgl.KeyS, turdgl.Instantaneous,
		func() {
			rect2.Move(turdgl.Vec{Y: 1})
			circle2.Move(turdgl.Vec{Y: 1})
		})
	win.RegisterKeybind(turdgl.KeyA, turdgl.Instantaneous,
		func() {
			rect2.Move(turdgl.Vec{X: -1})
			circle2.Move(turdgl.Vec{X: -1})
		})
	win.RegisterKeybind(turdgl.KeyD, turdgl.Instantaneous,
		func() {
			rect2.Move(turdgl.Vec{X: 1})
			circle2.Move(turdgl.Vec{X: 1})
		})

	for win.IsRunning() {
		win.SetBackground(color.Black)

		// Adjust shape colours to react to collisions
		rect1.SetStyle(turdgl.DefaultStyle)
		if turdgl.IsColliding(rect1, rect2) {
			rect1.SetStyle(turdgl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}
		if turdgl.IsColliding(circle2, rect1) {
			rect1.SetStyle(turdgl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}
		circle1.SetStyle(turdgl.DefaultStyle)
		if turdgl.IsColliding(circle1, circle2) {
			circle1.SetStyle(turdgl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}
		if turdgl.IsColliding(rect2, circle1) {
			circle1.SetStyle(turdgl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}

		// Draw shapes
		for _, shape := range []turdgl.Drawable{
			rect1,
			rect2,
			circle1,
			circle2,
			instruction,
		} {
			win.Draw(shape)
		}

		// Update the window
		win.Update()
	}
}
