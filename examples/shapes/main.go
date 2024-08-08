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
	rectSolid := tgl.NewRect(
		120, 90,
		tgl.Vec{X: 50, Y: 50},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{0, 0, 255, 1}, Thickness: 0}),
	)
	rectOutline := tgl.NewRect(
		120, 90,
		tgl.Vec{X: 50, Y: 50},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{255, 0, 0, 1}, Thickness: 2}),
	)
	circle := tgl.NewCircle(
		100,
		tgl.Vec{X: 300, Y: 100},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{255, 0, 0, 1}, Thickness: 30}),
	)
	circleButton := tgl.NewButton(circle, func() { fmt.Println("Circle pressed") })
	triangle := tgl.NewTriangle(
		tgl.Vec{X: 400, Y: 50},
		tgl.Vec{X: 580, Y: 120},
		tgl.Vec{X: 450, Y: 130},
	).SetStyle(tgl.Style{Colour: color.RGBA{100, 10, 100, 255}})
	txt := tgl.NewText("Hello there", tgl.Vec{X: 800, Y: 80}).
		SetColour(color.RGBA{255, 255, 255, 255})

	// Keybinds
	win.RegisterKeybind(tgl.KeyEscape, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyLCtrl, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyE, func() { rectSolid.Move(tgl.Vec{X: 2, Y: 2}) })

	for win.IsRunning() {
		win.Framebuffer.SetBackground(color.RGBA{39, 45, 53, 255})

		// Draw shapes
		win.Draw(rectSolid)
		win.Draw(rectOutline)
		win.Draw(circle)
		win.Draw(triangle)
		win.Draw(txt)

		circleButton.Update(win)

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
