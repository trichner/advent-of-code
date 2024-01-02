package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	hmap := readHeightMap("input.txt")
	partOne(hmap)
	partTwo(hmap)
}

func partOne(m *Map) {
	p, err := walk(m, m.Start, m.Target)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", len(p)-1)
}

func partTwo(m *Map) {
	var shortest []Vec2i
	for y, row := range m.HeightMap {
		for x, e := range row {
			if e == 0 {
				p, err := walk(m, Vec2i{x, y}, m.Target)
				if err != nil {
					continue
				}
				if shortest == nil || len(shortest) > len(p) {
					shortest = p
				}
			}
		}
	}
	fmt.Printf("%d\n", len(shortest)-1)
}

func walk(m *Map, start Vec2i, target Vec2i) ([]Vec2i, error) {
	h := &Queue{}

	visited := map[Vec2i]Vec2i{}
	h.Push(start, start)

	for {
		move, prev, ok := h.Pop()
		if !ok {
			return nil, fmt.Errorf("no path found")
		}
		_, ok = visited[move]
		if ok {
			continue
		}

		visited[move] = prev
		if move == target {
			break
		}
		h.PushAll(move, m.MovesFrom(move))
	}

	return calculatePath(visited, target), nil
}

func calculatePath(visited map[Vec2i]Vec2i, end Vec2i) []Vec2i {
	path := []Vec2i{end}
	pos := end
	for {
		prev, ok := visited[pos]
		if !ok {
			return nil
		}
		if pos == prev {
			break
		}
		pos = prev
		path = append(path, pos)
	}
	return path
}

func readHeightMap(fname string) *Map {
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var heightMap [][]int
	var row []int

	var start, target Vec2i
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			heightMap = append(heightMap, row)
			break
		} else if b == 'S' {
			start = Vec2i{
				X: len(row),
				Y: len(heightMap),
			}
			row = append(row, 0)
		} else if b == 'E' {
			target = Vec2i{
				X: len(row),
				Y: len(heightMap),
			}
			row = append(row, int('z'-'a'))
		} else if b == '\n' {
			heightMap = append(heightMap, row)
			row = nil
		} else {
			row = append(row, int(b)-'a')
		}
	}

	return &Map{
		HeightMap: heightMap,
		Start:     start,
		Target:    target,
	}
}

type Vec2i struct{ X, Y int }

func (v Vec2i) Add(summand Vec2i) Vec2i {
	return Vec2i{
		X: v.X + summand.X,
		Y: v.Y + summand.Y,
	}
}

type Map struct {
	HeightMap [][]int
	Start     Vec2i
	Target    Vec2i
}

func (m *Map) MovesFrom(p Vec2i) []Vec2i {
	targets := []Vec2i{
		p.Add(Vec2i{X: 1, Y: 0}),
		p.Add(Vec2i{-1, 0}),
		p.Add(Vec2i{X: 0, Y: 1}),
		p.Add(Vec2i{0, -1}),
	}
	var inReach []Vec2i
	for _, t := range targets {
		if t.X < 0 || t.X >= len(m.HeightMap[0]) || t.Y < 0 || t.Y >= len(m.HeightMap) {
			continue
		} else if m.HeightMap[t.Y][t.X]-m.HeightMap[p.Y][p.X] > 1 {
			continue
		} else {
			inReach = append(inReach, t)
		}
	}
	return inReach
}

type Queue [][]Vec2i

func (h *Queue) Push(position, previous Vec2i) {
	*h = append(*h, []Vec2i{position, previous})
}

func (h *Queue) PushAll(start Vec2i, elements []Vec2i) {
	for _, e := range elements {
		h.Push(e, start)
	}
}

func (h *Queue) Pop() (Vec2i, Vec2i, bool) {
	if len(*h) == 0 {
		return Vec2i{}, Vec2i{}, false
	}

	v := (*h)[0]
	*h = (*h)[1:]
	return v[0], v[1], true
}
