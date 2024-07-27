package turdgl

import "math"

// Vec represents Cartesian coordinates on a pixel grid.
type Vec struct{ X, Y int }

// Add adds a vector to v.
func (v Vec) Add(vec Vec) Vec {
	return Vec{
		X: v.X + vec.X,
		Y: v.Y + vec.Y,
	}
}

// Sub subtracts a vector to v.
func (v Vec) Sub(vec Vec) Vec {
	return Vec{
		X: v.X - vec.X,
		Y: v.Y - vec.Y,
	}
}

// Dist returns the distance between two vectors, assuming they both
// originate from the same point.
func Dist(v1, v2 Vec) float64 {
	aSqr := math.Pow(float64(v1.X-v2.X), 2)
	bSqr := math.Pow(float64(v1.Y-v2.Y), 2)
	return math.Sqrt(aSqr + bSqr)
}
