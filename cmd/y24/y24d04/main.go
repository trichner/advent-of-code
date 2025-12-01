package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"

	"aoc/pkg/in"
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
	words := parseInput(file)
	sum := 0
	for x := 0; x < len(words[0]); x++ {
		for y := 0; y < len(words); y++ {
			at := vec.Vec2i{x, y}
			if checkMasAt(words, at) {
				sum++
			}
		}
	}

	fmt.Printf("part two: %d\n", sum)
}

var mas = []byte("MAS")

func checkMasAt(words [][]byte, at vec.Vec2i) bool {
	if words[at.Y][at.X] != 'A' {
		return false
	}

	return (checkDirectionForMas(words, at, vec.Vec2i{-1, -1}) || checkDirectionForMas(words, at, vec.Vec2i{1, 1})) && (checkDirectionForMas(words, at, vec.Vec2i{-1, 1}) || checkDirectionForMas(words, at, vec.Vec2i{1, -1}))
}

func checkDirectionForMas(words [][]byte, at, dir vec.Vec2i) bool {
	boundary := vec.AABB{
		From: vec.Vec2i{0, 0},
		To:   vec.Vec2i{len(words[0]) - 1, len(words) - 1},
	}

	point := at.Sub(dir)
	for i := 0; i < len(mas); i++ {
		if !boundary.Contains(point) {
			return false
		}

		expected := mas[i]
		actual := words[point.Y][point.X]
		if actual != expected {
			return false
		}

		point = point.Add(dir)
	}
	return true
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()
	words := parseInput(file)

	sum := 0
	for x := 0; x < len(words[0]); x++ {
		for y := 0; y < len(words); y++ {
			c := words[y][x]
			if c != 'X' {
				continue
			}

			sum += countXmasAt(words, vec.Vec2i{x, y})
			// search
		}
	}

	fmt.Printf("part one: %d\n", sum)
}

var xmas = []byte("XMAS")

func countXmasAt(words [][]byte, start vec.Vec2i) int {
	c := words[start.Y][start.X]
	if c != xmas[0] {
		return 0
	}

	sum := 0

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			dir := vec.Vec2i{x, y}
			if checkDirectionForXmas(words, start, dir) {
				sum++
			}
		}
	}

	return sum
}

func checkDirectionForXmas(words [][]byte, start, dir vec.Vec2i) bool {
	boundary := vec.AABB{
		From: vec.Vec2i{0, 0},
		To:   vec.Vec2i{len(words[0]) - 1, len(words) - 1},
	}

	point := start
	for i := 0; i < len(xmas); i++ {
		if !boundary.Contains(point) {
			return false
		}

		expected := xmas[i]
		actual := words[point.Y][point.X]
		if actual != expected {
			return false
		}

		point = point.Add(dir)
	}

	return true
}

func parseInput(r io.Reader) [][]byte {
	scanner := bufio.NewScanner(r)

	var words [][]byte
	for scanner.Scan() {
		text := scanner.Text()
		wordLine := scanLine(text)
		words = append(words, wordLine)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return words
}

func scanLine(s string) []byte {
	words := make([]byte, len(s))
	for i, c := range s {
		words[i] = byte(c)
	}
	return words
}
