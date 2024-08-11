package turdgl

import (
	"testing"
)

func TestWithinFrame(t *testing.T) {
	type tc struct {
		vec      Vec
		padding  float64
		expected bool
	}

	w, h := 100, 200
	f := NewFrameBuffer(w, h)
	for n, tc := range []tc{
		{vec: Vec{50, 100}, padding: 0, expected: true},
		{vec: Vec{1000, 2000}, padding: 0, expected: false},
		{vec: Vec{100, 200}, padding: 0, expected: true},
		{vec: Vec{100, 200}, padding: 1, expected: false},
		{vec: Vec{10, 10}, padding: 20, expected: false},
	} {
		actual := f.WithinFrame(tc.vec, tc.padding)
		if tc.expected != actual {
			t.Errorf("Test: %d\nExpected: %v\nGot: %v", n+1, tc.expected, actual)
		}
	}
}

func BenchmarkAlphaBlend(b *testing.B) {
	src := Pixel{100, 80, 70, 60}
	dst := Pixel{200, 160, 140, 60}

	for n := 0; n < b.N; n++ {
		_ = alphaBlend(src, dst)
	}
}
