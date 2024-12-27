package gogl

import (
	"github.com/veandco/go-sdl2/sdl"
)

// MouseScrollCallback is executed when the user scrolls the mouse wheel in any direction.
//
// Positive X movement means scrolling to the right. Positive Y movement means scrolling up.
type MouseScrollCallback func(movement Vec)

// IsScrollLeft returns true if the mouse scroll vector is leftwards.
func (v Vec) IsScrollLeft() bool {
	return v.X < 0
}

// IsScrollRight returns true if the mouse scroll vector is rightwards.
func (v Vec) IsScrollRight() bool {
	return v.X > 0
}

// IsScrollUp returns true if the mouse scroll vector is upwards.
func (v Vec) IsScrollUp() bool {
	return v.Y > 0
}

// IsScrollDown returns true if the mouse scroll vector is downwards.
func (v Vec) IsScrollDown() bool {
	return v.Y < 0
}

// mouseScrollHandler handles mouse wheel events.
type mouseScrollHandler struct {
	Callback MouseScrollCallback
}

// newMouseScrollHandler constructs a mouse scroll handler with a blank callback function.
func newMouseScrollHandler() *mouseScrollHandler {
	return &mouseScrollHandler{
		Callback: func(_ Vec) {},
	}
}

// handleEvent handles a mouse scroll event.
func (m *mouseScrollHandler) handleEvent(event *sdl.MouseWheelEvent) {
	m.Callback(Vec{float64(event.PreciseX), float64(event.PreciseY)})
}
