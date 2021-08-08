package graph

import (
	"errors"

	"github.com/MaciejSzaflik/traffic/go_container"
)

var ErrNoPath = errors.New("No path found")

func AStar(start, goal int, d func(int, int) float32, getNeighbours func(int) []int, heuristic func(int, int) float32) error {
	openSet := make(map[int]struct{})
	cameFrom := make(map[int]int)

	openSet[start] = struct{}{}

	gScore := make(map[int]float32)
	gScore[start] = 0

	fScore := go_container.NewPriorityQueue()
	fScore.Push(
		&go_container.Item{
			Priority: heuristic(start, goal),
			Value:    start,
		},
	)

	for len(openSet) > 0 {
		current := fScore.Pop().(*go_container.Item).Value
		if current == goal {
			return nil
		}

		delete(openSet, current)
		for _, n := range getNeighbours(current) {
			tentativeGScore := gScore[current] + d(current, n)
			if v, ok := gScore[n]; ok && v > tentativeGScore || !ok {
				cameFrom[n] = current
				gScore[n] = tentativeGScore

				if _, ok := openSet[n]; !ok {
					fScore.Push(
						&go_container.Item{
							Priority: tentativeGScore + heuristic(n, goal),
							Value:    n,
						},
					)
					openSet[n] = struct{}{}
				}

			}
		}
	}

	return ErrNoPath
}
