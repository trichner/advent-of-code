package vec

import (
	"fmt"
)

type AABB struct{ From, To Vec2i }

func (v AABB) Contains(p Vec2i) bool {
	if p.X < v.From.X || p.X > v.To.X {
		return false
	}

	if p.Y < v.From.Y || p.Y > v.To.Y {
		return false
	}

	return true
}

func (v AABB) String() string {
	return fmt.Sprintf("AABB{%v, %v}", v.From, v.To)
}
