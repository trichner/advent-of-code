package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/bits"
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

		l := len(text)
		s1 := text[:(l / 2)]
		s2 := text[(l / 2):]
		total += findCommonPriority(s1, s2)
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
		r1 := scanner.Text()
		if !scanner.Scan() {
			panic(fmt.Errorf("expected more"))
		}
		r2 := scanner.Text()
		if !scanner.Scan() {
			panic(fmt.Errorf("expected more"))
		}
		r3 := scanner.Text()

		p := findCommonPriority(r1, r2, r3)
		total += p
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%d\n", total)
}

func seenField(s string) uint64 {
	var bitsSet uint64
	for _, b := range []byte(s) {
		bitsSet |= 1 << id(b)
	}
	return bitsSet
}

func findCommonPriority(lines ...string) int {
	var all uint64 = math.MaxUint64
	for _, s := range lines {
		all &= seenField(s)
	}

	return bits.TrailingZeros64(all) + 1
}

func id(c byte) int {
	if c >= 'a' && c <= 'z' {
		return int(c) - int('a')
	}
	if c >= 'A' && c <= 'Z' {
		return int(c) - int('A') + 26
	}
	panic(fmt.Errorf("invalid item: %d", c))
}
