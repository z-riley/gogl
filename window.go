package gogl

import (
	"fmt"
	"image/color"
	"os"
	"unsafe"

	"github.com/jupiterrider/purego-sdl3/img"
	"github.com/jupiterrider/purego-sdl3/sdl"
)

// WindowCfg contains adjustable configuration for a window.
type WindowCfg struct {
	// Title is the title of the window.
	Title string
	// Width, Height specifies the dimensions of the window, in pixels.
	Width, Height int
	// Icon is the image used for the window icon. Default if nil.
	Icon *os.File
	// Resizable can be set to true to allow the window to be resizable.
	Resizable bool
}

// Window represents an OS Window.
type Window struct {
	Framebuffer *FrameBuffer

	win      *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture

	engine *engine
	config WindowCfg
}

// NewWindow constructs a new window according to the provided configuration.
//
// Call Destroy to deallocate the window.
func NewWindow(cfg WindowCfg) (*Window, error) {
	if ok := sdl.Init(sdl.InitVideo); !ok {
		return nil, fmt.Errorf("failed to init sdl3: %s", sdl.GetError())
	}

	var resizableFlag3 sdl.WindowFlags
	if cfg.Resizable {
		resizableFlag3 = sdl.WindowResizable
	}

	w := sdl.CreateWindow(
		cfg.Title,
		int32(cfg.Width),
		int32(cfg.Height),
		resizableFlag3,
	)
	if w == nil {
		return nil, fmt.Errorf("failed to create sdl3 window: %s", sdl.GetError())
	}
	if !sdl.StartTextInput(w) {
		return nil, fmt.Errorf("failed to start text input: %s", sdl.GetError())
	}

	r := sdl.CreateRenderer(w, "")
	if r == nil {
		return nil, fmt.Errorf("failed to create sdl3 renderer: %s", sdl.GetError())
	}

	t := sdl.CreateTexture(
		r,
		sdl.PixelFormatRGBA8888,
		sdl.TextureAccessStreaming,
		int32(cfg.Width),
		int32(cfg.Height),
	)
	if t == nil {
		return nil, fmt.Errorf("failed to create sdl3 texture: %s", sdl.GetError())
	}

	// (Optional) set window icon
	if cfg.Icon != nil {
		iconSurface := img.Load(cfg.Icon.Name())
		sdl.SetWindowIcon(w, iconSurface)
	}

	return &Window{
		Framebuffer: NewFrameBuffer(cfg.Width, cfg.Height),

		win:      w,
		renderer: r,
		texture:  t,

		engine: newEngine(),
		config: cfg,
	}, nil
}

// Destroy deallocates the window's resources. Call it at the end of your application.
func (w *Window) Destroy() {
	sdl.DestroyTexture(w.texture)
	sdl.DestroyRenderer(w.renderer)
	sdl.DestroyWindow(w.win)
	sdl.Quit()
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
	var event sdl.Event
	for sdl.PollEvent(&event) {
		switch event.Type() {
		case sdl.EventQuit:
			w.engine.running = false
		case sdl.EventKeyDown, sdl.EventKeyUp:
			e := event.Key()
			w.engine.keyTracker.handleEvent(e)
			w.engine.textMutator.handleEvent(e)
		case sdl.EventTextInput:
			e := event.Text()
			w.engine.textMutator.Append(e.Text())
		case sdl.EventMouseWheel:
			e := event.Wheel()
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

	// Render to SDL window
	pixels := w.Framebuffer.Bytes()
	// Pitch is bytes per row: width * 4 bytes per pixel (RGBA8888)
	if !sdl.UpdateTexture(w.texture, nil, unsafe.Pointer(&pixels[0]), int32(w.config.Width*pxLen)) {
		fmt.Println("failed to update texture:", sdl.GetError())
	}
	if !sdl.RenderTexture(w.renderer, w.texture, nil, nil) {
		fmt.Println("failed to render texture:", sdl.GetError())
	}
	if !sdl.RenderPresent(w.renderer) {
		fmt.Println("failed to present render:", sdl.GetError())
	}
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
	var x, y float32
	_ = sdl.GetMouseState(&x, &y)
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
	return MouseState(sdl.GetMouseState(nil, nil))
}

// Width returns the width of the window in pixels.
func (w *Window) Width() int {
	var width int32
	if !sdl.GetWindowSize(w.win, &width, nil) {
		fmt.Println("failed to get window size: " + sdl.GetError())
	}
	return int(width)
}

// Height returns the height of the window in pixels.
func (w *Window) Height() int {
	var height int32
	if !sdl.GetWindowSize(w.win, nil, &height) {
		fmt.Println("failed to get window size: " + sdl.GetError())
	}
	return int(height)
}

// SetTitle sets the title of the window to the provided string.
func (w *Window) SetTitle(title string) {
	if !sdl.SetWindowTitle(w.win, title) {
		fmt.Println("failed to set window title: " + sdl.GetError())
	}
}

// engine contains constructs used to execute background logic.
type engine struct {
	drawQueue          []Drawable
	running            bool
	keyTracker         *keyTracker
	mouseScrollTracker *mouseScrollHandler
	textMutator        *textMutator
}

// newEngine constructs a new gogl engine.
func newEngine() *engine {
	return &engine{
		drawQueue:          []Drawable{},
		running:            true,
		keyTracker:         newKeyTracker(),
		mouseScrollTracker: newMouseScrollHandler(),
		textMutator:        newTextTracker(),
	}
}
