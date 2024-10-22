package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"net"
	"time"

	"github.com/charmbracelet/log"
	"github.com/z-riley/turdgl"
)

var gameClient Game

// pongClient is the entrypoint for a pong client instance.
func pongClient() {
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
	win, err := turdgl.NewWindow(turdgl.WindowCfg{
		Title:  "Pong Client",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// For measuring FPS
	frames := 0
	second := time.Tick(time.Second)

	// Shapes
	gameClient.paddleLeft = NewPaddle(turdgl.Vec{X: 50, Y: 200})
	gameClient.paddleRight = NewPaddle(turdgl.Vec{X: float64(win.GetConfig().Width) - 50, Y: 200})
	gameClient.ball = NewBall(turdgl.Vec{
		X: float64(win.GetConfig().Width / 2),
		Y: float64(win.GetConfig().Height / 2),
	})

	prevTime := time.Now()
	for win.IsRunning() {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// React to pressed keys
		if win.KeyIsPressed(turdgl.KeyUp) {
			gameClient.paddleRight.MovePos(dirUp, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(turdgl.KeyDown) {
			gameClient.paddleRight.MovePos(dirDown, dt, win.Framebuffer)
		}

		// Ball movement
		gameClient.ball.Update(dt, win.Framebuffer)
		if turdgl.IsColliding(gameClient.ball.body, gameClient.paddleLeft.body) ||
			turdgl.IsColliding(gameClient.ball.body, gameClient.paddleRight.body) {
			gameClient.ball.velocity.X *= -1
		}

		// Send paddle move update to server
		err := movePaddle(conn, gameClient.paddleRight.body.GetPos())
		if err != nil {
			log.Error("failed to move paddle: " + err.Error())
		}

		// Set background colour
		win.SetBackground(color.RGBA{39, 45, 53, 255})

		// Draw shapes
		win.Draw(gameClient.paddleLeft)
		win.Draw(gameClient.paddleRight)
		win.Draw(gameClient.ball)

		win.Update()

		// Count FPS
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", win.GetConfig().Title, frames))
			frames = 0
		default:
		}
	}
}

// clientUpdate contains data to send to the server to update the game state.
type ClientUpdate struct {
	RightPaddlePos turdgl.Vec `json:"rightPaddlePos"`
}

// movePaddle sends a request to move the paddle to the game server.
func movePaddle(conn *net.TCPConn, newPos turdgl.Vec) error {
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
