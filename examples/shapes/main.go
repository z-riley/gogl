package main

import (
	"fmt"
	"image/color"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/z-riley/gogl"
)

func main() {
	go func() { // for pprof
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	win, err := gogl.NewWindow(gogl.WindowCfg{
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
	rectSolid := gogl.NewRect(
		120, 90,
		gogl.Vec{X: 50, Y: 50},
	).SetStyle(gogl.Style{Colour: color.RGBA{0, 0, 255, 255}, Thickness: 0, Bloom: 30})

	rectOutline := gogl.NewRect(
		120, 90,
		gogl.Vec{X: 50, Y: 50},
	).SetStyle(gogl.Style{Colour: color.RGBA{255, 0, 0, 255}, Thickness: 2})

	// Buttons can be constructed from shapes
	styleHover := gogl.Style{Colour: color.RGBA{180, 180, 0, 255}, Thickness: 0, Bloom: 5}
	stylePressed := gogl.Style{Colour: color.RGBA{255, 0, 0, 255}, Thickness: 0, Bloom: 10}
	styleUnpressed := gogl.Style{Colour: color.RGBA{80, 0, 0, 255}, Thickness: 30}
	c := gogl.NewCircle(100, gogl.Vec{X: 300, Y: 100}).SetStyle(styleUnpressed)
	circleButton := gogl.NewButton(c, "../../fonts/arial.ttf").
		SetLabelText("Press me").
		SetLabelSize(16).
		SetLabelColour(color.White)
	circleButton.SetCallback(gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnHold},
		func() {
			c.SetStyle(styleHover)
			circleButton.SetLabelText("Hovering")
		}).
		SetCallback(gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnRelease},
			func() {
				c.SetStyle(styleUnpressed)
				circleButton.SetLabelText("Press me")
			}).
		SetCallback(gogl.ButtonTrigger{State: gogl.LeftClick, Behaviour: gogl.OnPress},
			func() {
				c.SetStyle(stylePressed)
				circleButton.SetLabelText("Pressed!")
			}).
		SetCallback(gogl.ButtonTrigger{State: gogl.LeftClick, Behaviour: gogl.OnRelease},
			func() {
				c.SetStyle(styleUnpressed)
				circleButton.SetLabelText("Press me")
			})
	circleButton.Label.SetAlignment(gogl.AlignBottomCentre)

	// More complex shapes with limited feature sets can also be created
	triangle := gogl.NewTriangle(
		gogl.Vec{X: 400, Y: 50},
		gogl.Vec{X: 580, Y: 120},
		gogl.Vec{X: 450, Y: 130},
	).SetStyle(gogl.Style{Colour: color.RGBA{100, 10, 100, 255}})

	txt := gogl.NewText("Hello there", gogl.Vec{X: 800, Y: 80}, "../../fonts/arial.ttf").
		SetColour(color.RGBA{255, 255, 255, 255}).
		SetSize(40)

	curvedRect := gogl.NewCurvedRect(
		120, 90, 20,
		gogl.Vec{X: 50, Y: 200},
	).SetStyle(gogl.Style{Colour: gogl.Orange, Thickness: 10, Bloom: 0})

	// Put shapes on the background layer to avoid interactions with other shapes
	bgRect := gogl.NewRect(
		120, 90,
		gogl.Vec{X: 250, Y: 200},
	).SetStyle(gogl.Style{Colour: color.RGBA{175, 136, 90, 255}, Thickness: 0})
	fgRect := gogl.NewRect(
		120, 90,
		gogl.Vec{X: 280, Y: 230},
	).SetStyle(gogl.Style{Colour: gogl.RosyBrown, Thickness: 0})
	txtBox := gogl.NewTextBox(fgRect, "Click to edit text", "../../fonts/arial.ttf").
		SetTextSize(46).
		SetTextColour(color.RGBA{100, 255, 100, 100})
	txtBox.SetSelectedCB(func() { txtBox.SetTextColour(gogl.Yellow) })
	txtBox.SetDeselectedCB(func() { txtBox.SetTextColour(color.RGBA{100, 255, 100, 200}) })

	// Set variable alpha values to blend shape colours
	circleRed := gogl.NewCircle(
		100,
		gogl.Vec{X: 550, Y: 210},
	).SetStyle(gogl.Style{Colour: color.RGBA{255, 0, 0, 100}})
	circleGreen := gogl.NewCircle(
		100,
		gogl.Vec{X: 580, Y: 260},
	).SetStyle(gogl.Style{Colour: color.RGBA{0, 255, 0, 100}})
	circleBlue := gogl.NewCircle(
		100,
		gogl.Vec{X: 520, Y: 260},
	).SetStyle(gogl.Style{Colour: color.RGBA{0, 0, 255, 100}})

	// Register window-level keybinds
	win.RegisterKeybind(gogl.KeyEscape, gogl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(gogl.KeyLCtrl, gogl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(gogl.KeyE, gogl.Instantaneous, func() { rectSolid.Move(gogl.Vec{X: 2, Y: 2}) })

	for win.IsRunning() {
		win.SetBackground(color.RGBA{35, 39, 46, 255})

		// Dynamic components' Update method must be called
		circleButton.Update(win)
		txtBox.Update(win)

		// Draw shapes
		for _, shape := range []gogl.Drawable{
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
		win.SetTitle(fmt.Sprintf("%s Colour: %x", loc, win.Framebuffer.GetPixel(int(loc.X), int(loc.Y))))

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
