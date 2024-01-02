package main

import (
	"bufio"
	"embed"
	"fmt"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

func main() {
	partOne()
	partTwo()
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
