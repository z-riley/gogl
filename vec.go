package turdgl

import "math"

// Vec represents Cartesian coordinates on a pixel grid.
type Vec struct{ X, Y float64 }

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

// Mag calculates the magnitude of a vector.
func (v Vec) Mag() float64 {
	return math.Sqrt(math.Pow(float64(v.X), 2) + math.Pow(float64(v.Y), 2))
}

// SetMag sets the magnitude of the vector whilst preserving the direction.
func (v Vec) SetMag(newMag float64) Vec {
	// Get unit vector in same direction
	mag := v.Mag()
	uv := Vec{X: float64(v.X) / mag, Y: float64(v.Y) / mag}
	// Scale the normalised vector
	return Vec{X: uv.X * newMag, Y: uv.Y * newMag}
}

// Dist returns the distance between two vectors, assuming they both
// originate from the same point.
func Dist(v1, v2 Vec) float64 {
	aSqr := math.Pow(float64(v1.X-v2.X), 2)
	bSqr := math.Pow(float64(v1.Y-v2.Y), 2)
	return math.Sqrt(aSqr + bSqr)
}

// Sub subtracts vector v2 from v1.
func Sub(v1, v2 Vec) Vec {
	return Vec{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
	}
}

// Add adds vector v2 to v1.
func Add(v1, v2 Vec) Vec {
	return Vec{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}
