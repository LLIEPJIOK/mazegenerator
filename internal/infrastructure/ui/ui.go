package ui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

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

type Generator interface {
	GenerateMaze(
		height, width int,
		start, end domain.Coord,
		drawingChan chan<- domain.CellRenderData,
	) (domain.Maze, error)
}

type Painter interface {
	PaintGeneration(
		ctx context.Context,
		mazeHeight, mazeWidth int,
		cellChan <-chan domain.CellRenderData,
	)

	PaintPath(path []domain.Coord, delay time.Duration)

	MoveCursor(rowID, colID int)
}

type PathFinder interface {
	FindPath(maze domain.Maze, start, end domain.Coord) ([]domain.Coord, bool)
}

type UI struct {
	in         io.Reader
	out        io.Writer
	mazeHeight int
	mazeWidth  int
	start      domain.Coord
	end        domain.Coord
	gen        Generator
	painter    Painter
	pathFinder PathFinder
}

func NewUI(in io.Reader, out io.Writer, gen Generator, painter Painter, pathFinder PathFinder) *UI {
	return &UI{
		in:         in,
		out:        out,
		gen:        gen,
		painter:    painter,
		pathFinder: pathFinder,
	}
}

func (ui *UI) getInt(scan *bufio.Scanner, rng domain.Range) (int, error) {
	for {
		if !scan.Scan() {
			return 0, ErrNoInputLines{}
		}

		inputLine := scan.Text()
		inputInt, err := strconv.Atoi(strings.TrimSpace(inputLine))

		switch {
		case err != nil:
			// ANSI code for red letters
			fmt.Fprintf(ui.out, "\033[31mError: %s.\033[0m\nType a single integer: ", err)
		case !rng.Contains(inputInt):
			// ANSI code for red letters
			fmt.Fprintf(
				ui.out,
				"\033[31mError: Integer should be in range %s.\033[0m\nType a valid integer: ",
				rng,
			)
		default:
			return inputInt, nil
		}
	}
}

func (ui *UI) getMazeDimension(scan *bufio.Scanner) error {
	fmt.Fprint(ui.out, "Enter maze height: ")

	rng, err := domain.NewRange(domain.NewRangePoint(2, true), domain.NewRangePoint(0, false))
	if err != nil {
		return fmt.Errorf("create range: %w", err)
	}

	ui.mazeHeight, err = ui.getInt(scan, rng)
	if err != nil {
		return fmt.Errorf("read height from input stream: %w", err)
	}

	fmt.Fprint(ui.out, "Enter maze width: ")

	rng, err = domain.NewRange(domain.NewRangePoint(2, true), domain.NewRangePoint(0, false))
	if err != nil {
		return fmt.Errorf("create range: %w", err)
	}

	ui.mazeWidth, err = ui.getInt(scan, rng)
	if err != nil {
		return fmt.Errorf("read width from input stream: %w", err)
	}

	return nil
}

func (ui *UI) getPoint(scan *bufio.Scanner, pointName string) (domain.Coord, error) {
	fmt.Fprintf(ui.out, "Enter maze %s point row id: ", pointName)

	var coord domain.Coord

	rng, err := domain.NewRange(
		domain.NewRangePoint(0, true),
		domain.NewRangePoint(ui.mazeHeight-1, true),
	)
	if err != nil {
		return domain.Coord{}, fmt.Errorf("create range: %w", err)
	}

	coord.RowID, err = ui.getInt(scan, rng)
	if err != nil {
		return domain.Coord{}, fmt.Errorf(
			"read %s point row id from input stream: %w",
			pointName,
			err,
		)
	}

	fmt.Fprintf(ui.out, "Enter maze %s point col id: ", pointName)

	rng, err = domain.NewRange(
		domain.NewRangePoint(0, true),
		domain.NewRangePoint(ui.mazeWidth-1, true),
	)
	if err != nil {
		return domain.Coord{}, fmt.Errorf("create range: %w", err)
	}

	coord.ColID, err = ui.getInt(scan, rng)
	if err != nil {
		return domain.Coord{}, fmt.Errorf(
			"read %s point col id from input stream: %w",
			pointName,
			err,
		)
	}

	return coord, nil
}

func (ui *UI) getStartAndEndPoints(scan *bufio.Scanner) error {
	var err error

	for {
		ui.start, err = ui.getPoint(scan, "start")
		if err != nil {
			return fmt.Errorf("get start point: %w", err)
		}

		if ui.start.RowID != 0 && ui.start.RowID != ui.mazeHeight-1 && ui.start.ColID != 0 &&
			ui.start.ColID != ui.mazeWidth-1 {
			fmt.Fprintln(
				ui.out,
				"\033[31mError: start point must lie on the boundary.\033[0m\nType correct start point!",
			)
		} else {
			break
		}
	}

EndPointLoop:
	for {
		ui.end, err = ui.getPoint(scan, "end")
		if err != nil {
			return fmt.Errorf("get end point: %w", err)
		}

		switch {
		case ui.end == ui.start:
			fmt.Fprintln(
				ui.out,
				"\033[31mError: start and end points are equal.\033[0m\nType correct end point!",
			)
		case ui.end.RowID != 0 && ui.end.RowID != ui.mazeWidth-1 && ui.end.ColID != 0 && ui.end.ColID != ui.mazeWidth-1:
			fmt.Fprintln(ui.out, "\033[31mError: end point must lie on the boundary.\033[0m\nType correct end point!")
		default:
			break EndPointLoop
		}
	}

	return nil
}

func (ui *UI) processInput() error {
	fmt.Fprint(ui.out, greetingMessage)

	scan := bufio.NewScanner(ui.in)

	if err := ui.getMazeDimension(scan); err != nil {
		return fmt.Errorf("getting dimension: %w", err)
	}

	if err := ui.getStartAndEndPoints(scan); err != nil {
		return fmt.Errorf("getting start and end points: %w", err)
	}

	return nil
}

func (ui *UI) Run() error {
	err := ui.processInput()
	if err != nil {
		return fmt.Errorf("ui.ProcessInput(): %w", err)
	}

	ch := make(chan domain.CellRenderData)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		ui.painter.PaintGeneration(context.Background(), ui.mazeHeight, ui.mazeWidth, ch)
	}()

	maze, err := ui.gen.GenerateMaze(
		ui.mazeHeight,
		ui.mazeWidth,
		ui.start,
		ui.end,
		ch,
	)
	if err != nil {
		return fmt.Errorf("generate maze: %w", err)
	}

	close(ch)
	wg.Wait()

	if path, ok := ui.pathFinder.FindPath(maze, ui.start, ui.end); ok {
		ui.painter.PaintPath(path, 20*time.Millisecond)
	} else {
		ui.painter.MoveCursor(ui.mazeHeight+1, 0)
		fmt.Println("No path in this maze")
	}

	ui.painter.MoveCursor(ui.mazeHeight+2, 0)

	return nil
}
