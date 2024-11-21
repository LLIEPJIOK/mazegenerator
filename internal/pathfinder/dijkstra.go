package pathfinder

import (
	"container/heap"

	"github.com/LLIEPJIOK/mazegenerator/internal/domain"
)

type dijkstraItem struct {
	curCoord  domain.Coord
	prevCoord domain.Coord
	distance  int
}

func newDijkstraItem(curCoord, prevCoord domain.Coord, distance int) dijkstraItem {
	return dijkstraItem{
		curCoord:  curCoord,
		prevCoord: prevCoord,
		distance:  distance,
	}
}

func (d dijkstraItem) dist() int {
	return d.distance
}

type Dijkstra struct {
	dir domain.Direction
}

func NewDijkstra() *Dijkstra {
	return &Dijkstra{
		dir: domain.DefaultDirection(),
	}
}

func (d *Dijkstra) ShortestPath(maze domain.Maze, pathChan chan<- []domain.Coord) ([]domain.Coord, bool) {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, newDijkstraItem(maze.Data.Start, domain.Coord{}, 0))

	prevCoords := make(map[domain.Coord]domain.Coord)

	for len(pq) != 0 {
		curItem := heap.Pop(&pq).(dijkstraItem)

		if _, ok := prevCoords[curItem.curCoord]; ok {
			continue
		}

		prevCoords[curItem.curCoord] = curItem.prevCoord

		curPath := getPath(prevCoords, curItem.curCoord, maze.Data.Start)
		pathChan <- curPath

		if curItem.curCoord == maze.Data.End {
			path := getPath(prevCoords, maze.Data.End, maze.Data.Start)

			return path, true
		}

		for i := range d.dir.Cols {
			dRow, dCol := d.dir.Rows[i], d.dir.Cols[i]
			newRow, newCol := curItem.curCoord.Row+dRow, curItem.curCoord.Col+dCol

			if min(newCol, newRow) >= 0 && newRow < maze.Data.Height && newCol < maze.Data.Width &&
				maze.Cells[newRow][newCol].IsTraversable() {
				heap.Push(
					&pq,
					newDijkstraItem(
						domain.NewCoord(newRow, newCol),
						curItem.curCoord,
						curItem.distance+maze.Cells[newRow][newCol].Cost(),
					),
				)
			}
		}
	}

	pathChan <- nil

	return nil, false
}
