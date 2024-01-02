package vec

import (
	"fmt"
	"math"
)

type Vec3i struct{ X, Y, Z int }

func (v Vec3i) Add(summand Vec3i) Vec3i {
	return Vec3i{
		X: v.X + summand.X,
		Y: v.Y + summand.Y,
		Z: v.Z + summand.Z,
	}
}

func (v Vec3i) Sub(subtrahend Vec3i) Vec3i {
	return Vec3i{
		X: v.X - subtrahend.X,
		Y: v.Y - subtrahend.Y,
		Z: v.Z - subtrahend.Z,
	}
}

func (v Vec3i) Abs() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vec3i) Norm1() int {
	return abs(v.X) + abs(v.Y) + abs(v.Z)
}

func (v Vec3i) String() string {
	return fmt.Sprintf("{%d,%d,%d}", v.X, v.Y, v.Z)
}
