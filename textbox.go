package turdgl

import (
	"image/color"
)

// TextBox is a shape that can be typed in.
type TextBox struct {
	Shape        hoverable // the base shape the text box is built on
	Text         *Text     // the text in the text box
	selectedCB   func()    // callback function for when the text box is selected (optional)
	deselectedCB func()    // callback function for when the text box is deselected (optional)
	modifiedCB   func()    // callback function for when the text changes (optional)

	isEditing bool
	prevText  string
}

// NewTextBox constructs a new text box from a hoverable shape.
func NewTextBox(shape hoverable, text, fontPath string) *TextBox {
	return &TextBox{
		Shape:        shape,
		Text:         NewText(text, shape.GetPos(), fontPath).SetAlignment(AlignCentre),
		selectedCB:   func() {},
		deselectedCB: func() {},
		modifiedCB:   func() {},
		isEditing:    false,
		prevText:     text,
	}
}

// Draw draws the text box onto the frame buffer.
func (t *TextBox) Draw(buf *FrameBuffer) {
	t.Shape.Draw(buf)

	// Align to centre of underlying shape
	t.Text.SetPos(func() Vec {
		switch t.Shape.(type) {
		case *Rect, *CurvedRect:
			p := t.Shape.GetPos()
			return Vec{p.X + t.Shape.Width()/2, p.Y + t.Shape.Height()/2}

		default:
			return t.Shape.GetPos()
		}
	}())

	t.Text.Draw(buf)

	if t.isEditing {
		// TODO: Draw cursor
		// NOTE: this is complicated. Text editing may need its own package
	}
}

// Update handles user interactions with the text box.
func (t *TextBox) Update(win *Window) {
	// Enter editing mode if text box is clicked
	if win.MouseButtonState() == LeftClick {
		if t.Shape.IsWithin(win.MouseLocation()) {
			t.isEditing = true
			win.engine.textMutator.Load(t.Text.Text())
			t.selectedCB() // user-defined callback
		} else {
			t.isEditing = false
			t.deselectedCB() // user-defined callback
		}
	}

	if t.isEditing {
		t.SetText(win.engine.textMutator.String())
	}

	if t.Text.body != t.prevText {
		t.modifiedCB()
	}

	t.prevText = t.Text.body
}

// SetSelectedCB sets the callback function which is triggered when the text box is selected.
func (t *TextBox) SetSelectedCB(fn func()) *TextBox {
	t.selectedCB = fn
	return t
}

// SetDeselectedCB sets the callback function which is triggered when the text box is deselected.
func (t *TextBox) SetDeselectedCB(fn func()) *TextBox {
	t.deselectedCB = fn
	return t
}

// SetModifiedCB sets the callback function which is triggered when the text is modified.
func (t *TextBox) SetModifiedCB(fn func()) *TextBox {
	t.modifiedCB = fn
	return t
}

// SetPos sets the position of the text box.
func (t *TextBox) SetPos(pos Vec) *TextBox {
	t.Shape.SetPos(pos)
	return t
}

// Move moves the text box by a given vector.
func (t *TextBox) Move(mov Vec) {
	t.Shape.Move(mov)
	t.Text.Move(mov)
}

// SetCallback configures a callback function to execute every time the text
// in the text box is updated by the user.
func (t *TextBox) SetCallback(callback func()) *TextBox {
	t.modifiedCB = callback
	return t
}

// IsEditing returns whether the text box is in edit mode.
func (t *TextBox) IsEditing() bool {
	return t.isEditing
}

// SetEditing sets whether the text box is in edit mode or not.
func (t *TextBox) SetEditing(editMode bool) *TextBox {
	t.isEditing = editMode
	return t
}

// SetText sets the text to the given string.
func (t *TextBox) SetText(s string) *TextBox {
	t.Text.SetText(s)
	return t
}

// SetTextAlignment sets the alignment of the text relative to the text box.
func (t *TextBox) SetTextAlignment(align Alignment) *TextBox {
	t.Text.SetAlignment(align)
	return t
}

// SetTextOffset manually sets the text's offset, providing the text is in AlignCustom mode.
func (t *TextBox) SetTextOffset(offset Vec) *TextBox {
	t.Text.SetOffset(offset)
	return t
}

// SetTextColour sets the text colour.
func (t *TextBox) SetTextColour(c color.Color) *TextBox {
	t.Text.SetColour(c)
	return t
}

// SetTextFont sets the path fo the .ttf file that is used to generate the text.
func (t *TextBox) SetTextFont(path string) error {
	return t.Text.SetFont(path)
}

// SetTextDPI sets the DPI of the text font.
func (t *TextBox) SetTextDPI(dpi float64) *TextBox {
	t.Text.SetDPI(dpi)
	return t
}

// SetTextSize sets the size of the text font.
func (t *TextBox) SetTextSize(size float64) *TextBox {
	t.Text.SetSize(size)
	return t
}

// SetTextSpacing sets the line spacing of the text.
func (t *TextBox) SetTextSpacing(spacing float64) *TextBox {
	t.Text.SetSpacing(spacing)
	return t
}
