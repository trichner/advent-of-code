package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"time"

	"aoc/pkg/sets"
	"aoc/pkg/vec"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

func main() {
	start := time.Now()
	partOne()
	partTwo()
	elapsed := time.Since(start)
	fmt.Printf("executed in: %s\n", elapsed)
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	antennas, bounds := readMap(file)

	antinodes := sets.New[vec.Vec2i]()
	for _, towers := range antennas {
		nodes := findAntinodesP2(towers, bounds)
		antinodes.PutAll(nodes.Keys())
	}

	fmt.Printf("part two: %d\n", antinodes.Size())
}

func findAntinodesP2(towers []vec.Vec2i, bounds vec.AABB) sets.Set[vec.Vec2i] {
	antinodes := sets.New[vec.Vec2i]()
	for i, t1 := range towers {
		for j := i + 1; j < len(towers); j++ {
			t2 := towers[j]

			d := t1.Sub(t2)

			t := gcd(d.X, d.Y)
			d = vec.Vec2i{d.X / t, d.Y / t}

			p := t1
			for bounds.Contains(p) {
				antinodes.Put(p)
				p = p.Add(d)
			}

			p = t1
			for bounds.Contains(p) {
				antinodes.Put(p)
				p = p.Sub(d)
			}
		}
	}
	return antinodes
}

func gcd(a, b int) int {
	for b != 0 {
		h := a % b
		a = b
		b = h
	}
	return a
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	antennas, bounds := readMap(file)

	antinodes := sets.New[vec.Vec2i]()
	for _, towers := range antennas {
		nodes := findAntinodesP1(towers, bounds)
		antinodes.PutAll(nodes.Keys())
	}
	fmt.Printf("part one: %d\n", antinodes.Size())
}

func findAntinodesP1(towers []vec.Vec2i, bounds vec.AABB) sets.Set[vec.Vec2i] {
	antinodes := sets.New[vec.Vec2i]()
	for i, t1 := range towers {
		for j := i + 1; j < len(towers); j++ {
			t2 := towers[j]

			d := t2.Sub(t1)

			p1 := t1.Sub(d)
			p2 := t2.Add(d)
			if bounds.Contains(p1) {
				antinodes.Put(p1)
			}
			if bounds.Contains(p2) {
				antinodes.Put(p2)
			}
		}
	}
	return antinodes
}

func readMap(r io.Reader) (map[byte][]vec.Vec2i, vec.AABB) {
	scanner := bufio.NewScanner(r)

	signals := make(map[byte][]vec.Vec2i)

	y := 0
	maxX := 0
	for scanner.Scan() {
		text := scanner.Text()

		row := []byte(text)
		maxX = max(maxX, len(row))
		for x, v := range row {
			if !isAlphaNumeric(v) {
				continue
			}
			p := vec.Vec2i{x, y}
			s := signals[v]
			s = append(s, p)
			signals[v] = s
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return signals, vec.AABB{To: vec.Vec2i{maxX - 1, y - 1}}
}

func isAlphaNumeric(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
