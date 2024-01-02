package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	partOne()
	partTwo()
}

func partOne() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	max := 0
	current := 0

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			if current > max {
				max = current
			}
			current = 0
			continue
		}
		calories, err := strconv.Atoi(text)
		if err != nil {
			log.Fatal(err)
		}
		current += calories
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d\n", max)
}

func partTwo() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	current := 0
	var elves []int

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			elves = append(elves, current)
			current = 0
			continue
		}
		calories, err := strconv.Atoi(text)
		if err != nil {
			log.Fatal(err)
		}
		current += calories
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(elves)))
	sum := elves[0] + elves[1] + elves[2]

	fmt.Printf("%d\n", sum)
}
