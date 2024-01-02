package main

import (
	"embed"
	"fmt"
)

//go:embed *.txt
var inputs embed.FS

const (
	sPart = -1
	sVoid = -2
)

func main() {
	fmt.Println("part one: ")
	partOne()
	fmt.Println("part two: ")
	partTwo()
}

func calculatePart(lines [][]int, i int, candidates map[int]int) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			x := i + dx
			y := 1 + dy
			c := lines[y][x]
			if c >= 0 {
				start, num := parsePartNumber(lines[y], x)
				candidates[10*start+y] = num
			}
		}
	}
}

// each number is represented by its start position and number
func parsePartNumber(line []int, x int) (int, int) {
	c := line[x]
	for {
		x--
		c = line[x]
		if c < 0 {
			break
		}
	}
	x++

	start := x

	num := 0
	c = line[x]
	for {
		num += c
		next := line[x+1]
		if next < 0 {
			break
		}
		c = next
		num *= 10
		x++
	}

	return start, num
}

func shiftLines(lines [][]int, width int) {
	if lines[0] == nil {
		lines[0] = makeVoid(width)
		lines[1] = makeVoid(width)
	} else {
		lines[0] = lines[1]
		lines[1] = lines[2]
	}
	lines[2] = makeVoid(width)
}

func printLines(lines [][]int) {
	for y := range lines {
		for x := range lines[y] {
			c := lines[y][x]
			if c == sVoid {
				fmt.Print(".")
			} else if c == sPart {
				fmt.Print("*")
			} else {
				fmt.Printf("%d", c)
			}
		}
		fmt.Println()
	}
}

func makeVoid(s int) []int {
	l := make([]int, s)
	for i := range l {
		l[i] = sVoid
	}
	return l
}
