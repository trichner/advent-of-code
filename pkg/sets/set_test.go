package sets

import (
	"fmt"
	"testing"

	"aoc/pkg/be"
)

func TestPut(t *testing.T) {
	set := New[string]()

	set.Put("a")
	set.Put("b")

	be.Equal(t, set.Has("a"), true)
}

func TestUnion(t *testing.T) {
	setA := New[string]()
	setA.Put("a")
	setA.Put("b")

	setB := New[string]()
	setB.Put("b")
	setB.Put("c")

	u := Intersect(setA, setB)

	be.Equal(t, u.Has("b"), true)
	be.Equal(t, len(u), 1)
}

func TestString(t *testing.T) {
	set := New[string]()

	set.Put("a")
	set.Put("b")

	fmt.Println(set)
}
