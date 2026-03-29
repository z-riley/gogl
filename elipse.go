package gogl

type Ellipse struct {
	Pos   Vec
	w, h  float64
	style Style
}

var _ Shape = (*Ellipse)(nil)

func NewEllipse(width, height float64, pos Vec) *Ellipse {
	return &Ellipse{
		Pos:   pos,
		w:     width,
		h:     height,
		style: DefaultStyle,
	}
}

func (e *Ellipse) Draw(buf *FrameBuffer) {
	a := e.w / 2
	b := e.h / 2

	bbBoxPos := Vec{e.Pos.X - a, e.Pos.Y - b}
	bbox := NewRect(e.w, e.h, bbBoxPos)

	for x := bbox.Pos.X; x <= bbox.Pos.X+bbox.w; x++ {
		for y := bbox.Pos.Y; y <= bbox.Pos.Y+bbox.h; y++ {

			p1 := (x - e.Pos.X) * (x - e.Pos.X) / (a * a)
			p2 := (y - e.Pos.Y) * (y - e.Pos.Y) / (b * b)

			if p1+p2 <= 1 {
				buf.SetPixel(int(x), int(y), NewPixel(e.style.Colour))
			}
		}
	}
}

func (e *Ellipse) GetPos() Vec {
	return e.Pos
}

func (e *Ellipse) GetStyle() Style {
	return e.style
}

func (e *Ellipse) Height() float64 {
	return e.h
}

func (e *Ellipse) SetHeight(px float64) *Ellipse {
	e.h = max(px, 0)
	return e
}

func (e *Ellipse) Move(px Vec) {
	e.Pos = Add(e.Pos, px)
}

func (e *Ellipse) SetPos(pos Vec) {
	e.Pos = pos
}

func (e *Ellipse) String() string {
	return "ellipse"
}

func (e *Ellipse) Width() float64 {
	return e.w
}

func (e *Ellipse) SetWidth(px float64) *Ellipse {
	e.w = max(px, 0)
	return e
}

func (e *Ellipse) SetStyle(style Style) *Ellipse {
	e.style = style
	return e
}
