package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	total := 0

	for scanner.Scan() {
		text := scanner.Text()

		var e1a, e1b, e2a, e2b int
		_, err := fmt.Sscanf(text, "%d-%d,%d-%d", &e1a, &e1b, &e2a, &e2b)
		if err != nil {
			log.Fatal(err)
		}
		if isSubRange([2]int{e1a, e1b}, [2]int{e2a, e2b}) {
			total += 1
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%d\n", total)
}

func partTwo() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	total := 0

	for scanner.Scan() {
		text := scanner.Text()

		var e1a, e1b, e2a, e2b int
		_, err := fmt.Sscanf(text, "%d-%d,%d-%d", &e1a, &e1b, &e2a, &e2b)
		if err != nil {
			log.Fatal(err)
		}
		if isIntersecting([2]int{e1a, e1b}, [2]int{e2a, e2b}) {
			total += 1
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%d\n", total)
}

func isIntersecting(i1 [2]int, i2 [2]int) bool {
	// fully intersecting or intersecting start or intersecting end
	return isSubRange(i1, i2) || halfOverlapsStart(i1, i2) || halfOverlapsStart(i2, i1)
}

func isSubRange(i1 [2]int, i2 [2]int) bool {
	return contained(i1, i2) || contained(i2, i1)
}

func contained(i1 [2]int, i2 [2]int) bool {
	return stab(i1, i2[0]) && stab(i1, i2[1])
}

func halfOverlapsStart(i1 [2]int, i2 [2]int) bool {
	return stab(i1, i2[0]) && (!stab(i1, i2[1]))
}

func stab(haystack [2]int, needle int) bool {
	return haystack[0] <= needle && haystack[1] >= needle
}
