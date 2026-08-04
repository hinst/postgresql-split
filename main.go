package main

import (
	"flag"
)

func main() {
	flag.Parse()
	var app = new(DumpSplitter)
	app.inputFilePath = flag.Arg(0)
	app.run()
}
