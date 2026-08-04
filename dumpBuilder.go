package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hinst/go-gophers"
)

type DumpBuilder struct {
	dataDirectory  string
	outputFilePath string
}

func (me *DumpBuilder) run() {
	if me.dataDirectory == "" || me.outputFilePath == "" {
		log.Fatal("Data directory and output file path must be specified")
	}

	listFilePath := filepath.Join(me.dataDirectory, "files.txt")
	listFile := gophers.AssertResultError(os.Open(listFilePath))
	defer listFile.Close()

	outputFile := gophers.AssertResultError(os.Create(me.outputFilePath))
	defer outputFile.Close()

	scanner := bufio.NewScanner(listFile)
	for scanner.Scan() {
		fileName := strings.TrimSpace(scanner.Text())
		if fileName == "" {
			continue
		}

		filePath := filepath.Join(me.dataDirectory, filepath.FromSlash(fileName))
		contentFile := gophers.AssertResultError(os.Open(filePath))
		
		_, err := io.Copy(outputFile, contentFile)
		contentFile.Close()
		gophers.AssertError(err)
	}
	gophers.AssertError(scanner.Err())
}
