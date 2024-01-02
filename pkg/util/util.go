package util

import (
	"slices"
)

func Must[K any](v K, err error) K {
	if err != nil {
		panic(err)
	}
	return v
}

func MustOk[K any](v K, ok bool) K {
	if !ok {
		panic("not ok")
	}
	return v
}

func Copy[S [][]E, E comparable](s S) S {
	copied := make(S, len(s))
	for y, r := range s {
		copied[y] = make([]E, len(r))
		copy(copied[y], r)
	}
	return copied
}

func Equal[S [][]E, E comparable](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}

	for y := range s1 {
		if !slices.Equal(s1[y], s2[y]) {
			return false
		}
	}
	return true
}
