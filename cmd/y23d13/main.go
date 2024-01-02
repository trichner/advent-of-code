package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"slices"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

type FloorTile byte

const (
	FloorAsh FloorTile = iota
	FloorRock
)

func main() {
	partOne()
	partTwo()
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	patterns := parse(file)

	sum := 0
	for _, original := range patterns {
		oh, ov := findMirrorAxis(original)
		for _, pattern := range permutePattern(original) {
			h, v := findMirrorAxis(pattern)
			if len(h) == 0 && len(v) == 0 {
				// no reflections at all
				continue
			}
			if slices.Equal(oh, h) && slices.Equal(ov, v) {
				// exactly same reflections
				continue
			}
			dh := sub(h, oh)
			if len(dh) > 0 {
				sum += dh[0] * 100
				break
			}

			dv := sub(v, ov)
			if len(dv) > 0 {
				sum += dv[0]
				break
			}
			panic("wut?")
		}

	}

	fmt.Printf("part two: %d\n", sum)
	if sum != 32728 {
		panic("wut?")
	}
}

func permutePattern(pattern [][]FloorTile) [][][]FloorTile {
	var permutations [][][]FloorTile
	for y, row := range pattern {
		for x := range row {
			copied := shallowCopy(pattern)

			newRow := make([]FloorTile, len(row))
			copy(newRow, row)
			v := newRow[x]
			if v == FloorAsh {
				v = FloorRock
			} else {
				v = FloorAsh
			}
			newRow[x] = v

			copied[y] = newRow
			permutations = append(permutations, copied)
		}
	}

	return permutations
}

func shallowCopy[S ~[]E, E any](s S) S {
	copied := make(S, len(s))
	copy(copied, s)
	return copied
}

func sub[S ~[]E, E comparable](s1 S, s2 S) S {
	for _, e := range s2 {
		i := slices.Index(s1, e)
		if i < 0 {
			continue
		}
		s1 = slices.Delete(s1, i, i+1)
	}
	return s1
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()
	patterns := parse(file)

	sum := 0
	for _, p := range patterns {
		h, v := findMirrorAxis(p)
		if len(h) > 0 {
			sum += 100 * h[0]
		} else if len(v) > 0 {
			sum += v[0]
		}
	}

	fmt.Printf("part one: %d\n", sum)
	// 21272
	if sum != 27742 {
		panic("bad result")
	}
}

func findMirrorAxis(tiles [][]FloorTile) ([]int, []int) {
	horizontalAxis := findAllMirrorHorizontal(tiles)

	pt := transpose(tiles)
	verticalAxis := findAllMirrorHorizontal(pt)
	return horizontalAxis, verticalAxis
}

func findAllMirrorHorizontal(tiles [][]FloorTile) []int {
	var axis []int

	for y := 1; y < len(tiles); y++ {
		for offset := 0; offset < len(tiles)-1; offset++ {
			if y-1-offset < 0 || y+offset >= len(tiles) {
				axis = append(axis, y)
				break
			}
			upper := tiles[y-1-offset]
			lower := tiles[y+offset]
			if !slices.Equal(upper, lower) {
				// doesn't seem to be a valid axis
				break
			}
		}
	}

	return axis
}

func transpose(tiles [][]FloorTile) [][]FloorTile {
	transposed := make([][]FloorTile, len(tiles[0]))

	for ny := range transposed {
		row := make([]FloorTile, len(tiles))
		for nx := range row {
			row[nx] = tiles[len(row)-nx-1][ny]
		}
		transposed[ny] = row
	}

	return transposed
}

func printPattern(tiles [][]FloorTile) {
	fmt.Println("---")
	for y, row := range tiles {
		fmt.Printf("%2d ", y+1)
		for _, tile := range row {
			if tile == FloorAsh {
				fmt.Print(".")
			} else if tile == FloorRock {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
	fmt.Println("---")
}

func parse(f fs.File) [][][]FloorTile {
	scanner := bufio.NewScanner(f)

	var patterns [][][]FloorTile

	var lines []string
	for scanner.Scan() {
		text := scanner.Text()

		if text == "" {
			p := parsePattern(lines)
			patterns = append(patterns, p)
			lines = nil
			continue
		}
		lines = append(lines, text)
	}
	if lines != nil {
		p := parsePattern(lines)
		patterns = append(patterns, p)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return patterns
}

func parsePattern(lines []string) [][]FloorTile {
	tiles := make([][]FloorTile, len(lines))

	for y := range tiles {
		l := lines[y]
		row := make([]FloorTile, len(l))
		for x := 0; x < len(l); x++ {
			c := l[x]
			if c == '.' {
				row[x] = FloorAsh
			} else if c == '#' {
				row[x] = FloorRock
			}
		}

		tiles[y] = row
	}

	return tiles
}
