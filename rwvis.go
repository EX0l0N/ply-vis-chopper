package main

import (
	"EXoloN/visio"
	"fmt"
	"os"
)

func main() {
	var infile, outfile *os.File
	var err error

	if len(os.Args) != 3 {
		panic("Wrong number of Args.")
	}

	if infile, err = os.Open(os.Args[1]); err != nil {
		panic(fmt.Sprint("Can't open input file:", err))
	}
	if outfile, err = os.OpenFile(os.Args[2], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664); err != nil {
		panic(fmt.Sprint("Can't open ouput file:", err))
	}

	vis := visio.ReadVis(infile)
	visio.WriteVis(vis, outfile)
}
