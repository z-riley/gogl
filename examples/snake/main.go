package main

import (
	"image/color"
	"time"

	"github.com/z-riley/gogl"
)

func main() {
	win, err := gogl.NewWindow(gogl.WindowCfg{
		Title:  "Moving Snake Example",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	snake := NewSnake(gogl.Vec{X: 400, Y: 100})
	instruction := gogl.NewText("Use WASD to move", gogl.Vec{X: 10}, "../../fonts/arial.ttf")

	win.RegisterKeybind(gogl.KeyEscape, gogl.KeyPress, func() { win.Quit() })

	prevTime := time.Now()

	for win.IsRunning() {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// React to pressed keys
		const speed = 1000
		if win.KeyIsPressed(gogl.KeyW) {
			snake.velocity = &gogl.Vec{Y: -speed}
			snake.Update(dt, win.Framebuffer)
		}
		if win.KeyIsPressed(gogl.KeyA) {
			snake.velocity = &gogl.Vec{X: -speed}
			snake.Update(dt, win.Framebuffer)
		}
		if win.KeyIsPressed(gogl.KeyS) {
			snake.velocity = &gogl.Vec{Y: speed}
			snake.Update(dt, win.Framebuffer)
		}
		if win.KeyIsPressed(gogl.KeyD) {
			snake.velocity = &gogl.Vec{X: speed}
			snake.Update(dt, win.Framebuffer)
		}

		win.SetBackground(color.RGBA{39, 45, 53, 0})

		win.Draw(snake)
		win.Draw(instruction)

		win.Update()
	}
}
