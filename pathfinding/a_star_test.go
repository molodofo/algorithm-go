package pathfinding

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
	"unsafe"
)

type testPosition struct {
	X int
	Y int
}

func (p testPosition) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

type placeholder struct{}

type testGraph struct {
	width, height int
	obstacle      map[testPosition]placeholder
	start, end    testPosition
}

func (graph *testGraph) RandObstacle(proportion float64) {
	obstacleNum := int(float64(graph.width*graph.height) * proportion)
	for i := 0; i < obstacleNum; i++ {
		x, y := rand.Intn(graph.width), rand.Intn(graph.height)
		graph.obstacle[testPosition{X: x, Y: y}] = placeholder{}
	}
}

func (graph *testGraph) SetObstacle(obstacle []testPosition) {
	for _, o := range obstacle {
		graph.obstacle[o] = placeholder{}
	}
}

func (graph *testGraph) Init(width, height int) {
	graph.width, graph.height = width, height
	graph.obstacle = make(map[testPosition]placeholder)
	graph.start = testPosition{X: 0, Y: 0}
	graph.end = testPosition{X: width - 1, Y: height - 1}
}

func (graph *testGraph) SetStart(position any) {
	graph.start = position.(testPosition)
}

func (graph *testGraph) GetStart() any {
	return graph.start
}

func (graph *testGraph) SetEnd(position any) {
	graph.end = position.(testPosition)
}

func (graph *testGraph) GetEnd() any {
	return graph.end
}

func (graph *testGraph) IsEnd(position any) bool {
	p := position.(testPosition)
	return p.X == graph.end.X && p.Y == graph.end.Y
}

func (graph *testGraph) overBoundary(position testPosition) bool {
	return position.X < 0 || position.X >= graph.width || position.Y < 0 || position.Y > graph.height
}

func (graph *testGraph) isObstacle(position testPosition) bool {
	_, ok := graph.obstacle[position]
	return ok
}

func (graph *testGraph) String() string {
	s := strings.Repeat(" * ", graph.width+2) + "\n"
	for w := 0; w < graph.width; w++ {
		s += fmt.Sprintf(" * ")
		for h := 0; h < graph.height; h++ {
			if graph.isObstacle(testPosition{w, h}) {
				s += fmt.Sprintf(" * ")
			} else {
				s += fmt.Sprintf("   ")
			}
		}
		s += fmt.Sprintf(" * \n")
	}
	s += strings.Repeat(" * ", graph.width+2) + "\n"
	return s
}

func (graph *testGraph) DisplayPath(path []any) {
	pathSet := make(map[testPosition]placeholder)
	for _, p := range path {
		pathSet[p.(testPosition)] = placeholder{}
	}
	n := 0
	fmt.Println(strings.Repeat(" * ", graph.width+2))
	for w := 0; w < graph.width; w++ {
		fmt.Printf(" * ")
		for h := 0; h < graph.height; h++ {
			if _, ok := pathSet[testPosition{w, h}]; ok {
				fmt.Printf("%3d", n)
				n += 1
			} else if graph.isObstacle(testPosition{w, h}) {
				fmt.Printf(" * ")
			} else {
				fmt.Printf("   ")
			}
		}
		fmt.Printf(" * \n")
	}
	fmt.Println(strings.Repeat(" * ", graph.width+2))
}

func (graph *testGraph) Neighbors(position any) []any {
	p := position.(testPosition)
	direction := []testPosition{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	neighbors := make([]any, 0)
	for _, d := range direction {
		current := testPosition{X: p.X + d.X, Y: p.Y + d.Y}
		if !graph.overBoundary(current) && !graph.isObstacle(current) {
			neighbors = append(neighbors, current)
		}
	}
	return neighbors
}

func abs(n int) int {
	size := unsafe.Sizeof(n)
	m := n>>size - 1
	return (n ^ m) - m
}

func manhattanDistance(p1, p2 testPosition) int {
	return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}

func (graph *testGraph) Cost(position1, position2 any) int {
	return manhattanDistance(position1.(testPosition), position2.(testPosition))
}

func (graph *testGraph) Heuristic(position any) int {
	p := position.(testPosition)
	return manhattanDistance(p, graph.end)
}

/*
=== RUN   TestAStar
  - *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *
  - *   0  1  2  3  4  5                                           *
  - *                  6                                           *
  - *                  7                                           *
  - *                  8                                           *
  - *                  9                                           *
  - *                 10                                           *
  - *                 11                                           *
  - *                 12                                           *
  - *                 13                                           *
  - *                 14 15 16                                     *
  - *                       17                         *           *
  - *                       18                         *           *
  - *                       19                         *           *
  - *                       20 21 22                   *           *
  - *                          *  23                   *           *
  - *                          *  24                   *           *
  - *                          *  25    *  *  *  *  *  *           *
  - *                          *  26 27 28 29 30 31 32             *
  - *                          *                    33             *
  - *                          *                    34 35 36 37 38 *
  - *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *

-
  - *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *
  - *   0                                                          *
  - *   1                                                          *
  - *   2                                                          *
  - *   3                                                          *
  - *   4                                                          *
  - *   5                                                          *
  - *   6                                                          *
  - *   7                                                          *
  - *   8                                                          *
  - *   9                                                          *
  - *  10                                              *           *
  - *  11                                              *           *
  - *  12                                              *           *
  - *  13 14 15 16 17 18 19 20 21 22                   *           *
  - *                          *  23                   *           *
  - *                          *  24                   *           *
  - *                          *  25    *  *  *  *  *  *           *
  - *                          *  26                               *
  - *                          *  27                               *
  - *                          *  28 29 30 31 32 33 34 35 36 37 38 *
  - *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *  *

--- PASS: TestAStar (0.00s)
*/
func TestAStar(t *testing.T) {
	graph := testGraph{}
	graph.Init(20, 20)
	graph.SetObstacle([]testPosition{
		{16, 16}, {16, 15}, {16, 14}, {16, 13}, {16, 12}, {16, 11},
		{15, 16}, {14, 16}, {13, 16}, {12, 16}, {11, 16}, {10, 16},
		{14, 8}, {15, 8}, {16, 8}, {17, 8}, {18, 8}, {19, 8},
	})
	path := AStar(&graph)
	if path != nil {
		graph.DisplayPath(path)
	}

	path = Bfs(&graph)
	if path != nil {
		graph.DisplayPath(path)
	}

	graph = testGraph{}
	graph.Init(0, 0)
	path = AStar(&graph)
}

/*
=== RUN   TestAStar_Performance
a* cost: 1.3871671s
bfs cost: 761.1499ms
--- PASS: TestAStar_Performance (3.71s)
*/
func TestAStar_Performance(t *testing.T) {
	graph := testGraph{}
	graph.Init(1000, 1000)
	graph.RandObstacle(0.4)
	path := AStar(&graph)
	for path == nil {
		graph.Init(1000, 1000)
		graph.RandObstacle(0.4)
		path = AStar(&graph)
	}

	start := time.Now()
	AStar(&graph)
	fmt.Printf("a* cost: %v\n", time.Since(start))

	start = time.Now()
	Bfs(&graph)
	fmt.Printf("bfs cost: %v\n", time.Since(start))
}
