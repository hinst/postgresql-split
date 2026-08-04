package main

import (
	"flag"
)

func main() {
	buildDir := flag.String("build", "", "Build the dump from split files in the specified folder")
	flag.Parse()

	if *buildDir != "" {
		outputFile := flag.Arg(0)
		if outputFile == "" {
			outputFile = "dump.sql"
		}
		var builder = &DumpBuilder{
			dataDirectory:  *buildDir,
			outputFilePath: outputFile,
		}
		builder.run()
	} else {
		var app = new(DumpSplitter)
		app.inputFilePath = flag.Arg(0)
		app.run()
	}
}
