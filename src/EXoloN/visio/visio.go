package visio

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

type plyvis [][]byte

func (pv plyvis) WritePoint(point int, w io.Writer) {
	if err := binary.Write(w, binary.LittleEndian, uint32(len(pv[point])/4)); err != nil {
		panic(fmt.Sprint("Could not write number of reference images for this point: ", err))
	}
	if _, err := w.Write(pv[point]); err != nil {
		panic(fmt.Sprint("Could not write image references: ", err))
	}
}

func ReadVis(in io.Reader) plyvis {
	var vislen uint64

	bin := bufio.NewReader(in)

	if err := binary.Read(bin, binary.LittleEndian, &vislen); err != nil {
		panic(fmt.Sprint("Unable to read vis: ", err))
	}

	vis := make(plyvis, vislen)

	for c := uint64(0); c < vislen; c++ {
		var num_ref_imgs uint32

		if err := binary.Read(bin, binary.LittleEndian, &num_ref_imgs); err != nil {
			panic(fmt.Sprint("Could not retrieve number of images for this point: ", err))
		}

		vis[c] = make([]byte, 4*num_ref_imgs)

		if _, err := io.ReadFull(bin, vis[c]); err != nil {
			panic(fmt.Sprint("Could not read images references: ", err))
		}
	}

	return vis
}

func (pv plyvis) WriteTo(out io.Writer) {
	but := bufio.NewWriter(out)
	defer but.Flush()

	if err := binary.Write(but, binary.LittleEndian, uint64(len(pv))); err != nil {
		panic(fmt.Sprint("Unable to write vis: ", err))
	}

	for p := range pv {
		pv.WritePoint(p, but)
	}
}

func (pv plyvis) WriteListTo(list []int, out io.Writer) {
	but := bufio.NewWriter(out)
	defer but.Flush()

	if err := binary.Write(but, binary.LittleEndian, uint64(len(list))); err != nil {
		panic(fmt.Sprint("Unable to write vis: ", err))
	}

	for _, p := range list {
		pv.WritePoint(p, but)
	}
}
