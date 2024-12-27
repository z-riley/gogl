package main

import (
	"image/color"

	"github.com/z-riley/gogl"
)

func main() {
	win, err := gogl.NewWindow(gogl.WindowCfg{
		Title:  "gogl Collision Example",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// Initialise shapes
	rect1 := gogl.NewRect(100, 60, gogl.Vec{X: 500, Y: 200})
	rect2 := gogl.NewRect(130, 50, gogl.Vec{X: 500, Y: 300}).
		SetStyle(gogl.Style{Colour: color.RGBA{0, 0, 255, 255}})
	circle1 := gogl.NewCircle(88, gogl.Vec{X: 500, Y: 600})
	circle2 := gogl.NewCircle(130, gogl.Vec{X: 600, Y: 500}).
		SetStyle(gogl.Style{Colour: color.RGBA{0, 0, 255, 255}})
	instruction := gogl.NewText("Use WASD to move the shapes", gogl.Vec{X: 10}, "../../fonts/arial.ttf")

	// Set up keybinds
	win.RegisterKeybind(gogl.KeyEscape, gogl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(gogl.KeyW, gogl.Instantaneous,
		func() {
			rect2.Move(gogl.Vec{Y: -1})
			circle2.Move(gogl.Vec{Y: -1})
		})
	win.RegisterKeybind(gogl.KeyS, gogl.Instantaneous,
		func() {
			rect2.Move(gogl.Vec{Y: 1})
			circle2.Move(gogl.Vec{Y: 1})
		})
	win.RegisterKeybind(gogl.KeyA, gogl.Instantaneous,
		func() {
			rect2.Move(gogl.Vec{X: -1})
			circle2.Move(gogl.Vec{X: -1})
		})
	win.RegisterKeybind(gogl.KeyD, gogl.Instantaneous,
		func() {
			rect2.Move(gogl.Vec{X: 1})
			circle2.Move(gogl.Vec{X: 1})
		})

	for win.IsRunning() {
		win.SetBackground(color.Black)

		// Adjust shape colours to react to collisions
		rect1.SetStyle(gogl.DefaultStyle)
		if gogl.IsColliding(rect1, rect2) {
			rect1.SetStyle(gogl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}
		if gogl.IsColliding(circle2, rect1) {
			rect1.SetStyle(gogl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}
		circle1.SetStyle(gogl.DefaultStyle)
		if gogl.IsColliding(circle1, circle2) {
			circle1.SetStyle(gogl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}
		if gogl.IsColliding(rect2, circle1) {
			circle1.SetStyle(gogl.Style{Colour: color.RGBA{255, 0, 0, 255}})
		}

		// Draw shapes
		for _, shape := range []gogl.Drawable{
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
