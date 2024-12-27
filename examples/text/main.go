package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/z-riley/gogl"
)

func main() {
	win, err := gogl.NewWindow(gogl.WindowCfg{
		Title:     "Basic Shapes Example",
		Width:     1200,
		Height:    768,
		Resizable: true,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// For measuring FPS
	frames := 0
	second := time.Tick(time.Second)
	fpsCounter := gogl.NewText("Measuring FPS...", gogl.Vec{X: 1000, Y: 50}, "../../fonts/arial.ttf").
		SetColour(color.RGBA{255, 255, 255, 255}).
		SetSize(20)

	dynamicText := gogl.NewText("Testing", gogl.Vec{X: 300, Y: 100}, "../../fonts/arial.ttf")
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

	bottomRight := gogl.NewText("Bottom right alignment", gogl.Vec{X: 300, Y: 200}, "../../fonts/arial.ttf").
		SetAlignment(gogl.AlignBottomRight)

	dynamicAlignment := gogl.NewText("Scroll to change\nalignment", gogl.Vec{X: 1000, Y: 400}, "../../fonts/arial.ttf").
		SetAlignment(gogl.AlignCentre).
		SetSize(40)
	marker := gogl.NewRect(5, 5, dynamicAlignment.Pos())
	type alignPair struct {
		alignment gogl.Alignment
		label     string
	}
	// Change alignment using mouse scroll wheel
	alignments := []alignPair{
		{gogl.AlignTopLeft, "Top left alignment"},
		{gogl.AlignTopCentre, "Top centre alignment"},
		{gogl.AlignTopRight, "Top right alignment"},
		{gogl.AlignCentreLeft, "Centre left alignment"},
		{gogl.AlignCentre, "Centre alignment"},
		{gogl.AlignCentreRight, "Centre right alignment"},
		{gogl.AlignBottomLeft, "Bottom left alignment"},
		{gogl.AlignBottomCentre, "Bottom centre alignment"},
		{gogl.AlignBottomRight, "Bottom right alignment"},
	}
	i := 0
	win.SetMouseScrollCallback(func(movement gogl.Vec) {
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
	win.RegisterKeybind(gogl.KeyEscape, gogl.KeyPress, func() { win.Quit() })
	win.RegisterKeybind(gogl.KeyLCtrl, gogl.KeyPress, func() { win.Quit() })

	for win.IsRunning() {
		win.SetBackground(color.RGBA{35, 39, 46, 255})

		// Draw shapes
		for _, shape := range []gogl.Drawable{
			fpsCounter,
			dynamicText,
			bottomRight,
			marker,
			dynamicAlignment,
		} {
			win.Draw(shape)
		}

		loc := win.MouseLocation()
		win.SetTitle(
			fmt.Sprintf("Location: %s | Colour: %v", loc, win.Framebuffer.GetPixel(int(loc.X), int(loc.Y))),
		)

		win.Update()

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
