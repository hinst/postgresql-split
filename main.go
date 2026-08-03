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

	AssertError(os.MkdirAll("./data", 0755))
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

func (me *App) openFile(fileName string) {
	me.outputFile.Close()
	me.outputFilePath = fileName
	me.outputFile = AssertResultError(os.Create(fileName))
}

func (me *App) flushOngoingHeader() {
	if checkContainsToc(me.ongoingHeader) {
		sectionName := getSectionName(me.ongoingHeader)
		if sectionName == "" {
			panic("Cannot find section name in TOC header")
		}
		folder, filename := splitSectionName(sectionName)
		if folder != "" {
			AssertError(os.MkdirAll("./data/"+folder, 0755))
			me.openFile("./data/" + folder + "/" + filename + ".sql")
		} else {
			me.openFile("./data/" + filename + ".sql")
		}
	}
	for _, header := range me.ongoingHeader {
		me.writeOutput(me.outputFile, header)
	}
	me.ongoingHeader = nil
}

func (me *App) receiveLine(line string) {
	if "" == me.outputFilePath {
		me.openFile("./data/init.sql")
	}

	if strings.HasPrefix(line, "--") {
		me.ongoingHeader = append(me.ongoingHeader, line)
	} else {
		me.flushOngoingHeader()
		me.writeOutput(me.outputFile, line)
	}
}

func splitSectionName(sectionName string) (folder string, rest string) {
	i := 0
	for i < len(sectionName) && ((sectionName[i] >= 'A' && sectionName[i] <= 'Z') || (sectionName[i] >= 'a' && sectionName[i] <= 'z')) {
		i++
	}
	if i == 0 || i == len(sectionName) {
		rest = sanitizeFilename(sectionName)
	} else {
		folder = sanitizeFilename(sectionName[:i])
		rest = sanitizeFilename(sectionName[i:])
	}
	return
}

func main() {
	flag.Parse()
	var app = new(App)
	app.inputFilePath = flag.Arg(0)
	app.run()
}
