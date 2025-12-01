package main

import (
	"bufio"
	"embed"
	"fmt"
	"strconv"
	"strings"

	"aoc/pkg/sio"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

// d = (a * t1) * t2
// t = t1 + t2
// t2 = t - t1

// d = (a * t1) * (t - t1)
// d = (a * t1) * t - (a * t1) * (t - t1)
// d = (a * t1) * t - (a * t1) * t1

func main() {
	partOne()
	partTwo()
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	times := scanner.Text()
	times = strings.TrimPrefix(times, "Time: ")
	t := parseJoined(times)

	scanner.Scan()
	durations := scanner.Text()
	durations = strings.TrimPrefix(durations, "Distance: ")
	minD := parseJoined(durations)

	candidates := countWinningGames(t, minD)

	fmt.Printf("part two: %d\n", candidates)
	if candidates != 23654842 {
		panic("bad result")
	}
}

func parseJoined(l string) int {
	durationsSplits := strings.Fields(l)
	r, err := strconv.Atoi(strings.Join(durationsSplits, ""))
	if err != nil {
		panic(err)
	}
	return r
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	times := scanner.Text()
	times = strings.TrimPrefix(times, "Time: ")
	timesNum := sio.IntFields(times)

	scanner.Scan()
	durations := scanner.Text()
	durations = strings.TrimPrefix(durations, "Distance: ")
	durationsNum := sio.IntFields(durations)

	var games [][]int
	for i := range timesNum {
		games = append(games, []int{timesNum[i], durationsNum[i]})
	}

	result := 1
	for _, g := range games {
		t := g[0]
		minD := g[1]
		candidates := countWinningGames(t, minD)
		result *= candidates
	}

	fmt.Printf("part one: %d\n", result)
	if result != 303600 {
		panic("bad result")
	}
}

func countWinningGames(t, minD int) int {
	candidates := 0
	for ta := 1; ta < t; ta++ {
		r := playGame(t, ta)
		if r > minD {
			candidates++
		}
	}
	return candidates
}

func playGame(t, ta int) int {
	// d = t1 * t - t1 * t1
	return ta*t - ta*ta
}
