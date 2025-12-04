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
	field := readInput()
	newField := copyField(field)

	sum := 0

	for {
		removed := 0
		for i, row := range field {
			for j, e := range row {
				if e != '@' {
					continue
				}

				count := 0
				for dx := -1; dx <= 1; dx++ {
					for dy := -1; dy <= 1; dy++ {
						if dx == 0 && dy == 0 {
							continue
						}
						x := i + dx
						y := j + dy
						if x < 0 || y < 0 || x >= len(field) || y >= len(field[0]) {
							continue
						}
						if field[x][y] == '@' {
							count++
						}
					}
				}
				if count < 4 {
					newField[i][j] = '.'
					removed++
					sum++
				}

			}
		}
		if removed == 0 {
			break
		}
		field = newField
		newField = copyField(field)
	}
	fmt.Printf("part two: %d\n", sum)
}

func copyField(field [][]byte) [][]byte {
	newField := make([][]byte, len(field))
	for i, row := range field {
		newField[i] = make([]byte, len(row))
		copy(newField[i], row)
	}
	return newField
}

func partOne() {
	field := readInput()

	sum := 0
	for i, row := range field {
		for j, e := range row {
			if e != '@' {
				continue
			}

			count := 0
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					x := i + dx
					y := j + dy
					if x < 0 || y < 0 || x >= len(field) || y >= len(field[0]) {
						continue
					}
					if field[x][y] == '@' {
						count++
					}
				}
			}
			if count < 4 {
				sum++
			}

		}
	}

	// 1474
	fmt.Printf("part one: %d\n", sum)
}

func readInput() [][]byte {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var field [][]byte
	for scanner.Scan() {
		text := scanner.Text()

		field = append(field, []byte(text))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return field
}
