package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"slices"

	"aoc/pkg/vec"

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

	galaxies := expand(parseStarmap(file), 1000000)

	sum := sumPairwiseDistances(galaxies)

	fmt.Printf("part two: %d\n", sum)
	if sum != 611998089572 {
		panic("bad result")
	}
}

type Galaxy struct {
	ID         int
	Coordinate vec.Vec2i
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	galaxies := expand(parseStarmap(file), 2)

	sum := sumPairwiseDistances(galaxies)
	fmt.Printf("part one: %d\n", sum)
	if sum != 10313550 {
		panic("bad result")
	}
}

func sumPairwiseDistances(galaxies []Galaxy) int {
	sum := 0
	for i, g1 := range galaxies {
		if i+1 >= len(galaxies) {
			break
		}
		others := galaxies[(i + 1):]
		for _, g2 := range others {
			sub := g1.Coordinate.Sub(g2.Coordinate)
			dist := sub.Norm1()
			sum += dist
		}
	}
	return sum
}

func locateGalaxies(sm [][]int) []Galaxy {
	var galaxies []Galaxy
	for y, row := range sm {
		for x, e := range row {
			if e > 0 {
				galaxies = append(galaxies, Galaxy{
					ID:         e,
					Coordinate: vec.Vec2i{X: x, Y: y},
				})
			}
		}
	}
	return galaxies
}

func parseStarmap(f fs.File) [][]int {
	scanner := bufio.NewScanner(f)

	id := 1
	var starmap [][]int
	for scanner.Scan() {
		text := scanner.Text()

		row := make([]int, len(text))
		for x := 0; x < len(text); x++ {
			if text[x] == '#' {
				row[x] = id
				id++
			}
		}
		starmap = append(starmap, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return starmap
}

func expand(sm [][]int, factor int) []Galaxy {
	emptyRows := findEmptyRows(sm)
	emptyColumns := findEmptyColumns(sm)

	galaxies := locateGalaxies(sm)

	expanded := make([]Galaxy, len(galaxies))
	for i, g := range galaxies {

		shift := vec.Vec2i{
			X: (factor - 1) * countLessThan(emptyColumns, g.Coordinate.X),
			Y: (factor - 1) * countLessThan(emptyRows, g.Coordinate.Y),
		}
		expanded[i] = Galaxy{
			ID:         g.ID,
			Coordinate: g.Coordinate.Add(shift),
		}
	}

	return expanded
}

func countLessThan(vals []int, d int) int {
	count := 0
	for _, v := range vals {
		if v < d {
			count++
		} else {
			break
		}
	}
	return count
}

func findEmptyColumns(sm [][]int) []int {
	var emptyColumns []int
	for x := 0; x < len(sm[0]); x++ {
		isEmpty := true
		for _, row := range sm {
			if row[x] != 0 {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			emptyColumns = append(emptyColumns, x)
		}
	}
	slices.Sort(emptyColumns)
	return emptyColumns
}

func findEmptyRows(sm [][]int) []int {
	var emptyRows []int
	for y, row := range sm {
		isEmpty := true
		for _, p := range row {
			if p != 0 {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			emptyRows = append(emptyRows, y)
		}
	}
	slices.Sort(emptyRows)
	return emptyRows
}
