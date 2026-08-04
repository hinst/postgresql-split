package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/hinst/go-gophers"
	"github.com/hinst/go-gophers/file_mode"
)

type App struct {
	inputFilePath  string
	outputFilePath string
	outputFile     *os.File
	ongoingHeader  []string
	processedFiles []string
}

const default_directory_permission = file_mode.USER_RWX | file_mode.GROUP_R | file_mode.GROUP_X | file_mode.OTH_R | file_mode.OTH_X
const line_length_limit = 32 * 1024 * 1024
const data_directory = "data"

func (me *App) run() {
	if me.inputFilePath == "" {
		log.Fatal("Input file is not specified")
	}

	file := gophers.AssertResultError(os.Open(me.inputFilePath))
	defer file.Close()

	gophers.AssertError(os.RemoveAll(data_directory))
	gophers.AssertError(os.MkdirAll(data_directory, default_directory_permission))
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, line_length_limit), line_length_limit)
	for scanner.Scan() {
		line := scanner.Text()
		me.receiveLine(line)
	}
	me.flushOngoingHeader()
	gophers.AssertError(scanner.Err())
	me.saveProcessedFiles()
}

func (me *App) writeOutput(output *os.File, content string) {
	gophers.AssertResultError(output.WriteString(content + "\n"))
}

func (me *App) openFile(fileName string) {
	if me.outputFile != nil {
		me.outputFile.Close()
	}
	me.outputFilePath = fileName
	me.outputFile = gophers.AssertResultError(os.Create(fileName))
	me.processedFiles = append(me.processedFiles, strings.TrimPrefix(fileName, data_directory+"/"))
}

func (me *App) flushOngoingHeader() {
	if checkContainsToc(me.ongoingHeader) {
		sectionName := getSectionName(me.ongoingHeader)
		if sectionName == "" {
			panic("Cannot find section name in TOC header")
		}
		folder, filename := getQualifiedName(sectionName)
		if folder != "" {
			gophers.AssertError(os.MkdirAll(data_directory+"/"+folder, default_directory_permission))
		}
		filename = data_directory + "/" + folder + gophers.IfElse(len(folder) > 0, "/", "") + filename + ".sql"
		me.openFile(filename)
	}
	for _, header := range me.ongoingHeader {
		me.writeOutput(me.outputFile, header)
	}
	me.ongoingHeader = nil
}

func (me *App) saveProcessedFiles() {
	filePath := data_directory + "/files.txt"
	file := gophers.AssertResultError(os.Create(filePath))
	defer file.Close()
	for _, f := range me.processedFiles {
		me.writeOutput(file, f)
	}
}

func (me *App) receiveLine(line string) {
	if "" == me.outputFilePath {
		me.openFile(data_directory + "/init.sql")
	}

	if strings.HasPrefix(line, "--") {
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
