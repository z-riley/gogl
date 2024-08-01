package turdgl

import (
	"reflect"
	"testing"
)

func TestBytes(t *testing.T) {
	f := FrameBuffer{
		[]Pixel{{0, 0, 0, 0}, {1, 1, 1, 1}, {2, 2, 2, 2}, {3, 3, 3, 3}},
		[]Pixel{{4, 4, 4, 4}, {5, 5, 5, 5}, {6, 6, 6, 6}, {7, 7, 7, 7}},
	}
	expected := []byte{
		4, 4, 4, 4,
		5, 5, 5, 5,
		6, 6, 6, 6,
		7, 7, 7, 7,
		0, 0, 0, 0,
		1, 1, 1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3,
	}
	actual := f.BytesReverse()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected: %v\nGot: %v", expected, actual)
	}
}

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
