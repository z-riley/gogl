package turdgl

import (
	"image/color"
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// WindowCfg contains adjustable configuration for a window.
type WindowCfg struct {
	Title         string
	Width, Height int
}

// engine contains constructs used to execute background logic.
type engine struct {
	foregroundDrawQueue []Drawable
	backgroundDrawQueue []Drawable
	running             bool
	keyTracker          *keyTracker
	textMutator         *textMutator
}

// newEngine constructs a new turdgl engine.
func newEngine() *engine {
	return &engine{
		running:     true,
		keyTracker:  newKeyTracker(),
		textMutator: newTextTracker(),
	}
}

// Window represents an OS Window.
type Window struct {
	Framebuffer *FrameBuffer
	win         *sdl.Window
	renderer    *sdl.Renderer
	texture     *sdl.Texture
	engine      *engine
	config      WindowCfg
}

// NewWindow constructs a new window according to the provided configuration.
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
		Framebuffer: NewFrameBuffer(cfg.Width, cfg.Height),
		win:         w,
		renderer:    r,
		texture:     t,
		engine:      newEngine(),
		config:      cfg,
	}, nil
}

// Destroy deallocates the window's resources. Call it at the end of your application.
func (w *Window) Destroy() {
	sdl.Quit()
	w.win.Destroy()
	w.renderer.Destroy()
	w.texture.Destroy()
}

// Draw is an alias for DrawForeground.
func (w *Window) Draw(s Drawable) {
	w.DrawForeground(s)
}

// DrawForeground draws a shape to the foreground layer.
func (w *Window) DrawForeground(s Drawable) {
	w.engine.foregroundDrawQueue = append(w.engine.foregroundDrawQueue, s)
}

// DrawBackground draws a shape to the background layer.
func (w *Window) DrawBackground(s Drawable) {
	w.engine.backgroundDrawQueue = append(w.engine.backgroundDrawQueue, s)
}

// RegisterKeybind sets a callback function which is executed when a key is pressed.
// The callback can be executed always while the key is pressed, on press, or on release.
func (w *Window) RegisterKeybind(key sdl.Keycode, mode KeybindMode, cb func()) {
	w.engine.keyTracker.registerKeybind(key, mode, cb)
}

// RegisterKeybind removes a keybind for a specific key mode combination.
func (w *Window) UnregisterKeybind(key sdl.Keycode, mode KeybindMode) {
	w.engine.keyTracker.unregisterKeybind(key, mode)
}

// DropKeybinds unregisters all keybinds.
func (w *Window) DropKeybinds() {
	w.engine.keyTracker.dropKeybinds()
}

// KeyIsPressed returns whether a given key is currently pressed.
func (w *Window) KeyIsPressed(key sdl.Keycode) bool {
	return w.engine.keyTracker.isPressed(key)
}

// SetBackground sets the background to a uniform colour.
func (w *Window) SetBackground(c color.Color) {
	// Alpha must be 255 for bloom effects to work
	r, g, b, _ := RGBA8(c)
	w.Framebuffer.Fill(color.RGBA{r, g, b, 255})
}

func (w *Window) Update() {
	// Handle internal events
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			w.engine.running = false
		case *sdl.KeyboardEvent:
			w.engine.keyTracker.handleEvent(e)
			w.engine.textMutator.handleEvent(e)
		case *sdl.TextInputEvent:
			w.engine.textMutator.Append(e.GetText())
		}
	}

	// React to key presses
	w.engine.keyTracker.update()

	// Draw shapes to frame buffer
	for _, shape := range w.engine.backgroundDrawQueue {
		shape.Draw(w.Framebuffer)
	}
	w.engine.backgroundDrawQueue = nil
	for _, shape := range w.engine.foregroundDrawQueue {
		shape.Draw(w.Framebuffer)
	}
	w.engine.foregroundDrawQueue = nil

	// Render latest frame buffer to window
	pixels := w.Framebuffer.Bytes()
	w.texture.Update(nil, unsafe.Pointer(&pixels[0]), w.config.Width*pxLen)
	w.renderer.Copy(w.texture, nil, nil)
	w.renderer.Present()
}

// IsRunning returns true while the window is running.
func (w *Window) IsRunning() bool {
	return w.engine.running
}

// Quit signals the window to exit.
func (w *Window) Quit() {
	w.engine.running = false
}

// GetConfig returns a copy of the window's config.
func (w *Window) GetConfig() WindowCfg {
	return w.config
}

// MouseLocation returns the location of the mouse cursor, relative to
// the origin of the window.
func (w *Window) MouseLocation() Vec {
	x, y, _ := sdl.GetMouseState()
	return Vec{X: float64(x), Y: float64(y)}
}

// MouseState represents the state of the mouse buttons.
type MouseState int

const (
	NoClick           MouseState = 0
	LeftClick         MouseState = 1
	RightClick        MouseState = 4
	LeftAndRightClick MouseState = 5
)

// MouseButtonState returns the current state of the mouse buttons.
func (w *Window) MouseButtonState() MouseState {
	_, _, s := sdl.GetMouseState()
	return MouseState(s)
}

// Width returns the width of the window in pixels.
func (w *Window) Width() int {
	width, _ := w.win.GetSize()
	return int(width)
}

// Height returns the height of the window in pixels.
func (w *Window) Height() int {
	_, height := w.win.GetSize()
	return int(height)
}

// SetTitle sets the title of the window to the provided string.
func (w *Window) SetTitle(title string) {
	w.win.SetTitle(title)
}
