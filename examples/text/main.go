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
	fpsCounter := turdgl.NewText("Measuing FPS...", turdgl.Vec{X: 1000, Y: 50}, "../../fonts/arial.ttf").
		SetColour(color.RGBA{255, 255, 255, 255}).
		SetSize(20)

	dynamicText := turdgl.NewText("Testing", turdgl.Vec{X: 300, Y: 100}, "../../fonts/arial.ttf")
	go func() {
		for {
			for i := range 5 {
				dynamicText.SetText(fmt.Sprint(dynamicText.Text(), " ", i))
				time.Sleep(time.Second)
			}
			dynamicText.SetText("Testing")
			time.Sleep(time.Second)
		}
	}()

	bottomRight := turdgl.NewText("Bottom-right alignment", turdgl.Vec{X: 300, Y: 200}, "../../fonts/arial.ttf").
		SetAlignment(turdgl.AlignBottomRight)

	dynamicAlignment := turdgl.NewText("Dynamic alignment", turdgl.Vec{X: 1000, Y: 400}, "../../fonts/arial.ttf").
		SetAlignment(turdgl.AlignBottomRight).
		SetSize(50)
	dynamicMarker := turdgl.NewRect(5, 5, dynamicAlignment.Pos())
	type alignPair struct {
		alignment turdgl.Alignment
		label     string
	}
	// Change alignment using mouse scroll wheel
	alignments := []alignPair{
		{turdgl.AlignTopLeft, "Top left alignment"},
		{turdgl.AlignTopCentre, "Top centre alignment"},
		{turdgl.AlignTopRight, "Top right alignment"},
		{turdgl.AlignCentreLeft, "Centre left alignment"},
		{turdgl.AlignCentre, "Centre alignment"},
		{turdgl.AlignCentreRight, "Centre right alignment"},
		{turdgl.AlignBottomLeft, "Bottom left alignment"},
		{turdgl.AlignBottomCentre, "Bottom centre alignment"},
		{turdgl.AlignBottomRight, "Bottom right alignment"},
	}
	i := 0
	win.SetMouseScrollCallback(func(movement turdgl.Vec) {
		if movement.IsScrollUp() {
			i++
			if i > len(alignments)-1 {
				i = 0
			}

		} else if movement.IsScrollDown() {
			i--
			if i < 0 {
				i = len(alignments) - 1
			}
		}
		dynamicAlignment.SetAlignment(alignments[i].alignment)
		dynamicAlignment.SetText(alignments[i].label)
	})

	// Register window-level keybinds
	win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(turdgl.KeyLCtrl, turdgl.KeyPress, func() { win.Quit() })

	for win.IsRunning() {
		win.SetBackground(color.RGBA{35, 39, 46, 255})

		// Draw foreground shapes
		for _, shape := range []turdgl.Drawable{
			fpsCounter,
			dynamicText,
			bottomRight,
			dynamicMarker,
			dynamicAlignment,
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
			fpsCounter.SetText(fmt.Sprintf("FPS: %d", frames))
			frames = 0
		default:
		}
	}
}
