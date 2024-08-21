package main

import (
	"fmt"
	"image/color"
	"time"

	tgl "github.com/z-riley/turdgl"
)

func main() {
	win, err := tgl.NewWindow(tgl.WindowCfg{
		Title:  "Shape Collision Example",
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
	rect1 := tgl.NewRect(100, 60, tgl.Vec{X: 500, Y: 200})
	rect2 := tgl.NewRect(130, 50, tgl.Vec{X: 500, Y: 300},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{0, 0, 255, 255}}),
	)
	circle1 := tgl.NewCircle(88, tgl.Vec{X: 500, Y: 600})
	circle2 := tgl.NewCircle(130, tgl.Vec{X: 600, Y: 500},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{0, 0, 255, 255}}),
	)

	// Keybinds
	win.RegisterKeybind(tgl.KeyEscape, tgl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyLCtrl, tgl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyW, tgl.Instantaneous,
		func() {
			rect2.Move(tgl.Vec{Y: -1})
			circle2.Move(tgl.Vec{Y: -1})
		})
	win.RegisterKeybind(tgl.KeyS, tgl.Instantaneous,
		func() {
			rect2.Move(tgl.Vec{Y: 1})
			circle2.Move(tgl.Vec{Y: 1})
		})
	win.RegisterKeybind(tgl.KeyA, tgl.Instantaneous,
		func() {
			rect2.Move(tgl.Vec{X: -1})
			circle2.Move(tgl.Vec{X: -1})
		})
	win.RegisterKeybind(tgl.KeyD, tgl.Instantaneous,
		func() {
			rect2.Move(tgl.Vec{X: 1})
			circle2.Move(tgl.Vec{X: 1})
		})

	for win.IsRunning() {
		win.SetBackground(color.Black)

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

		// Draw shapes
		win.Draw(rect1)
		win.Draw(rect2)
		win.Draw(circle1)
		win.Draw(circle2)

		win.Update()

		// Count FPS
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", win.GetConfig().Title, frames))
			frames = 0
		default:
		}
	}
}
