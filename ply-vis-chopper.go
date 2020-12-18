package main

import (
	"EXoloN/cloudcompare"
	"EXoloN/plyreader"
	"EXoloN/visio"
	"fmt"
	"os"
)

func open_files(args []string) (*os.File, *os.File, *os.File, *os.File) {
	var err error
	var out [4]*os.File

	for c := 0; c < 3; c++ {
		if out[c], err = os.Open(args[c]); err != nil {
			panic(fmt.Sprint("Can't open input file:", err))
		}
	}

	if out[3], err = os.OpenFile(args[3], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664); err != nil {
		panic(fmt.Sprint("Can't open ouput file:", err))
	}

	return out[0], out[1], out[2], out[3]
}

func main() {
	var original, cropped cloudcompare.PointCloud

	if len(os.Args) != 5 {
		panic("Must have four Args!")
	}

	original_ply, cropped_ply, ply_vis, chopped_vis := open_files(os.Args[1:])

	original = plyreader.ReadPLY(original_ply)
	cropped = plyreader.ReadPLY(cropped_ply)
	vis := visio.ReadVis(ply_vis)

	positions := make([]int, cropped.Elements())

	for pos := range positions {
		if newpos, exists := original.GetPosition(cropped.GetPointAt(pos)); exists {
			positions[pos] = newpos
		} else {
			panic(fmt.Sprint("Point at position ", pos, "can not be found in original PLY file."))
		}
	}

	vis.WriteListTo(positions, chopped_vis)
}
