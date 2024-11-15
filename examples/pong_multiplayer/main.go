package main

import (
	"log"
	"os"
)

func main() {
	const helpMsg = "Must specify arg \"left\" or \"right\""
	if len(os.Args[1:]) == 0 {
		log.Fatal(helpMsg)
	}

	switch os.Args[1] {
	case "left":
		log.Println("Left side selected. Running app in server mode...")
		pongServer()
	case "right":
		log.Println("Right side selected. Running app in client mode...")
		pongClient()
	default:
		log.Fatal(helpMsg)
	}
}
