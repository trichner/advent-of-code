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
	file := in.MustOpenExampleTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()

		_ = text
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part two: %d\n", 0)
}

func partOne() {
	file := in.MustOpenExampleTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()

		_ = text
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part one: %d\n", 0)
}
