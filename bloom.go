package turdgl

import (
	"image/color"
	"math"
)

type bloom struct {
	Pos   Vec
	reach float64
	power float64
}

func NewBloom(pos Vec) *bloom {
	return &bloom{
		Pos:   pos,
		reach: 800,
		power: 0.14,
	}
}

// Note: the alpha channel of the background must be a non zero value
func (g *bloom) Draw(buf *FrameBuffer) {
	// Origin pixel
	buf.SetPixel(
		int(math.Round(g.Pos.Y)), int(math.Round(g.Pos.X)),
		NewPixel(color.RGBA{100, 255, 100, 255}),
	)

	width := g.reach
	// Construct bounding box
	radius := width / 2
	bbBoxPos := Vec{g.Pos.X - (radius), g.Pos.Y - (radius)}
	bbox := NewRect(width, width, bbBoxPos)

	// Iterate over every pixel in the bounding box
	for i := bbox.Pos.X; i <= bbox.Pos.X+bbox.w; i++ {
		for j := bbox.Pos.Y; j <= bbox.Pos.Y+bbox.h; j++ {
			dist := Dist(g.Pos, Vec{i, j})
			if dist <= float64(radius) {
				brightness := math.Exp(-(1 / (g.power * g.reach)) * dist)
				c := color.RGBA{
					uint8(brightness * 100),
					uint8(brightness * 255),
					uint8(brightness * 100),
					uint8(brightness * 255),
				}
				buf.addToPixel(int(math.Round(j)), int(math.Round(i)), NewPixel(c))
			}
		}
	}
}
