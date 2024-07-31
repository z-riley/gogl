package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"net"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/pixelgl"
	tgl "github.com/zac460/turdgl"
)

// pongClient is the entrypoint for a pong client instance.
func pongClient() {
	pixelgl.Run(runClient)
}

var gameClient Game

func runClient() {
	// Set up TCP client
	const addr, port = "localhost", 3333
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal("ResolveTCPAddr failed: " + err.Error())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("Dial failed: " + err.Error())
	}
	defer conn.Close()

	// Run the game locally
	gameLoop(conn)
}

// gameLoop is the main loop for the client.
func gameLoop(conn *net.TCPConn) {

	cfg := pixelgl.WindowConfig{
		Title:     "Pong",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Screen state
	screenWidth := win.Canvas().Texture().Width()
	screenHeight := win.Canvas().Texture().Height()
	framebuf := tgl.NewFrameBuffer(screenWidth, screenHeight)
	prevSize := win.Bounds().Size()

	// Shapes
	gameClient.paddleLeft = NewPaddle(tgl.Vec{X: 50, Y: 200})
	gameClient.paddleRight = NewPaddle(tgl.Vec{X: float64(screenWidth) - 50, Y: 200})
	gameClient.ball = NewBall(tgl.Vec{X: win.Bounds().Center().X, Y: win.Bounds().Center().Y})

	prevTime := time.Now()
	for {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// Handle user input
		if win.Closed() || win.JustPressed(pixelgl.KeyLeftControl) || win.JustPressed(pixelgl.KeyEscape) {
			return
		}
		if win.Pressed(pixelgl.KeyW) {
			gameClient.paddleLeft.MovePos(dirUp, dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyS) {
			gameClient.paddleLeft.MovePos(dirDown, dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyUp) {
			gameClient.paddleRight.MovePos(dirUp, dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyDown) {
			gameClient.paddleRight.MovePos(dirDown, dt, framebuf)
		}
		gameClient.ball.Update(dt, framebuf)
		if tgl.IsColliding(gameClient.ball.body, gameClient.paddleLeft.body) ||
			tgl.IsColliding(gameClient.ball.body, gameClient.paddleRight.body) {
			gameClient.ball.velocity.X *= -1
		}

		// Send paddle move update to server
		err := movePaddle(conn, gameClient.paddleRight.body.GetPos())
		if err != nil {
			log.Error("failed to move paddle: " + err.Error())
		}

		// Adjust frame buffer size if window size changes
		if !prevSize.Eq(win.Bounds().Size()) {
			framebuf = tgl.NewFrameBuffer(win.Canvas().Texture().Width(), win.Canvas().Texture().Height())
		}

		// Set background colour
		framebuf.SetBackground(color.RGBA{39, 45, 53, 255})

		// Modify frame buffer
		gameClient.paddleLeft.Draw(framebuf)
		gameClient.paddleRight.Draw(framebuf)
		gameClient.ball.Draw(framebuf)

		// Render screen
		win.Canvas().SetPixels(framebuf.Bytes())
		win.Update()

		// Count FPS
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

// clientUpdate contains data to send to the server to update the game state.
type ClientUpdate struct {
	RightPaddlePos tgl.Vec `json:"rightPaddlePos"`
}

// movePaddle sends a request to move the paddle to the game server.
func movePaddle(conn *net.TCPConn, newPos tgl.Vec) error {
	b, err := json.Marshal(ClientUpdate{RightPaddlePos: newPos})
	if err != nil {
		return fmt.Errorf("failed to marshal client update %w", err)
	}

	// Add delimiter to message server knows when to stop reading
	b = append(b, ';')

	_, err = conn.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write to server: %w", err)
	}

	resp := make([]byte, 512)
	_, err = conn.Read(resp)
	if err != nil {
		return fmt.Errorf("failed to read server reply to server failed: %w", err)
	}
	// Remove delimiter from message
	i := bytes.IndexByte(resp, ';')
	resp = bytes.Trim(resp, ";")[:i]

	log.Info("New game state: " + string(resp))

	var gs GameState
	err = json.Unmarshal(resp, &gs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal game state from server: %w", err)
	}

	// Update local game with latest game state from server
	gameClient.paddleLeft.body.SetPos(gs.LeftPaddlePos)
	gameClient.paddleRight.body.SetPos(gs.RightPaddlePos)
	gameClient.ball.body.SetPos(gs.BallPos)

	return nil
}
