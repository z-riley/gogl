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
	actual := f.Bytes()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected: %v\nGot: %v", expected, actual)
	}
}
