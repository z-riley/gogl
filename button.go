package turdgl

type hoverable interface {
	IsHovering(Vec) bool
}

type Button struct {
	shape          hoverable
	cb             func()     // the callback function to execute on press
	prevMouseState MouseState // the previous mouse button state
}

// NewButton constructs a new button from any shape that satisfies
// the hoverable interface.
func NewButton(shape hoverable, callback func()) *Button {
	return &Button{
		shape:          shape,
		cb:             callback,
		prevMouseState: NoClick,
	}
}

// Update examines button state and executes behaviour accordingly.
func (b *Button) Update(win *Window) {

	// TODO: implement different behaviour models
	// - Execute once on left/right click
	// - Execute continously on left/right click

	// THIS LOGIC WORKS FOR SINGLE CLICK BEHAVIOUR
	leftClicked := b.prevMouseState == NoClick && win.MouseButtonState() == LeftClick
	if leftClicked {
		b.prevMouseState = LeftClick
		b.cb()
	} else if win.MouseButtonState() == NoClick {
		b.prevMouseState = NoClick
	}
}
