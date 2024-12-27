package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net"
	"time"

	"github.com/z-riley/gogl"
)

// Game contains the game's assets.
type Game struct {
	paddleLeft  *paddle
	paddleRight *paddle
	ball        *ball
}

var gameServer Game

// GameState contains the game's current state.
type GameState struct {
	LeftPaddlePos  gogl.Vec `json:"leftPaddlePos"`
	RightPaddlePos gogl.Vec `json:"rightPaddlePos"`
	BallPos        gogl.Vec `json:"ballPos"`
}

// pongServer is the entrypoint for a pong server instance.
func pongServer() {
	go NewServer("0.0.0.0", 3333).Run()

	win, err := gogl.NewWindow(gogl.WindowCfg{
		Title:  "Pong Server",
		Width:  1024,
		Height: 768,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// Initialise shapes
	gameServer.paddleLeft = NewPaddle(gogl.Vec{X: 50, Y: 200})
	gameServer.paddleRight = NewPaddle(gogl.Vec{X: float64(win.GetConfig().Width) - 50, Y: 200})
	gameServer.ball = NewBall(gogl.Vec{
		X: float64(win.GetConfig().Width / 2),
		Y: float64(win.GetConfig().Height / 2),
	})

	win.RegisterKeybind(gogl.KeyEscape, gogl.KeyPress, func() { win.Quit() })

	prevTime := time.Now()
	for win.IsRunning() {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// React to pressed keys
		if win.KeyIsPressed(gogl.KeyW) {
			gameServer.paddleLeft.MovePos(dirUp, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(gogl.KeyS) {
			gameServer.paddleLeft.MovePos(dirDown, dt, win.Framebuffer)
		}

		gameServer.ball.Update(dt, win.Framebuffer)
		if gogl.IsColliding(gameServer.ball.body, gameServer.paddleLeft.body) ||
			gogl.IsColliding(gameServer.ball.body, gameServer.paddleRight.body) {
			gameServer.ball.velocity.X *= -1
		}

		// Set background colour
		win.SetBackground(color.RGBA{39, 45, 53, 255})

		// Draw shapes
		win.Draw(gameServer.paddleLeft)
		win.Draw(gameServer.paddleRight)
		win.Draw(gameServer.ball)

		win.Update()
	}
}

type Server struct {
	host string
	port int
}

func NewServer(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}
		go client.handleRequest()
	}
}

type Client struct {
	conn net.Conn
}

func (client *Client) handleRequest() {
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString(';')
		if err != nil {
			client.conn.Close()
			return
		}

		// Remove delimiter from message
		message = message[:len(message)-1]

		var clientUpdate ClientUpdate
		err = json.Unmarshal([]byte(message), &clientUpdate)
		if err != nil {
			log.Println("Failed to unmarshal client update:", err)
		}

		// Update game state from client data
		gameServer.paddleRight.body.SetPos(clientUpdate.RightPaddlePos)

		// Reply with current game state
		b, err := json.Marshal(GameState{
			LeftPaddlePos:  gameServer.paddleLeft.body.GetPos(),
			RightPaddlePos: gameServer.paddleRight.body.GetPos(),
			BallPos:        gameServer.ball.body.GetPos(),
		})
		if err != nil {
			log.Println("Failed to marshal new game state:", err)
		}

		_, err = client.conn.Write(append(b, ';')) // add delimiter for easier client parsing
		if err != nil {
			log.Println("Failed to reply to client:", err)
		}
	}
}
