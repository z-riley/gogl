package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"
	"net"
	"time"

	"github.com/charmbracelet/log"
	"github.com/z-riley/turdgl"
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
	LeftPaddlePos  turdgl.Vec `json:"leftPaddlePos"`
	RightPaddlePos turdgl.Vec `json:"rightPaddlePos"`
	BallPos        turdgl.Vec `json:"ballPos"`
}

// pongServer is the entrypoint for a pong server instance.
func pongServer() {
	go NewServer("0.0.0.0", 3333).Run()

	win, err := turdgl.NewWindow(turdgl.WindowCfg{
		Title:  "Pong Server",
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
	gameServer.paddleLeft = NewPaddle(turdgl.Vec{X: 50, Y: 200})
	gameServer.paddleRight = NewPaddle(turdgl.Vec{X: float64(win.GetConfig().Width) - 50, Y: 200})
	gameServer.ball = NewBall(turdgl.Vec{
		X: float64(win.GetConfig().Width / 2),
		Y: float64(win.GetConfig().Height / 2),
	})

	prevTime := time.Now()
	for win.IsRunning() {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// React to pressed keys
		if win.KeyIsPressed(turdgl.KeyW) {
			gameServer.paddleLeft.MovePos(dirUp, dt, win.Framebuffer)
		}
		if win.KeyIsPressed(turdgl.KeyS) {
			gameServer.paddleLeft.MovePos(dirDown, dt, win.Framebuffer)
		}

		gameServer.ball.Update(dt, win.Framebuffer)
		if turdgl.IsColliding(gameServer.ball.body, gameServer.paddleLeft.body) ||
			turdgl.IsColliding(gameServer.ball.body, gameServer.paddleRight.body) {
			gameServer.ball.velocity.X *= -1
		}

		// Set background colour
		win.SetBackground(color.RGBA{39, 45, 53, 255})

		// Draw shapes
		win.Draw(gameServer.paddleLeft)
		win.Draw(gameServer.paddleRight)
		win.Draw(gameServer.ball)

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
		log.Info("Received from client: " + message)

		var clientUpdate ClientUpdate
		err = json.Unmarshal([]byte(message), &clientUpdate)
		if err != nil {
			log.Error("Failed to unmarshal client update: " + err.Error())
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
			log.Error("Failed to marshal new game state: " + err.Error())
		}

		_, err = client.conn.Write(append(b, ';')) // add delimiter for easier client parsing
		if err != nil {
			log.Error("Failed to reply to client:", err)
		}
	}
}
