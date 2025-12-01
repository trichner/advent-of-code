package main

import (
	"bufio"
	"embed"
	"fmt"
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
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	start := 50
	pw := 0

	for scanner.Scan() {
		text := scanner.Text()

		var dir rune
		var count int
		_, err := fmt.Sscanf(text, "%c%d", &dir, &count)
		if err != nil {
			panic(err)
		}

		next := start

		if dir == 'R' {
			next = mod100(next + count)

			//wraparound?
			if next < start && start != 0 {
				pw++
			} else if next == 0 {
				pw++
			}

		} else if dir == 'L' {
			next = mod100(next - count)

			//wraparound?
			if start < next && start != 0 {
				pw++
			} else if next == 0 {
				pw++
			}
		} else {
			panic("wut?")
		}

		if count > 100 {
			pw += count / 100
		}

		start = next
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// 6561
	fmt.Printf("part two: %d\n", pw)
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	start := 50

	pw := 0
	for scanner.Scan() {
		text := scanner.Text()

		var dir rune
		var count int
		_, err := fmt.Sscanf(text, "%c%d", &dir, &count)
		if err != nil {
			panic(err)
		}

		if dir == 'R' {
			start += count
		} else if dir == 'L' {
			start -= count
		} else {
			panic("wut?")
		}
		start = mod100(start)

		if start == 0 {
			pw++
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part one: %d\n", pw)
}

func mod100(n int) int {
	v := n % 100
	if v < 0 {
		v += 100
	}
	return v
}
