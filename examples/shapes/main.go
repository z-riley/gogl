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
	rectSolid := turdgl.NewRect(
		120, 90,
		turdgl.Vec{X: 50, Y: 50},
	).SetStyle(turdgl.Style{Colour: color.RGBA{0, 0, 255, 255}, Thickness: 0, Bloom: 30})

	rectOutline := turdgl.NewRect(
		120, 90,
		turdgl.Vec{X: 50, Y: 50},
	).SetStyle(turdgl.Style{Colour: color.RGBA{255, 0, 0, 255}, Thickness: 2})

	// Buttons can be constructed from shapes
	styleHover := turdgl.Style{Colour: color.RGBA{180, 180, 0, 255}, Thickness: 0, Bloom: 5}
	stylePressed := turdgl.Style{Colour: color.RGBA{255, 0, 0, 255}, Thickness: 0, Bloom: 10}
	styleUnpressed := turdgl.Style{Colour: color.RGBA{80, 0, 0, 255}, Thickness: 30}
	c := turdgl.NewCircle(100, turdgl.Vec{X: 300, Y: 100}).SetStyle(styleUnpressed)
	circleButton := turdgl.NewButton(c, "../../fonts/arial.ttf").
		SetLabelText("Press me").
		SetLabelSize(16).
		SetLabelColour(color.White)
	circleButton.SetCallback(turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			c.SetStyle(styleHover)
			circleButton.SetLabelText("Hovering")
		}).
		SetCallback(turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
			func() {
				c.SetStyle(styleUnpressed)
				circleButton.SetLabelText("Press me")
			}).
		SetCallback(turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnPress},
			func() {
				c.SetStyle(stylePressed)
				circleButton.SetLabelText("Pressed!")
			}).
		SetCallback(turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
			func() {
				c.SetStyle(styleUnpressed)
				circleButton.SetLabelText("Press me")
			})
	circleButton.Label.SetAlignment(turdgl.AlignBottomCentre)

	// More complex shapes with limited feature sets can also be created
	triangle := turdgl.NewTriangle(
		turdgl.Vec{X: 400, Y: 50},
		turdgl.Vec{X: 580, Y: 120},
		turdgl.Vec{X: 450, Y: 130},
	).SetStyle(turdgl.Style{Colour: color.RGBA{100, 10, 100, 255}})

	txt := turdgl.NewText("Hello there", turdgl.Vec{X: 800, Y: 80}, "../../fonts/arial.ttf").
		SetColour(color.RGBA{255, 255, 255, 255}).
		SetSize(40)

	curvedRect := turdgl.NewCurvedRect(
		120, 90, 20,
		turdgl.Vec{X: 50, Y: 200},
	).SetStyle(turdgl.Style{Colour: turdgl.Orange, Thickness: 10, Bloom: 0})

	// Put shapes on the background layer to avoid interactions with other shapes
	bgRect := turdgl.NewRect(
		120, 90,
		turdgl.Vec{X: 250, Y: 200},
	).SetStyle(turdgl.Style{Colour: color.RGBA{175, 136, 90, 255}, Thickness: 0})
	fgRect := turdgl.NewRect(
		120, 90,
		turdgl.Vec{X: 280, Y: 230},
	).SetStyle(turdgl.Style{Colour: turdgl.RosyBrown, Thickness: 0})
	txtBox := turdgl.NewTextBox(fgRect, "Click to edit text", "../../fonts/arial.ttf").
		SetTextSize(46).
		SetTextColour(color.RGBA{100, 255, 100, 100})
	txtBox.SetSelectedCB(func() { txtBox.SetTextColour(turdgl.Yellow) })
	txtBox.SetDeselectedCB(func() { txtBox.SetTextColour(color.RGBA{100, 255, 100, 200}) })

	// Set variable alpha values to blend shape colours
	circleRed := turdgl.NewCircle(
		100,
		turdgl.Vec{X: 550, Y: 210},
	).SetStyle(turdgl.Style{Colour: color.RGBA{255, 0, 0, 100}})
	circleGreen := turdgl.NewCircle(
		100,
		turdgl.Vec{X: 580, Y: 260},
	).SetStyle(turdgl.Style{Colour: color.RGBA{0, 255, 0, 100}})
	circleBlue := turdgl.NewCircle(
		100,
		turdgl.Vec{X: 520, Y: 260},
	).SetStyle(turdgl.Style{Colour: color.RGBA{0, 0, 255, 100}})

	// Register window-level keybinds
	win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(turdgl.KeyLCtrl, turdgl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(turdgl.KeyE, turdgl.Instantaneous, func() { rectSolid.Move(turdgl.Vec{X: 2, Y: 2}) })

	for win.IsRunning() {
		win.SetBackground(color.RGBA{35, 39, 46, 255})

		// Dynamic components' Update method must be called
		circleButton.Update(win)
		txtBox.Update(win)

		// Draw shapes
		for _, shape := range []turdgl.Drawable{
			bgRect,
			rectSolid,
			rectOutline,
			curvedRect,
			triangle,
			txt,
			circleButton,
			txtBox,
			circleRed,
			circleBlue,
			circleGreen,
		} {
			win.Draw(shape)
		}

		// Lastly, the window must be updated
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
