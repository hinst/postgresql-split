package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/hinst/go-gophers"
	"github.com/hinst/go-gophers/file_mode"
)

type DumpSplitter struct {
	inputFilePath  string
	outputFilePath string
	outputFile     *os.File
	ongoingHeader  []string
	listFile       *os.File
}

const default_directory_permission = file_mode.USER_RWX | file_mode.GROUP_R | file_mode.GROUP_X | file_mode.OTH_R | file_mode.OTH_X
const line_length_limit = 32 * 1024 * 1024
const data_directory = "data"

func (me *DumpSplitter) run() {
	if me.inputFilePath == "" {
		log.Fatal("Input file is not specified")
	}

	file := gophers.AssertResultError(os.Open(me.inputFilePath))
	defer file.Close()

	gophers.AssertError(os.RemoveAll(data_directory))
	gophers.AssertError(os.MkdirAll(data_directory, default_directory_permission))

	me.listFile = gophers.AssertResultError(os.Create(data_directory + "/files.txt"))
	defer me.listFile.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, line_length_limit), line_length_limit)
	for scanner.Scan() {
		line := scanner.Text()
		me.receiveLine(line)
	}
	me.flushOngoingHeader()
	gophers.AssertError(scanner.Err())

	me.verifyRoundTrip()
}

func (me *DumpSplitter) verifyRoundTrip() {
	log.Println("Verifying round-trip integrity...")

	tmpOutput := data_directory + "/_dump.sql"
	var builder = &DumpBuilder{
		dataDirectory:  data_directory,
		outputFilePath: tmpOutput,
	}
	builder.run()
	defer os.Remove(tmpOutput)

	if !gophers.CheckFilesEqual(me.inputFilePath, tmpOutput) {
		log.Fatalf("Round-trip verification FAILED: %s", me.inputFilePath)
	}
	log.Println("Round-trip verified.")
}

func (me *DumpSplitter) writeOutput(output *os.File, content string) {
	gophers.AssertResultError(output.WriteString(content + "\n"))
}

func (me *DumpSplitter) openFile(fileName string) {
	if me.outputFile != nil {
		me.outputFile.Close()
	}
	me.outputFilePath = fileName
	me.outputFile = gophers.AssertResultError(os.Create(fileName))

	relativeName := strings.TrimPrefix(fileName, data_directory+"/")
	gophers.AssertResultError(me.listFile.WriteString(relativeName + "\n"))
}

func (me *DumpSplitter) flushOngoingHeader() {
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

func (me *DumpSplitter) receiveLine(line string) {
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
