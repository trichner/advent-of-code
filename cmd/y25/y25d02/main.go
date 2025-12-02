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

func partTwo() {
	ivals := readInputIntervals()

	sum := 0
	for _, i := range ivals {
		start := i[0]
		end := i[1]

		for j := start; j <= end; j++ {
			if hasRepetitions(j) {
				sum += j
			}
		}

	}

	// 15704845910
	fmt.Printf("part two: %d\n", sum)
}

func readInputIntervals() [][]int {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	ivals := make([][]int, 0)
	for scanner.Scan() {
		text := scanner.Text()
		splits := strings.Split(text, ",")
		for _, s := range splits {
			isplits := strings.Split(s, "-")
			if len(isplits) != 2 {
				continue
			}
			l, _ := strconv.Atoi(isplits[0])
			r, _ := strconv.Atoi(isplits[1])
			ivals = append(ivals, []int{l, r})
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return ivals
}

func partOne() {
	ivals := readInputIntervals()
	sum := 0
	for _, i := range ivals {
		start := i[0]
		end := i[1]

		for j := start; j <= end; j++ {
			if isDoubled(j) {
				sum += j
			}
		}

	}

	// 5398419778
	fmt.Printf("part one: %d\n", sum)
}

func isDoubled(n int) bool {
	s := strconv.Itoa(n)
	w := len(s)
	if w%2 != 0 {
		return false
	}

	l := s[:(w / 2)]
	r := s[(w / 2):]
	return l == r
}

func hasRepetitions(n int) bool {
	s := strconv.Itoa(n)
	w := len(s)

	for i := 1; i <= w/2; i++ {
		if hasRepetitionsWithPartSize(s, i) {
			return true
		}
	}

	return false
}

func hasRepetitionsWithPartSize(num string, partWidth int) bool {
	w := len(num)
	if w%partWidth != 0 {
		return false
	}

	var parts []string
	for j := 0; j < w; j += partWidth {
		parts = append(parts, num[j:j+partWidth])
	}

	for j := 0; j < len(parts)-1; j++ {
		if parts[j] != parts[j+1] {
			return false
		}
	}
	return true
}
