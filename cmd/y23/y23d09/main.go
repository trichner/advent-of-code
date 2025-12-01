package main

import (
	"bufio"
	"embed"
	"fmt"

	"aoc/pkg/sio"

	"aoc/pkg/in"
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
	sum := 0
	for scanner.Scan() {
		text := scanner.Text()

		fields := sio.IntFields(text)

		stack := stackUp(fields)
		p := predictFirst(stack)
		sum += p

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part two: %d\n", sum)
	if sum != 977 {
		panic("bad result")
	}
}

func predictFirst(stack [][]int) int {
	e := 0
	for _, row := range stack {
		first := row[0]
		e = first - e
	}
	return e
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0
	for scanner.Scan() {
		text := scanner.Text()

		fields := sio.IntFields(text)

		stack := stackUp(fields)
		p := predictLast(stack)
		sum += p

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part one: %d\n", sum)
	if sum != 1980437560 {
		panic("bad result")
	}
}

func predictLast(stack [][]int) int {
	e := 0
	for _, row := range stack {
		last := row[len(row)-1]
		e += last
	}
	return e
}

func printStack(stack [][]int) {
	indent := ""
	for i := len(stack) - 1; i >= 0; i-- {
		fmt.Print(indent)
		indent += " "
		l := stack[i]
		for _, e := range l {
			fmt.Printf("%2d ", e)
		}
		fmt.Println()
	}
}

func stackUp(fields []int) [][]int {
	stack := goDeeper(fields)
	stack = append(stack, fields)
	return stack
}

func goDeeper(fields []int) [][]int {
	if len(fields) == 1 {
		panic("wtf?")
	}

	var next []int
	for i := 0; i < len(fields)-1; i++ {
		v := fields[i+1] - fields[i]
		next = append(next, v)
	}

	for _, v := range next {
		if v != 0 {
			deeper := goDeeper(next)
			return append(deeper, next)
		}
	}
	return [][]int{next}
}
