package main

import (
	"bufio"
	"embed"
	"fmt"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

var spelled = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

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

		first := -1
		last := -1

		line := []byte(text)
		for i := 0; i < len(line); i++ {
			sub := line[i:]
			d := findDigit(sub)
			if d >= 0 {
				first = d
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			sub := line[i:]
			d := findDigit(sub)
			if d >= 0 {
				last = d
				break
			}
		}
		sum += first*10 + last
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("%d\n", sum)
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0

	for scanner.Scan() {
		text := scanner.Text()

		first := -1
		last := -1

		line := []byte(text)
		for i := 0; i < len(line); i++ {
			c := line[i]
			if c >= '0' && c <= '9' {
				first = int(c - '0')
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			c := line[i]
			if c >= '0' && c <= '9' {
				last = int(c - '0')
				break
			}
		}
		sum += first*10 + last
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("%d\n", sum)
}

var lookup = func() map[string]int {
	lookup := map[string]int{}
	for c := '0'; c <= '9'; c++ {
		lookup[string(c)] = int(c - '0')
	}
	for i, s := range spelled {
		lookup[s] = i
	}
	return lookup
}()

func findDigit(s []byte) int {
	for i := 0; i < min(5, len(s)); i++ {
		sub := string(s[:(i + 1)])
		d, ok := lookup[sub]
		if !ok {
			continue
		}
		return d
	}
	return -1
}
