package main

import (
	"bufio"
	"fmt"

	"aoc/pkg/in"
)

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0

	lines := make([][]int, 3)
	for scanner.Scan() {
		text := scanner.Text()
		width := len(text) + 2 // border to make parsing easier
		shiftLines(lines, width)
		scanLineForGears(lines[2], text)

		sum += checkLinePartTwo(lines)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("%d\n", sum)
	if sum != 81463996 {
		panic("wrong result!")
	}
}

func checkLinePartTwo(lines [][]int) int {
	sum := 0

	width := len(lines[0])
	for i := 1; i < width-1; i++ {
		if lines[1][i] != sPart {
			continue
		}

		candidates := map[int]int{}

		calculatePart(lines, i, candidates)
		if len(candidates) != 2 {
			continue
		}

		mul := 1
		for _, v := range candidates {
			mul *= v
		}
		sum += mul
	}

	return sum
}

func scanLineForGears(line []int, text string) {
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c >= '0' && c <= '9' {
			line[i+1] = int(c - '0')
		} else if c == '*' {
			line[i+1] = sPart
		}
	}
}
