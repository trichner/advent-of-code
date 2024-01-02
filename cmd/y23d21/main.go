package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"

	"aoc/pkg/in"
	"aoc/pkg/progress"
	"aoc/pkg/queue"
	"aoc/pkg/vec"
)

//go:embed *.txt
var inputs embed.FS

func main() {
	partOne()
	partTwo()
}

var (
	UP    = vec.Vec2i{0, -1}
	DOWN  = vec.Vec2i{0, 1}
	LEFT  = vec.Vec2i{-1, 0}
	RIGHT = vec.Vec2i{1, 0}
)

var headings = []vec.Vec2i{UP, DOWN, LEFT, RIGHT}

type Step struct {
	Pos  vec.Vec2i
	Cost int
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	garden := parse(file)

	pointXs := []int{
		65,
		65 + garden.GridSize,
		65 + 2*garden.GridSize,
	}

	points := []vec.Vec2i{}
	for _, maxCost := range pointXs {
		count := walkBfs(garden, maxCost)
		points = append(points, vec.Vec2i{maxCost, count})
	}

	steps := 26501365
	count := solveQuadratic(garden, points, steps)

	fmt.Printf("part two: %d\n", count)
	if count != 636350496972143 {
		panic("bad")
	}
}

func solveQuadratic(garden *InfGarden, points []vec.Vec2i, steps int) int {
	// with a lot of inspiration from
	// https://colab.research.google.com/github/derailed-dash/Advent-of-Code/blob/master/src/AoC_2023/Dazbo's_Advent_of_Code_2023.ipynb#scrollTo=1kgtFCEOXyaH

	c := points[0].Y
	b := (4*points[1].Y - 3*points[0].Y - points[2].Y) / 2
	a := points[1].Y - points[0].Y - b

	x := (steps - garden.GridSize/2) / garden.GridSize

	// fmt.Printf("%d * %d^2 + %d * %d + %d\n", a, x, b, x, c)
	return a*x*x + b*x + c
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	garden := parse(file)

	steps := 64

	count := walkBfs(garden, steps)

	fmt.Printf("part one: %d\n", count)

	if count != 3858 {
		panic("bad result")
	}
}

func walkBfs(garden *InfGarden, maxCost int) int {
	pending := &queue.Queue[Step]{}

	pending.Push(Step{Pos: garden.Start})

	visited := make(map[vec.Vec2i]int)

	visited[garden.Start] = 0

	highest := 0
	var progressTicker progress.ProgressTicker

	for {
		current, ok := pending.Pop()
		if !ok {
			break
		}

		cost := current.Cost

		if cost >= maxCost {
			continue
		}

		cost++

		if cost%512 == 0 && progressTicker.ShouldUpdate() {
			highest = max(highest, cost)
			fmt.Printf("cost: %8d %f%%\n", cost, float32(cost)/float32(maxCost)*100)
		}

		for _, h := range headings {
			next := current.Pos.Add(h)

			if garden.Get(next) == TileRock {
				continue
			}

			_, alreadyVisited := visited[next]
			if alreadyVisited {
				continue
			}

			visited[next] = cost

			pending.Push(Step{
				Pos:  next,
				Cost: cost,
			})
		}

	}

	return countOptions(visited, maxCost)
}

type InfGarden struct {
	Tiles     [][]Tile
	Start     vec.Vec2i
	GridSize  int
	MaxRadius int
}

func (g *InfGarden) Get(p vec.Vec2i) Tile {
	return g.Tiles[mod(p.Y, g.GridSize)][mod(p.X, g.GridSize)]
}

func mod(a, m int) int {
	r := a % m
	if r >= 0 {
		return r
	}
	return r + m
}

func countOptions(visited map[vec.Vec2i]int, target int) int {
	count := 0
	for _, v := range visited {
		if v > target {
			continue
		} else if v == target {
			count++
		} else if (target-v)%2 == 0 {
			count++
		}
	}

	return count
}

type Tile byte

const (
	TileGarden Tile = iota
	TileRock
	TileStart
)

func parse(r io.Reader) *InfGarden {
	scanner := bufio.NewScanner(r)

	var tiles [][]Tile

	for scanner.Scan() {
		text := scanner.Text()
		tiles = append(tiles, parseRow(text))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	start := findStart(tiles)

	if len(tiles) != len(tiles[0]) {
		panic("can only deal with squares")
	}

	return &InfGarden{
		Tiles:    tiles,
		Start:    start,
		GridSize: len(tiles),
	}
}

func findStart(tiles [][]Tile) vec.Vec2i {
	for y, row := range tiles {
		for x, t := range row {
			if t == TileStart {
				return vec.Vec2i{x, y}
			}
		}
	}
	panic("no start found")
}

func parseRow(s string) []Tile {
	fields := []byte(s)
	row := make([]Tile, len(fields))
	for i, b := range fields {
		if b == '#' {
			row[i] = TileRock
		} else if b == '.' {
			row[i] = TileGarden
		} else if b == 'S' {
			row[i] = TileStart
		} else {
			panic("wut?")
		}
	}
	return row
}
