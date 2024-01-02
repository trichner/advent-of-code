package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"

	"aoc/pkg/in"
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

	raw := parse(file)
	tiles := NewTiles(raw)

	totalCycles := 1_000_000_000
	seen := map[string]int{}
	for i := 0; i < totalCycles; i++ {

		tiles = TiltCycle(tiles)
		k := tiles.Hash()
		prev, ok := seen[k]
		if ok {
			w := WeighTiles(tiles)
			loopLen := i - prev
			if (totalCycles-i-1)%loopLen == 0 {
				fmt.Printf("part two: %d\n", w)
				if w != 104533 {
					panic("bad result")
				}
				return
			}
		}
		seen[k] = i
	}
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	raw := parse(file)
	tiles := NewTiles(raw)

	tiles = TiltNorth(tiles)
	sum := WeighTiles(tiles)

	fmt.Printf("part one: %d\n", sum)
	if sum != 108813 {
		panic("bad result")
	}
}

func parse(f fs.File) [][]Tile {
	scanner := bufio.NewScanner(f)

	var tiles [][]Tile

	for scanner.Scan() {
		text := scanner.Text()
		tiles = append(tiles, parseLine(text))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return tiles
}

func parseLine(l string) []Tile {
	return []Tile(l)
}
