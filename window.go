package turdgl

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// WindowCfg contains adjustable configuration for a window.
type WindowCfg struct {
	// Title is the title of the window.
	Title string
	// Width, Height specifies the dimensions of the window, in pixels.
	Width, Height int
	// Icon is the image used for the window. Default if nil.
	Icon *os.File
	// Resizable can be set to true to allow the window to be resizable.
	Resizable bool
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

	var resizableFlag uint32 = sdl.WINDOW_RESIZABLE
	if cfg.Resizable {
		resizableFlag = 0
	}

	w, err := sdl.CreateWindow(
		cfg.Title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(cfg.Width),
		int32(cfg.Height),
		sdl.WINDOW_SHOWN|resizableFlag,
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

	// (Optional) set window icon
	if cfg.Icon != nil {
		iconSurface, err := img.Load(cfg.Icon.Name())
		if err != nil {
			return nil, fmt.Errorf("could not load icon: %w", err)
		}
		defer iconSurface.Free()
		w.SetIcon(iconSurface)
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
	_ = w.win.Destroy()
	_ = w.renderer.Destroy()
	_ = w.texture.Destroy()
}

// Draw draws a shape to the window.
func (w *Window) Draw(s Drawable) {
	w.engine.drawQueue = append(w.engine.drawQueue, s)
}

// RegisterKeybind sets a callback function which is executed when a key is pressed.
// The callback can be executed always while the key is pressed, on press, or on release.
func (w *Window) RegisterKeybind(key sdl.Keycode, mode KeybindMode, callback func()) {
	w.engine.keyTracker.registerKeybind(key, mode, callback)
}

// RegisterKeybind removes a keybind for a specific key mode combination.
func (w *Window) UnregisterKeybind(key sdl.Keycode, mode KeybindMode) {
	w.engine.keyTracker.unregisterKeybind(key, mode)
}

// SetMouseScrollCallback sets a callback to be executed on mouse scroll events.
func (w *Window) SetMouseScrollCallback(cb MouseScrollCallback) {
	w.engine.mouseScrollTracker.Callback = cb
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
		case *sdl.MouseWheelEvent:
			w.engine.mouseScrollTracker.handleEvent(e)
		}
	}

	// React to key presses
	w.engine.keyTracker.update()

	// Draw shapes to frame buffer
	for _, shape := range w.engine.drawQueue {
		shape.Draw(w.Framebuffer)
	}
	w.engine.drawQueue = nil

	// Render latest frame buffer to window
	pixels := w.Framebuffer.Bytes()
	if err := w.texture.Update(nil, unsafe.Pointer(&pixels[0]), w.config.Width*pxLen); err != nil {
		fmt.Println("SDL texture failed to update from frame buffer")
	}
	if err := w.renderer.Copy(w.texture, nil, nil); err != nil {
		fmt.Println("SDL renderer failed to copy texture")
	}
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

// String returns the mouse state as a string.
func (m MouseState) String() string {
	switch m {
	case NoClick:
		return "no click"
	case LeftClick:
		return "left click"
	case RightClick:
		return "right click"
	case LeftAndRightClick:
		return "left and right click"
	default:
		return "invalid"
	}
}

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

// engine contains constructs used to execute background logic.
type engine struct {
	drawQueue          []Drawable
	running            bool
	keyTracker         *keyTracker
	mouseScrollTracker *mouseScrollHandler
	textMutator        *textMutator
}

// newEngine constructs a new turdgl engine.
func newEngine() *engine {
	return &engine{
		drawQueue:          []Drawable{},
		running:            true,
		keyTracker:         newKeyTracker(),
		mouseScrollTracker: newMouseScrollHandler(),
		textMutator:        newTextTracker(),
	}
}
