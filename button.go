package turdgl

import "fmt"

// buttonable is an interface for shapes buttons can be built with.
type buttonable interface {
	Draw(*FrameBuffer)
	IsWithin(Vec) bool
	GetPos() Vec
}

// Button can be build on top of shapes to create pressable buttons.
type Button struct {
	Shape     buttonable       // the base shape the button is built on
	Label     *Text            // the text to display on the button (if any)
	CB        func(MouseState) // the callback function to execute on press
	Trigger   MouseState       // which mouse button must be used to press the button
	Behaviour ButtonBehaviour  // how the button responds to being pressed

	prevMouseState MouseState
	prevMouseLoc   Vec
	prevLabel      string
}

// NewButton constructs a new button from any shape that satisfies
// the hoverable interface.
func NewButton(shape buttonable, fontPath string) *Button {
	return &Button{
		Shape:     shape,
		Label:     NewText("", shape.GetPos(), fontPath),
		CB:        func(MouseState) { fmt.Println("Warning: Button callback not configured") },
		Trigger:   LeftClick,
		Behaviour: OnAll,
	}
}

// SetCallback configures a callback function to execute every time a press
// or unpress event occurs. The type of event (left-click, right-click, etc...)
// is passed into the function so the callback can take appropriate action.
func (b *Button) SetCallback(callback func(MouseState)) *Button {
	b.CB = callback
	return b
}

// SetText sets the text label to the given string.
func (b *Button) SetText(s string) *Button {
	b.Label.SetText(s)
	return b
}

// Draw draws the button onto the frame buffer.
func (b *Button) Draw(buf *FrameBuffer) {
	b.Shape.Draw(buf)
	b.Label.Draw(buf)
}

// ButtonBehaviour represents how a button responds to being pressed.
type ButtonBehaviour int

const (
	OnAll             ButtonBehaviour = iota // execute behaviour every time Update() is called
	OnPress                                  // execute behaviour on press
	OnRelease                                // execute behaviour on release
	OnPressAndRelease                        // execute behaviour on press and release
	OnHold                                   // execute behaviour as long as button is held down
	OnHover                                  // execute behaviour if cursor is over button
)

// Update examines button state and executes behaviour accordingly.
func (b *Button) Update(win *Window) {

	currentMouseState := win.MouseButtonState()
	hovering := b.Shape.IsWithin(win.MouseLocation())

	switch b.Behaviour {
	case OnAll:
		b.CB(currentMouseState)
	case OnHover:
		if hovering {
			b.CB(currentMouseState)
		}
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
	b.prevMouseLoc = win.MouseLocation()
}

// IsHovering returns whether the cursor is hovering over the button.
func (b *Button) IsHovering() bool {
	return b.Shape.IsWithin(b.prevMouseLoc)
}
