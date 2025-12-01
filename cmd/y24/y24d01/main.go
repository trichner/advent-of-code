package main

import (
	"bufio"
	"embed"
	"fmt"
	"slices"

	"aoc/pkg/in"
	"aoc/pkg/util"
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

	scanner := bufio.NewScanner(file)

	var left []int
	var right []int

	for scanner.Scan() {
		text := scanner.Text()

		var l, r int
		util.Must(fmt.Sscanf(text, "%d %d", &l, &r))

		left = append(left, l)
		right = append(right, r)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	frequencyMap := map[int]int{}
	for _, r := range right {
		current, ok := frequencyMap[r]
		if !ok {
			current = 0
		}
		frequencyMap[r] = current + 1
	}

	similarity := 0
	for _, l := range left {
		f := frequencyMap[l]
		similarity += l * f
	}

	fmt.Printf("part two: %d\n", similarity)
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var left []int
	var right []int

	for scanner.Scan() {
		text := scanner.Text()

		var l, r int
		util.Must(fmt.Sscanf(text, "%d %d", &l, &r))

		left = append(left, l)
		right = append(right, r)
	}

	slices.Sort(left)
	slices.Sort(right)

	diff := 0
	for i := range left {
		diff += abs(left[i] - right[i])
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part one: %d\n", diff)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
