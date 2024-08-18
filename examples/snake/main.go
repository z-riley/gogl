package main

import (
	"image/color"
	"time"

	tgl "github.com/z-riley/turdgl"
)

func main() {
	win, err := tgl.NewWindow(tgl.WindowCfg{
		Title:  "Moving Snake Example",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	snake := NewSnake(tgl.Vec{X: 400, Y: 100})

	win.RegisterKeybind(tgl.KeyEscape, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyLCtrl, func() { win.Quit() })

	prevTime := time.Now()

	for win.IsRunning() {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// React to pressed keys
		const speed = 1000
		if win.KeyIsPressed(tgl.KeyW) {
			snake.velocity = &tgl.Vec{Y: -speed}
			snake.Update(dt, win.Framebuffer)
		}
		if win.KeyIsPressed(tgl.KeyA) {
			snake.velocity = &tgl.Vec{X: -speed}
			snake.Update(dt, win.Framebuffer)
		}
		if win.KeyIsPressed(tgl.KeyS) {
			snake.velocity = &tgl.Vec{Y: speed}
			snake.Update(dt, win.Framebuffer)
		}
		if win.KeyIsPressed(tgl.KeyD) {
			snake.velocity = &tgl.Vec{X: speed}
			snake.Update(dt, win.Framebuffer)
		}

		// Set background colour
		win.SetBackground(color.RGBA{39, 45, 53, 0})

		// Modify frame buffer
		win.Draw(snake)

		win.Update()
	}
}
