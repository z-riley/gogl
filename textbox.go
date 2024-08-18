package turdgl

import (
	"fmt"
	"image/color"
	"time"
)

// TextBox is a shape that can be typed in.
type TextBox struct {
	Shape        hoverable // the base shape the text box is built on
	Body         *Text     // the text in the text box
	selectedCB   func()    // callback function for when the text box is selected (optional)
	deselectedCB func()    // callback function for when the text box is deselected (optional)
	modifiedCB   func()    // callback function for when the text changes (optional)

	isEditing bool
	prevText  string
}

// NewTextBox constructs a new text box from a shape.
func NewTextBox(shape hoverable, fontPath string) *TextBox {
	return &TextBox{
		Shape:        shape,
		Body:         NewText("", shape.GetPos(), fontPath),
		selectedCB:   func() {},
		deselectedCB: func() {},
		modifiedCB:   func() {},
	}
}

// Draw draws the text box onto the frame buffer.
func (t *TextBox) Draw(buf *FrameBuffer) {
	t.Shape.Draw(buf)

	// Align to centre of underlying shape
	t.Body.SetPos(func() Vec {
		switch t.Shape.(type) {
		case *Rect:
			p := t.Shape.GetPos()
			return Vec{p.X + t.Shape.Width()/2, p.Y + t.Shape.Height()/2}
		default:
			return t.Shape.GetPos()
		}
	}())

	t.Body.Draw(buf)

	if t.isEditing {
		// Draw cursor
		// TODO: this is complicated. Text editing should ideally be a separate package
	}
}

// Update executes the callback function if the text has changed.
func (t *TextBox) Update(win *Window) {
	// Enter editing mode if text box is clicked
	if win.MouseButtonState() == LeftClick {
		if t.Shape.IsWithin(win.MouseLocation()) {
			t.isEditing = true
			win.engine.textMutator.Load(t.Body.Text())
			t.selectedCB() // user-defined callback
			fmt.Println(time.Now().Nanosecond(), t.Shape, "select")
		} else {
			t.isEditing = false
			fmt.Println(win.engine.textMutator.buffer)
			t.deselectedCB() // user-defined callback
			fmt.Println(time.Now().Nanosecond(), t.Shape, "deselect")
		}
	}

	if t.isEditing {
		t.SetText(win.engine.textMutator.String())
	}

	if t.Body.body != t.prevText {
		t.modifiedCB()
	}

	t.prevText = t.Body.body
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

// Move moves the text box by a given vector.
func (t *TextBox) Move(mov Vec) {
	t.Shape.Move(mov)
	t.Body.Move(mov)
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
	t.Body.SetText(s)
	return t
}

// SetTextAlignment sets the alignment of the text relative to the text box.
func (t *TextBox) SetTextAlignment(align Alignment) *TextBox {
	t.Body.SetAlignment(align)
	return t
}

// SetTextOffset manually sets the text's offset, providing the text is in AlignCustom mode.
func (t *TextBox) SetTextOffset(offset Vec) *TextBox {
	t.Body.SetOffset(offset)
	return t
}

// SetTextColour sets the text colour.
func (t *TextBox) SetTextColour(c color.Color) *TextBox {
	t.Body.SetColour(c)
	return t
}

// SetTextFont sets the path fo the .ttf file that is used to generate the text.
func (t *TextBox) SetTextFont(path string) *TextBox {
	t.Body.SetFont(path)
	return t
}

// SetTextDPI sets the DPI of the text font.
func (t *TextBox) SetTextDPI(dpi float64) *TextBox {
	t.Body.SetDPI(dpi)
	return t
}

// SetTextSize sets the size of the text font.
func (t *TextBox) SetTextSize(size float64) *TextBox {
	t.Body.SetSize(size)
	return t
}

// SetTextSpacing sets the line spacing of the text.
func (t *TextBox) SetTextSpacing(spacing float64) *TextBox {
	t.Body.SetSpacing(spacing)
	return t
}

// SetTextMaskSize sets the size of the mask used to generate the text.
func (t *TextBox) SetTextMaskSize(w, h int) *TextBox {
	t.Body.SetMaskSize(w, h)
	return t
}
