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

	file, fileError := os.Open(me.inputFilePath)
	if fileError != nil {
		log.Fatalf("Failed to open input file: %v", fileError)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, BUFFER_SIZE), BUFFER_SIZE)
	for scanner.Scan() {
		line := scanner.Text()
		me.receiveLine(line)
	}
	if scannerError := scanner.Err(); scannerError != nil {
		log.Fatalf("Error reading input file: %v", scannerError)
	}
}

func (me *App) receiveLine(line string) {
	if "" == me.outputFilePath {
		me.outputFilePath = "./init.sql"
	}
	if nil == me.outputFile {
		outputFile, err := os.Create(me.outputFilePath)
		if err != nil {
			log.Fatalf("Failed to open output file %s: %v", me.outputFilePath, err)
		}
		me.outputFile = outputFile
	}
}

func main() {
	flag.Parse()
	var app = new(App)
	app.inputFilePath = flag.Arg(0)
	app.run()
}
