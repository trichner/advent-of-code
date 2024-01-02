package vec

import (
	"cmp"
	"fmt"
	"math"
	"slices"
)

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
	return abs(v.X) + abs(v.Y)
}

func abs(i int) int {
	if i < 0 {
		i = -i
	}
	return i
}

func (v Vec2i) String() string {
	return fmt.Sprintf("{%d,%d}", v.X, v.Y)
}

func Compare2i(a, b Vec2i) int {
	d := cmp.Compare(a.X, b.X)
	if d != 0 {
		return d
	}
	return cmp.Compare(a.Y, b.Y)
}

func BoundingBox2i(points []Vec2i) AABB {
	cmpX := func(a, b Vec2i) int {
		return cmp.Compare(a.X, b.X)
	}

	minX := slices.MinFunc(points, cmpX).X
	maxX := slices.MaxFunc(points, cmpX).X

	cmpY := func(a, b Vec2i) int {
		return cmp.Compare(a.Y, b.Y)
	}

	minY := slices.MinFunc(points, cmpY).Y
	maxY := slices.MaxFunc(points, cmpY).Y

	return AABB{From: Vec2i{minX, minY}, To: Vec2i{maxX, maxY}}
}
