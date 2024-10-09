package pathfinder

import (
	"container/heap"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
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

func (d *Dijkstra) ShortestPath(maze domain.Maze) ([]domain.Coord, bool) {
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

		if curItem.curCoord == maze.Data.End {
			path := getPath(prevCoords, maze.Data.End, maze.Data.Start)

			return path, true
		}

		for i := range d.dir.Cols {
			dRow, dCol := d.dir.Rows[i], d.dir.Cols[i]
			newRow, newCol := curItem.curCoord.Row+dRow, curItem.curCoord.Col+dCol

			if min(newCol, newRow) >= 0 && newRow < maze.Data.Height && newCol < maze.Data.Width &&
				maze.Cells[newRow][newCol] != domain.Wall {
				heap.Push(
					&pq,
					newDijkstraItem(
						domain.NewCoord(newRow, newCol),
						curItem.curCoord,
						curItem.distance+int(maze.Cells[newRow][newCol]),
					),
				)
			}
		}
	}

	return nil, false
}
