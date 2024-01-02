package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	ropeWalk(2)
	ropeWalk(10)
}

type Vec2i struct{ X, Y int }

func ropeWalk(k int) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	knots := make([]Vec2i, k)

	visited := map[string]struct{}{}
	visited[knots[len(knots)-1].String()] = struct{}{}

	for scanner.Scan() {
		text := scanner.Text()

		var v string
		var n int
		_, err := fmt.Sscanf(text, "%s %d", &v, &n)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < n; i++ {
			knots[0] = knots[0].Move(v)
			for j := 1; j < len(knots); j++ {
				knots[j] = CatchUp(knots[j-1], knots[j])
			}
			visited[knots[len(knots)-1].String()] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d\n", len(visited))
}

func CatchUp(head, tail Vec2i) Vec2i {
	diff := head.Sub(tail)
	if diff.Abs() < 2 {
		return tail
	}

	return tail.Add(Vec2i{sign(diff.X), sign(diff.Y)})
}

func (v Vec2i) Move(dir string) Vec2i {
	switch dir {
	case "R":
		return Vec2i{v.X + 1, v.Y}
	case "L":
		return Vec2i{v.X - 1, v.Y}
	case "U":
		return Vec2i{v.X, v.Y + 1}
	case "D":
		return Vec2i{v.X, v.Y - 1}
	}
	panic(fmt.Errorf("unexpected direction: %q", dir))
}

func (v Vec2i) Add(summand Vec2i) Vec2i {
	return Vec2i{
		X: v.X + summand.X,
		Y: v.Y + summand.Y,
	}
}

func (v Vec2i) Sub(subtrahend Vec2i) Vec2i {
	return Vec2i{
		X: v.X - subtrahend.X,
		Y: v.Y - subtrahend.Y,
	}
}

func (v Vec2i) Abs() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v Vec2i) String() string {
	return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

func sign(x int) int {
	if x == 0 {
		return 0
	}
	return 1 | (x >> (strconv.IntSize - 1))
}
