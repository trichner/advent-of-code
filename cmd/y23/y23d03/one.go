package main

import (
	"bufio"
	"fmt"

	"aoc/pkg/in"
)

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0

	lines := make([][]int, 3)
	for scanner.Scan() {
		text := scanner.Text()
		width := len(text) + 2 // border to make parsing easier
		shiftLines(lines, width)
		scanLineForParts(lines[2], text)

		sum += checkLinePartOne(lines)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("%d\n", sum)

	if sum != 527144 {
		panic("wrong result!")
	}
}

func checkLinePartOne(lines [][]int) int {
	candidates := map[int]int{}

	width := len(lines[0])
	for i := 1; i < width-1; i++ {
		if lines[1][i] != sPart {
			continue
		}
		calculatePart(lines, i, candidates)
	}
	sum := 0
	for _, v := range candidates {
		sum += v
	}
	return sum
}

func scanLineForParts(line []int, text string) {
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c >= '0' && c <= '9' {
			line[i+1] = int(c - '0')
		} else if c != '.' {
			line[i+1] = sPart
		}
	}
}
