package main

import (
	"bufio"
	"embed"
	"fmt"
	"math"
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
	banks := readBanks()

	sum := 0
	for _, bank := range banks {
		cm := maxBattery(bank, 11)
		sum += cm
	}

	// 171518260283767
	fmt.Printf("part two: %d\n", sum)
}

func partOne() {
	banks := readBanks()

	sum := 0
	for _, bank := range banks {
		sum += maxBattery(bank, 1)
	}

	fmt.Printf("part one: %d\n", sum)
}

func maxBattery(bank []int, depth int) int {
	memo := make(map[[2]int]int)
	return maxBatteryRec(bank, 0, depth, memo)
}

func maxBatteryRec(bank []int, start, depth int, memo map[[2]int]int) int {
	if start >= len(bank) {
		return -1
	}

	key := [2]int{start, depth}
	if cached, ok := memo[key]; ok {
		return cached
	}

	var candidate int
	if depth == 0 {
		candidate = -1
		for i := start; i < len(bank); i++ {
			candidate = max(bank[i], candidate)
		}
	} else {
		candidate = -1
		for i := start; i < len(bank); i++ {
			maxRemaining := maxBatteryRec(bank, i+1, depth-1, memo)
			if maxRemaining < 0 {
				continue
			}
			candidate = max(int(math.Pow10(depth))*bank[i]+maxRemaining, candidate)
		}
	}

	memo[key] = candidate
	return candidate
}

func readBanks() [][]int {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var banks [][]int
	for scanner.Scan() {
		text := scanner.Text()

		batteries := make([]int, len(text))

		for i, c := range text {
			batteries[i] = int(c - '0')
		}

		banks = append(banks, batteries)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return banks
}
