package presentation

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

const greetingMessage = `Welcome to the Maze Generator! The program generates a maze and a path for the shortest passage.

Before we start, please keep the following in mind:
 - To display the maze correctly, you must enter the dimensions so that the maze fits into the console
 - The start and end points must be on the boundaries of the maze
 - Maze width and height must be >= 2. For smaller values, we get a simple labyrinth in which the path 
   from the start point to the end point is clearly found

Enjoy the program!

`

type dimension struct {
	height int
	width  int
}

type Presentation struct {
	in  io.Reader
	out io.Writer
}

func New(in io.Reader, out io.Writer) *Presentation {
	return &Presentation{
		in:  in,
		out: out,
	}
}

func (p *Presentation) getInt(scan *bufio.Scanner, rng rangeNumber) (int, error) {
	for {
		if !scan.Scan() {
			return 0, ErrNoInputLines{}
		}

		inputLine := scan.Text()
		inputInt, err := strconv.Atoi(strings.TrimSpace(inputLine))

		switch {
		case err != nil:
			// ANSI code for red letters
			fmt.Fprintf(p.out, "\033[31mError: %s.\033[0m\nType a single integer: ", err)
		case !rng.Contains(inputInt):
			// ANSI code for red letters
			fmt.Fprintf(
				p.out,
				"\033[31mError: Integer should be in range %s.\033[0m\nType a valid integer: ",
				rng,
			)
		default:
			return inputInt, nil
		}
	}
}

func (p *Presentation) mazeDimension(scan *bufio.Scanner) (dimension, error) {
	dim := dimension{}

	fmt.Fprint(p.out, "Enter maze height: ")

	rng, err := newRange(newRangePoint(2, true), newRangePoint(0, false))
	if err != nil {
		return dimension{}, fmt.Errorf("create range: %w", err)
	}

	dim.height, err = p.getInt(scan, rng)
	if err != nil {
		return dimension{}, fmt.Errorf("read height from input stream: %w", err)
	}

	fmt.Fprint(p.out, "Enter maze width: ")

	rng, err = newRange(newRangePoint(2, true), newRangePoint(0, false))
	if err != nil {
		return dimension{}, fmt.Errorf("create range: %w", err)
	}

	dim.width, err = p.getInt(scan, rng)
	if err != nil {
		return dimension{}, fmt.Errorf("read width from input stream: %w", err)
	}

	return dim, nil
}

func (p *Presentation) point(scan *bufio.Scanner, pointName string, dim dimension) (domain.Coord, error) {
	fmt.Fprintf(p.out, "Enter maze %s point row id: ", pointName)

	var coord domain.Coord

	rng, err := newRange(
		newRangePoint(0, true),
		newRangePoint(dim.height-1, true),
	)
	if err != nil {
		return domain.Coord{}, fmt.Errorf("create range: %w", err)
	}

	coord.Row, err = p.getInt(scan, rng)
	if err != nil {
		return domain.Coord{}, fmt.Errorf(
			"read %s point row id from input stream: %w",
			pointName,
			err,
		)
	}

	fmt.Fprintf(p.out, "Enter maze %s point col id: ", pointName)

	rng, err = newRange(
		newRangePoint(0, true),
		newRangePoint(dim.width-1, true),
	)
	if err != nil {
		return domain.Coord{}, fmt.Errorf("create range: %w", err)
	}

	coord.Col, err = p.getInt(scan, rng)
	if err != nil {
		return domain.Coord{}, fmt.Errorf(
			"read %s point col id from input stream: %w",
			pointName,
			err,
		)
	}

	return coord, nil
}

func (p *Presentation) startCoord(scan *bufio.Scanner, dim dimension) (domain.Coord, error) {
	var start domain.Coord

	var err error

	for {
		start, err = p.point(scan, "start", dim)
		if err != nil {
			return domain.Coord{}, fmt.Errorf("get start point: %w", err)
		}

		if start.Row != 0 && start.Row != dim.height-1 && start.Col != 0 &&
			start.Col != dim.width-1 {
			fmt.Fprintln(
				p.out,
				"\033[31mError: start point must lie on the boundary.\033[0m\nType correct start point!",
			)
		} else {
			break
		}
	}

	return start, nil
}

func (p *Presentation) endCoord(scan *bufio.Scanner, dim dimension, start domain.Coord) (domain.Coord, error) {
	var end domain.Coord

	var err error

CoordLoop:
	for {
		end, err = p.point(scan, "end", dim)
		if err != nil {
			return domain.Coord{}, fmt.Errorf("get end point: %w", err)
		}

		switch {
		case end == start:
			fmt.Fprintln(
				p.out,
				"\033[31mError: start and end points are equal.\033[0m\nType correct end point!",
			)

		case end.Row != 0 && end.Row != dim.width-1 && end.Col != 0 && end.Col != dim.width-1:
			fmt.Fprintln(p.out, "\033[31mError: end point must lie on the boundary.\033[0m\nType correct end point!")

		default:
			break CoordLoop
		}
	}

	return end, nil
}

func (p *Presentation) ProcessInput() (domain.MazeData, error) {
	fmt.Fprint(p.out, greetingMessage)

	scan := bufio.NewScanner(p.in)

	dim, err := p.mazeDimension(scan)
	if err != nil {
		return domain.MazeData{}, fmt.Errorf("getting dimension: %w", err)
	}

	start, err := p.startCoord(scan, dim)
	if err != nil {
		return domain.MazeData{}, fmt.Errorf("getting start coord: %w", err)
	}

	end, err := p.endCoord(scan, dim, start)
	if err != nil {
		return domain.MazeData{}, fmt.Errorf("getting end coord: %w", err)
	}

	return domain.NewMazeData(dim.height, dim.width, start, end), nil
}
