package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

type App struct {
	inputFilePath  string
	outputFilePath string
	outputFile     *os.File
}

const BUFFER_SIZE = 32 * 1024 * 1024 // 32 MB

func (me *App) run() {
	if me.inputFilePath == "" {
		log.Fatal("Input file is not specified")
	}

	file, err := os.Open(me.inputFilePath)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, BUFFER_SIZE), BUFFER_SIZE)
	for scanner.Scan() {
		line := scanner.Text()
		me.receiveLine(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}
}

func (me *App) receiveLine(line string) {
	if "" == me.outputFilePath {
		me.outputFilePath = "./init.sql"
	}
	if nil == me.outputFile {

	}
}

func main() {
	flag.Parse()
	var app = new(App)
	app.inputFilePath = flag.Arg(0)
	app.run()
}
