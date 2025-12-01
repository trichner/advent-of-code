package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"time"

	"aoc/pkg/in"
	"aoc/pkg/vec"
)

//go:embed *.txt
var inputs embed.FS

func main() {
	start := time.Now()
	partOne()
	elapsed := time.Since(start)
	fmt.Printf("executed in: %s\n", elapsed)
	start = time.Now()
	partTwo()
	elapsed = time.Since(start)
	fmt.Printf("executed in: %s\n", elapsed)
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	room := readRoom(file)
	start := findStart(room)
	sum := 0
	for x := 0; x < len(room[0]); x++ {
		for y := 0; y < len(room); y++ {
			if isLoop(room, start, vec.Vec2i{x, y}) {
				sum++
			}
		}
	}

	fmt.Printf("part two: %d\n", sum)
}

type loopKey struct {
	pos, dir vec.Vec2i
}

func isLoop(room [][]byte, start vec.Vec2i, obstruction vec.Vec2i) bool {
	if start == obstruction {
		return false
	}

	if room[obstruction.Y][obstruction.X] == '#' {
		return false
	}

	visited := map[loopKey]struct{}{}

	boundary := vec.AABB{From: vec.Vec2i{}, To: vec.Vec2i{len(room[0]) - 1, len(room) - 1}}
	dir := vec.Vec2i{0, -1}

	pos := start
	for {
		_, ok := visited[loopKey{pos, dir}]
		if ok {
			return true
		}

		visited[loopKey{pos, dir}] = struct{}{}
		next := pos.Add(dir)
		if !boundary.Contains(next) {
			return false
		}
		for room[next.Y][next.X] == '#' || next == obstruction {
			dir = vec.NewRotCW().Mul(dir)
			next = pos.Add(dir)
		}
		pos = pos.Add(dir)
	}
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	room := readRoom(file)

	sum := walk(room)

	fmt.Printf("part one: %d\n", sum)
}

func walk(room [][]byte) int {
	pos := findStart(room)

	freq := make([][]byte, len(room))
	for i := range freq {
		freq[i] = make([]byte, len(room[i]))
	}

	boundary := vec.AABB{From: vec.Vec2i{}, To: vec.Vec2i{len(room[0]) - 1, len(room) - 1}}
	dir := vec.Vec2i{0, -1}

	for {
		freq[pos.Y][pos.X]++
		next := pos.Add(dir)
		if !boundary.Contains(next) {
			break
		}
		if room[next.Y][next.X] == '#' {
			dir = vec.NewRotCW().Mul(dir)
			pos = pos.Add(dir)
		} else {
			pos = next
		}
	}

	sum := 0
	for _, row := range freq {
		for _, f := range row {
			if f > 0 {
				sum++
			}
		}
	}
	return sum
}

func readRoom(r io.Reader) [][]byte {
	scanner := bufio.NewScanner(r)

	var room [][]byte

	for scanner.Scan() {
		text := scanner.Text()
		room = append(room, []byte(text))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return room
}

func findStart(room [][]byte) vec.Vec2i {
	for y, row := range room {
		for x, e := range row {
			if e == '^' {
				return vec.Vec2i{x, y}
			}
		}
	}
	panic("start not found")
}
