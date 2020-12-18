package main

import (
	"EXoloN/plyreader"
	"fmt"
	"os"
)

func main() {
	var file *os.File
	var err error

	if len(os.Args) != 2 {
		panic("Must have one arg.")
	}

	if file, err = os.Open(os.Args[1]); err != nil {
		panic(fmt.Sprint("Can't open input file:", err))
	}

	plyreader.ReadPLY(file)
}
