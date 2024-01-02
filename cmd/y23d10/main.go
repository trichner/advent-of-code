package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"slices"

	"aoc/pkg/in"
	"aoc/pkg/sets"
	"aoc/pkg/vec"
)

//go:embed *.txt
var inputs embed.FS

const (
	VERTICAL   = '|'
	HORIZONTAL = '-'
	BEND_NE    = 'L'
	BEND_NW    = 'J'
	BEND_SW    = '7'
	BEND_SE    = 'F'
	GROUND     = '.'
	START      = 'S'
)

func main() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()
	field := parseField(file)

	loop := findLoop(field)

	partOne(loop)
	partTwo(loop)
}

type Node struct {
	Key     vec.Vec2i
	Next    []vec.Vec2i
	IsStart bool
}

func (n *Node) String() string {
	return fmt.Sprintf("{%v: %v}", n.Key, n.Next)
}

func parseField(file fs.File) [][]byte {
	scanner := bufio.NewScanner(file)

	field := [][]byte{}

	for scanner.Scan() {
		text := scanner.Text()
		field = append(field, []byte(text))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return field
}

func findLoop(field [][]byte) []vec.Vec2i {
	graph := buildGraph(field)

	start := findStart(field)
	node := graph[start]

	visited := make(sets.Set[vec.Vec2i])
	visited.Put(start)
	for _, n := range node.Next {
		path, ok := walk(graph, start, n, sets.NewFrom(visited))
		if ok {
			path = append(path, start)
			return path
		}
	}
	panic("no loop found?!")
}

func walk(graph map[vec.Vec2i]*Node, parent, current vec.Vec2i, visited sets.Set[vec.Vec2i]) ([]vec.Vec2i, bool) {
	node := graph[current]
	if node.IsStart {
		panic("no loop!")
	}

	for _, c := range node.Next {
		if c == parent {
			continue
		}
		if visited.Has(c) { // needed?
			dst := graph[c]
			if !dst.IsStart {
				panic("not start? wtf?")
			}
			return []vec.Vec2i{current}, true
		}
		visited.Put(c)

		p, ok := walk(graph, current, c, sets.NewFrom(visited))
		if ok {
			p = append(p, current)
			return p, true
		}
	}

	return nil, false
}

func buildGraph(field [][]byte) map[vec.Vec2i]*Node {
	graph := make(map[vec.Vec2i]*Node)

	height := len(field)
	width := len(field[0])

	for y := range field {
		for x := range field[y] {
			from := vec.Vec2i{x, y}
			a := field[from.Y][from.X]

			var connections []vec.Vec2i

			for dx := -1; dx <= 1; dx++ {
				if dx == 0 {
					continue
				}
				d := vec.Vec2i{X: dx, Y: 0}
				to := from.Add(d)
				if to.X < 0 || to.X >= width {
					continue
				}
				b := field[to.Y][to.X]
				if isConnected(a, b, toDirection(d)) {
					connections = append(connections, to)
				}
			}
			for dy := -1; dy <= 1; dy++ {
				if dy == 0 {
					continue
				}
				d := vec.Vec2i{X: 0, Y: dy}
				to := from.Add(d)
				if to.Y < 0 || to.Y >= height {
					continue
				}
				b := field[to.Y][to.X]
				if isConnected(a, b, toDirection(d)) {
					connections = append(connections, to)
				}
			}
			graph[from] = &Node{
				Key:     from,
				Next:    connections,
				IsStart: a == START,
			}
		}
	}

	return graph
}

type Direction int

const (
	DirectionNone Direction = iota
	DirectionNorth
	DirectionEast
	DirectionWest
	DirectionSouth
)

var bendConnectivity = map[byte][]Direction{
	GROUND:     {},
	START:      {DirectionNorth, DirectionEast, DirectionSouth, DirectionWest},
	VERTICAL:   {DirectionNorth, DirectionSouth},
	HORIZONTAL: {DirectionEast, DirectionWest},
	BEND_NW:    {DirectionNorth, DirectionWest},
	BEND_NE:    {DirectionNorth, DirectionEast},
	BEND_SE:    {DirectionSouth, DirectionEast},
	BEND_SW:    {DirectionSouth, DirectionWest},
}

var inversion = map[Direction]Direction{DirectionNorth: DirectionSouth, DirectionEast: DirectionWest, DirectionWest: DirectionEast, DirectionSouth: DirectionNorth}

func invert(d Direction) Direction {
	return inversion[d]
}

var vec2dir = map[vec.Vec2i]Direction{
	{-1, 0}: DirectionWest,
	{1, 0}:  DirectionEast,
	{0, 1}:  DirectionSouth,
	{0, -1}: DirectionNorth,
}

func toDirection(v vec.Vec2i) Direction {
	return vec2dir[v]
}

func isConnected(a, b byte, direction Direction) bool {
	aConnectivity := bendConnectivity[a]
	bConnectivity := bendConnectivity[b]

	if slices.Index(aConnectivity, direction) < 0 {
		return false
	}

	if slices.Index(bConnectivity, invert(direction)) < 0 {
		return false
	}

	return true
}

func findStart(field [][]byte) vec.Vec2i {
	for y := range field {
		for x := range field[y] {
			if field[y][x] == START {
				return vec.Vec2i{X: x, Y: y}
			}
		}
	}
	panic("no start")
}
