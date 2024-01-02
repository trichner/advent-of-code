package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	partOne()
	partTwo()
}

type Material int

const (
	AIR Material = iota
	ROCK
	SAND
	VOID
)

func (m Material) IsSolid() bool {
	return m == SAND || m == ROCK
}

func partOne() {
	paths := parsePaths("input.txt")
	cave := NewCavePartOne(paths)

	// fmt.Printf("%s\n", cave.String())
	dropped := 0
	for {
		if !dropSand(cave) {
			break
		}
		// fmt.Printf("%s\n", cave)
		dropped++
	}
	fmt.Printf("%d\n", dropped)
}

func partTwo() {
	paths := parsePaths("input.txt")
	cave := NewCavePartTwo(paths)
	dropped := 0
	for {
		if !dropSandPartTwo(cave) {
			break
		}
		dropped++
	}
	fmt.Printf("%d\n", dropped)
}

func parsePaths(fname string) [][]Vec2i {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	paths := [][]Vec2i{}
	for scanner.Scan() {

		text := scanner.Text()

		splits := strings.Split(text, " -> ")
		var path []Vec2i
		for _, s := range splits {
			var x, y int
			_, err := fmt.Sscanf(s, "%d,%d", &x, &y)
			if err != nil {
				panic(err)
			}

			path = append(path, Vec2i{x, y})
		}
		paths = append(paths, path)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return paths
}

type AABB struct {
	Origin, Size Vec2i
}

func (a *AABB) InBounds(p Vec2i) bool {
	if p.X < a.Origin.X || p.X > a.Origin.X+a.Size.X {
		return false
	}
	if p.Y < a.Origin.Y || p.Y > a.Origin.Y+a.Size.Y {
		return false
	}
	return true
}

type Cave struct {
	aabb AABB
	cave [][]Material
}

func (c *Cave) Get(p Vec2i) Material {
	if !c.aabb.InBounds(p) {
		return VOID
	}
	adjusted := p.Sub(c.aabb.Origin)
	return c.cave[adjusted.Y][adjusted.X]
}

func (c *Cave) Set(p Vec2i, m Material) {
	if !c.aabb.InBounds(p) {
		panic("out of bounds")
	}
	adjusted := p.Sub(c.aabb.Origin)
	c.cave[adjusted.Y][adjusted.X] = m
}

func (c *Cave) DrawRocks(start, end Vec2i) {
	if !c.aabb.InBounds(start) || !c.aabb.InBounds(end) {
		panic(fmt.Errorf("vectors out of bounds: %s -> %s", start, end))
	}

	if start.X != end.X {
		b := start.X
		e := end.X
		if start.X > end.X {
			b, e = e, b
		}
		for dx := b; dx <= e; dx++ {
			c.Set(Vec2i{dx, start.Y}, ROCK)
		}
	} else if start.Y != end.Y {
		b := start.Y
		e := end.Y
		if start.Y > end.Y {
			b, e = e, b
		}
		for dy := b; dy <= e; dy++ {
			c.Set(Vec2i{start.X, dy}, ROCK)
		}
	} else {
		panic("invalid path")
	}
}

func (c *Cave) String() string {
	var buf strings.Builder
	for _, row := range c.cave {
		for _, p := range row {
			switch p {
			case AIR:
				buf.WriteRune('.')
			case ROCK:
				buf.WriteRune('#')
			case SAND:
				buf.WriteRune('O')
			}
		}
		buf.WriteRune('\n')
	}

	return buf.String()
}

func NewCavePartOne(paths [][]Vec2i) *Cave {
	maxVec := Vec2i{math.MinInt, math.MinInt}
	minVec := Vec2i{math.MaxInt, 0}
	for _, path := range paths {
		for _, point := range path {
			maxVec = max2i(maxVec, point)
			minVec = min2i(minVec, point)
		}
	}

	origin := minVec
	size := maxVec.Sub(minVec).Add(Vec2i{1, 1})

	caveMaterial := make([][]Material, size.Y)
	for y := 0; y < len(caveMaterial); y++ {
		caveMaterial[y] = make([]Material, size.X)
	}
	cave := &Cave{
		aabb: AABB{
			Origin: origin,
			Size:   size,
		},
		cave: caveMaterial,
	}
	for _, path := range paths {
		start := path[0]
		for i := 1; i < len(path); i++ {
			end := path[i]
			cave.DrawRocks(start, end)
			start = end
		}
	}
	return cave
}

func NewCavePartTwo(paths [][]Vec2i) *Cave {
	maxVec := Vec2i{math.MinInt, math.MinInt}
	minVec := Vec2i{math.MaxInt, 0}
	for _, path := range paths {
		for _, point := range path {
			maxVec = max2i(maxVec, point)
			minVec = min2i(minVec, point)
		}
	}

	origin := Vec2i{0, 0}
	// extra layer
	size := maxVec.Sub(origin).Add(Vec2i{1, 1 + 1})
	size = size.Add(Vec2i{size.Y, 0})

	// add some extra X

	caveMaterial := make([][]Material, size.Y)
	for y := 0; y < len(caveMaterial); y++ {
		caveMaterial[y] = make([]Material, size.X)
	}
	cave := &Cave{
		aabb: AABB{
			Origin: origin,
			Size:   size,
		},
		cave: caveMaterial,
	}
	for _, path := range paths {
		start := path[0]
		for i := 1; i < len(path); i++ {
			end := path[i]
			cave.DrawRocks(start, end)
			start = end
		}
	}
	return cave
}

func dropSand(cave *Cave) bool {
	pos := Vec2i{500, cave.aabb.Origin.Y}
	m := cave.Get(pos)
	if m != AIR {
		panic("can't drop sand")
	}
	for {
		if cave.Get(pos) == VOID {
			return false
		}

		if !cave.Get(pos.Add(Vec2i{0, 1})).IsSolid() {
			pos = pos.Add(Vec2i{0, 1})
		} else if !cave.Get(pos.Add(Vec2i{-1, 1})).IsSolid() {
			pos = pos.Add(Vec2i{-1, 1})
		} else if !cave.Get(pos.Add(Vec2i{1, 1})).IsSolid() {
			pos = pos.Add(Vec2i{1, 1})
		} else {
			cave.Set(pos, SAND)
			return true
		}
	}
}

func dropSandPartTwo(cave *Cave) bool {
	pos := Vec2i{500, cave.aabb.Origin.Y}
	m := cave.Get(pos)
	if m == SAND {
		return false
	}
	floorY := cave.aabb.Origin.Y + cave.aabb.Size.Y
	for {
		straightDown := pos.Add(Vec2i{0, 1})
		if straightDown.Y >= floorY {
			cave.Set(pos, SAND)
			return true
		}

		if !cave.Get(pos.Add(Vec2i{0, 1})).IsSolid() {
			pos = pos.Add(Vec2i{0, 1})
		} else if !cave.Get(pos.Add(Vec2i{-1, 1})).IsSolid() {
			pos = pos.Add(Vec2i{-1, 1})
		} else if !cave.Get(pos.Add(Vec2i{1, 1})).IsSolid() {
			pos = pos.Add(Vec2i{1, 1})
		} else {
			cave.Set(pos, SAND)
			return true
		}
	}
}

func max2i(a, b Vec2i) Vec2i {
	return Vec2i{
		X: max(a.X, b.X),
		Y: max(a.Y, b.Y),
	}
}

func min2i(a, b Vec2i) Vec2i {
	return Vec2i{
		X: min(a.X, b.X),
		Y: min(a.Y, b.Y),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Vec2i struct {
	X, Y int
}

func (v Vec2i) String() string {
	return fmt.Sprintf("(%d/%d)", v.X, v.Y)
}

func (v Vec2i) Add(summand Vec2i) Vec2i {
	return Vec2i{X: v.X + summand.X, Y: v.Y + summand.Y}
}

func (v Vec2i) Sub(subtrahend Vec2i) Vec2i {
	return Vec2i{X: v.X - subtrahend.X, Y: v.Y - subtrahend.Y}
}
