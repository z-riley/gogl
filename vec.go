package gogl

import (
	"fmt"
	"math"
)

// Vec is a Cartesian vector.
type Vec struct{ X, Y float64 }

// String returns a human-readable version of the vector.
func (v Vec) String() string {
	return fmt.Sprintf("{%.2f, %.2f}", v.X, v.Y)
}

// Mag calculates the magnitude of a vector.
func (v Vec) Mag() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// SetMag sets the magnitude of the vector whilst preserving the direction.
func (v Vec) SetMag(newMag float64) Vec {
	// Get unit vector in same direction, and scale it by the new magnitude
	uv := Normalise(v)
	return Vec{X: uv.X * newMag, Y: uv.Y * newMag}
}

// Rotate rotates a vector clockwise by theta radians, preserving the magnitude.
func (v Vec) Rotate(theta float64) Vec {
	// Apply rotation matrix
	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)

	return Vec{
		X: v.X*cosTheta + v.Y*sinTheta,
		Y: -v.X*sinTheta + v.Y*cosTheta,
	}
}

// Round rounds the components of the vector to the nearest whole numbers.
func (v Vec) Round() Vec {
	return Vec{math.Round(v.X), math.Round(v.Y)}
}

// Normalise returns a unit vector with the same direction.
func Normalise(v Vec) Vec {
	mag := v.Mag()
	return Vec{X: v.X / mag, Y: v.Y / mag}
}

// Dist returns the distance between two vectors, assuming they both
// originate from the same point.
func Dist(v1, v2 Vec) float64 {
	aSqr := (v1.X - v2.X) * (v1.X - v2.X)
	bSqr := (v1.Y - v2.Y) * (v1.Y - v2.Y)
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

// Dot calculates the dot product between two vectors.
func Dot(v1, v2 Vec) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// Cross calculates the cross product of two vectors.
func Cross(v1, v2 Vec) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// Theta calculates the angle between two vectors, in radians.
func Theta(v1, v2 Vec) float64 {
	return math.Atan2(v1.Y, v1.X) - math.Atan2(v2.Y, v2.X)
}
