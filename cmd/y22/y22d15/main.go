package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

func main() {
	partOne()
	partTwo()
}

type Sensor struct {
	Position      Vec2i
	ClosestBeacon Vec2i
}

func partTwo() {
	sensors := parseInput("input.txt")
	maxSize := 4000000

	aabb := AABB{
		Origin: Vec2i{0, 0},
		Size:   Vec2i{maxSize, maxSize},
	}

	searchIntervalX := [2]int{aabb.Origin.X, aabb.Origin.X + aabb.Size.X}
	for y := aabb.Origin.Y; y < aabb.Origin.Y+aabb.Size.Y; y++ {
		freeInterval := scanLine(sensors, searchIntervalX, y)
		if len(freeInterval) > 0 {
			fmt.Printf("%d\n", tune(Vec2i{freeInterval[0][0], y}))
			return
		}
	}
}

func scanLine(sensors []Sensor, searchInterval [2]int, y int) [][2]int {
	var edges [][2]int

	edges = append(edges, [2]int{searchInterval[0], -1})
	edges = append(edges, [2]int{searchInterval[1], 1})

	for _, s := range sensors {
		ival := intersectSensor(s, y)
		if ival != nil {
			edges = append(edges, [2]int{ival[0].X, 1})
			edges = append(edges, [2]int{ival[1].X, -1})
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		a := edges[i]
		b := edges[j]
		return a[0]-b[0] < 0
	})

	layers := 1
	start := 0
	var intervals [][2]int
	for _, edge := range edges {
		prev := layers
		layers += edge[1]
		if prev > 0 && layers == 0 {
			start = edge[0]
		}
		if prev == 0 && layers > 0 {
			end := edge[0]
			if start != end {
				intervals = append(intervals, [2]int{start, end})
			}
		}
	}
	return intervals
}

func intersectSensor(s Sensor, y int) []Vec2i {
	r := s.Position.DistanceTo(s.ClosestBeacon)
	return intersectCircle(s.Position, r, y)
}

func intersectCircle(s Vec2i, r int, y int) []Vec2i {
	if s.Y-r > y {
		return nil
	}
	if s.Y+r < y {
		return nil
	}
	if s.Y+r == y || s.Y-r == y {
		return []Vec2i{{s.X, y}, {s.X + 1, y}}
	}

	a := r - absInt(y-s.Y)
	x1 := s.X - a
	x2 := s.X + a
	return []Vec2i{{x1, y}, {x2 + 1, y}}
}

func tune(p Vec2i) int {
	return p.X*4000000 + p.Y
}

func partOne() {
	sensors := parseInput("input.txt")

	aabb := NewAabbContaining(sensors)

	total := 0
	y := 2000000
	for x := aabb.Origin.X; x <= aabb.Origin.X+aabb.Size.X; x++ {
		hasSensorInRange := false
		isOccupied := false
		for _, s := range sensors {
			p := Vec2i{x, y}
			if p == s.ClosestBeacon {
				isOccupied = true
				break
			}
			if p == s.Position {
				isOccupied = true
				break
			}
			dist := s.Position.DistanceTo(p)
			closest := s.Position.DistanceTo(s.ClosestBeacon)
			if dist <= closest {
				hasSensorInRange = true
			}
		}
		if hasSensorInRange && !isOccupied {
			total++
		}
	}
	fmt.Printf("%d\n", total)
}

func parseInput(fname string) []Sensor {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var sensors []Sensor

	for scanner.Scan() {
		text := scanner.Text()

		var sx, sy, bx, by int
		_, err := fmt.Sscanf(text, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sx, &sy, &bx, &by)
		if err != nil {
			log.Fatal(err)
		}
		sensors = append(sensors, Sensor{
			Position:      Vec2i{sx, sy},
			ClosestBeacon: Vec2i{bx, by},
		})

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return sensors
}

func NewAabbContaining(vecs []Sensor) AABB {
	maxVec := Vec2i{math.MinInt, math.MinInt}
	minVec := Vec2i{math.MaxInt, math.MaxInt}
	for _, v := range vecs {
		d := v.Position.DistanceTo(v.ClosestBeacon)
		maxVec = max2i(maxVec, v.Position.Add(Vec2i{d, d}))
		minVec = min2i(minVec, v.Position.Sub(Vec2i{d, d}))

		maxVec = max2i(maxVec, v.ClosestBeacon)
		minVec = min2i(minVec, v.ClosestBeacon)
	}

	origin := minVec
	size := maxVec.Sub(minVec)

	return AABB{
		Origin: origin,
		Size:   size,
	}
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

type Vec2i struct{ X, Y int }

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

func (v Vec2i) Norm1() int {
	return absInt(v.X) + absInt(v.Y)
}

func (v Vec2i) DistanceTo(o Vec2i) int {
	return v.Sub(o).Norm1()
}

func (v Vec2i) String() string {
	return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

func area(r int) int {
	a := 0
	for i := 0; i < r; i++ {
		a += (i*2 + 1) * 2
	}

	a += 2*r + 1
	return a
}

func absInt(a int) int {
	if a >= 0 {
		return a
	}
	return -a
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
