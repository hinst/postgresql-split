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

func sanitizeFilename(name string) string {
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	return b.String()
}

func (me *App) flushOngoingHeader() {
	if checkContainsToc(me.ongoingHeader) {
		sectionName := getSectionName(me.ongoingHeader)
		if "" == sectionName {
			panic("Cannot find section name in TOC header")
		}
		filename := sanitizeFilename(sectionName)
		me.outputFile.Close()
		outputFile := AssertResultError(os.Create(filename))
		me.outputFile = outputFile
	}
	for _, header := range me.ongoingHeader {
		me.writeOutput(me.outputFile, header)
	}
	me.ongoingHeader = nil
}

func (me *App) receiveLine(line string) {
	if "" == me.outputFilePath {
		me.outputFilePath = "./data/init.sql"
	}
	if nil == me.outputFile {
		outputFile := AssertResultError(os.Create(me.outputFilePath))
		me.outputFile = outputFile
	}

	if strings.HasPrefix(line, "-- ") {
		me.ongoingHeader = append(me.ongoingHeader, line)
	} else {
		me.flushOngoingHeader()
		me.writeOutput(me.outputFile, line)
	}
}

func main() {
	flag.Parse()
	var app = new(App)
	app.inputFilePath = flag.Arg(0)
	app.run()
}
