package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"aoc/pkg/in"
)

func partTwo() {
	// TOO LOW: 5932134731224
	// TOO LOW: 7996215336396
	input := readInput()
	input = transpose(input)

	for _, line := range input {
		fmt.Println(string(line))
	}

	sum := 0
	var rows [][]byte
	for _, line := range input {
		if isEmpty(line) {
			sum += solveLines(rows)
			rows = nil
		} else {
			rows = append(rows, line)
		}
	}
	if len(rows) > 0 {
		sum += solveLines(rows)
	}

	fmt.Printf("part two: %d\n", sum)
}

func solveLines(lines [][]byte) int {
	op := lines[0][len(lines[0])-1]
	lines[0][len(lines[0])-1] = ' '

	if op == '+' {
		sum := 0
		for _, line := range lines {
			v, err := strconv.Atoi(strings.TrimSpace(string(line)))
			if err != nil {
				panic(err)
			}
			sum += v
		}
		return sum
	} else if op == '*' {
		sum := 1
		for _, line := range lines {
			v, err := strconv.Atoi(strings.TrimSpace(string(line)))
			if err != nil {
				panic(err)
			}
			sum *= v
		}
		return sum
	}
	panic("unknown op")
}

func isEmpty(line []byte) bool {
	for _, c := range line {
		if c != ' ' {
			return false
		}
	}
	return true
}

func transpose(grid [][]byte) [][]byte {
	if len(grid) == 0 {
		return nil
	}

	rows := len(grid)
	cols := len(grid[0])

	result := make([][]byte, cols)
	for i := range result {
		result[i] = make([]byte, rows)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			result[c][r] = grid[r][c]
		}
	}

	return result
}

func readInput() [][]byte {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines [][]byte
	maxLen := 0

	for scanner.Scan() {
		line := []byte(scanner.Text())
		lines = append(lines, line)
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Pad all lines to equal length
	for i, line := range lines {
		if len(line) < maxLen {
			padded := make([]byte, maxLen)
			copy(padded, line)
			for j := len(line); j < maxLen; j++ {
				padded[j] = ' '
			}
			lines[i] = padded
		}
	}

	return lines
}
