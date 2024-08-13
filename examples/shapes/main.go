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
		Width:  1200,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// For measuring FPS
	frames := 0
	second := time.Tick(time.Second)

	// Shapes can be filled or have an outline of specified thickness
	rectSolid := tgl.NewRect(
		120, 90,
		tgl.Vec{X: 50, Y: 50},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{0, 0, 255, 255}, Thickness: 0, Bloom: 30}),
	)
	rectOutline := tgl.NewRect(
		120, 90,
		tgl.Vec{X: 50, Y: 50},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{255, 0, 0, 255}, Thickness: 2}),
	)

	// Buttons can be constructed from shapes
	styleUnpressed := tgl.Style{Colour: color.RGBA{80, 0, 0, 255}, Thickness: 30}
	stylePressed := tgl.Style{Colour: color.RGBA{255, 0, 0, 255}, Thickness: 0, Bloom: 10}
	c := tgl.NewCircle(100, tgl.Vec{X: 300, Y: 100}, tgl.WithStyle(styleUnpressed))
	circleButton := tgl.NewButton(c, "../../fonts/arial.ttf").SetText("Press me")
	circleButton.Label.SetSize(16)
	circleButton.Label.SetColour(color.White)
	circleButton.SetCallback(func(m tgl.MouseState) {
		if m == tgl.LeftClick {
			c.SetStyle(stylePressed)
			circleButton.SetText("Pressed!")
		} else {
			c.SetStyle(styleUnpressed)
			circleButton.SetText("Press me")
		}
	})
	circleButton.Label.SetAlignment(tgl.AlignBottomCentre)

	// More complex shapes with limited feature sets can also be created
	triangle := tgl.NewTriangle(
		tgl.Vec{X: 400, Y: 50},
		tgl.Vec{X: 580, Y: 120},
		tgl.Vec{X: 450, Y: 130},
	).SetStyle(tgl.Style{Colour: color.RGBA{100, 10, 100, 255}})

	txt := tgl.NewText("Hello there", tgl.Vec{X: 800, Y: 80}, "../../fonts/opensans.ttf").
		SetColour(color.RGBA{255, 255, 255, 255}).
		SetSize(40)

	curvedRect := tgl.NewCurvedRect(
		120, 90, 12,
		tgl.Vec{X: 50, Y: 200},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{240, 170, 90, 255}, Thickness: 8}),
	)

	// Put shapes on the background layer to avoid interactions with other shapes
	bgRect := tgl.NewRect(
		120, 90,
		tgl.Vec{X: 250, Y: 200},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{175, 136, 90, 255}, Thickness: 0}),
	)
	fgRect := tgl.NewRect(
		120, 90,
		tgl.Vec{X: 280, Y: 230},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{90, 65, 48, 255}, Thickness: 0}),
	)

	// Set variable alpha values to blend shape colours
	circleRed := tgl.NewCircle(
		100,
		tgl.Vec{X: 550, Y: 210},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{255, 0, 0, 100}}),
	)
	circleGreen := tgl.NewCircle(
		100,
		tgl.Vec{X: 580, Y: 260},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{0, 255, 0, 100}}),
	)
	circleBlue := tgl.NewCircle(
		100,
		tgl.Vec{X: 520, Y: 260},
		tgl.WithStyle(tgl.Style{Colour: color.RGBA{0, 0, 255, 100}}),
	)

	// Register window-level keybinds
	win.RegisterKeybind(tgl.KeyEscape, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyLCtrl, func() { win.Quit() })
	win.RegisterKeybind(tgl.KeyE, func() { rectSolid.Move(tgl.Vec{X: 2, Y: 2}) })

	for win.IsRunning() {
		win.SetBackground(color.RGBA{35, 39, 46, 255})

		// Draw foreground shapes
		for _, shape := range []tgl.Drawable{
			rectSolid,
			rectOutline,
			curvedRect,
			triangle,
			txt,
			circleButton,
			fgRect,
			circleRed,
			circleBlue,
			circleGreen,
		} {
			win.Draw(shape)
		}

		// Shapes drawn to the background appear behind foreground shapes
		for _, shape := range []tgl.Drawable{
			bgRect,
		} {
			win.Draw(shape)
		}

		circleButton.Update(win)

		win.Update()

		loc := win.MouseLocation()
		win.SetTitle(fmt.Sprint(loc, win.Framebuffer.GetPixel(int(loc.X), int(loc.Y))))

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
