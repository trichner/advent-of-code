package main

import (
	"bufio"
	"embed"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func partOne() {
	worksheet := readWorksheetPartOne()

	sum := calcWorksheet(worksheet)
	fmt.Printf("part one: %d\n", sum)
}

func calcWorksheet(worksheet []*column) int {
	sum := 0
	for _, w := range worksheet {
		if w.op == '+' {
			r := 0
			for _, n := range w.numbers {
				r += n
			}
			sum += r
		} else if w.op == '*' {
			r := 1
			for _, n := range w.numbers {
				r *= n
			}
			sum += r
		}
	}
	return sum
}

type column struct {
	op      byte
	numbers []int
}

func readWorksheetPartOne() []*column {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if len(lines) == 0 {
		return nil
	}

	// Last line contains operators, determine column width from it
	opLine := lines[len(lines)-1]
	numLines := lines[:len(lines)-1]

	ops := splitTrim(opLine)

	var parsed [][]int

	for _, line := range numLines {
		parsed = append(parsed, splitParse(line))
	}

	columns := make([]*column, len(parsed[0]))
	for _, row := range parsed {
		for i, e := range row {
			if columns[i] == nil {
				columns[i] = &column{
					op:      ops[i][0],
					numbers: nil,
				}
			}
			columns[i].numbers = append(columns[i].numbers, e)
		}

	}

	return columns
}

func splitTrim(s string) []string {
	splits := strings.Split(s, " ")

	var elements []string
	for _, e := range splits {
		e = strings.TrimSpace(e)
		if len(e) > 0 {
			elements = append(elements, e)
		}
	}
	return elements
}

func splitParse(s string) []int {
	splits := splitTrim(s)

	nums := make([]int, len(splits))
	for i, s := range splits {
		var err error
		nums[i], err = strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
	}
	return nums
}
