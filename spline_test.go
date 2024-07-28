package turdgl

import (
	"fmt"
	"testing"
)

func TestGenerateCatmullRomSpline(t *testing.T) {
	points := []Vec{
		{0, 0},
		{1, 2},
		{2, 3},
		{3, 5},
		{4, 4},
		{5, 2},
	}
	splinePoints := GenerateCatmullRomSpline(points, 100)
	for _, p := range splinePoints {
		fmt.Printf("x: %.2f, y: %.2f\n", p.X, p.Y)
	}
}
