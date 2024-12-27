package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net"
	"time"

	"github.com/z-riley/gogl"
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
	win, err := gogl.NewWindow(gogl.WindowCfg{
		Title:  "Pong Client",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// Initialise shapes
	gameClient.paddleLeft = NewPaddle(gogl.Vec{X: 50, Y: 200})
	gameClient.paddleRight = NewPaddle(gogl.Vec{X: float64(win.GetConfig().Width) - 50, Y: 200})
	gameClient.ball = NewBall(gogl.Vec{
		X: float64(win.GetConfig().Width / 2),
		Y: float64(win.GetConfig().Height / 2),
	})

	win.RegisterKeybind(gogl.KeyEscape, gogl.KeyPress, func() { win.Quit() })

	prevTime := time.Now()
	for win.IsRunning() {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// React to pressed keys
		if win.KeyIsPressed(gogl.KeyUp) {
			gameClient.paddleRight.MovePos(dirUp, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(gogl.KeyDown) {
			gameClient.paddleRight.MovePos(dirDown, dt, win.Framebuffer)
		}

		// Ball movement
		gameClient.ball.Update(dt, win.Framebuffer)
		if gogl.IsColliding(gameClient.ball.body, gameClient.paddleLeft.body) ||
			gogl.IsColliding(gameClient.ball.body, gameClient.paddleRight.body) {
			gameClient.ball.velocity.X *= -1
		}

		// Send paddle move update to server
		err := movePaddle(conn, gameClient.paddleRight.body.GetPos())
		if err != nil {
			log.Println("failed to move paddle:", err)
		}

		win.SetBackground(color.RGBA{39, 45, 53, 255})

		// Draw shapes
		win.Draw(gameClient.paddleLeft)
		win.Draw(gameClient.paddleRight)
		win.Draw(gameClient.ball)

		win.Update()
	}
}

// clientUpdate contains data to send to the server to update the game state.
type ClientUpdate struct {
	RightPaddlePos gogl.Vec `json:"rightPaddlePos"`
}

// movePaddle sends a request to move the paddle to the game server.
func movePaddle(conn *net.TCPConn, newPos gogl.Vec) error {
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
