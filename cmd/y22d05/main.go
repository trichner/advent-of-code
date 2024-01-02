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

	stacks := [9][]byte{{}, {}, {}, {}, {}, {}, {}, {}, {}}

	maxHeight := 8
	for i := 0; i < maxHeight; i++ {
		if !scanner.Scan() {
			panic("can't scan")
		}
		l := scanner.Text()
		crates := parseLine(l)
		for i, e := range crates {
			if e == 0 {
				continue
			}
			stacks[i] = append([]byte{e}, stacks[i]...)

		}
	}

	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		text := scanner.Text()

		var n, from, to int
		_, err := fmt.Sscanf(text, "move %d from %d to %d", &n, &from, &to)
		if err != nil {
			panic(err)
		}
		from -= 1
		to -= 1
		moveCrates(&stacks[from], &stacks[to], n)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, s := range stacks {
		fmt.Printf("%c", s[len(s)-1])
	}
	fmt.Println()
}

func partTwo() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	stacks := [9][]byte{{}, {}, {}, {}, {}, {}, {}, {}, {}}

	maxHeight := 8
	for i := 0; i < maxHeight; i++ {
		if !scanner.Scan() {
			panic("can't scan")
		}
		l := scanner.Text()
		crates := parseLine(l)
		for i, e := range crates {
			if e == 0 {
				continue
			}
			stacks[i] = append([]byte{e}, stacks[i]...)

		}
	}

	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		text := scanner.Text()

		var n, from, to int
		_, err := fmt.Sscanf(text, "move %d from %d to %d", &n, &from, &to)
		if err != nil {
			panic(err)
		}
		from -= 1
		to -= 1
		moveCratesOver9000(&stacks[from], &stacks[to], n)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, s := range stacks {
		fmt.Printf("%c", s[len(s)-1])
	}
	fmt.Println()
}

func moveCratesOver9000(from *[]byte, to *[]byte, n int) {
	crates := (*from)[len(*from)-n:]
	*from = (*from)[:(len(*from) - n)]

	*to = append(*to, crates...)
}

func moveCrates(from *[]byte, to *[]byte, n int) {
	for i := 0; i < n; i++ {
		crate := (*from)[len(*from)-1]
		*from = (*from)[:(len(*from) - 1)]

		*to = append(*to, crate)
	}
}

func parseLine(s string) []byte {
	var packedCrate string
	var crates []byte
	for {
		if len(s) < 3 {
			return crates
		} else if len(s) == 3 {
			packedCrate = s
			s = ""
		} else {
			packedCrate = s[:4]
			s = s[4:]
		}

		if packedCrate[0] == ' ' {
			crates = append(crates, 0)
			continue
		}

		var contents byte
		_, err := fmt.Sscanf(packedCrate, "[%c] ", &contents)
		if err != nil {
			panic(err)
		}
		crates = append(crates, contents)
	}
}
