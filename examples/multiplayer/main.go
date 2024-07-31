package main

import (
	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/log"
)

var args struct {
	Side string `arg:"-s, --side" help:"\"left\" or \"right\" side"`
}

func main() {
	log.SetLevel(log.DebugLevel)
	arg.MustParse(&args)
	if args.Side == "left" {
		log.Info("Left side selected. Running app in server mode")
		pongServer()
	} else {
		log.Info("Right side selected. Running app in client mode")
		pongClient()
	}
}
