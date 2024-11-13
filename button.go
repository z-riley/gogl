package turdgl

import (
	"image/color"
)

// hoverable is an interface for shapes can detect cursor hovering.
type hoverable interface {
	Shape
	IsWithin(Vec) bool
}

// Button can be build on top of shapes to create pressable buttons.
type Button struct {
	Shape     hoverable                // the base shape the button is built on
	Label     *Text                    // the text to display on the button (if any)
	Callbacks map[ButtonTrigger]func() // mapping of triggers to callback functions
	IsEnabled bool                     // stops callbacks from executing if false

	prevMouseState MouseState
	prevHovering   bool
}

// NewButton constructs a new button from any shape that satisfies the buttonable interface.
func NewButton(shape hoverable, fontPath string) *Button {
	return &Button{
		Shape:     shape,
		Label:     NewText("", shape.GetPos(), fontPath).SetAlignment(AlignCentre),
		Callbacks: make(map[ButtonTrigger]func()),
		IsEnabled: true,
	}
}

// Draw draws the button onto the frame buffer.
func (b *Button) Draw(buf *FrameBuffer) {
	b.Shape.Draw(buf)

	// Align to centre of underlying shape
	b.Label.SetPos(func() Vec {
		switch b.Shape.(type) {
		case *Rect, *CurvedRect:
			p := b.Shape.GetPos()
			return Vec{p.X + b.Shape.Width()/2, p.Y + b.Shape.Height()/2}
		default:
			return b.Shape.GetPos()
		}
	}())

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
)

// Update examines button state and executes behaviour accordingly.
func (b *Button) Update(win *Window) {
	mouseState := win.MouseButtonState()
	hovering := b.Shape.IsWithin(win.MouseLocation())

	if b.IsEnabled {
		for trigger, cb := range b.Callbacks {
			stateMatchesTrigger := func() bool {
				switch {
				case trigger.Behaviour == OnAll:
					return true

				case (hovering && b.prevMouseState != trigger.State && mouseState == trigger.State),
					(trigger.State == (NoClick) && !b.prevHovering && hovering):
					return trigger.Behaviour == OnPress || trigger.Behaviour == OnPressAndRelease

				case (hovering && b.prevMouseState == trigger.State && mouseState != trigger.State),
					(trigger.State == (NoClick) && b.prevHovering && !hovering):
					return trigger.Behaviour == OnRelease || trigger.Behaviour == OnPressAndRelease

				case hovering && mouseState == trigger.State:
					return trigger.Behaviour == OnHold

				default:
					return false
				}
			}()

			if stateMatchesTrigger {
				cb()
			}
		}
	}

	b.prevMouseState = win.MouseButtonState()
	b.prevHovering = b.Shape.IsWithin(win.MouseLocation())
}

// ButtonTrigger is a trigger for executing a button callback function.
type ButtonTrigger struct {
	// State is the state of the mouse buttons in a given moment.
	State MouseState
	// Behaviour describes at what point of the mouse click to execute the callback on.
	Behaviour ButtonBehaviour
}

// SetCallback configures a callback function to execute if the conditions described
// by the trigger are met.
func (b *Button) SetCallback(trigger ButtonTrigger, callback func()) *Button {
	b.Callbacks[trigger] = callback
	return b
}

// UnsetCallback removes a callback for a specified trigger. If no callback is configured,
// nothing will happen.
func (b *Button) UnsetCallback(trigger ButtonTrigger) *Button {
	delete(b.Callbacks, trigger)
	return b
}

// Disable disables the button. If it is already disabled, nothing happens.
func (b *Button) Disable() *Button {
	b.IsEnabled = false
	return b
}

// Enable enables the button. If it is already enabled, nothing happens.
func (b *Button) Enable() *Button {
	b.IsEnabled = true
	return b
}

// Move moves the button by a given vector.
func (b *Button) Move(mov Vec) {
	b.Shape.Move(mov)
	b.Label.Move(mov)
}

// IsHovering returns whether the cursor is hovering over the button.
func (b *Button) IsHovering() bool {
	return b.prevHovering
}

// SetLabelText sets the text label to the given string.
func (b *Button) SetLabelText(s string) *Button {
	b.Label.SetText(s)
	return b
}

// SetLabelAlignment sets the alignment of the text label relative to the centre of the shape.
func (b *Button) SetLabelAlignment(align Alignment) *Button {
	b.Label.SetAlignment(align)
	return b
}

// SetLabelOffset manually sets the label's offset, providing the text is in AlignCustom mode.
func (b *Button) SetLabelOffset(offset Vec) *Button {
	b.Label.SetOffset(offset)
	return b
}

// SetLabelPos sets the label's position on the screen.
func (b *Button) SetLabelPos(pos Vec) *Button {
	b.Label.SetPos(pos)
	return b
}

// SetLabelColour sets the label text colour.
func (b *Button) SetLabelColour(c color.Color) *Button {
	b.Label.SetColour(c)
	return b
}

// SetLabelFont sets the path of the .ttf file that is used to generate the label.
func (b *Button) SetLabelFont(path string) error {
	return b.Label.SetFont(path)
}

// SetLabelDPI sets the DPI of the label font.
func (b *Button) SetLabelDPI(dpi float64) *Button {
	b.Label.SetDPI(dpi)
	return b
}

// SetLabelSize sets the size of the label font.
func (b *Button) SetLabelSize(size float64) *Button {
	b.Label.SetSize(size)
	return b
}

// SetLabelSpacing sets the line spacing of the label.
func (b *Button) SetLabelSpacing(spacing float64) *Button {
	b.Label.SetSpacing(spacing)
	return b
}
