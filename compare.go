package main

import (
	"EXoloN/cloudcompare"
	"EXoloN/plyreader"
	"fmt"
	"os"
)

func main() {
	var first, second *os.File
	var err error
	var one, two cloudcompare.PointCloud

	if len(os.Args) != 3 {
		panic("Must have two args.")
	}

	if first, err = os.Open(os.Args[1]); err != nil {
		panic(fmt.Sprint("Can't open input file:", err))
	}

	if second, err = os.Open(os.Args[2]); err != nil {
		panic(fmt.Sprint("Can't open input file:", err))
	}

	one = plyreader.ReadPLY(first)
	two = plyreader.ReadPLY(second)

	for c := 0; c < one.Elements(); c++ {
		fmt.Println(two.GetPosition(one.GetPointAt(c)))
	}
}
