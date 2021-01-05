package visio

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

type plyvis [][]byte

func (pv plyvis) WritePoint(point int, w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, uint32(len(pv[point])/4)); err != nil {
		return fmt.Errorf("Could not write number of reference images for this point: %w", err)
	}
	if _, err := w.Write(pv[point]); err != nil {
		return fmt.Errorf("Could not write image references: %w", err)
	}

	return nil
}

func ReadVis(in io.Reader) (plyvis, error) {
	var vislen uint64

	bin := bufio.NewReader(in)

	if err := binary.Read(bin, binary.LittleEndian, &vislen); err != nil {
		return nil, fmt.Errorf("Unable to read vis: %w", err)
	}

	vis := make(plyvis, vislen)

	for c := uint64(0); c < vislen; c++ {
		var num_ref_imgs uint32

		if err := binary.Read(bin, binary.LittleEndian, &num_ref_imgs); err != nil {
			return nil, fmt.Errorf("Could not retrieve number of images for this point: %w", err)
		}

		vis[c] = make([]byte, 4*num_ref_imgs)

		if _, err := io.ReadFull(bin, vis[c]); err != nil {
			return nil, fmt.Errorf("Could not read images references: %w", err)
		}
	}

	return vis, nil
}

func (pv plyvis) WriteTo(out io.Writer) error {
	but := bufio.NewWriter(out)
	defer but.Flush()

	if err := binary.Write(but, binary.LittleEndian, uint64(len(pv))); err != nil {
		return fmt.Errorf("Unable to write vis: %w", err)
	}

	for p := range pv {
		pv.WritePoint(p, but)
	}

	return nil
}

func (pv plyvis) WriteListTo(list []int, out io.Writer) error {
	but := bufio.NewWriter(out)
	defer but.Flush()

	if err := binary.Write(but, binary.LittleEndian, uint64(len(list))); err != nil {
		return fmt.Errorf("Unable to write vis: %w", err)
	}

	for _, p := range list {
		if err := pv.WritePoint(p, but); err != nil {
			return err
		}
	}

	return nil
}
