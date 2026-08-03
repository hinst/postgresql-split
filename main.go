package main

import (
	"flag"
	"log"
)

type App struct {
	inputFile string
}

func (me *App) run() {
	if me.inputFile == "" {
		log.Fatal("Input file is not specified")
	}
}

func main() {
	flag.Parse()
	var app = new(App)
	app.inputFile = flag.Arg(0)
	app.run()
}
