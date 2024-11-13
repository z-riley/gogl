package turdgl

import "golang.org/x/exp/constraints"

// Clamp constraints a variable between lower and upper bounds.
func Clamp[T constraints.Ordered](x, lower, upper T) T {
	switch {
	case x > upper:
		return upper
	case x < lower:
		return lower
	default:
		return x
	}
}

func UNUSED(...any) {}
