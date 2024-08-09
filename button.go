package turdgl

type button interface {
	Draw(*FrameBuffer)
	IsHovering(Vec) bool
}

// Button can be build on top of shapes to create pressable buttons.
type Button struct {
	Shape          button           // the base shape the button is built on
	CB             func(MouseState) // the callback function to execute on press
	Trigger        MouseState       // which mouse button must be used to press the button
	Behaviour      ButtonBehaviour  // how the button responds to being pressed
	prevMouseState MouseState       // the previous mouse button state
}

// NewButton constructs a new button from any shape that satisfies
// the hoverable interface.
func NewButton(shape button, callback func(MouseState)) *Button {
	return &Button{
		Shape:     shape,
		CB:        callback,
		Trigger:   LeftClick,
		Behaviour: OnPressAndRelease,
	}
}

// Draw draws the button's base shape to the provided frame buffer.
func (b *Button) Draw(buf *FrameBuffer) {
	b.Shape.Draw(buf)
}

// ButtonBehaviour represents how a button responds to being pressed.
type ButtonBehaviour int

const (
	OnPress           ButtonBehaviour = iota // execute behaviour on press
	OnRelease                                // execute behaviour on release
	OnPressAndRelease                        // execute behaviour on press and release
	OnHold                                   // execute behaviour as long as button is held down
)

// Update examines button state and executes behaviour accordingly.
func (b *Button) Update(win *Window) {

	currentMouseState := win.MouseButtonState()
	hovering := b.Shape.IsHovering(win.MouseLocation())

	switch b.Behaviour {
	case OnPress:
		if hovering && (b.prevMouseState == NoClick && currentMouseState == b.Trigger) {
			b.CB(currentMouseState)
		}
	case OnRelease:
		if hovering && (b.prevMouseState == b.Trigger && currentMouseState == NoClick) {
			b.CB(currentMouseState)
		}
	case OnPressAndRelease:
		if hovering && (b.prevMouseState != currentMouseState) {
			b.CB(currentMouseState)
		}
	case OnHold:
		if hovering && (currentMouseState == b.Trigger) {
			b.CB(currentMouseState)
		}
	default:
		panic("unsupported button behaviour")
	}

	b.prevMouseState = win.MouseButtonState()
}
