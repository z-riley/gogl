package turdgl

import (
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type WindowCfg struct {
	Title         string
	Width, Height int
}

type engine struct {
	running bool
	keys    *keyTracker
}

type Window struct {
	KeyBindings map[sdl.Keycode]func()
	Framebuffer *FrameBuffer
	win         *sdl.Window
	renderer    *sdl.Renderer
	texture     *sdl.Texture
	engine      *engine
	config      WindowCfg
}

func NewWindow(cfg WindowCfg) (*Window, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	w, err := sdl.CreateWindow(
		cfg.Title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(cfg.Width),
		int32(cfg.Height),
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE,
	)
	if err != nil {
		return nil, err
	}

	r, err := sdl.CreateRenderer(w, -1,
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC,
	)
	if err != nil {
		log.Fatalf("failed to create renderer: %s", err)
	}

	t, err := r.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(cfg.Width),
		int32(cfg.Height),
	)
	if err != nil {
		log.Fatalf("failed to create texture: %s", err)
	}

	return &Window{
		KeyBindings: make(map[sdl.Keycode]func()),
		Framebuffer: NewFrameBuffer(cfg.Width, cfg.Height),
		win:         w,
		renderer:    r,
		texture:     t,
		engine: &engine{
			running: true,
			keys:    newKeyTracker(),
		},
		config: cfg,
	}, nil
}

func (w *Window) Destroy() {
	sdl.Quit()
	w.win.Destroy()
	w.renderer.Destroy()
	w.texture.Destroy()
}

func (w *Window) RegisterKeybind(key sdl.Keycode, cb func()) {
	w.KeyBindings[key] = cb
}

func (w *Window) KeyIsPressed(key sdl.Keycode) bool {
	return w.engine.keys.IsPressed(key)
}

func (w *Window) Draw(s Drawable) {
	s.Draw(w.Framebuffer)
}

func (w *Window) Update() {
	// Handle internal events
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			w.engine.running = false
		case *sdl.KeyboardEvent:
			w.engine.keys.eventHandler(e)
		}
	}

	// React to key presses
	for key, fn := range w.KeyBindings {
		if w.engine.keys.IsPressed(key) {
			fn()
		}
	}

	// Present latest frame buffer
	pixels := w.Framebuffer.Bytes()
	w.texture.Update(nil, unsafe.Pointer(&pixels[0]), w.config.Width*pxLen)
	w.renderer.Copy(w.texture, nil, nil)
	w.renderer.Present()
}

func (w *Window) IsRunning() bool {
	return w.engine.running
}

func (w *Window) Quit() {
	w.engine.running = false
}

// GetConfig returns a copy of the window's config.
func (w *Window) GetConfig() WindowCfg {
	return w.config
}

func (w *Window) Width() int {
	width, _ := w.win.GetSize()
	return int(width)
}

func (w *Window) Height() int {
	_, height := w.win.GetSize()
	return int(height)
}

func (w *Window) SetTitle(title string) {
	w.win.SetTitle(title)
}
