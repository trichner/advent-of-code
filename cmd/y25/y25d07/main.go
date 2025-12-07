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

type manifold struct {
	next []*manifold
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var manifolds [][]*manifold

	start := &manifold{}

	for scanner.Scan() {
		line := scanner.Bytes()
		if manifolds == nil {
			manifolds = make([][]*manifold, len(line))
			for i, c := range line {
				if c == 'S' {
					manifolds[i] = append(manifolds[i], start)
				}
			}
			continue
		}

		for i, b := range line {
			if b == '^' && manifolds[i] != nil {
				prev := manifolds[i]
				manifolds[i] = nil

				next := &manifold{}
				for _, p := range prev {
					p.next = append(p.next, next)
				}
				manifolds[i-1] = append(manifolds[i-1], next)
				manifolds[i+1] = append(manifolds[i+1], next)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	memo := make(map[*manifold]int)
	fmt.Printf("part two: %d\n", countPaths(start, memo)+1)
}

func countPaths(start *manifold, memo map[*manifold]int) int {
	if cached, ok := memo[start]; ok {
		return cached
	}

	sum := 0
	for _, m := range start.next {
		sum += countPaths(m, memo) + 1
	}

	memo[start] = sum
	return sum
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var splits int

	var beams []bool
	for scanner.Scan() {
		line := scanner.Bytes()
		if beams == nil {
			beams = make([]bool, len(line))
			for i, c := range line {
				if c == 'S' {
					beams[i] = true
				}
			}
			continue
		}

		for i, b := range line {
			if b == '^' && beams[i] {
				splits++
				beams[i-1] = true
				beams[i] = false
				beams[i+1] = true
			}
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part one: %d\n", splits)
}
