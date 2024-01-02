package main

import (
	"fmt"

	"aoc/pkg/vec"
)

func partOne(loop []vec.Vec2i) {
	steps := (len(loop) + 1) / 2
	fmt.Printf("part one: %v\n", steps)
	if steps != 7005 {
		panic("bad result")
	}
}
