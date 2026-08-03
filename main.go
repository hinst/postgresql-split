package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

type App struct {
	inputFile string
}

func (me *App) run() {
	if me.inputFile == "" {
		log.Fatal("Input file is not specified")
	}

	file, err := os.Open(me.inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer file.Close()

	const bufSize = 32 * 1024 * 1024 // 32 MB
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, bufSize), bufSize)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		// TODO: insert logic here to process each line
	}
}

func main() {
	flag.Parse()
	var app = new(App)
	app.inputFile = flag.Arg(0)
	app.run()
}
