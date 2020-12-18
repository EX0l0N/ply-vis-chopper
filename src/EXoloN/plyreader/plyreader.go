package plyreader

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	REQ_X = iota
	REQ_Y
	REQ_Z
	REQ_RED
	REQ_GREEN
	REQ_BLUE
	OPT_NX
	OPT_NY
	OPT_NZ
	REQ_ALPHA
	REQ_VERTEX
	REQ_FORMAT
	REQ_FIELD_LEN
)

type ply_header struct {
	num_vertices           int64
	has_alpha, has_normals bool
	field_order            []byte
}

func parse_header(in *bufio.Reader) ply_header {
	var checked_fields [REQ_FIELD_LEN]bool
	var header ply_header
	header.field_order = make([]byte, 0, REQ_FIELD_LEN-2)

	tokens := func() []string {
		if txt, err := in.ReadString('\n'); err != nil {
			panic(fmt.Sprint("Unable to get line: ", err))
		} else {
			return strings.Split(strings.TrimRight(txt, " \r\n"), " ")
		}
	}

	if magic := tokens(); len(magic) != 1 || magic[0] != "ply" {
		panic(fmt.Sprint("Magic line corrupted:", magic))
	}

	for run := true; run; {
		line := tokens()

		switch line[0] {
		case "comment":
			fmt.Println("Ignoring comment:", line[1:])
		case "format":
			checked_fields[REQ_FORMAT] = true
			if line[1] != "binary_little_endian" || line[2] != "1.0" {
				fmt.Println(`Format needs to be exactly "binary_little_endian 1.0"`)
				fmt.Printf("Got %q instead.", line[1:])
				panic("Can't parse format")
			}
		case "element":
			if line[1] == "vertex" {
				checked_fields[REQ_VERTEX] = true
				fmt.Print("Parsing vertex count: ")
				if i, err := strconv.ParseInt(line[2], 10, 64); err != nil {
					panic(err)
				} else {
					header.num_vertices = i
				}
				fmt.Printf("Setting up for %d vertices.\n", header.num_vertices)
			} else {
				fmt.Printf("I hope it's ok to totally ignore %q.\n", line)
			}
		case "property":
			switch line[1] {
			case "float":
				switch line[2] {
				case "x":
					checked_fields[REQ_X] = true
					header.field_order = append(header.field_order, REQ_X)
				case "y":
					checked_fields[REQ_Y] = true
					header.field_order = append(header.field_order, REQ_Y)
				case "z":
					checked_fields[REQ_Z] = true
					header.field_order = append(header.field_order, REQ_Z)
				case "nx":
					checked_fields[OPT_NX] = true
					header.has_normals = true
					header.field_order = append(header.field_order, OPT_NX)
				case "ny":
					checked_fields[OPT_NY] = true
					header.has_normals = true
					header.field_order = append(header.field_order, OPT_NY)
				case "nz":
					checked_fields[OPT_NZ] = true
					header.has_normals = true
					header.field_order = append(header.field_order, OPT_NZ)
				default:
					fmt.Println(line)
					panic("Uknown float property.")
				}
			case "uchar":
				switch line[2] {
				case "red":
					checked_fields[REQ_RED] = true
					header.field_order = append(header.field_order, REQ_RED)
				case "green":
					checked_fields[REQ_GREEN] = true
					header.field_order = append(header.field_order, REQ_GREEN)
				case "blue":
					checked_fields[REQ_BLUE] = true
					header.field_order = append(header.field_order, REQ_BLUE)
				case "alpha":
					checked_fields[REQ_ALPHA] = true
					header.field_order = append(header.field_order, REQ_ALPHA)
					header.has_alpha = true
					fmt.Println("Please note that alpha values will be ignored.")
				default:
					fmt.Println(line)
					panic("Unknown uchar property.")
				}
			case "list":
				fmt.Printf("Ignoring unknown porperty list %q.\n", line)
			default:
				fmt.Println(line)
				panic("Can't read that.")
			}
		case "end_header":
			fmt.Println("That's it. From now on I expect binary data.")
			run = false
		default:
			fmt.Println(line)
			panic("Can't read that.")
		}

	}

	if !(checked_fields[REQ_FORMAT] &&
		checked_fields[REQ_VERTEX] &&
		checked_fields[REQ_X] &&
		checked_fields[REQ_Y] &&
		checked_fields[REQ_Z] &&
		checked_fields[REQ_RED] &&
		checked_fields[REQ_GREEN] &&
		checked_fields[REQ_BLUE]) {

		fmt.Println(checked_fields)
		panic("Did not see all the required fields in header. Giving up.")
	}

	if header.has_normals && !(checked_fields[OPT_NX] && checked_fields[OPT_NY] && checked_fields[OPT_NZ]) {
		panic("If normals are used than they need to exist for all three dimensions.")
	}

	return header
}

type colors struct {
	r, g, b byte
}

type positions struct {
	x, y, z float32
}

type normals struct {
	nx, ny, nz float32
}

type colordata struct {
	pos positions
	col colors
}

func (p *colordata) get_pos() *positions {
	return &p.pos
}

func (p *colordata) get_col() *colors {
	return &p.col
}

func (p *colordata) get_norm() *normals {
	return nil
}

type colordata_with_normals struct {
	pos  positions
	col  colors
	norm normals
}

func (p *colordata_with_normals) get_pos() *positions {
	return &p.pos
}

func (p *colordata_with_normals) get_col() *colors {
	return &p.col
}

func (p *colordata_with_normals) get_norm() *normals {
	return &p.norm
}

type point interface {
	get_pos() *positions
	get_col() *colors
	get_norm() *normals
}

type pointcloud struct {
	hasher map[positions]int
	points []point
}

func read_pointcloud(in *bufio.Reader, header ply_header) pointcloud {
	var pc pointcloud

	pc.hasher = make(map[positions]int)
	pc.points = make([]point, int(header.num_vertices))

	read_float32 := func() float32 {
		var f float32

		if err := binary.Read(in, binary.LittleEndian, &f); err != nil {
			panic(fmt.Sprint("binary.Read failed:", err))
		}

		return f
	}

	read_byte := func() byte {
		var b byte

		if err := binary.Read(in, binary.LittleEndian, &b); err != nil {
			panic(fmt.Sprint("binary.Read failed:", err))
		}

		return b
	}

	for c := range pc.points {
		var p point

		if header.has_normals {
			p = new(colordata_with_normals)
		} else {
			p = new(colordata)
		}

		for i := 0; i < len(header.field_order); i++ {
			switch header.field_order[i] {
			case REQ_X:
				p.get_pos().x = read_float32()
			case REQ_Y:
				p.get_pos().y = read_float32()
			case REQ_Z:
				p.get_pos().z = read_float32()
			case REQ_RED:
				p.get_col().r = read_byte()
			case REQ_GREEN:
				p.get_col().g = read_byte()
			case REQ_BLUE:
				p.get_col().b = read_byte()
			case REQ_ALPHA:
				if d, err := (*in).Discard(1); err != nil || d != 1 {
					panic("Unable to discard one byte.")
				}
			case OPT_NX:
				p.get_norm().nx = read_float32()
			case OPT_NY:
				p.get_norm().ny = read_float32()
			case OPT_NZ:
				p.get_norm().nz = read_float32()
			default:
				panic("Wrong use of field order.")
			}

			pc.hasher[*p.get_pos()] = c
			pc.points[c] = p
		}
	}

	return pc
}

func ReadPLY(file io.Reader) pointcloud {
	br := bufio.NewReader(file)

	head := parse_header(br)

	return read_pointcloud(br, head)
}
