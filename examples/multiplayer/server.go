package main

import (
	"bufio"
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

// Game contains the game's assets.
type Game struct {
	paddleLeft  *paddle
	paddleRight *paddle
	ball        *ball
}

var gameServer Game

// GameState contains the game's current state.
type GameState struct {
	LeftPaddlePos  tgl.Vec `json:"leftPaddlePos"`
	RightPaddlePos tgl.Vec `json:"rightPaddlePos"`
	BallPos        tgl.Vec `json:"ballPos"`
}

var (
	frames = 0
	second = time.Tick(time.Second)
)

// pongServer is the entrypoint for a pong server instance.
func pongServer() {
	pixelgl.Run(runServer)
}

func runServer() {
	go NewServer("0.0.0.0", 3333).Run()

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
	gameServer.paddleLeft = NewPaddle(tgl.Vec{X: 50, Y: 200})
	gameServer.paddleRight = NewPaddle(tgl.Vec{X: float64(screenWidth) - 50, Y: 200})
	gameServer.ball = NewBall(tgl.Vec{X: win.Bounds().Center().X, Y: win.Bounds().Center().Y})

	prevTime := time.Now()
	for {
		dt := time.Since(prevTime)
		prevTime = time.Now()

		// Handle user input
		if win.Closed() || win.JustPressed(pixelgl.KeyLeftControl) || win.JustPressed(pixelgl.KeyEscape) {
			return
		}
		if win.Pressed(pixelgl.KeyW) {
			gameServer.paddleLeft.MovePos(dirUp, dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyS) {
			gameServer.paddleLeft.MovePos(dirDown, dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyUp) {
			gameServer.paddleRight.MovePos(dirUp, dt, framebuf)
		}
		if win.Pressed(pixelgl.KeyDown) {
			gameServer.paddleRight.MovePos(dirDown, dt, framebuf)
		}
		gameServer.ball.Update(dt, framebuf)
		if tgl.IsColliding(gameServer.ball.body, gameServer.paddleLeft.body) ||
			tgl.IsColliding(gameServer.ball.body, gameServer.paddleRight.body) {
			gameServer.ball.velocity.X *= -1
		}

		// Adjust frame buffer size if window size changes
		if !prevSize.Eq(win.Bounds().Size()) {
			framebuf = tgl.NewFrameBuffer(win.Canvas().Texture().Width(), win.Canvas().Texture().Height())
		}

		// Set background colour
		framebuf.SetBackground(color.RGBA{39, 45, 53, 255})

		// Modify frame buffer
		gameServer.paddleLeft.Draw(framebuf)
		gameServer.paddleRight.Draw(framebuf)
		gameServer.ball.Draw(framebuf)

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
