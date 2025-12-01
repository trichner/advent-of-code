package main

import (
	"fmt"
	"slices"
	"strings"

	"aoc/pkg/vec"
)

func printFieldWithPath(field [][]byte, path []vec.Vec2i) {
	for y, line := range field {
		fmt.Printf("%2d ", y)
		for x, c := range line {
			if slices.Index(path, vec.Vec2i{x, y}) >= 0 {
				fmt.Print("*")
			} else {
				fmt.Printf("%s", string([]byte{c}))
			}
		}
		fmt.Println()
	}
}

func printPath(loop []vec.Vec2i) {
	bb := vec.BoundingBox2i(loop)
	for y := 0; y <= bb.To.Y; y++ {
		fmt.Printf("%2d ", y)
		for x := 0; x <= bb.To.X; x++ {

			onPath := slices.Contains(loop, vec.Vec2i{x, y})
			if onPath {
				fmt.Print("*")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func printHeader(ident, max int) {
	fmt.Print(strings.Repeat(" ", ident))
	for x := 0; x <= max; x++ {
		d := (x / 100) % 10
		if d > 0 {
			fmt.Printf("%d", d)
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	fmt.Print(strings.Repeat(" ", ident))
	for x := 0; x <= max; x++ {
		d := (x / 10) % 10
		if d > 0 || x >= 100 {
			fmt.Printf("%d", d)
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	fmt.Print(strings.Repeat(" ", ident))
	for x := 0; x <= max; x++ {
		d := x % 10
		fmt.Printf("%d", d)
	}
	fmt.Println()
}

func printSegments(segments []*Segment) {
	aabb := boundingBox(segments)
	fmt.Printf("%v\n", aabb)

	printHeader(3, aabb.To.X)
	for y := 0; y <= aabb.To.Y; y++ {
		fmt.Printf("%2d ", y)
		for x := 0; x <= aabb.To.X; x++ {

			intersection := stabsAnySegment(segments, vec.Vec2i{x, y})

			if intersection == IntersectsInside {
				fmt.Print("|")
			} else if intersection == IntersectsEnd {
				fmt.Print("v")
			} else if intersection == IntersectsStart {
				fmt.Print("^")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
