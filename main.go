package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
)

type App struct {
	inputFilePath  string
	outputFilePath string
	outputFile     *os.File
	ongoingHeader  []string
}

const BUFFER_SIZE = 32 * 1024 * 1024

func (me *App) run() {
	if me.inputFilePath == "" {
		log.Fatal("Input file is not specified")
	}

	file := AssertResultError(os.Open(me.inputFilePath))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, BUFFER_SIZE), BUFFER_SIZE)
	for scanner.Scan() {
		line := scanner.Text()
		me.receiveLine(line)
	}
	me.flushOngoingHeader()
	AssertError(scanner.Err())
}

func (me *App) writeOutput(output *os.File, content string) {
	AssertResultError(output.WriteString(content + "\n"))
}

func (me *App) flushOngoingHeader() {
	for _, header := range me.ongoingHeader {
		me.writeOutput(me.outputFile, header)
	}
	me.ongoingHeader = nil
}

func (me *App) receiveLine(line string) {
	if "" == me.outputFilePath {
		me.outputFilePath = "./init.sql"
	}
	if nil == me.outputFile {
		outputFile := AssertResultError(os.Create(me.outputFilePath))
		me.outputFile = outputFile
	}

	if strings.HasPrefix(line, "-- ") {
		me.ongoingHeader = append(me.ongoingHeader, line)
	} else {
		me.flushOngoingHeader()
	}
	if me.outputFile != nil {
		me.writeOutput(me.outputFile, line)
	}
}

func main() {
	flag.Parse()
	var app = new(App)
	app.inputFilePath = flag.Arg(0)
	app.run()
}
