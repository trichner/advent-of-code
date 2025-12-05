package main

import (
	"bufio"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"aoc/pkg/in"
)

type interval struct {
	Start, End int
}

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
	ivals, _ := readInput()

	ivals = compactIntervals(ivals)

	sum := 0
	for _, ival := range ivals {
		sum += ival.End - ival.Start + 1
	}

	fmt.Printf("part two: %d\n", sum)
}

func compactIntervals(ivals []*interval) []*interval {
	if len(ivals) == 0 {
		return nil
	}

	sort.Slice(ivals, func(i, j int) bool {
		return ivals[i].Start < ivals[j].Start
	})

	result := []*interval{{Start: ivals[0].Start, End: ivals[0].End}}

	for i := 1; i < len(ivals); i++ {
		curr := ivals[i]
		last := result[len(result)-1]

		if curr.Start <= last.End+1 {
			last.End = max(last.End, curr.End)
		} else {
			result = append(result, &interval{Start: curr.Start, End: curr.End})
		}
	}

	return result
}

func partOne() {
	ivals, ids := readInput()

	sum := 0

	for _, id := range ids {
		for _, ival := range ivals {
			if ival.Start <= id && id <= ival.End {
				sum += 1
				break
			}
		}
	}

	fmt.Printf("part one: %d\n", sum)
}

func readInput() ([]*interval, []int) {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var intervals []*interval
	var numbers []int
	parsingIntervals := true

	for scanner.Scan() {
		text := scanner.Text()

		if text == "" {
			parsingIntervals = false
			continue
		}

		if parsingIntervals {
			parts := strings.Split(text, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			intervals = append(intervals, &interval{Start: start, End: end})
		} else {
			n, _ := strconv.Atoi(text)
			numbers = append(numbers, n)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return intervals, numbers
}
