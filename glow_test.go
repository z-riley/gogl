package turdgl

import "testing"

func TestDraw(t *testing.T) {
	g := NewBloom(Vec{0, 0})
	g.Draw(NewFrameBuffer(100, 100))
}
