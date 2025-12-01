package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"math"
	"strings"

	"aoc/pkg/util"

	"aoc/pkg/in"
	"aoc/pkg/queue"
	"aoc/pkg/vec"
)

//go:embed *.txt
var inputs embed.FS

func main() {
	partOne()
	partTwo()
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	island := parse(file)
	heat := walk(island, 3, 10)

	fmt.Printf("part two: %d\n", heat)
	if heat != 1411 {
		panic("bad value")
	}
}

type key struct {
	Pos, Heading vec.Vec2i
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	island := parse(file)
	heat := walk(island, 0, 3)

	fmt.Printf("part one: %d\n", heat)
	if heat != 1263 {
		panic("bad")
	}
}

func walk(island [][]uint8, skip int, maxSteps int) int {
	size := vec.Vec2i{len(island[0]), len(island)}
	boundingBox := vec.AABB{To: size.Sub(vec.Vec2i{1, 1})}

	best := make(map[key]int, 1024)

	var todo queue.Queue[[]vec.Vec2i]

	k := key{Pos: vec.Vec2i{}, Heading: HeadRight}
	best[k] = 0

	k = key{Pos: vec.Vec2i{}, Heading: HeadDown}
	best[k] = 0

	todo.Push([]vec.Vec2i{{}, HeadRight})
	todo.Push([]vec.Vec2i{{}, HeadDown})

	for {
		current, ok := todo.Pop()
		if !ok {
			break
		}

		cpos := current[0]
		cheading := current[1]

		k := key{Pos: cpos, Heading: cheading}
		cost, ok := best[k]
		cost = util.MustOk(cost, ok)

		if cpos == boundingBox.To {
			continue
		}

		headings := nextHeadings(cheading)
		for i := 0; i < maxSteps; i++ {

			cpos = cpos.Add(cheading)
			if !boundingBox.Contains(cpos) {
				break
			}
			cost += int(island[cpos.Y][cpos.X])
			if i < skip {
				// need to travel at least skip + 1
				continue
			}
			for _, h := range headings {

				k := key{Pos: cpos, Heading: h}
				currCost, ok := best[k]
				if ok && currCost <= cost {
					continue
				}
				best[k] = cost

				todo.Push([]vec.Vec2i{cpos, h})
			}
		}
	}

	heat := findMin(best, boundingBox.To)

	return heat
}

func findMin(best map[key]int, dst vec.Vec2i) int {
	heat := math.MaxInt
	for _, h := range []vec.Vec2i{HeadUp, HeadLeft, HeadRight, HeadDown} {
		b, ok := best[key{
			Pos:     dst,
			Heading: h,
		}]
		if ok {
			heat = min(heat, b)
		}
	}
	return heat
}

var (
	HeadUp    = vec.Vec2i{0, -1}
	HeadDown  = vec.Vec2i{0, 1}
	HeadLeft  = vec.Vec2i{-1, 0}
	HeadRight = vec.Vec2i{1, 0}
)

var lut = map[vec.Vec2i][]vec.Vec2i{
	HeadUp:    {HeadRight, HeadLeft},
	HeadDown:  {HeadRight, HeadLeft},
	HeadLeft:  {HeadUp, HeadDown},
	HeadRight: {HeadUp, HeadDown},
}

func nextHeadings(heading vec.Vec2i) []vec.Vec2i {
	return lut[heading]
}

func stringIsland(island [][]uint8) string {
	var buf strings.Builder

	for y, row := range island {
		fmt.Fprintf(&buf, "%3d |", y)
		for _, v := range row {
			fmt.Fprintf(&buf, "%d", v)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func stringBest(size vec.Vec2i, best map[key]int) string {
	var buf strings.Builder

	for y := 0; y < size.Y; y++ {
		fmt.Fprintf(&buf, "%3d |", y)
		for x := 0; x < size.X; x++ {
			heat := math.MaxInt
			for _, h := range []vec.Vec2i{HeadUp, HeadLeft, HeadRight, HeadDown} {
				b, ok := best[key{
					Pos:     vec.Vec2i{x, y},
					Heading: h,
				}]
				if ok {
					heat = min(heat, b)
				}
			}
			if heat < math.MaxInt {
				fmt.Fprintf(&buf, " %3d", heat)
			} else {
				fmt.Fprintf(&buf, "  - ")
			}
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func parse(r io.Reader) [][]uint8 {
	scanner := bufio.NewScanner(r)

	var island [][]uint8
	for scanner.Scan() {
		text := scanner.Text()

		row := make([]uint8, len(text))
		for i := range row {
			n := text[i] - '0'
			row[i] = n
		}

		island = append(island, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return island
}
