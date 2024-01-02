package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	partOne()
	partTwo()
}

func partOne() {
	forest := readForest()
	isVisible := map[uint64]struct{}{}

	top := newBorder(len(forest))
	left := newBorder(len(forest[0]))

	for y := 0; y < len(forest); y++ {
		for x := 0; x < len(forest[0]); x++ {
			tree := forest[y][x]

			if top[x] < tree || left[y] < tree {
				isVisible[(uint64(y)<<32)|uint64(x)] = struct{}{}
			}

			top[x] = max(top[x], tree)
			left[y] = max(left[y], tree)
		}
	}

	bottom := newBorder(len(forest))
	right := newBorder(len(forest[0]))

	for y := len(forest) - 1; y >= 0; y-- {
		for x := len(forest[0]) - 1; x >= 0; x-- {
			tree := forest[y][x]

			if bottom[x] < tree || right[y] < tree {
				isVisible[(uint64(y)<<32)|uint64(x)] = struct{}{}
			}

			bottom[x] = max(bottom[x], tree)
			right[y] = max(right[y], tree)
		}
	}

	fmt.Printf("%d\n", len(isVisible))
}

func partTwo() {
	maxScore := 0

	forest := readForest()
	for y := 0; y < len(forest); y++ {
		for x := 0; x < len(forest[0]); x++ {
			tree := forest[y][x]
			score := 1

			distance := 0
			for dx := 1; dx+x < len(forest[0]); dx++ {
				if tree > forest[y][x+dx] {
					distance++
				} else {
					distance++
					break
				}
			}
			score *= distance

			distance = 0
			for dx := -1; dx+x >= 0; dx-- {
				if tree > forest[y][x+dx] {
					distance++
				} else {
					distance++
					break
				}
			}
			score *= distance

			distance = 0
			for dy := 1; dy+y < len(forest); dy++ {
				if tree > forest[y+dy][x] {
					distance++
				} else {
					distance++
					break
				}
			}
			score *= distance

			distance = 0
			for dy := -1; dy+y >= 0; dy-- {
				if tree > forest[y+dy][x] {
					distance++
				} else {
					distance++
					break
				}
			}
			score *= distance

			if score > maxScore {
				maxScore = score
			}
		}
	}

	fmt.Printf("%d", maxScore)
}

func readForest() [][]int8 {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var forest [][]int8
	var row []int8

	for {

		b, err := reader.ReadByte()
		if err == io.EOF {
			forest = append(forest, row)
			row = nil
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if b == '\n' {
			forest = append(forest, row)
			row = []int8{}
		} else if b >= '0' && b <= '9' {
			n := b - '0'
			row = append(row, int8(n))
		} else {
			log.Fatalf("unexpected char: %q", b)
		}
	}
	return forest
}

func max(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}

func newBorder(l int) []int8 {
	s := make([]int8, l)
	for i := range s {
		s[i] = -1
	}
	return s
}
